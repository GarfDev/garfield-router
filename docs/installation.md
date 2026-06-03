# Installation Guide

## Prerequisites

- Go 1.22+ (for building from source)
- Docker (optional, for containerised deployment)
- Kubernetes + Helm (optional, for Helm chart deployment)

## Option 1: Download Binary

Download the latest release for your platform:

```bash
# Linux (amd64)
curl -sL https://github.com/GarfDev/garfield-router/releases/latest/download/garfield-router-linux-amd64.tar.gz | tar xz
sudo mv garfield-router /usr/local/bin/

# macOS (arm64)
curl -sL https://github.com/GarfDev/garfield-router/releases/latest/download/garfield-router-darwin-arm64.tar.gz | tar xz
sudo mv garfield-router /usr/local/bin/

# Verify
garfield-router --version
```

## Option 2: Build from Source

```bash
git clone https://github.com/GarfDev/garfield-router.git
cd garfield-router
go build -o garfield-router .
./garfield-router
```

## Option 3: Docker

```bash
docker build -t garfield-router .
docker run -p 8050:8050 -v $(pwd)/config.yaml:/app/config.yaml garfield-router
```

Or with Docker Compose:

```yaml
services:
  garfield-router:
    build: ./garfield-router
    ports:
      - "8050:8050"
    volumes:
      - ./config.yaml:/app/config.yaml
    environment:
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - DATABASE_URL=postgres://user:pass@db:5432/mydb?sslmode=disable
```

## Option 4: Helm (Kubernetes)

```bash
helm install garfield-router ./helm/garfield-router \
  --set env.GEMINI_API_KEY=$GEMINI_API_KEY \
  --set env.ROUTER_API_TOKEN=$ROUTER_API_TOKEN
```

See [Deployment Guide](deployment.md) for detailed Kubernetes setup.

## Option 5: Python SDK

```bash
pip install garfield-router
```

```python
from garfield_router import GarfieldRouter
router = GarfieldRouter("http://localhost:8050", service="my-app")
response = router.chat("Hello, world!")
```

## Option 6: TypeScript SDK

```bash
npm install garfield-router
```

```typescript
import { GarfieldRouter } from 'garfield-router';
const router = new GarfieldRouter('http://localhost:8050', { service: 'my-app' });
const response = await router.chat('Hello, world!');
```

## First Run

On first run without a config file, the router auto-generates a default `config.yaml`:

```bash
./garfield-router
# [router] no config file at config.yaml, generating default
# [router] loaded config: 1 backends, 1 rules
# [router] listening on :8050
```

Open http://localhost:8050 to access the web UI.

## Verify Installation

```bash
curl http://localhost:8050/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "garfield-router",
  "version": "1.0.0",
  "backends_total": 1,
  "backends_healthy": 0
}
```

## Next Steps

1. [Configure backends and routing rules](configuration.md)
2. [Point your services at the router](user-guide.md)
3. [Set up monitoring](monitoring.md)
