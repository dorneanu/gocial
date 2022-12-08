##
# gocial
#
# @file
# @version 0.1



# end
current_dir = $(shell pwd)

build: go-binary tailwind

go-binary:
	go build -o gocial ./cli/main.go

netlify:
	mkdir -p functions
	mkdir -p public
	GOOS=linux
	GOARCH=amd64
	GO111MODULE=on
	GOBIN=${PWD}/functions go get ./...

tailwind:
	cd server/html && \
	npx tailwindcss build -i tailwind.css -o static/main.css

go-lambda:
	go build -o gocial-lambda ./lambda/main.go
