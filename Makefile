build:
	docker build -t app:local . && docker-compose stop app &&  docker-compose up app