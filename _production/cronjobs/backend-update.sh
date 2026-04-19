#!/bin/bash

# exit immediately if any command fails
set -e

# set the repo directory
REPO_DIR="/c/Users/leach/Documents/ResearchProject/drafty3"

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
# git pull origin dev

# read latest commit hash again
AFTER_COMMIT="$(git rev-parse HEAD)"

# if hashes match then no new changes so exit
# if [ "$BEFORE_COMMIT" = "$AFTER_COMMIT" ]; then
#   echo "No source changes pulled. Exiting."
#   exit 0
# fi

# make sure temp directory exists
mkdir -p "$REPO_DIR/tmp"

# path to the currently deployed backend binary
CURRENT_BIN="$REPO_DIR/bin/drafty-backend"

# temporary location for the newly built binary
TEMP_BIN="$REPO_DIR/tmp/drafty-backend.new"

# script that actually swaps in the new binary and restarts the service
DEPLOY_SCRIPT="$REPO_DIR/_production/deploy.sh"

# if hashes no match then build the new binary in a temp location
cd "$REPO_DIR/backend/endpoints"
go build -o "$TEMP_BIN"

# make sure current binary exists and compare silently with the new one
# if they match then no deploy is needed so clean up temp binary and exit
if [ -f "$CURRENT_BIN" ] && cmp -s "$TEMP_BIN" "$CURRENT_BIN"; then
  echo "Built binary matches deployed binary. No deploy needed."
  rm -f "$TEMP_BIN"
  exit 0
fi

# if no match then run deploy and still remove the temp binary after 
#"$DEPLOY_SCRIPT"
# rm -f "$TEMP_BIN"
echo "Deployment successful"