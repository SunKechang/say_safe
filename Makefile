HOMEDIR := $(shell pwd)

pre-build:
	cd D:\\GolangProjects\\safeWeb
	$Env:GOARCH="386";$Env:GOOS="linux"
	go build -o bjfu main.go
	#scp bjfu root@39.107.25.37:/usr/projects

compile:build
build:
	docker rm -f safeweb
	docker rmi saysafe
	docker build -f Dockerfile -t saysafe .

run:
	docker run -d -it --network=host --name safeweb saysafe:latest

mysql-start:

	docker run -tid --name mysql --network=host -e MYSQL_ROOT_PASSWORD=Pgone3123 daff57b7d2d1