.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f xuhaidong/webook:v0.0.2
	@docker build -t xuhaidong/webook:v0.0.2 .
