/*
Copyright 2024 The Kubernetes Authors.

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

package cache

import (
	"k8s.io/client-go/tools/cache"
	k8smetrics "k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

// Metrics subsystem and keys used by the workqueue.
const (
	CacheSubsystem          = "cache"
	PendingNotificationsKey = "pending_notifications"
)

var (
	pendingNotifications = k8smetrics.NewGaugeVec(&k8smetrics.GaugeOpts{
		Subsystem:      CacheSubsystem,
		Name:           PendingNotificationsKey,
		StabilityLevel: k8smetrics.ALPHA,
		Help:           "Current number of pending notifications",
	}, []string{"name", "resources"})

	metrics = []k8smetrics.Registerable{
		pendingNotifications,
	}
)

func init() {
	for _, m := range metrics {
		legacyregistry.MustRegister(m)
	}
	cache.SetMetricsProvider(prometheusMetricsProvider{})
}

type prometheusMetricsProvider struct{}

func (prometheusMetricsProvider) NewPendingNotificationsMetric(name, resources string) cache.GaugeMetric {
	return pendingNotifications.WithLabelValues(name, resources)
}
