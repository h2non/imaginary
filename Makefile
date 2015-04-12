docker-build:
	docker build --no-cache -t h2non/imaginary:$(VERSION) .

docker-push:
	docker pull h2non/imaginary

docker:
	docker-build
	docker-push
