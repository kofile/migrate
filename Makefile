all: dist/darwin_x86/migrate dist/linux_x86/migrate dist/alpine_x86/migrate

dist/darwin_x86/migrate: main.go
	gox -osarch="darwin/amd64" -output="./dist/darwin_x86/migrate"

dist/linux_x86/migrate: main.go
	gox -osarch="linux/amd64" -output="./dist/linux_x86/migrate"

dist/alpine_x86/migrate: main.go
	mkdir -p ./dist/alpine_x86
	docker build -t migrate/alpine-builder .
	docker run --name migrate-alpine-builder migrate/alpine-builder build -o dist/migrate-alpine-amd64 -ldflags "-s" -a -installsuffix cgo .
	docker cp migrate-alpine-builder:/go/src/migrate/dist/migrate-alpine-amd64 ./dist/alpine_x86/migrate
	docker rm migrate-alpine-builder
