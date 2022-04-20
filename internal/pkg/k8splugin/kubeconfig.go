package k8splugin

import (
	"encoding/base64"
	"io/ioutil"
	"os"
)

func TransferKubeConfigToBase64(kubeconfigPath string) (string, error) {
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return "", err
	}

	kubeconfigBytes, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(kubeconfigBytes), nil
}
