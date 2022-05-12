##
# Gomation
#
# @file
# @version 0.1



# end

build: go-binary tailwind

go-binary:
	go build -o gomation ./cli/main.go

tailwind:
	cd server/html && \
	npx tailwindcss build -i tailwind.css -o static/main.css
