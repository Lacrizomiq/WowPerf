# Comment setup WoW Perf dans votre propre environnement

Afin de setup WoW Perf dans votre propre environnement, voici la documentation et les étapes à suivre :

## Prérequis

- Docker
- Docker Compose
- Go
- Git

## Setup

1. Clone le repository
2. Avoir Go installé [(https://go.dev/dl/)](https://go.dev/dl/)
3. Avoir Docker et Docker Compose installé [(https://docs.docker.com/get-docker/)](https://docs.docker.com/get-docker/)
4. Copier `.env.example` en `.env` et le remplir avec vos propres valeurs

## Setup backend

1. Run `go mod tidy`

## Setup frontend

1. Run `npm install`
2. Run `npm run dev`

## Setup des variables d'environnement

1. Copier `.env.example` en `.env` et le remplir avec vos propres valeurs
2. Créer un client Blizzard API
3. Créer un client WarcraftLogs API
4. Créer un client Google Console API
5. Remplir les variables d'environnement avec les bonnes valeurs
6. Créer des secrets avec des hash générés avec `openssl rand -hex 32` pour :
   - `JWT_SECRET`
   - `CSRF_SECRET`
   - `ENCRYPTION_KEY`
7. Créer un mot de passe pour traefik auth avec `htpasswd -c auth <username>`

## Setup Traefik

1. Créer les certificats auto-signés dans `/deploy/local/certs` avec `openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /deploy/local/certs/tls.key -out /deploy/local/certs/tls.crt -subj "/CN=localhost"`
2. Copier le reste des fichiers dans `/deploy/local/traefik` si pas déjà fait

## Run le projet

1. Run `docker compose up -d`
