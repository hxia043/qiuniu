package k8splugin

import (
	"fmt"
	"github/hxia043/qiuniu/internal/pkg/helmclient"
	"io/ioutil"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

func ValidateKubeconfig(kubeconfig []byte) error {
	config, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return err
	}

	if err = clientcmd.Validate(*config); err != nil {
		return err
	}

	return nil
}

func ParseKubeconfig(kubeconfig []byte) (string, string, string, error) {
	fmt.Println("Info: parse kubeconfig start...")

	clientGetter, err := helmclient.NewRESTClientGetter(kubeconfig)
	if err != nil {
		return "", "", "", err
	}

	clientConfig, err := clientGetter.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		return "", "", "", err
	}

	namespace, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return "", "", "", err
	}

	fmt.Println("Info: parse kubeconfig finished.")

	return namespace, clientConfig.Host, clientConfig.BearerToken, nil
}

func TransferKubeconfigFromPath(path string) ([]byte, error) {
	fmt.Println("Info: transfer kubeconfig start...")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	kubeconfigBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.Load(kubeconfigBytes)
	if err != nil {
		return nil, err
	}

	if err = clientcmd.Validate(*config); err != nil {
		return nil, err
	}

	fmt.Println("Info: transfer kubeconfig finished.")

	return kubeconfigBytes, nil
}
