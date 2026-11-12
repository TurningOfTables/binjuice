init:
	cp ./.env.dist ./.env

build-for-pi:
	env GOOS=linux GOARCH=arm GOARM=7 go build .
