# Launch Wow Perf

## Créez d'abord le réseau traefik-public

```zsh
docker network create traefik-public
```

## Supprimer le reseau traefik-public

```zsh
docker network rm traefik-public
```

## Lancer le systeme local

```zsh
./deploy.sh up local
```

## Arreter le system local

```zsh
./deploy.sh down local
```

## Restart les containeurs

```zsh
./deploy.sh restart local
```
