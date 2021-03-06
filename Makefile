build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/garmin-stress2pixela garmin-stress2pixela/main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
