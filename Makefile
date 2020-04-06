build:
	docker-compose build
.PHONY: build

up:
	docker-compose up
.PHONY: up

client:
	docker-compose run --rm client yarn build
.PHONY: client

set:
	docker-compose run --rm set go build
.PHONY: set

clean:
	docker-compose stop
	docker-compose rm -fv
.PHONY: clean
