#!/bin/bash
set -e

DROPLET="deploy@167.172.131.3"
SSH_KEY="$HOME/.ssh/id_ed25519_landscape_droplet"
REPO="/Users/sanay_d/Documents/Personal/LandscapeForm"

# ── SSH agent ──────────────────────────────────────────────────────────────────
eval "$(ssh-agent -s)" > /dev/null
ssh-add "$SSH_KEY"

# ── Stop backend ───────────────────────────────────────────────────────────────
echo "→ Stopping backend..."
ssh -i "$SSH_KEY" "$DROPLET" "pm2 stop landscapeform-api || true"

# ── Backend ────────────────────────────────────────────────────────────────────
echo "→ Building backend..."
cd "$REPO/backend"
GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

echo "→ Uploading backend..."
scp -i "$SSH_KEY" api "$DROPLET:/var/www/landscapeform/backend/api"
rm api

# ── Frontend ───────────────────────────────────────────────────────────────────
echo "→ Building frontend..."
cd "$REPO/frontend"
npm run build

echo "→ Uploading frontend..."
rsync -az --delete \
  -e "ssh -i $SSH_KEY" \
  .next \
  "$DROPLET:/var/www/landscapeform/frontend/"

rsync -az \
  -e "ssh -i $SSH_KEY" \
  public package.json package-lock.json \
  "$DROPLET:/var/www/landscapeform/frontend/"

# ── Restart services ───────────────────────────────────────────────────────────
echo "→ Restarting services..."
ssh -i "$SSH_KEY" "$DROPLET" "cd /var/www/landscapeform/frontend && npm ci --omit=dev && pm2 restart landscapeform-api && pm2 restart landscapeform-frontend"

echo "✓ Done"
