test:
	go test

test-bench:
	go test -bench=.

benchmark:
	bash benchmark.sh

docker-build:
	docker build --no-cache=true -t h2non/imaginary:$(VERSION) .

docker-push:
	docker push h2non/imaginary:$(VERSION)

docker: docker-build docker-push

.PHONY: test test-bench benchmark docker-build docker-push docker
