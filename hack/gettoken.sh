#!/bin/bash
kubectl get serviceAccounts -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | while read -r line; do
  kubectl get secret -l serviceAccountName=$line -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | while read -r line; do
    kubectl describe secret $line | grep token | awk '{print $2}'
  done
done

kubectl config view --minify -o jsonpath=={.clusters[0].cluster.server}


kubectl get serviceAccounts svc-acct -n dev -o=jsonpath={.secrets[*].name}
kubectl get secret -n dev -l serviceAccountName=svc-acct -o json

kubectl get secret -n dev <svc-acct-secret-name> -n dev -o json