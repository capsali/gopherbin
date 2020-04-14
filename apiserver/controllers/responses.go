// Copyright 2019 Gabriel-Adrian Samfira
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package controllers

// ErrorResponse holds any errors generated during
// a request
type ErrorResponse struct {
	Errors map[string]string
}

// APIErrorResponse holds information about an error, returned by the API
type APIErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

var (
	notFoundResponse = APIErrorResponse{
		Error:   "Not Found",
		Details: "The resource you are looking for was not found",
	}
	unauthorizedResponse = APIErrorResponse{
		Error:   "Not Authorized",
		Details: "You do not have the required permissions to access this resource",
	}
)
