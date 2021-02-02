go-build:
	GOOS=linux go build -o ./app .
docker-build: go-build
	docker build -t luebken/cd-auto .
docker-push: docker-build
	docker push luebken/cd-auto

k8s-run: docker-push
	kubectl run --rm -i cd-auto --image=luebken/cd-auto
k8s-delete:
	kubectl delete pods cd-auto

k8s-example-run:
	kubectl run nginx --image=nginx --labels="app=nginx,instana_customdashboard_auto=true"
k8s-example-delete:
	kubectl delete pod/nginx

configmap-create:
	kubectl create configmap instana-dashboard --from-file=example.json
	kubectl label configmap instana-dashboard instana_customdashboard_auto=true
configmap-delete:
	kubectl delete configmap instana-dashboard