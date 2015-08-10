test:
	go test

bench:
	go test -bench=.

test-load:
	bash bench.sh

docker-build:
	docker build --no-cache=true -t h2non/imaginary:$(VERSION) .

docker-push:
	docker push h2non/imaginary:$(VERSION)

docker: docker-build docker-push
