/*
Copyright 2022.

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

package core

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	customcorev1 "github.com/Shaad7/bookstore-controller-kubebuilder/apis/core/v1"
)

// BookstoreReconciler reconciles a Bookstore object
type BookstoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.gopher.com,resources=bookstores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.gopher.com,resources=bookstores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.gopher.com,resources=bookstores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bookstore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *BookstoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here

	fmt.Println("Reconciling")

	bookstore := &customcorev1.Bookstore{}
	err := r.Client.Get(ctx, req.NamespacedName, bookstore)
	if err != nil {
		log.Error(err, "Bookstore Not found")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploymentName := bookstore.Spec.Name + "-deployment"
	deployment := &appsv1.Deployment{}

	depNamespaceName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      deploymentName,
	}
	fmt.Println(req.NamespacedName)
	fmt.Println(depNamespaceName)

	err = r.Client.Get(ctx, depNamespaceName, deployment)
	if errors.IsNotFound(err) {
		// Create deployment
		fmt.Printf("Creating Deployment\n")
		if err = r.Client.Create(ctx, newDeployment(bookstore)); err != nil {
			log.Error(err, "error creating deployment")
			return ctrl.Result{}, err
		} else {
			fmt.Println(deploymentName + "created")
		}
	} else if err != nil {
		log.Error(err, "Unable to fetch Deployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else if *bookstore.Spec.Replicas != *deployment.Spec.Replicas {

		fmt.Printf("Bookstore Replicas %d , Deployment Replicas %d. .. Updating", *bookstore.Spec.Replicas, *deployment.Spec.Replicas)
		if err = r.Client.Update(ctx, newDeployment(bookstore)); err != nil {
			log.Error(err, "error updating deployment")
			return ctrl.Result{}, err
		} else {
			fmt.Println("Updated deployment : " + deploymentName)
		}
	}

	// Update Status
	if err = r.updateBookstoreStatus(ctx, bookstore, deployment); err != nil {
		log.Error(err, "Cannot update status")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

func (r *BookstoreReconciler) updateBookstoreStatus(ctx context.Context, bookstore *customcorev1.Bookstore, deployment *appsv1.Deployment) error {
	bookstoreCopy := bookstore.DeepCopy()
	bookstoreCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	if err := r.Status().Update(ctx, bookstoreCopy); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookstoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&customcorev1.Bookstore{}).
		Complete(r)
}
