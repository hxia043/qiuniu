package resource

import (
	"encoding/json"
	"fmt"
	"loki/internal/loki/resource/configmap"
	"loki/internal/loki/resource/deployment"
	"loki/internal/loki/resource/endpoint"
	"loki/internal/loki/resource/event"
	"loki/internal/loki/resource/imagestream"
	"loki/internal/loki/resource/job"
	"loki/internal/loki/resource/node"
	"loki/internal/loki/resource/pod"
	"loki/internal/loki/resource/pv"
	"loki/internal/loki/resource/pvc"
	"loki/internal/loki/resource/role"
	"loki/internal/loki/resource/rolebinding"
	"loki/internal/loki/resource/service"
	"loki/internal/loki/resource/serviceaccount"
	"loki/internal/loki/resource/statefulset"
	"loki/internal/loki/resource/storageclass"
	"os"
)

type Resource interface {
	Log() error
}

func GenerateLog(resp []byte, logDir string) error {
	ep := make(map[string]interface{})
	err := json.Unmarshal(resp, &ep)
	if err != nil {
		return err
	}

	_, ok := ep["items"]
	if !ok {
		return nil
	}

	items := ep["items"].([]interface{})
	for _, item := range items {
		itemElem := item.(map[string]interface{})

		metadata := itemElem["metadata"].(map[string]interface{})
		name := metadata["name"].(string)
		logPath := fmt.Sprintf("%s/%s.json", logDir, name)
		f, err := os.OpenFile(logPath, os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		data, err := json.MarshalIndent(itemElem, "", "    ")
		if err != nil {
			return err
		}

		//_, err = f.Write(data)
		err = os.WriteFile(f.Name(), data, os.ModePerm)
		if err != nil {
			return nil
		}
	}

	return nil
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
