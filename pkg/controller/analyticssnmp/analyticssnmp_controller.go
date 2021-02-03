package analyticssnmp

import (
	"context"
	"fmt"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// InstanceType is a string value for AnalyticsSnmp
var InstanceType = "analyticssnmp"

// Log is a default logger for AnalyticsSnmp
var Log = logf.Log.WithName("controller_" + InstanceType)

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.AnalyticsSnmpList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
		UpdateFunc: func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.MetaNew.GetNamespace()}
			list := &v1alpha1.AnalyticsSnmpList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.MetaNew.GetNamespace(),
					}})
				}
			}
		},
		DeleteFunc: func(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.AnalyticsSnmpList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
		GenericFunc: func(e event.GenericEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.AnalyticsSnmpList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
	}
	return appHandler
}

// Add adds the AnalyticsSnmp controller to the manager.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAnalyticsSnmp{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Manager:    mgr,
		Kubernetes: k8s.New(mgr.GetClient(), mgr.GetScheme()),
	}
}
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller.
	c, err := controller.New(InstanceType+"-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AnalyticsSnmp.
	if err = c.Watch(&source.Kind{Type: &v1alpha1.AnalyticsSnmp{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}
	serviceMap := map[string]string{"contrail_manager": InstanceType}
	srcPod := &source.Kind{Type: &corev1.Pod{}}
	podHandler := resourceHandler(mgr.GetClient())
	predInitStatus := utils.PodInitStatusChange(serviceMap)
	predPodIPChange := utils.PodIPChange(serviceMap)
	predInitRunning := utils.PodInitRunning(serviceMap)

	if err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Config{},
	}); err != nil {
		return err
	}

	if err = c.Watch(srcPod, podHandler, predPodIPChange); err != nil {
		return err
	}
	if err = c.Watch(srcPod, podHandler, predInitStatus); err != nil {
		return err
	}
	if err = c.Watch(srcPod, podHandler, predInitRunning); err != nil {
		return err
	}

	/*
		cassandraServiceMap := map[string]string{"contrail_manager": "cassandra"}
		predCassandraPodIPChange := utils.PodIPChange(cassandraServiceMap)
		if err = c.Watch(srcPod, podHandler, predCassandraPodIPChange); err != nil {
			return err
		}

		srcRabbitmq := &source.Kind{Type: &v1alpha1.Rabbitmq{}}
		rabbitmqHandler := resourceHandler(mgr.GetClient())
		predRabbitmqSizeChange := utils.RabbitmqActiveChange()
		if err = c.Watch(srcRabbitmq, rabbitmqHandler, predRabbitmqSizeChange); err != nil {
			return err
		}

		srcZookeeper := &source.Kind{Type: &v1alpha1.Zookeeper{}}
		zookeeperHandler := resourceHandler(mgr.GetClient())
		predZookeeperSizeChange := utils.ZookeeperActiveChange()
		if err = c.Watch(srcZookeeper, zookeeperHandler, predZookeeperSizeChange); err != nil {
			return err
		}

		srcSTS := &source.Kind{Type: &appsv1.StatefulSet{}}
		stsHandler := &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &v1alpha1.Config{},
		}
		stsPred := utils.STSStatusChange(utils.ConfigGroupKind())
		if err = c.Watch(srcSTS, stsHandler, stsPred); err != nil {
			return err
		}
	*/
	return nil
}

// blank assignment to verify that ReconcileAnalyticsSnmp implements reconcile.Reconciler.
var _ reconcile.Reconciler = &ReconcileAnalyticsSnmp{}

// ReconcileAnalyticsSnmp reconciles a AnalyticsSnmp object.
type ReconcileAnalyticsSnmp struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver.
	Client     client.Client
	Scheme     *runtime.Scheme
	Manager    manager.Manager
	Kubernetes *k8s.Kubernetes
}

