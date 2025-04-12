APP_NAME=mockzilla
build:
	go build -o bin/$(APP_NAME) main.go

run: build
	./bin/$(APP_NAME)

clean:
	rm -rf bin/

test:
	./scripts/test.sh
