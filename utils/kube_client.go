package utils

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	clientSet *kubernetes.Clientset
	config    *viper.Viper
	env       string
}

func NewKubeClient(env string) *KubeClient {
	conf := "/home/hhoroot/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", conf)
	if err != nil {
		panic(err)
	}
	c := LoadYaml()

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &KubeClient{
		clientSet: clientSet,
		config:    c,
		env:       env,
	}
}

func (k *KubeClient) ListDeployment() (replicas, availableReplicas int32) {
	fmt.Println("检查deployment replicas 状态")
	deploys, err := k.clientSet.AppsV1().Deployments(fmt.Sprintf("%s-%s", k.config.GetString("proj"), k.env)).
		List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("name=%s", k.config.GetString("app"))})
	if err != nil {
		panic(err)
	}

	return deploys.Items[0].Status.Replicas, deploys.Items[0].Status.AvailableReplicas
}
