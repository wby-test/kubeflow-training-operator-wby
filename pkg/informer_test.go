package pkg

import (
	"fmt"
	"testing"
	"time"

	kubefloworgv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	clientset "github.com/kubeflow/training-operator/pkg/client/clientset/versioned"
	versioned "github.com/kubeflow/training-operator/pkg/client/clientset/versioned"
	"github.com/kubeflow/training-operator/pkg/client/informers/externalversions"
	v1 "github.com/kubeflow/training-operator/pkg/client/informers/externalversions/kubeflow.org/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type TfJob struct {
}

func (p *TfJob) OnAdd(obj interface{}) {
	fmt.Println("TfJob OnAdd")
}

func (p *TfJob) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("TfJob OnUpdate")
}

func (p *TfJob) OnDelete(obj interface{}) {
	fmt.Println("TfJob OnDelete")
}

type PytorchJob struct {
}

func (p *PytorchJob) OnAdd(obj interface{}) {
	fmt.Println("PytorchJob OnAdd")
}

func (p *PytorchJob) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("PytorchJob OnUpdate")
}

func (p *PytorchJob) OnDelete(obj interface{}) {
	fmt.Println("PytorchJob OnDelete")
}

func common(v versioned.Interface, t time.Duration) cache.SharedIndexInformer {
	return v1.NewFilteredTFJobInformer(v, "public-resource", t, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, nil)
}

func common1(v versioned.Interface, t time.Duration) cache.SharedIndexInformer {
	return v1.NewFilteredPyTorchJobInformer(v, "public-resource", t, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, nil)
}
func RunInformer() {
	client, err := clientcmd.BuildConfigFromFlags("", "/Users/wby/.kube/config")
	if err != nil {
		panic(err)
	}

	clientSet, err := clientset.NewForConfig(client)
	if err != nil {
		panic(err)
	}

	stopTfCh := make(chan struct{})
	defer close(stopTfCh)
	pytorchCh := make(chan struct{})
	defer close(pytorchCh)

	sharedInformers := externalversions.NewSharedInformerFactory(clientSet, time.Minute)
	tfInformer := sharedInformers.InformerFor(&kubefloworgv1.TFJob{}, common)
	tfInformer.AddEventHandler(&TfJob{})

	sharedInformers1 := externalversions.NewSharedInformerFactory(clientSet, time.Minute)
	pytorchInformer := sharedInformers1.InformerFor(&kubefloworgv1.PyTorchJob{}, common)
	pytorchInformer.AddEventHandler(&PytorchJob{})

	tfInformer.Run(stopTfCh)
	pytorchInformer.Run(pytorchCh)
}

func RunPytorchInformer() {
	client, err := clientcmd.BuildConfigFromFlags("", "/Users/wby/.kube/config")
	if err != nil {
		panic(err)
	}

	clientSet, err := clientset.NewForConfig(client)
	if err != nil {
		panic(err)
	}

	pytorchCh := make(chan struct{})
	defer close(pytorchCh)

	sharedInformers1 := externalversions.NewSharedInformerFactory(clientSet, time.Minute)
	pytorchInformer := sharedInformers1.InformerFor(&kubefloworgv1.PyTorchJob{}, common1)
	pytorchInformer.AddEventHandler(&PytorchJob{})

	pytorchInformer.Run(pytorchCh)
}

func Test_RunInformer(t *testing.T) {
	go RunInformer()
	go RunPytorchInformer()
	time.Sleep(time.Hour)
}
