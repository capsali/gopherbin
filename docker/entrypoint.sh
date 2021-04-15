#!/bin/bash

set -e

gopher_dir="/etc/gopherbin/"
gopher_conf="${gopher_dir}/gopherbin-config.toml"
gopher_tmpl="/templates/gopherbin-config.toml.tmpl"

# Declare default values for variables if they are empty
export API_PORT=${API_PORT:-"9997"}
GEN_SESSION_SECRET=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 64 | head -n 1)
if [ ! -f "/opt/api_session_secret" ]; then
	echo $GEN_SESSION_SECRET | tee /secrets/api_session_secret
fi
RANDOM_SESSION_SECRET=$(cat /secrets/api_session_secret)
export API_SESSION_SECRET=${API_SESSION_SECRET:-"$RANDOM_SESSION_SECRET"}
GEN_JWT_SECRET=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 64 | head -n 1)
if [ ! -f "/opt/api_jwt_secret" ]; then
        echo $GEN_JWT_SECRET | tee /secrets/api_jwt_secret
fi
RANDOM_JWT_SECRET=$(cat /secrets/api_jwt_secret)
export API_JWT_SECRET=${API_JWT_SECRET:-"$RANDOM_JWT_SECRET"}
export DB_BACKEND=${DB_BACKEND:-"mysql"}
export DB_USER=${DB_USER:-"gopherbin"}
GEN_DB_PASSWORD=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
if [ ! -f "/opt/db_password" ]; then
        echo $GEN_DB_PASSWORD | tee /secrets/db_password
fi
RANDOM_DB_PASSWORD=$(cat /secrets/db_password)
export DB_PASSWORD=${DB_PASSWORD:-"$RANDOM_DB_PASSWORD"}
export DB_HOST=${DB_HOST:-"127.0.0.1"}
export DB_NAME=${DB_NAME:-"gopherbin"}

# TO DO: Include db server with the image

# Create gopher user
export PUID=${PUID:-"10001"}
export PGID=${PGID:-"10001"}
addgroup -g $PGID -S gopher && adduser -u $PUID -S gopher -G gopher

# Check if config file exists. If not, render it.
if [ -f "/etc/gopherbin/gopherbin-config.toml" ]; then
	echo "Config file already exists"
else
	envsubst < $gopher_tmpl > $gopher_conf
fi
chown -R gopher:gopher $gopher_dir

exec su-exec gopher:gopher "$@"

