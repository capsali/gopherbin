package controllers

import (
	"encoding/json"
	adminCommon "gopherbin/admin/common"
	"gopherbin/auth"
	gErrors "gopherbin/errors"
	"gopherbin/paste/common"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// NewAPIController returns a new APIController
func NewAPIController(paster common.Paster, mgr adminCommon.UserManager) *APIController {
	return &APIController{
		paster:  paster,
		manager: mgr,
	}
}

// APIController implements handlers for the REST API
type APIController struct {
	paster  common.Paster
	manager adminCommon.UserManager
}

func handleError(w http.ResponseWriter, err error) {
	apiErr := APIErrorResponse{
		Details: err.Error(),
	}
	switch errors.Cause(err) {
	case gErrors.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		apiErr.Error = "Not Found"
	case gErrors.ErrUnauthorized:
		w.WriteHeader(http.StatusUnauthorized)
		apiErr.Error = "Not Authorized"
	default:
		w.WriteHeader(http.StatusInternalServerError)
		apiErr.Error = "Server error"
	}
	json.NewEncoder(w).Encode(apiErr)
	return
}

// PasteViewHandler returns details about a single paste
func (p *APIController) PasteViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	pasteID, ok := vars["pasteID"]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pasteInfo, err := p.paster.Get(ctx, pasteID)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(pasteInfo)
}

// PasteListHandler returns a list of pastes
func (p *APIController) PasteListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := r.URL.Query().Get("page")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	maxResultsOpt := r.URL.Query().Get("max_results")
	maxResults, _ := strconv.ParseInt(maxResultsOpt, 10, 64)
	if maxResults == 0 {
		maxResults = 50
	}

	res, err := p.paster.List(ctx, pageInt, maxResults)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(res)
}

// DeletePasteHandler deletes a single paste
func (p *APIController) DeletePasteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	pasteID, ok := vars["pasteID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIErrorResponse{
			Error:   "Bad Request",
			Details: "No paste ID specified",
		})
		return
	}
	if err := p.paster.Delete(ctx, pasteID); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UserListHandler handles the list of pastes
func (p *APIController) UserListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if auth.IsSuperUser(ctx) == false && auth.IsAdmin(ctx) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(unauthorizedResponse)
		return
	}

	page := r.URL.Query().Get("page")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	maxResultsOpt := r.URL.Query().Get("max_results")
	maxResults, _ := strconv.ParseInt(maxResultsOpt, 10, 64)
	if maxResults == 0 {
		maxResults = 50
	}

	res, err := p.manager.List(ctx, pageInt, maxResults)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(res)
	return
}

// NotFoundHandler is returned when an invalid URL is acccessed
func (p *APIController) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(notFoundResponse)
}
