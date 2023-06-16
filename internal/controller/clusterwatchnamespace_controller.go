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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "cdx.foc/clusterwatch/api/v1"
)

const NamespaceAutomationMarker string = "backstage-namespace-pr"
const MaxAllowedDaysWithoutRaisingPR float64 = 20

// ClusterWatchNamespaceReconciler reconciles a ClusterWatchNamespace object
type ClusterWatchNamespaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	log    logr.Logger
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

	r.log = log.FromContext(ctx)
	r.log.Info("Reconciling clusterwatchnamespace")

	var cns clusterv1.ClusterWatchNamespace

	if err := r.Get(ctx, req.NamespacedName, &cns); err != nil {
		r.log.Error(err, "Unable to obtain crd created for cluster watcher instance")
	} else {
		r.GetNamespaceWithRequiredPRTag()
	}

	var syncPeriod = 300 * time.Second
	scheduledResult := ctrl.Result{RequeueAfter: syncPeriod}
	return scheduledResult, nil
}

func (r *ClusterWatchNamespaceReconciler) GetNamespaceWithRequiredPRTag() {

	// InClusterConfig
	config, err := rest.InClusterConfig()
	if err != nil {
		r.log.Error(err, "Fault in rest.InClusterConfig")
	}

	// home, _ := os.UserHomeDir()
	// kubeConfigPath := filepath.Join(home, ".kube", "config")
	// config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	// if err != nil {
	// 	panic(err)
	// }

	client := kubernetes.NewForConfigOrDie(config)

	// creates the clientset
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Error(err, "Error setup clientset")
	// }

	// Get the namespace with proper annotations //
	nslist, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		r.log.Error(err, "Error listing namespace in the cluster.")
	}

	for _, s := range nslist.Items {

		r.log.Info(fmt.Sprintf("Performing checks on namespace label: %s", s.Name))

		labelValue, ok := s.Labels[NamespaceAutomationMarker]

		if ok {

			r.log.Info(fmt.Sprintf("Found a namespace with label: backstage-pr-automation"))
			if labelValue == "required" {
				// Get creation time //
				r.log.Info(s.CreationTimestamp.String())
				timeDrift := time.Now().Sub(s.CreationTimestamp.Time).Hours() / 24
				r.log.Info(fmt.Sprintf("TimeDiff %f", timeDrift))

				if timeDrift > MaxAllowedDaysWithoutRaisingPR {
					r.log.Info("Namespace has exceeded max number of days")
				} else {
					r.log.Info("Namespace with matching label was found but does not match the max days:")
				}
			}
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterWatchNamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.ClusterWatchNamespace{}).
		Complete(r)
}
