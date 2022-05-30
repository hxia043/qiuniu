# Qiuniu

## 0. Overview

qiuniu is a log collector for application based on Kubernetes.

![avatar](./image/qiuniu.png)

qiuniu can support user to collect these logs of application:
- configmap
- container
- deployment
- endpoint
- event
- image
- job
- node
- pod
- persistent volume
- persistent volume claim
- roles
- rolebinding
- service
- serviceaccount
- statefulset
- storageclass
- helm

## 1. Install qiuniu with CLI
To install qiuniu, you need to install Go and set your Go workspace first.

1. Clone qiuniu from github
```
git clone https://github.com/hxia043/qiuniu.git
```

2. Build qiuniu
```
$ cd cmd
$ go build -o qiuniu main.go
```

3. Move qiuniu to system env
```
mv qiuniu /usr/bin
```

### 1.1 Quickly start
1. Find help from help command
```
$ qiuniu help
The log collector for Kubernetes application

Available Commands:
  help                  print help information
  version               print the version information
  env                   print host env information
  log                   collect the kubernetes application log
  zip                   compress the log
  clean                 clean the log

Options:
  log
    -n, --namespace             kubernetes cluster namespace
    -w, --workspace             the workspace of qiuniu
    -k, --kubeconfig            local kubeconfig path of kubernetes cluster
  zip
    -d, --dir                   the dir of compress the log
  clean
    -w, --workspace             the workspace for qiuniu
    -i, --interval              the time interval between log collect time and current time, unit(h)
```

2. Get qiuniu version
```
$ qiuniu version
qiuniu version: 1.0 Release  go version: go1.17.7
```

## 2. Instanll qiuniu with Helm Chart
The helm chart under the directory of deployment, can deploy the qiuniu by `helm install`:
```
$ cd /deployment
$ helm install qiuniu .
```

Before deployment, please update the `repository` under the image field of `values.yaml`.

## Docs
1. Get start for command line mode of qiuniu.
2. Get start for service mode of qiuniu.

## Q&A
1. raise issue from [issues](https://github.com/hxia043/qiuniu/issues)
