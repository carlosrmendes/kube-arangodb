#!/usr/bin/fish

source helper.fish

set -g TESTNAME test4a
set -g TESTDESC "Deployment of mode cluster (development, enterprise, local storage)"
set -g YAMLFILE generated/cluster-local-storage-enterprise-dev.yaml
set -g YAMLFILESTORAGE generated/local-storage-community-dev.yaml
set -g DEPLOYMENT acceptance-cluster
printheader

# Deploy local storage:
kubectl apply -f $YAMLFILESTORAGE
and waitForKubectl "get storageclass" "acceptance.*arangodb.*localstorage" "" 1 60
or fail "Local storage could not be deployed."

# Deploy and check
kubectl apply -f $YAMLFILE
and waitForKubectl "get pod" "$DEPLOYMENT-prmr" "1/1 *Running" 3 120
and waitForKubectl "get pod" "$DEPLOYMENT-agnt" "1/1 *Running" 3 120
and waitForKubectl "get pod" "$DEPLOYMENT-crdn" "1/1 *Running" 3 120
and waitForKubectl "get service" "$DEPLOYMENT *ClusterIP" 8529 1 120
and waitForKubectl "get service" "$DEPLOYMENT-ea *LoadBalancer" "-v;pending" 1 180
and waitForKubectl "get pvc" "$DEPLOYMENT" "RWO *acceptance" 6 120
or fail "Deployment did not get ready."

# Automatic check
set ip (getLoadBalancerIP "$DEPLOYMENT-ea")
testArangoDB $ip 120
or fail "ArangoDB was not reachable."

# Manual check
output "Work" "Now please check external access on this URL with your browser:" "  https://$ip:8529/" "then type the outcome followed by ENTER."
inputAndLogResult

# Cleanup
kubectl delete -f $YAMLFILE
waitForKubectl "get pod" $DEPLOYMENT "" 0 120
or fail "Could not delete deployment."

kubectl delete -f $YAMLFILESTORAGE
kubectl delete storageclass acceptance
waitForKubectl "get storageclass" "acceptance.*arangodb.*localstorage" "" 0 120
or fail "Could not delete deployed storageclass."

output "Ready" ""
