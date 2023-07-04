/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/client-go/kubernetes"

	clusterv1 "cdx.foc/clusterwatch/api/v1"
)

// ClusterWatchNamespaceReconciler reconciles a ClusterWatchNamespace object
type BuildAgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	log    logr.Logger
	client *kubernetes.Clientset
}

//+kubebuilder:rbac:groups=cluster.cdx.foc,resources=clusterwatchnamespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.cdx.foc,resources=clusterwatchnamespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.cdx.foc,resources=clusterwatchnamespaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterWatchNamespace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *BuildAgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.log = log.FromContext(ctx)
	r.log.Info("Starting up BuildAgentReconciler")

	var cns clusterv1.ClusterWatchNamespace

	if err := r.Get(ctx, req.NamespacedName, &cns); err != nil {
		r.log.Error(err, "Unable to obtain crds created for cdx cluster watcher instance.")
	}

	r.SetupClientConfig()

	r.InitRestartPods(cns)

	var syncPeriod = 300 * time.Second
	scheduledResult := ctrl.Result{RequeueAfter: syncPeriod}
	return scheduledResult, nil
}

func (r *BuildAgentReconciler) SetupClientConfig() {
	configInstance := ClientConfig{}
	config := configInstance.GetConfig()
	r.client = kubernetes.NewForConfigOrDie(config)
}

func (r *BuildAgentReconciler) InitRestartPods(cns clusterv1.ClusterWatchNamespace) {

	for _, targetedNamespace := range cns.Spec.BuildAgentNamespaces {

		pods, err := r.client.CoreV1().Pods(targetedNamespace).List(context.TODO(), v1.ListOptions{})

		if err != nil {
			r.log.Info(fmt.Sprintf("Unable get pods from namespace %s ", targetedNamespace))
		}

		for _, pod := range pods.Items {
			restartCount := pod.Status.ContainerStatuses[0].RestartCount
			if restartCount > int32(cns.Spec.BuildAgentRestartMaxCount) {
				r.log.Info(fmt.Sprintf("Removing build agent: %s in namespace: %s ", pod.Name, targetedNamespace))
				//r.client.CoreV1().Pods(targetedNamespace).Delete(context.TODO(), pod.Name, v1.DeleteOptions{})
			}
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuildAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.ClusterWatchNamespace{}).
		Complete(r)
}
