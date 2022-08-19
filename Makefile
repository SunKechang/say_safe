HOMEDIR := $(shell pwd)

pre-build:
	cd $(HOMEDIR)
	env GOOS=linux GOARCH=386 go build -o bjfu main.go
	#scp bjfu root@39.107.25.37:/usr/projects

compile:build
build:
	docker build -f Dockerfile -t saysafe .

run:
	docker run -d -it -p 8080:8080 saysafe:latest