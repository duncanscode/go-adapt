# Deployment Guide

## Prerequisites
- Linux server (Digital Ocean: 64.227.100.21 / 10.124.0.2)
- Nginx configured as reverse proxy
- `.env` file with `ANTHROPIC_API_KEY`

## Build for Production

```bash
./build.sh
```

This creates `go-adapt-linux` - a single binary with embedded frontend.

## Deploy to Server

### 1. Upload Files

```bash
# Upload binary and .env
scp go-adapt-linux root@64.227.100.21:/opt/go-adapt/
scp .env root@64.227.100.21:/opt/go-adapt/
```

### 2. On Server

```bash
# SSH into server
ssh root@64.227.100.21

# Navigate to app directory
cd /opt/go-adapt

# Make executable
chmod +x go-adapt-linux

# Test run (Ctrl+C to stop)
GIN_MODE=release ./go-adapt-linux

# Should see: "Server starting in release mode on port 1234"
```

### 3. Run Permanently

**Option A: Screen session** (simplest)
```bash
screen -S go-adapt
GIN_MODE=release ./go-adapt-linux
# Press Ctrl+A then D to detach

# To reattach: screen -r go-adapt
```

**Option B: Keep-alive cron**
Create `/opt/go-adapt/keep-alive.sh`:
```bash
#!/bin/bash
if ! pgrep -f go-adapt-linux > /dev/null; then
    cd /opt/go-adapt
    nohup GIN_MODE=release ./go-adapt-linux >> app.log 2>&1 &
fi
```

Then:
```bash
chmod +x keep-alive.sh
crontab -e
# Add: */5 * * * * /opt/go-adapt/keep-alive.sh
```

### 4. Configure Nginx

Add to your nginx config:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:1234;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://127.0.0.1:1234/health;
    }
}
```

Then:
```bash
nginx -t  # Test config
systemctl reload nginx
```

## Environment Variables

Create `.env` file on server:
```bash
ANTHROPIC_API_KEY=your_api_key_here
```

Optional variables:
- `GIN_MODE=release` (can also set in shell)
- `PORT=1234` (default, can change if needed)

## Updating

```bash
# 1. Build new version locally
./build.sh

# 2. Upload
scp go-adapt-linux root@64.227.100.21:/opt/go-adapt/

# 3. On server, kill old process
pkill go-adapt-linux

# 4. Start new version (screen or cron will handle)
GIN_MODE=release ./go-adapt-linux
```

## Monitoring

```bash
# Check if running
pgrep -f go-adapt-linux

# View logs (if using nohup)
tail -f /opt/go-adapt/app.log

# Test health endpoint
curl http://localhost:1234/health
```

## Troubleshooting

**Port already in use:**
```bash
lsof -ti:1234 | xargs kill -9
```

**Frontend not loading:**
- Frontend is embedded - no separate files needed
- Check nginx proxy headers

**LLM mode not working:**
- Verify `ANTHROPIC_API_KEY` in `.env`
- Check: `cat /opt/go-adapt/.env`

## Security Notes

- `.env` file contains API key - keep it secure (chmod 600)
- App runs on localhost:1234 - only accessible via nginx
- No database - sessions stored in memory (lost on restart)
- Trusted proxies configured for 127.0.0.1 and 10.124.0.2
