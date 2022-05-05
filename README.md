# Qiuniu
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

## Installation
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

## Quickly start
1. Find help from help command
```
$qiuniu help
The log collector for Kubernetes application

Available Commands:
  help                  help for qiuniu
  version               print the qiuniu version information
  env                   qiuniu client env information
  zip                   compress the log
  clean                 clean the log
  helm                  collect the helm release log for application
  log                   collect the kubernetes application log

Options:
  log
    -h, --host                  kubernetes cluster hostname or ip
    -p, --port                  kubernetes cluster port
    -t, --token                 kubernetes cluster token
    -n, --namespace             kubernetes cluster namespace
    -w, --workspace             the workspace for qiuniu
  helm
    -k, --kubeconfig            local kubeconfig path of kubernetes cluster
    -n, --namespace             kubernetes cluster namespace
    -w, --workspace             the workspace for qiuniu
  zip
    -d, --dir                   the dir of compress the log
  clean
    -w, --workspace             the workspace for qiuniu
    -i, --interval              the time interval between log collect time and current time
```

2. Get qiuniu version
```
$ qiuniu version
qiuniu version: 0.1 Debug  go version: go1.17.7
```

## Docs
1. Get start from the user guide
2. Get start from the developer guide

## Q&A
1. raise issue from [issues](https://github.com/hxia043/qiuniu/issues)
