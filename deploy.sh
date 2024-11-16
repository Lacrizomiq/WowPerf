#!/bin/bash

# Terminal colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Help function
show_help() {
    echo -e "${BLUE}WoWPerf Deployment Script${NC}"
    echo -e "${GREEN}Usage: $0 <command> <environment>${NC}"
    echo ""
    echo "Commands:"
    echo "  up          Start the environment"
    echo "  down        Stop the environment"
    echo "  restart     Restart the environment"
    echo "  logs        Show logs (use -f for follow)"
    echo "  ps          Show running containers"
    echo "  pull        Pull latest images"
    echo "  build       Build containers"
    echo ""
    echo "Environments:"
    echo "  local       Local development environment"
    echo "  test        Test environment"
    echo "  prod        Production environment"
    echo ""
    echo "Options:"
    echo "  -f, --follow    Follow log output (with logs command)"
    echo ""
    echo "Examples:"
    echo "  $0 up local"
    echo "  $0 logs -f test"
    echo "  $0 build prod"
}

# Argument validation
if [ $# -lt 2 ]; then
    show_help
    exit 1
fi

COMMAND=$1
ENV=$2
FOLLOW=""
if [ "$3" == "-f" ] || [ "$3" == "--follow" ]; then
    FOLLOW="-f"
fi

# Validate environment
if [[ ! "$ENV" =~ ^(local|test|prod)$ ]]; then
    echo -e "${RED}Invalid environment. Must be local, test, or prod${NC}"
    exit 1
fi

# Check .env file with simple path
if [ ! -f "deploy/${ENV}/.env" ]; then
    echo -e "${RED}Error: deploy/${ENV}/.env file not found${NC}"
    echo -e "${YELLOW}Please ensure your .env file exists in deploy/${ENV}/ directory${NC}"
    exit 1
fi

# Load environment variables
set -a
source "deploy/${ENV}/.env"
set +a

# Check traefik-public network
if ! docker network ls | grep -q "traefik-public"; then
    echo -e "${YELLOW}Creating traefik-public network...${NC}"
    docker network create traefik-public
fi

# Docker compose configuration with simple paths
COMPOSE_FILES="-f deploy/docker-compose.base.yml -f deploy/${ENV}/docker-compose.yml"

# Execution log function
log_exec() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

# Command execution
case $COMMAND in
    "up")
        log_exec "Starting $ENV environment..."
        cd "$ENV"
        docker-compose $COMPOSE_FILES up
        log_exec "Environment $ENV started successfully"
        docker-compose $COMPOSE_FILES ps
        ;;
    "down")
        log_exec "Stopping $ENV environment..."
        cd "$ENV"
        docker-compose $COMPOSE_FILES down
        log_exec "Environment $ENV stopped"
        ;;
    "restart")
        log_exec "Restarting $ENV environment..."
        cd "$ENV"
        docker-compose $COMPOSE_FILES down
        docker-compose $COMPOSE_FILES up
        log_exec "Environment $ENV restarted"
        docker-compose $COMPOSE_FILES ps
        ;;
    "logs")
        cd "$ENV"
        docker-compose $COMPOSE_FILES logs $FOLLOW
        ;;
    "ps")
        cd "$ENV"
        docker-compose $COMPOSE_FILES ps
        ;;
    "pull")
        log_exec "Pulling latest images for $ENV..."
        cd "$ENV"
        docker-compose $COMPOSE_FILES pull
        log_exec "Images updated"
        ;;
    "build")
        log_exec "Building containers for $ENV..."
        cd "$ENV"
        docker-compose $COMPOSE_FILES build --no-cache
        log_exec "Build completed"
        ;;
    *)
        echo -e "${RED}Invalid command${NC}"
        show_help
        exit 1
        ;;
esac