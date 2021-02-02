package main

import (
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"

	kubeinformers "k8s.io/client-go/informers"
)

func main() {
	fmt.Println("Start") // doesn't work

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(clientset, time.Second*30)
	confgMapInformer := kubeInformerFactory.Core().V1().ConfigMaps().Informer()
	confgMapInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		DeleteFunc: onDelete,
		UpdateFunc: onUpdate,
	})

	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)
	for {
		time.Sleep(time.Second)
	}
}

func onAdd(obj interface{}) {
	onUpdate(obj, nil)
}

func onUpdate(obj interface{}, obj2 interface{}) {
	cm := obj.(*corev1.ConfigMap)
	fmt.Printf("ConfigMap was added/updated: %s. With labels: %s.\n", cm.Name, cm.Labels)
	fmt.Printf("Data: %s", cm.Data)
	customdashboardAuto := cm.Labels["instana_customdashboard_auto"] == "true"
	appLabel := cm.Labels["app"] != ""
	if customdashboardAuto && appLabel {
		fmt.Printf("===>: Will create dashboard for %s with %s \n", cm.Name, cm.Labels)
	}
}
func onDelete(obj interface{}) {
	//pod := obj.(*corev1.Pod)
	//fmt.Printf("pod delete: %s \n", pod.Name)
}
