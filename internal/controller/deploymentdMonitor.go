package controller

import (
	"fmt"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/tools/cache"
)

type DeploymentMonitor struct {
	deploymentInformer coreinformers.DeploymentInformer
	informerFactory    informers.SharedInformerFactory
}

func (r *DeploymentMonitor) deploymentAdd(obj interface{}) {
}

func (r *DeploymentMonitor) deploymentUpdate(obj interface{}, newObj interface{}) {
}

func (r *DeploymentMonitor) deploymentDelete(obj interface{}) {
}

func (c *DeploymentMonitor) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.deploymentInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to sync")
	}
	return nil
}

func NewDeploymentMonitor(informerFactory informers.SharedInformerFactory) (*DeploymentMonitor, error) {

	deploymentInformer := informerFactory.Apps().V1().Deployments()
	c := &DeploymentMonitor{
		informerFactory:    informerFactory,
		deploymentInformer: deploymentInformer,
	}

	_, err := deploymentInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.deploymentAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.deploymentUpdate,
			// Called on resource deletion.
			DeleteFunc: c.deploymentDelete,
		},
	)
	return c, err
}
