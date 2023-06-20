package controller

// factory := informers.NewSharedInformerFactory(client, time.Hour*24)
// 	controller, err := NewPodMonitor(factory)

// 	stop := make(chan struct{})
// 	err = controller.Run(stop)

// 	if err != nil {
// 		r.log.Error(err, "error")
// 	}
// 	select {}

import (
	"fmt"

	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type PodMonitor struct {
	podInformer     coreinformers.PodInformer
	informerFactory informers.SharedInformerFactory
}

func (c *PodMonitor) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to sync")
	}
	return nil
}

func (r *PodMonitor) podAdd(obj interface{}) {
}

func (r *PodMonitor) podUpdate(obj interface{}, newObj interface{}) {
}

func (r *PodMonitor) podDelete(obj interface{}) {
}

func NewPodMonitor(informerFactory informers.SharedInformerFactory) (*PodMonitor, error) {

	podInformer := informerFactory.Core().V1().Pods()

	c := &PodMonitor{
		informerFactory: informerFactory,
		podInformer:     podInformer,
	}

	_, err := podInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.podUpdate,
			// Called on resource deletion.
			DeleteFunc: c.podDelete,
		},
	)
	return c, err
}
