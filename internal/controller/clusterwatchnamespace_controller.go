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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	clusterv1 "cdx.foc/clusterwatch/api/v1"
)

// ClusterWatchNamespaceReconciler reconciles a ClusterWatchNamespace object
type ClusterWatchNamespaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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
func (r *ClusterWatchNamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("we are good")

	var cns clusterv1.ClusterWatchNamespace

	if err := r.Get(ctx, req.NamespacedName, &cns); err != nil {
		log.Error(err, "unable to fetch cluster watcher instance")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err, "Fault in rest.InClusterConfig")
	}
	// creates the clientsets
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "Error settign up clientset")
	}

	// Get the namespace with proper annotations //
	pods, err := clientset.CoreV1().Pods("").Lisst(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error(err, "Error getting namespace")
	}

	if pods != nil {
		log.Info("pod info ")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterWatchNamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.ClusterWatchNamespace{}).
		Complete(r)
}
