test:
	go test

bench:
	go test -bench=.

docker-build:
	docker build -t h2non/imaginary:$(VERSION) .

docker-push:
	docker push h2non/imaginary

docker: docker-build docker-push
