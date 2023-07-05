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
const MaxAllowedDaysWithoutRaisingPR float64 = 3
const LabelMatchingValue string = "required"

// ClusterWatchNamespaceReconciler reconciles a ClusterWatchNamespace object
type ClusterWatchNamespaceReconciler struct {
	client.Client
	Scheme                   *runtime.Scheme
	log                      logr.Logger
	namespaceNeedsReviewList []string
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

	var cns clusterv1.ClusterWatchNamespace

	if err := r.Get(ctx, req.NamespacedName, &cns); err != nil {
		r.log.Error(err, "Unable to obtain crds created for cdx cluster watcher instance.")
	} else {
		r.GetNamespaceWithRequiredPRTag(cns)
	}

	// Next cycle wait time //
	var syncPeriod = 300 * time.Second
	scheduledResult := ctrl.Result{RequeueAfter: syncPeriod}
	return scheduledResult, nil
}

func (r *ClusterWatchNamespaceReconciler) GetNamespaceWithRequiredPRTag(cns clusterv1.ClusterWatchNamespace) {

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

		labelValue, ok := s.Labels[NamespaceAutomationMarker]

		if ok {

			if labelValue == LabelMatchingValue {
				r.log.Info(fmt.Sprintf("Matching namespace found, checking namespace creation time: %s", s.Name))
				// Get creation time //
				r.log.Info(s.CreationTimestamp.String())
				timeDriftDays := time.Now().Sub(s.CreationTimestamp.Time).Hours() / 24

				var maxDaysPR float64 = 0

				if cns.Spec.RequiredPRNamespaceMaxWaitDays == 0 {
					maxDaysPR = MaxAllowedDaysWithoutRaisingPR
				} else {
					maxDaysPR = cns.Spec.RequiredPRNamespaceMaxWaitDays
				}

				if timeDriftDays > maxDaysPR {
					logMessageNamespaceExceeded := fmt.Sprintf("%s namespace has exceeded max number of days. TimeDiff %f", s.Name, timeDriftDays)

					r.log.Info(logMessageNamespaceExceeded)
					r.namespaceNeedsReviewList = append(r.namespaceNeedsReviewList, s.Name)

					// send message to a webhook //
					sm := NewSlackMessenger(cns.Spec.NotificationWebHookEndpoint)
					nw := NewNotificationWorker(sm)
					nw.SendMessage(logMessageNamespaceExceeded)
				} else {
					r.log.Info(fmt.Sprintf("Namespace found but does not meet the PR requirement criteria of %f", MaxAllowedDaysWithoutRaisingPR))
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
