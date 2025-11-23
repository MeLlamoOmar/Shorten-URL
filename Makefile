GO := go
APP := urlShortener
MIGRATIONSDIR := migrations
DBURL = sqlite3://urlShortener.db

.PHONY: build test fmt create up down migrate-status dev

build:
	go build -o urlShortener main.go	

test:
	go test ./..

fmt:
	go fmt ./..

dev:
	go run .
	
create:
	migrate create -ext sql -dir ${MIGRATIONSDIR} -seq ${NAME}

up:
	migrate -path ${MIGRATIONSDIR} -database ${DBURL} up

down:
	migrate -path ${MIGRATIONSDIR} -database ${DBURL} down

migrate-status:
	migrate -path ${MIGRATIONSDIR} -database ${DBURL} version
