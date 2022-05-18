package k8splugin

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

func ValidateKubeconfig(kubeconfig string) error {
	config, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		return err
	}

	if err = clientcmd.Validate(*config); err != nil {
		return err
	}

	return nil
}

func TransferKubeconfigToBase64(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	kubeconfigBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	config, err := clientcmd.Load(kubeconfigBytes)
	if err != nil {
		return "", err
	}

	if err = clientcmd.Validate(*config); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(kubeconfigBytes), nil
}
