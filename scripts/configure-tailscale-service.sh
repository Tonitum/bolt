#!/bin/bash
set -euo pipefail

sudo tailscale serve --service=svc:blt --http=80 127.0.0.1:8081
sudo tailscale serve --service=svc:blt --https=443 127.0.0.1:8080
