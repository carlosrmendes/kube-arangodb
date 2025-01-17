# Introduction

Kubernetes ArangoDB Operator.

# Chart Details

Chart will install fully operational ArangoDB Kubernetes Operator. CRD are moved to different Helm package.

# Prerequisites

To be able to work with Operator, Custom Resource Definitions needs to be installed. More details can be found in `kube-arangodb-crd` chart.

# Resources Required

In default installation deployment with 2 pods will be created. Each default pod require 256MB of ram and 250m of CPU.

# Installing the Chart

Chart can be installed in two methods:
- With all Operators in single Helm Release
- One Helm Release per Operator

Possible Operators:
- `ArangoDeployment` - enabled by default
- `ArangoDeploymentReplications` - enabled by default
- `ArangoLocalStorage` - disabled by default

To install Operators in mode "One per Helm Release" we can use:

```
helm install --name arango-deployment kube-arangodb.tar.gz \
             --set operator.features.deployment=true \
             --set operator.features.deploymentReplications=false \
             --set operator.features.storage=false


helm install --name arango-deployment-replications kube-arangodb.tar.gz \
             --set operator.features.deployment=false \
             --set operator.features.deploymentReplications=true \
             --set operator.features.storage=false


helm install --name arango-storage kube-arangodb.tar.gz \
             --set operator.features.deployment=false \
             --set operator.features.deploymentReplications=false \
             --set operator.features.storage=true
```


# Configuration

### `operator.image`

Image used for the ArangoDB Operator.

Default: `arangodb/kube-arangodb:latest`

### `operator.imagePullPolicy`

Image pull policy for Operator images.

Default: `IfNotPresent`

### `operator.imagePullSecrets`

List of the Image Pull Secrets for Operator images.

Default: `[]string`

### `operator.scope`

Scope on which Operator will be configured.

Default: `legacy`

Supported modes:
- `legacy` - mode with limited cluster scope access
- `namespaced` - mode with namespace access only

### `operator.service.type`

Type of the Operator service.

Default: `ClusterIP`

### `operator.annotations`

Annotations passed to the Operator Deployment definition.

Default: `[]string`

### `operator.resources.limits.cpu`

CPU limits for operator pods.

Default: `1`

### `operator.resources.limits.memory`

Memory limits for operator pods.

Default: `256Mi`

### `operator.resources.requested.cpu`

Requested CPI by Operator pods.

Default: `250m`

### `operator.resources.requested.memory`

Requested memory for operator pods.

Default: `256Mi`

### `operator.nodeSelector`

NodeSelector for Deployment pods.

Default: `{}`

### `operator.tolerations`

Tolerations for Deployment pods.

Default: 
```
  - key: node.kubernetes.io/unreachable
    operator: Exists
    effect: NoExecute
    tolerationSeconds: 5
  - key: node.kubernetes.io/not-ready
    operator: Exists
    effect: NoExecute
    tolerationSeconds: 5
```

### `operator.replicaCount`

Replication count for Operator deployment.

Default: `2`

### `operator.updateStrategy`

Update strategy for operator pod.

Default: `Recreate`

### `operator.features.deployment`

Define if ArangoDeployment Operator should be enabled.

Default: `true`

### `operator.features.deploymentReplications`

Define if ArangoDeploymentReplications Operator should be enabled.

Default: `true`

### `operator.features.storage`

Define if ArangoLocalStorage Operator should be enabled.

Default: `false`

### `operator.features.backup`

Define if ArangoBackup Operator should be enabled.

Default: `false`

### `rbac.enabled`

Define if RBAC should be enabled.

Default: `true`

# Limitations

N/A