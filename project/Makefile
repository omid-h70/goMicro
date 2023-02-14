SHELL=cmd.exe
FRONT_END_BINARY=frontEndApp.exe
BROKER_BINARY=brokerApp.exe
AUTH_SERVICE_BINARY=authApp.exe
LOGGER_SERVICE_BINARY=loggerService.exe
##
## Attention
## All MakeFile Command Lines Must Contain Tab Not Spaces
##
##

##
## CROSS Compiling Golang
## By Setting These Options we're cross Compiling Go From Windows To Linux
## GOOS=linux && set GOARCH=amd64
##

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

### up_build: stops docker-compose (if running), builds all projects and starts docker compose
#up_build: build_broker
#    @echo Stopping docker images (if running...)
#    docker-compose down
#    @echo Building (when required) and starting docker images...
#    docker-compose up --build -d
#    @echo Docker images built and started!
#
### down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"
#
# build Front End As Well
build_frontend:
	@echo "Building frontend binary..."
	chdir ..\frontend && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"
### build_broker: builds the broker binary as a linux executable
## Append && to GOOS=linux File !!!!!!!!!!! with no spaces !!!!!!!!!!!!
## && means execute second command if first was successful
build_broker:
	@echo "Building broker binary..."
	chdir ..\broker-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"
#
#
build_logger:
	@echo "Logger Service binary..."
	chdir ..\logger-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${LOGGER_SERVICE_BINARY} ./cmd/api
	@echo "Done!"
#
### build_auth: builds the authentication service
build_auth:
	@echo "Building Authentication Service..."
	chdir ..\authentication-service && set CGO_ENABLED=0&& set GOOS=windows&& go build -o ${AUTH_SERVICE_BINARY} ./cmd/api
	@echo "Done!"
#
### start: starts the front end
#start: build_front
#    @echo Starting front end
#    chdir ..\front-end && start /B ${FRONT_END_BINARY} &
#
### stop: stop the front end
#stop:
#    @echo Stopping front end...
#    @taskkill /IM "${FRONT_END_BINARY}" /F
#    @echo "Stopped front end!"
