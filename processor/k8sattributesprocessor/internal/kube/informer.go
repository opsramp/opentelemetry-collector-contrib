// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kube // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/kube"

import (
	"context"
	"time"

	api_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/tools/cache"
)

const kubeSystemNamespace = "kube-system"

// InformerProvider defines a function type that returns a new SharedInformer. It is used to
// allow passing custom shared informers to the watch client.
type InformerProvider func(
	client kubernetes.Interface,
	namespace string,
	labelSelector labels.Selector,
	fieldSelector fields.Selector,
) cache.SharedInformer

// InformerProviderNamespace defines a function type that returns a new SharedInformer. It is used to
// allow passing custom shared informers to the watch client for fetching namespace objects.
type InformerProviderNamespace func(
	client metadata.Interface,
) cache.SharedInformer

// InformerProviderWorkload defines a function type that returns a new SharedInformer. It is used to
// allow passing custom shared informers to the watch client.
// It's used for high-level workloads such as ReplicaSets, Deployments, DaemonSets, StatefulSets or Jobs
type InformerProviderWorkload func(
	client metadata.Interface,
	namespace string,
) cache.SharedInformer

func newSharedInformer(
	client kubernetes.Interface,
	namespace string,
	ls labels.Selector,
	fs fields.Selector,
	watchSyncPeriod time.Duration,
) cache.SharedInformer {
	return newSharedInformerWithObserver(client, namespace, ls, fs, watchSyncPeriod, nil)
}

func newSharedInformerWithObserver(
	client kubernetes.Interface,
	namespace string,
	ls labels.Selector,
	fs fields.Selector,
	watchSyncPeriod time.Duration,
	observer *k8sAPICallObserver,
) cache.SharedInformer {
	informer := cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc:  informerListFuncWithSelectorsWithObserver(client, namespace, ls, fs, observer),
			WatchFuncWithContext: informerWatchFuncWithSelectorsWithObserver(client, namespace, ls, fs, observer),
		},
		&api_v1.Pod{},
		watchSyncPeriod,
	)
	return informer
}

func informerListFuncWithSelectors(client kubernetes.Interface, namespace string, ls labels.Selector, fs fields.Selector) cache.ListWithContextFunc {
	return informerListFuncWithSelectorsWithObserver(client, namespace, ls, fs, nil)
}

func informerListFuncWithSelectorsWithObserver(client kubernetes.Interface, namespace string, ls labels.Selector, fs fields.Selector, observer *k8sAPICallObserver) cache.ListWithContextFunc {
	return func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		observer.Record("list", "pods", namespace)
		opts.LabelSelector = ls.String()
		opts.FieldSelector = fs.String()
		return client.CoreV1().Pods(namespace).List(ctx, opts)
	}
}

func informerWatchFuncWithSelectors(client kubernetes.Interface, namespace string, ls labels.Selector, fs fields.Selector) cache.WatchFuncWithContext {
	return informerWatchFuncWithSelectorsWithObserver(client, namespace, ls, fs, nil)
}

func informerWatchFuncWithSelectorsWithObserver(client kubernetes.Interface, namespace string, ls labels.Selector, fs fields.Selector, observer *k8sAPICallObserver) cache.WatchFuncWithContext {
	return func(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
		observer.Record("watch", "pods", namespace)
		opts.LabelSelector = ls.String()
		opts.FieldSelector = fs.String()
		return client.CoreV1().Pods(namespace).Watch(ctx, opts)
	}
}

// newKubeSystemSharedInformer watches only kube-system namespace
func newKubeSystemSharedInformer(
	client metadata.Interface,
	watchSyncPeriod time.Duration,
) cache.SharedInformer {
	return newKubeSystemSharedInformerWithObserver(client, watchSyncPeriod, nil)
}

func newKubeSystemSharedInformerWithObserver(
	client metadata.Interface,
	watchSyncPeriod time.Duration,
	observer *k8sAPICallObserver,
) cache.SharedInformer {
	gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
	return cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc: func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
				observer.Record("list", gvr.Resource, kubeSystemNamespace)
				opts.FieldSelector = fields.OneTermEqualSelector("metadata.name", kubeSystemNamespace).String()
				return client.Resource(gvr).List(ctx, opts)
			},
			WatchFuncWithContext: func(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
				observer.Record("watch", gvr.Resource, kubeSystemNamespace)
				opts.FieldSelector = fields.OneTermEqualSelector("metadata.name", kubeSystemNamespace).String()
				return client.Resource(gvr).Watch(ctx, opts)
			},
		},
		&metav1.PartialObjectMetadata{},
		watchSyncPeriod,
	)
}

func newNamespaceSharedInformer(
	client metadata.Interface,
	watchSyncPeriod time.Duration,
) cache.SharedInformer {
	return newNamespaceSharedInformerWithObserver(client, watchSyncPeriod, nil)
}

func newNamespaceSharedInformerWithObserver(
	client metadata.Interface,
	watchSyncPeriod time.Duration,
	observer *k8sAPICallObserver,
) cache.SharedInformer {
	gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
	return cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc:  metadataListFuncWithObserver(client, gvr, "", observer),
			WatchFuncWithContext: metadataWatchFuncWithObserver(client, gvr, "", observer),
		},
		&metav1.PartialObjectMetadata{},
		watchSyncPeriod,
	)
}

func newReplicaSetSharedInformer(client metadata.Interface, namespace string, watchSyncPeriod time.Duration) cache.SharedInformer {
	return newReplicaSetSharedInformerWithObserver(client, namespace, watchSyncPeriod, nil)
}

