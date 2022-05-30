## 0. Qiuniu

Here let's using qiuniu to collect the application log.

## 1. Help
Get help of qiuniu:
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
  service               provide the log collect service by restful api

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
  service
    -si, --service-ip           listening ip of qiuniu
    -sp, --service-port         listening port of qiuniu
```

Let's using the command one by one.

## 2. version
Get the version of qiuniu by `qiuniu version`:
```
qiuniu version: 2.0 Release  go version: go1.17.7
```

## 3. env
Get the env of qiuniu by `qiuniu env`:
```
QIUNIU_NAMESPACE=
QIUNIU_WORKSPACE=
QIUNIU_KUBECONFIG=
QIUNIU_SERVICE_IP=
QIUNIU_SERVICE_PORT=
```

qiuniu get the env from os, the default env is empty, it can set by user manually.

For example, set env:
```
$ export QIUNIU_NAMESPACE=demo
$ export QIUNIU_WORKSPACE=/home
```

Get the env:
```
$ qiuniu env
QIUNIU_NAMESPACE=demo
QIUNIU_WORKSPACE=/home
```

*Notice: The env is to simplify the options of log command, if no log options specific the env will be used. Otherwise, the log options will be used.*

## 4. log
Collect the application log by `qiuniu log`.

### 4.1 collect log by options
```
$ qiuniu log -k [kubeconfig_path] -n [namespace] -w [workspace]
```

### 4.2 collect log by env
1) Set env as chapter `2. env`

2) collect log
```
$ qiuniu log
```

The log saved in the path of [workspace]/qiuniu, for example:
```
$ ls /var/home/core/qiuniu/
2022-04-29T08:09:21Z

$ ls -l 2022-04-29T08\:09\:21Z/
total 372
drwxrwxr-x.   59 core core   4096 Apr 29 08:09 configmap
drwxrwxr-x.    6 core core    116 Apr 29 08:09 deployment
drwxrwxr-x.   49 core core   4096 Apr 29 08:09 endpoint
drwxrwxr-x. 5228 core core 266240 Apr 29 08:09 event
drwxrwxr-x.   50 core core   4096 Apr 29 08:09 image_stream
drwxrwxr-x.    7 core core    186 Apr 29 08:09 node
drwxrwxr-x.   39 core core   4096 Apr 29 08:09 pod
drwxrwxr-x.   61 core core   4096 Apr 29 08:09 pv
drwxrwxr-x.   19 core core   4096 Apr 29 08:09 pvc
-rwxrwxr-x.    1 core core    211 Apr 29 08:09 qiuniu_description.json
drwxrwxr-x.    5 core core     79 Apr 29 08:09 role
drwxrwxr-x.   10 core core    207 Apr 29 08:09 helm
drwxrwxr-x.   10 core core    207 Apr 29 08:09 role_binding
drwxrwxr-x.   49 core core   4096 Apr 29 08:09 service
drwxrwxr-x.    9 core core    137 Apr 29 08:09 service_account
drwxrwxr-x.   20 core core   4096 Apr 29 08:09 statefulset
drwxrwxr-x.    7 core core    162 Apr 29 08:09 storage_class
```

Here the pod log include the previous and current container log under the path of pod, the pod log is not only the pod, but also container log:
```
$ ls
current-container  demo-0.json  previous-container
```

- demo-0.json is the pod log include the pod info, status, etc.
- current-container and previours-container is the container log under the pod.

The log is naming by the time format, so qiuniu can collect different log in same workspace.

## 5. clean
`clean` is to clean the log according to the time interval between the collect time and the current time, for example:
```
$ qiuniu clean -w /var/home/core -i 24
```

It means clean the log over `24h` under workspace `/var/home/core`.

If not specific the workspace and interval, qiuniu will clean all logs under the default workspace `$HOME/qiuniu`.

## 6. zip
`zip` command is used to compress the log to zip, for example:
```
$ qiuniu zip -d /var/home/core/qiuniu
```

It means to compress the log under `/var/home/core/qiuniu` to zip package `qiuniu.zip`.

If not specific the dir of zip, qiuniu will zip the directory under `$HOME/qiuniu`.

## 7. service
`service` is to provide the RestfulAPI mode of qiuniu.

```
$ ./qiuniu service
[GIN-debug] Listening and serving HTTP on 0.0.0.0:9189
```

qiuniu will listening to the localhost and port `9189` by default, can change the serviceIP and servicePort with options `-si` and `-sp`.

For the detailed info of service OpenAPI, please refer here. 
