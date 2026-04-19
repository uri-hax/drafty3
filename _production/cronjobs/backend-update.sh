#!/bin/bash

# exit immediately if any command fails
set -e

# set the repo directory
REPO_DIR="/vol/drafty3"
BACKEND_DIR="$REPO_DIR/backend"

# make sure we're in the correct directory
cd "$REPO_DIR"

# abort if repo has local changes so we don't overwrite work in progress
if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Repo has uncommitted changes. Aborting."
  exit 1
fi

# read latest commit hash
BEFORE_COMMIT="$(git rev-parse HEAD)"

# pull in latest changes from remote
git pull origin main

# read latest commit hash again
AFTER_COMMIT="$(git rev-parse HEAD)"

# if hashes match then no new changes so exit
if [ "$BEFORE_COMMIT" = "$AFTER_COMMIT" ]; then
  echo "No source changes pulled. Exiting."
  exit 0
fi

# if hashes don't match then check if there are changes in the backend endpoints directory 
# if no changes in the endpoints directory then no deploy needed so exit
if git diff --quiet "$BEFORE_COMMIT" "$AFTER_COMMIT" -- backend/endpoints/; then
  echo "No endpoint changes pulled. Exiting."
  exit 0
fi

# make sure temp directory exists
mkdir -p "$REPO_DIR/tmp"

# temporary location for the newly built binary
TEMP_BIN="$REPO_DIR/tmp/drafty-backend.new"

# script that actually swaps in the new binary and restarts the service
DEPLOY_SCRIPT="$REPO_DIR/_production/deploy.sh"

# build the new binary in a temp location
cd "$BACKEND_DIR/endpoints"
go build -o "$TEMP_BIN"

# if no match then run deploy and still remove the temp binary after 
"$DEPLOY_SCRIPT"
rm -f "$TEMP_BIN"
echo "Deployment successful"