func newReplicaSetSharedInformerWithObserver(client metadata.Interface, namespace string, watchSyncPeriod time.Duration, observer *k8sAPICallObserver) cache.SharedInformer {
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"}
	return cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc:  metadataListFuncWithObserver(client, gvr, namespace, observer),
			WatchFuncWithContext: metadataWatchFuncWithObserver(client, gvr, namespace, observer),
		},
		&metav1.PartialObjectMetadata{},
		watchSyncPeriod,
	)
}

func newDeploymentSharedInformer(client metadata.Interface, namespace string, watchSyncPeriod time.Duration) cache.SharedInformer {
	return newDeploymentSharedInformerWithObserver(client, namespace, watchSyncPeriod, nil)
}

func newDeploymentSharedInformerWithObserver(client metadata.Interface, namespace string, watchSyncPeriod time.Duration, observer *k8sAPICallObserver) cache.SharedInformer {
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	return cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc:  metadataListFuncWithObserver(client, gvr, namespace, observer),
			WatchFuncWithContext: metadataWatchFuncWithObserver(client, gvr, namespace, observer),
		},
		&metav1.PartialObjectMetadata{},
		watchSyncPeriod,
	)
}

func newStatefulSetSharedInformer(client metadata.Interface, namespace string, watchSyncPeriod time.Duration) cache.SharedInformer {
	return newStatefulSetSharedInformerWithObserver(client, namespace, watchSyncPeriod, nil)
}

func newStatefulSetSharedInformerWithObserver(client metadata.Interface, namespace string, watchSyncPeriod time.Duration, observer *k8sAPICallObserver) cache.SharedInformer {
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}
	return cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc:  metadataListFuncWithObserver(client, gvr, namespace, observer),
			WatchFuncWithContext: metadataWatchFuncWithObserver(client, gvr, namespace, observer),
		},
		&metav1.PartialObjectMetadata{},
		watchSyncPeriod,
	)
}

func newDaemonSetSharedInformer(client metadata.Interface, namespace string, watchSyncPeriod time.Duration) cache.SharedInformer {
	return newDaemonSetSharedInformerWithObserver(client, namespace, watchSyncPeriod, nil)
}

func newDaemonSetSharedInformerWithObserver(client metadata.Interface, namespace string, watchSyncPeriod time.Duration, observer *k8sAPICallObserver) cache.SharedInformer {
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}
	return cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc:  metadataListFuncWithObserver(client, gvr, namespace, observer),
			WatchFuncWithContext: metadataWatchFuncWithObserver(client, gvr, namespace, observer),
		},
		&metav1.PartialObjectMetadata{},
		watchSyncPeriod,
	)
}

func newJobSharedInformer(client metadata.Interface, namespace string, watchSyncPeriod time.Duration) cache.SharedInformer {
	return newJobSharedInformerWithObserver(client, namespace, watchSyncPeriod, nil)
}

func newJobSharedInformerWithObserver(client metadata.Interface, namespace string, watchSyncPeriod time.Duration, observer *k8sAPICallObserver) cache.SharedInformer {
	gvr := schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}
	return cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc:  metadataListFuncWithObserver(client, gvr, namespace, observer),
			WatchFuncWithContext: metadataWatchFuncWithObserver(client, gvr, namespace, observer),
		},
		&metav1.PartialObjectMetadata{},
		watchSyncPeriod,
	)
}

func newNodeSharedInformer(client metadata.Interface, nodeName string, watchSyncPeriod time.Duration) cache.SharedInformer {
	return newNodeSharedInformerWithObserver(client, nodeName, watchSyncPeriod, nil)
}

func newNodeSharedInformerWithObserver(client metadata.Interface, nodeName string, watchSyncPeriod time.Duration, observer *k8sAPICallObserver) cache.SharedInformer {
	gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}
	return cache.NewSharedInformer(
		&cache.ListWatch{
			ListWithContextFunc: func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
				observer.Record("list", gvr.Resource, nodeName)
				if nodeName != "" {
					opts.FieldSelector = fields.OneTermEqualSelector("metadata.name", nodeName).String()
				}
				return client.Resource(gvr).List(ctx, opts)
			},
			WatchFuncWithContext: func(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
				observer.Record("watch", gvr.Resource, nodeName)
				if nodeName != "" {
					opts.FieldSelector = fields.OneTermEqualSelector("metadata.name", nodeName).String()
				}
				return client.Resource(gvr).Watch(ctx, opts)
			},
		},
		&metav1.PartialObjectMetadata{},
		watchSyncPeriod,
	)
}

func metadataListFunc(mc metadata.Interface, gvr schema.GroupVersionResource, namespace string) cache.ListWithContextFunc {
	return metadataListFuncWithObserver(mc, gvr, namespace, nil)
}

func metadataListFuncWithObserver(mc metadata.Interface, gvr schema.GroupVersionResource, namespace string, observer *k8sAPICallObserver) cache.ListWithContextFunc {
	return func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		observer.Record("list", gvr.Resource, namespace)
		return mc.Resource(gvr).Namespace(namespace).List(ctx, opts)
	}
}

func metadataWatchFunc(mc metadata.Interface, gvr schema.GroupVersionResource, namespace string) cache.WatchFuncWithContext {
	return metadataWatchFuncWithObserver(mc, gvr, namespace, nil)
}

func metadataWatchFuncWithObserver(mc metadata.Interface, gvr schema.GroupVersionResource, namespace string, observer *k8sAPICallObserver) cache.WatchFuncWithContext {
	return func(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
		observer.Record("watch", gvr.Resource, namespace)
		return mc.Resource(gvr).Namespace(namespace).Watch(ctx, opts)
	}
}
