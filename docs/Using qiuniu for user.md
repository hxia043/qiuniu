## 0. Qiuniu

Here let's using qiuniu to collect the application log.

## 1. Help
Get help from qiuniu by `qiuniu help`:
```
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

Let's using the command one by one.

## 2. version
Get the version of qiuniu by `qiuniu version`:
```
$ qiuniu version
qiuniu version: 1.0 Release  go version: go1.17.7
```

## 3. env
Get the env of qiuniu by `qiuniu env`:
```
$ qiuniu env
QIUNIU_PORT=
QIUNIU_TOKEN=
QIUNIU_NAMESPACE=
QIUNIU_WORKSPACE=
QIUNIU_HOST=
```

qiuniu get the env from system, the default env is empty, it can set by user manully.

For example, set env:
```
$ export QIUNIU_HOST=127.0.0.1
$ export QIUNIU_PORT=6443
$ export QIUNIU_NAMESPACE=demo
$ export QIUNIU_WORKSPACE=/home
```

Get the env:
```
$ qiuniu env
QIUNIU_PORT=6443
QIUNIU_TOKEN=
QIUNIU_NAMESPACE=demo
QIUNIU_WORKSPACE=/home
QIUNIU_HOST=127.0.0.1
```

*Notice: The env is to simplify the options of log command, if no log options specific the env will be used. Otherwise, the log options will be used.*

## 4. log
Collect the application log by `qiuniu log`.

### 4.1 collect log by options
```
$ qiuniu log -h [host] -p [port] -t [token] -n [namespace] -w [workspace]
```

### 4.2 collect log by env
1) Set env as chapter `2. env`

2) collect log
```
$ qiuniu log
```


For the host, port and token info can be checked from the kubernetes cluster kubeconfig, more detailed can refer [here](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/).

If there is no error raised up, it means log collect finished:
```
$ qiuniu log -h 10.10.xxx.xxx -p 6443 -t sha256~x8C4ouzPt3tqyEAwMxLtPZ8Mz5rgtHEQOS72Ie6fZJY -n demo -w /var/home/core
$
```

The log saved in the path of [workspace]/qiuniu:
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
drwxrwxr-x.   10 core core    207 Apr 29 08:09 role_binding
drwxrwxr-x.   49 core core   4096 Apr 29 08:09 service
drwxrwxr-x.    9 core core    137 Apr 29 08:09 service_account
drwxrwxr-x.   20 core core   4096 Apr 29 08:09 statefulset
drwxrwxr-x.    7 core core    162 Apr 29 08:09 storage_class
```

Here need highlight the pod log cause it include the previous and current container log under the pod, the pod is not only the "pod":
```
$ ls
current-container  demo-0.json  previous-container
```

- demo-0.json is the pod log include the pod info, status, etc.
- current-container and previours-container is the container log under the pod.

The log has been collect in the package which naming with time format, so qiuniu can collect the log in same workspace.

## 5. helm
Collect the helm release log by `qiuniu helm`.
```
qiuniu helm -k [kubeconfig_path] -n [namespace] -w [workspace]
```

For the detail info of kubeconfig, please refer [here](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/)

The log saved in the path of [workspace]/qiuniu/helm, for example:
```
$ qiuniu helm -k /var/home/core/kubeconfig -n demo -w /var/home/core

$ pwd
/var/home/core/qiuniu/helm

$ ls
2022-05-05T03:01:19Z  2022-05-05T03:02:29Z
```
