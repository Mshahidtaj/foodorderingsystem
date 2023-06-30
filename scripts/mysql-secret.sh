#!/bin/bash
set -euox

NAMESPACE="food-order"
SECRET_NAME="mysql-secret"

# Check if the namespace exists, create it if it doesn't
kubectl get namespace $NAMESPACE > /dev/null 2>&1
if [ $? -ne 0 ]; then
  kubectl create namespace $NAMESPACE
fi

# Get the MySQL root password from the mysql secret
MYSQL_ROOT_PASSWORD=$(kubectl get secret --namespace mysql mysql -o jsonpath="{.data.mysql-root-password}" | base64 -d)

# Create the Kubernetes secret
kubectl create secret generic $SECRET_NAME \
  --namespace $NAMESPACE \
  --from-literal=password="$MYSQL_ROOT_PASSWORD"
