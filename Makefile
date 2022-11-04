bootstrap:
	kind version || brew install kind
	kind get clusters |grep kind || ./deploy/kind-registry.sh
	helm version || brew install helm
	./deploy/install-localstack.sh
	make build-demo-image
	make build-capturer-image
	./deploy/deploy.sh
	kubectl logs -f -l app=s3test -c capturer

build-capturer-image:
	docker build . -f build/capturer.Dockerfile -t capturer
	docker tag capturer localhost:5001/capturer
	docker push localhost:5001/capturer

log:
	./deploy/log.sh

update-rbac:
	kubectl apply -f deploy/rbac.yaml

redeploy:
	./deploy/deploy.sh

build-demo-image:
	docker build . -f demo-app/s3.Dockerfile -t demo
	docker tag demo localhost:5001/demo
	docker push localhost:5001/demo

rebuild-demo:
	make build-demo-image
	make redeploy

rebuild-capturer:
	make build-capturer-image
	make redeploy