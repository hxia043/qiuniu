package resource

import (
	"github/hxia043/qiuniu/internal/app/resource/configmap"
	"github/hxia043/qiuniu/internal/app/resource/deployment"
	"github/hxia043/qiuniu/internal/app/resource/endpoint"
	"github/hxia043/qiuniu/internal/app/resource/event"
	"github/hxia043/qiuniu/internal/app/resource/helm"
	"github/hxia043/qiuniu/internal/app/resource/imagestream"
	"github/hxia043/qiuniu/internal/app/resource/job"
	"github/hxia043/qiuniu/internal/app/resource/node"
	"github/hxia043/qiuniu/internal/app/resource/pod"
	"github/hxia043/qiuniu/internal/app/resource/pv"
	"github/hxia043/qiuniu/internal/app/resource/pvc"
	"github/hxia043/qiuniu/internal/app/resource/role"
	"github/hxia043/qiuniu/internal/app/resource/rolebinding"
	"github/hxia043/qiuniu/internal/app/resource/service"
	"github/hxia043/qiuniu/internal/app/resource/serviceaccount"
	"github/hxia043/qiuniu/internal/app/resource/statefulset"
	"github/hxia043/qiuniu/internal/app/resource/storageclass"
)

type Resource interface {
	Log() error
}

func NewResources(host, token, dir string, kubeconfig []byte) []Resource {
	resources := []Resource{
		pod.NewPod(host, token, dir),
		endpoint.NewEndpoint(host, token, dir),
		deployment.NewDeployment(host, token, dir),
		statefulset.NewStatefulset(host, token, dir),
		event.NewEvent(host, token, dir),
		node.NewNode(host, token, dir),
		job.NewJob(host, token, dir),
		pvc.NewPersistentVolumeClaims(host, token, dir),
		configmap.NewConfigMap(host, token, dir),
		storageclass.NewStorageClass(host, token, dir),
		pv.NewPersistentVolume(host, token, dir),
		role.NewRole(host, token, dir),
		imagestream.NewImageStream(host, token, dir),
		rolebinding.NewRoleBinding(host, token, dir),
		serviceaccount.NewServiceAccount(host, token, dir),
		service.NewService(host, token, dir),
	}

	helm := helm.NewHelm(kubeconfig, dir)
	resources = append(resources, helm)

	return resources
}
