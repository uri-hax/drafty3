#!/bin/bash

# exit immediately if any command fails
set -e

# root repo directory
REPO_DIR="/vol/drafty3"
BACKEND_DIR="$REPO_DIR/backend"

# make sure we're in the correct directory
cd "$REPO_DIR"

# abort if repo has local changes so we don't overwrite work in progress
if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Repo has uncommitted changes. Aborting."
  exit 1
fi

# make sure temp directory exists
mkdir -p "$REPO_DIR/tmp"

# sqlite database file (csprofs)
DB_FILE="$REPO_DIR/backend/db/drafty_new_gorm.db"

# temporary generated csv
TEMP_CSV="$REPO_DIR/tmp/generated.csv"

# csv file tracked in git that the frontend uses (csprofs)
CURRENT_CSV="$REPO_DIR/public/suggestions.csv"

# build a fresh csv from the sqlite database (csprofs)
echo "Generating CSV..."
cd "$BACKEND_DIR"
go run "$BACKEND_DIR/csv/build_csv.go" --db "$DB_FILE" --out "$TEMP_CSV" --csv_type "csprofs"

# make sure the csv was actually created
if [ ! -f "$TEMP_CSV" ]; then
  echo "CSV was not created. Aborting."
  exit 1
fi

# make sure current csv exists and if generated csv matches the current tracked csv then do nothing
if [ -f "$CURRENT_CSV" ] && cmp -s "$TEMP_CSV" "$CURRENT_CSV"; then
  echo "CSV unchanged. Exiting."
  rm -f "$TEMP_CSV"
  exit 0
fi

# overwrite tracked csv with the new one
cp "$TEMP_CSV" "$CURRENT_CSV"

# clean up temp file
rm -f "$TEMP_CSV"

# stage the updated csv
git add "$CURRENT_CSV"

# commit and push so github updates the frontend
git commit -m "Update generated CSV"
git push origin main

echo "Frontend CSV updated successfully."