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
	v1beta1 "github.com/Neaj-Morshad-101/crd-kubebuilder-controller/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// KlusterReconciler reconciles a Kluster object
type KlusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=neajmorshad.dev.neajmorshad.com,resources=klusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=neajmorshad.dev.neajmorshad.com,resources=klusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=neajmorshad.dev.neajmorshad.com,resources=klusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *KlusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logs := log.FromContext(ctx)
	logs.WithValues("ReqName", req.Name, "ReqNamesapce", req.Namespace)

	// TODO(user): your logic here
	/*
		### 1: Load the Aadee by name

		We'll fetch the Aadee using our client.  All client methods take a
		context (to allow for cancellation) as their first argument, and the object
		in question as their last.  Get is a bit special, in that it takes a
		[`NamespacedName`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client?tab=doc#ObjectKey)
		as the middle argument (most don't have a middle argument, as we'll see
		below).

		Many client methods also take variadic options at the end.
	*/

	var kluster v1beta1.Kluster

	if err := r.Get(ctx, req.NamespacedName, &kluster); err != nil {
		fmt.Println(err, "unable to fetch kluster")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, nil
	}

	fmt.Println("Kluster Name", kluster.Name)

	// deploymentObject carry the all data of deployment in specific namespace and name
	var deploymentObject appsv1.Deployment
	deploymentName := "deployment-" + kluster.Name
	serviceName := "service-" + kluster.Name

	objectKey := client.ObjectKey{
		Namespace: req.Namespace,
		Name:      deploymentName,
	}

	if err := r.Get(ctx, objectKey, &deploymentObject); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("could not find existing Deployment for ", kluster.Name, " creating one...")
			err := r.Client.Create(ctx, newDeployment(&kluster, deploymentName))
			if err != nil {
				fmt.Printf("error while creteing depluoyment %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s deployment created...\n", kluster.Name)
			}

		} else {
			fmt.Printf("error fetchiong deploymewnt %s\n", err)
			return ctrl.Result{}, err
		}
	} else {
		if kluster.Spec.Replicas != nil && *kluster.Spec.Replicas != *deploymentObject.Spec.Replicas {
			fmt.Println(*kluster.Spec.Replicas, *deploymentObject.Spec.Replicas)
			fmt.Println("deployemnts replicas don't match...")
			deploymentObject.Spec.Replicas = kluster.Spec.Replicas
			if err := r.Update(ctx, &kluster); err != nil {
				fmt.Printf("error updating deployment %s\n", err)
				return ctrl.Result{}, err
			}
			fmt.Println("deployment updated")
		}
	}

	var serviceObject corev1.Service
	objectKey = client.ObjectKey{
		Namespace: req.Namespace,
		Name:      serviceName,
	}

	if err := r.Get(ctx, objectKey, &serviceObject); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("could not find existing Service for ", kluster.Name, " creating one...")
			err := r.Client.Create(ctx, newService(&kluster, serviceName, deploymentName))
			if err != nil {
				fmt.Printf("error while creating Service %s\n", err)
				return ctrl.Result{}, err
			} else {
				fmt.Printf("%s Service created...\n", kluster.Name)
			}

		} else {
			fmt.Printf("error fetchiong service %s\n", err)
			return ctrl.Result{}, err
		}
	} else {
		fmt.Println("Service exists")
	}

	if kluster.Spec.Replicas != nil && *kluster.Spec.Replicas != kluster.Status.AvailableReplicas {
		fmt.Println(*kluster.Spec.Replicas, kluster.Status.AvailableReplicas)
		kluster.Status.AvailableReplicas = *kluster.Spec.Replicas
		if err := r.Status().Update(ctx, &kluster); err != nil {
			fmt.Printf("error while updating the replica count %s\n", err)
			return ctrl.Result{}, err
		}
		fmt.Println("replica count updated")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KlusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Kluster{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
