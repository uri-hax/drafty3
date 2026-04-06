#!/bin/bash
set -e
### fails fast - so if something breaks
### in the middle we do not hose oursleves

echo "STARTING Drafty Backend Build/Deployment..."

SCRIPT_DIR=$(pwd)
BACKEND_DIR="/vol/drafty3"

BACKEND_SERVICE="drafty-backend.service"

echo "...Stopping existing services (if running)"
systemctl stop drafty-backend || true

echo "...building go backend"
cd $BACKEND_DIR/endpoints
go build -o drafty-backend
mv drafty-backend ../bin/

echo "...installing systemd service files from cwd"
cp "$SCRIPT_DIR/$BACKEND_SERVICE" /etc/systemd/system/

echo "...reloading systemd"
systemctl daemon-reload

echo "...enabling services"
systemctl enable drafty-backend

echo "...starting services"
systemctl start drafty-backend

echo "...status"
systemctl status drafty-backend --no-pager

echo "::: Complete! :) :::"