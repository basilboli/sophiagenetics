TIMESTAMP :=$(shell /bin/date +"%Y%m%d.%H%M%S")
IMAGE :=sg
COMMIT_HASH := $(shell git log -1 --pretty=format:%h || echo 'master')

# Default target : Do nothing

default:build

.PHONY: all
all: build

test-locally:
	echo "You should be running local database. Use the following handy command ;)"
	echo "docker run -p 27017:27017 -d mongo:3.4"	
	cd backend; go test -v ./...

test:
	cd backend; docker-compose -f docker-compose.test.yml build
	cd backend; docker-compose -f docker-compose.test.yml up --abort-on-container-exit

build-backend:test
	cd backend; docker build --file Dockerfile --build-arg BUILD_TIME=$(TIMESTAMP) --build-arg COMMIT_HASH=$(COMMIT_HASH) -t sg-api .
	docker tag sg-api basilboli/sg-api

build-frontend:
	cd frontend; docker build --file Dockerfile -t sg-web .
	docker tag sg-web basilboli/sg-web

build: build-backend build-frontend

push:
	docker push basilboli/sg-api
	docker push basilboli/sg-web

pull-latest:
	docker pull basilboli/sg-api
	docker pull basilboli/sg-web

run-locally-backend:
	echo "You should be running local database. Use the following handy command ;)"
	echo "docker run -p 27017:27017 -d mongo:3.4"	
	cd backend; go run main.go

run-locally-web:
	cd frontend; yarn start

run:
	docker-compose up

stop:
	docker-compose down	

clean:
	docker-compose stop
	$(eval CONTAINERS=$(shell docker ps -a -q))		
	docker rmi basilboli/sg-api || true
	docker rmi basilboli/sg-web || true
	docker rm $(CONTAINERS)	