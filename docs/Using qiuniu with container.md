## Download qiuniu image

The qiuniu image has been stored in [docker hub](https://hub.docker.com/repository/docker/hxia043/qiuniu).

pull qiuniu to local by:
```
$ docker pull hxia043/qiuniu:latest

$ docker images | grep qiuniu
hxia043/qiuniu         latest              710da4187a5c   14 minutes ago   380MB
```

## helm install qiuniu
After push image to repository, can deploy qiuniu with `helm install`:
```
$ cd deployment
$ helm install qiuniu .
```

*Notice: before helm install, need modify the values.yaml according to the configure of infra, like image repository, servcie-ip and service-port*

## visit qiuniu
the service type of qiuniu is NodePort, you can visit qiuniu by `http://$NODE_IP:$NODE_PORT`.

The variable $NODE_IP and $NODE_PORT can be found from NOTES after `helm install`, For example:
```
$ helm install qiuniu .
NAME: qiuniu
LAST DEPLOYED: Sun May 29 16:10:45 2022
NAMESPACE: ci1
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
1. Get the application URL by running these commands:
  export NODE_PORT=$(kubectl get --namespace ci1 -o jsonpath="{.spec.ports[0].nodePort}" services qiuniu)
  export NODE_IP=$(kubectl get nodes --namespace ci1 -o jsonpath="{.items[0].status.addresses[0].address}")
  echo "Visit http://$NODE_IP:$NODE_PORT to use qiuniu"
```

For the OpenAPI of qiuniu, please refer [here](https://github.com/hxia043/qiuniu/blob/main/api/openapi.yaml).
