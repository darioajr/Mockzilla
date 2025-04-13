APP_NAME=mockzilla
build:
	go build -o bin/$(APP_NAME) main.go

run: build
	./bin/$(APP_NAME) -response-code=201 -cert=server.crt -key=server.key -port=8443

clean:
	rm -rf bin/

test:
	./scripts/test.sh
