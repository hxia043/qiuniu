package resource

import (
	"github/hxia043/qiuniu/internal/app/resource/configmap"
	"github/hxia043/qiuniu/internal/app/resource/deployment"
	"github/hxia043/qiuniu/internal/app/resource/endpoint"
	"github/hxia043/qiuniu/internal/app/resource/event"
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

func NewResources(dir string) []Resource {
	resources := []Resource{
		pod.NewPod(dir),
		endpoint.NewEndpoint(dir),
		deployment.NewDeployment(dir),
		statefulset.NewStatefulset(dir),
		event.NewEvent(dir),
		node.NewNode(dir),
		job.NewJob(dir),
		pvc.NewPersistentVolumeClaims(dir),
		configmap.NewConfigMap(dir),
		storageclass.NewStorageClass(dir),
		pv.NewPersistentVolume(dir),
		role.NewRole(dir),
		imagestream.NewImageStream(dir),
		rolebinding.NewRoleBinding(dir),
		serviceaccount.NewServiceAccount(dir),
		service.NewService(dir),
	}

	return resources
}
