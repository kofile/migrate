.PHONEY: archive

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

dist/migrate_darwin_x86.tar.gz: dist/darwin_x86/migrate
	(cd dist && tar czf migrate_darwin_x86.tar.gz darwin_x86/migrate)

dist/migrate_linux_x86.tar.gz: dist/linux_x86/migrate
	(cd dist && tar czf migrate_linux_x86.tar.gz linux_x86/migrate)

dist/migrate_alpine_x86.tar.gz: dist/alpine_x86/migrate
	(cd dist && tar czf migrate_alpine_x86.tar.gz alpine_x86/migrate)

archive: dist/migrate_darwin_x86.tar.gz dist/migrate_linux_x86.tar.gz dist/migrate_alpine_x86.tar.gz