// Reconcile reconciles AnalyticsSnmp.
func (r *ReconcileAnalyticsSnmp) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := Log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AnalyticsSnmp")

	// Get instance
	instance := &v1alpha1.AnalyticsSnmp{}
	if err := r.Client.Get(context.TODO(), request.NamespacedName, instance); err != nil && errors.IsNotFound(err) {
		reqLogger.Error(err, "Instance not found.")
		return reconcile.Result{}, nil
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		reqLogger.Info("Instance is deleting, skip reconcile.")
		return reconcile.Result{}, nil
	}

	// TODO: Add check cassandra, zookeeper, rabitmq active

	// Get or create configmaps
	_, _, err := r.GetOrCreateConfigMap(FullName("configmap", request), instance, request)
	if err != nil {
		reqLogger.Error(err, "ConfigMap not created.")
		return reconcile.Result{}, err
	}

	_, err = v1alpha1.CreateSecret(request.Name+"-secret-certificates", r.Client, r.Scheme, request, InstanceType, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet, err := v1alpha1.QuerySTS(FullName("statefulset", request), request.Namespace, r.Client)
	if err != nil {
		reqLogger.Error(err, "StatefulSet not found.")
		return reconcile.Result{}, nil
	}

	updateStatefulSetFlag := false
	if statefulSet != nil {
		// StatefulSet has been already created
		// TODO: check if some data changed: updateStatefulSetFlag = true|false
	}

	if statefulSet == nil || updateStatefulSetFlag {
		// StatefulSet haven't been created or need to be updated

		// Get basic stateful set
		statefulSet, err := GetStatefulsetFromYaml()
		if err != nil {
			reqLogger.Error(err, "Cant load the stateful set from yaml.")
			return reconcile.Result{}, nil
		}

		// Add common configuration to stateful set
		if err := v1alpha1.PrepareSTS(statefulSet, &instance.Spec.CommonConfiguration, InstanceType, request, r.Scheme, instance, true); err != nil {
			reqLogger.Error(err, "Cant prepare the stateful set.")
		}

		// Add volumes to stateful set
		v1alpha1.AddVolumesToIntendedSTS(statefulSet, map[string]string{
			FullName("configmap", request):     FullName("volume", request),
			certificates.SignerCAConfigMapName: request.Name + "-csr-signer-ca",
		})
		v1alpha1.AddSecretVolumesToIntendedSTS(statefulSet, map[string]string{
			request.Name + "-secret-certificates": request.Name + "-secret-certificates",
		})

		// Don't know what is it
		statefulSet.Spec.Template.Spec.Affinity = &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{{
					LabelSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{{
							Key:      InstanceType,
							Operator: "In",
							Values:   []string{request.Name},
						}},
					},
					TopologyKey: "kubernetes.io/hostname",
				}},
			},
		}

		// Manual settings for containers
		for idx := range statefulSet.Spec.Template.Spec.Containers {
			container := &statefulSet.Spec.Template.Spec.Containers[idx]
			instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
			if instanceContainer == nil {
				reqLogger.Info(fmt.Sprintf("There is no %s container in the manifect", container.Name))
				continue
			}

			// Add sleep command to container
			sleepCommand := []string{"bash", "-c", "while true; do echo Working; sleep 1; done"}
			container.Command = sleepCommand

			// Add image from manifest
			container.Image = instanceContainer.Image

			// Add volume mounts to container
			volumeMountList := []corev1.VolumeMount{}
			if len(container.VolumeMounts) > 0 {
				volumeMountList = container.VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      FullName("volume", request),
				MountPath: "/etc/contrailconfigmaps",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-secret-certificates",
				MountPath: "/etc/certificates",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-csr-signer-ca",
				MountPath: certificates.SignerCAMountPath,
			}
			volumeMountList = append(volumeMountList, volumeMount)
			container.VolumeMounts = volumeMountList
		}
		
		v1alpha1.CreateSTS(statefulSet, InstanceType, request, r.Client)
		v1alpha1.UpdateSTS(statefulSet, InstanceType, request, r.Client, "rolling")
	}

	// Check secrets to start containers (init containers running until pod `status` label not equal to `ready`)
	podIpList, podIpMap, err := v1alpha1.PodIPListAndIPMapFromInstance(InstanceType,
		&instance.Spec.CommonConfiguration,
		request,
		r.Client,
		true, true, false, false, false, false,
	)
	if err != nil {
		reqLogger.Error(err, "Pod list not found")
		return reconcile.Result{}, err
	}
	if len(podIpMap) > 0 {
		// Ensure certificates exists
		certSubjects := v1alpha1.PodsCertSubjects(podIpList,
			instance.Spec.CommonConfiguration.HostNetwork,
			v1alpha1.PodAlternativeIPs{},
		)
		crt := certificates.NewCertificate(r.Client, r.Scheme, instance, certSubjects, InstanceType)
		if err := crt.EnsureExistsAndIsSigned(); err != nil {
			reqLogger.Error(err, "Certificates for pod not exist.")
			return reconcile.Result{Requeue: true}, nil
		}
		// Set pod `status` label to `ready`
		if err = v1alpha1.SetPodsToReady(podIpList, r.Client); err != nil {
			reqLogger.Error(err, "Failed to set pods to ready")
			return reconcile.Result{}, err
		}		
	}

	return reconcile.Result{}, nil
}

// FullName ...
func FullName(name string, request reconcile.Request) string {
	return request.Name + "-" + InstanceType + "-" + name
}

// GetOrCreateConfigMap ...
func (r *ReconcileAnalyticsSnmp) GetOrCreateConfigMap(name string,
	instance *v1alpha1.AnalyticsSnmp,
	request reconcile.Request,
) (configMap *corev1.ConfigMap, isCreated bool, err error) {

	configMap, isCreated, err = v1alpha1.GetOrCreateConfigMap(name,
		r.Client,
		r.Scheme,
		request,
		InstanceType,
		instance,
	)
	return
}
