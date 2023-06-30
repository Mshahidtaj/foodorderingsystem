# Variables
DB_CONTAINER_NAME = food-ordering-db
APP_CONTAINER_NAME = mshahidtaj/food-order-app
UI_CONTAINER_NAME = mshahidtaj/food-order-ui

.PHONY: parse

help:
	@echo "************************** Food Order System ******************************* "

	@echo "1.  deploy-mysql              - Deploy MYSQL to the Kubernetes"
	@echo "2.  create-mysql-secret       - Create MySQL Secret"
	@echo "3.  deploy-food-order-app     - Deploy food order app to your current K8s Cluster"
	@echo "4.  deploy-food-order-ui      - Deploy food order UI to your current K8s Cluster"
	@echo "5.  start-food-order-app-ui   - Start Food Order App UI"
	@echo "6.  build-food-order-app      - build food order app"
	@echo "7.  Run Test                  - Testing food order app"
	@echo "8.  build-food-order-ui       - Build food order UI"
	@echo "8.  push-food-order-app       - build food order app"
	@echo "10.  push-food-order-ui        - Build food order UI"
	@echo "11. uninstall-food-system     - Uninstall Food System"
	@echo "12. uninstall-mysql-db        - Uninstall MYSQL DB"

deploy-mysql:
	helm upgrade --install mysql oci://registry-1.docker.io/bitnamicharts/mysql

create-mysql-secret:
	./scripts/mysql-secret.sh

deploy-food-order-app:
	helm upgrade --install food-order-app --namespace food-order charts/food-order-app --create-namespace

deploy-food-order-ui:
	helm upgrade --install food-order-ui --namespace food-order charts/food-order-ui

start-food-order-app-ui:
	kubectl port-forward svc/food-order-app 8080:8080 &
	kubectl port-forward svc/food-order-ui 8081:80 &

build-food-order-app:
    docker buildx build --tag $(APP_CONTAINER_NAME):$(TAG) .

push-food-order-app:
    docker push $(APP_CONTAINER_NAME):$(TAG)

build-food-order-ui:
    docker buildx build --tag $(UI_CONTAINER_NAME):$(TAG) ui/

push-food-order-ui:
    docker push $(UI_CONTAINER_NAME):$(TAG)

uninstall-food-system:
	helm uninstall food-order-app --namespace food-order
	helm uninstall food-order-ui --namespace food-order

uninstall-mysql-db:
	helm uninstall mysql --namespace mysql
