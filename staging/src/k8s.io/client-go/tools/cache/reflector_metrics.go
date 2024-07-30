/*
Copyright 2016 The Kubernetes Authors.

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

// This file provides abstractions for setting the provider (e.g., prometheus)
// of metrics.

package cache

import (
	"sync"
)

// GaugeMetric represents a single numerical value that can arbitrarily go up
// and down.
type GaugeMetric interface {
	Set(float64)
}

// CounterMetric represents a single numerical value that only ever
// goes up.
type CounterMetric interface {
	Inc()
}

// HistrgramMetric captures individual observations.
type HistrgramMetric interface {
	Observe(float64)
}

type noopMetric struct{}

func (noopMetric) Inc()            {}
func (noopMetric) Dec()            {}
func (noopMetric) Observe(float64) {}
func (noopMetric) Set(float64)     {}

type MetricsProvider interface {
	// the informer name
	// NewStoredItemMetric(name string) GaugeMetric
	// NewQueuedItemMetric(name string) GaugeMetric

	// the eventHandler name
	NewPendingNotificationsMetric(name, resources string) GaugeMetric
	// NewRingGrowingMetric(name string) GaugeMetric
	// NewPrcoessDurationMetric(name string) HistogramMetric
}

type informerMetrics struct {
	// clock clock.Clock

	// // total number of item in store
	// numbernOfStoredItem GaugeMetric
	// // total number of item in queue
	// numberOfQueuedItem GaugeMetric

	// each eventHandler metrics
	eventHandlerMetrics map[string]eventHandlerMetrics

	mutex sync.Mutex
}

type eventHandlerMetrics struct {

	// number of pending data
	numberOfPendingNotifications GaugeMetric

	// 	// size of RingGrowring data
	// 	sizeOfRingGrowing GaugeMetric

	// // how long processing an item from informer reflector
	// prcoessDuration HistogramMetric
}

type noopMetricsProvider struct{}

// func (noopMetricsProvider) NewListsMetric(name string) CounterMetric         { return noopMetric{} }
// func (noopMetricsProvider) NewListDurationMetric(name string) SummaryMetric  { return noopMetric{} }
// func (noopMetricsProvider) NewItemsInListMetric(name string) SummaryMetric   { return noopMetric{} }
// func (noopMetricsProvider) NewWatchesMetric(name string) CounterMetric       { return noopMetric{} }
// func (noopMetricsProvider) NewShortWatchesMetric(name string) CounterMetric  { return noopMetric{} }
// func (noopMetricsProvider) NewWatchDurationMetric(name string) SummaryMetric { return noopMetric{} }
// func (noopMetricsProvider) NewItemsInWatchMetric(name string) SummaryMetric  { return noopMetric{} }
//
//	func (noopMetricsProvider) NewLastResourceVersionMetric(name string) GaugeMetric {
//		return noopMetric{}
//	}
func (noopMetricsProvider) NewPendingNotificationsMetric(name, resources string) GaugeMetric {
	return noopMetric{}
}

func (m *eventHandlerMetrics) updatePendingNotificationsCount(count int64) {
	m.numberOfPendingNotifications.Set(float64(count))
}

var metricsFactory = struct {
	metricsProvider MetricsProvider
	setProviders    sync.Once
}{
	metricsProvider: noopMetricsProvider{},
}

// SetMetricsProvider sets the metrics provider
func SetMetricsProvider(metricsProvider MetricsProvider) {
	metricsFactory.setProviders.Do(func() {
		metricsFactory.metricsProvider = metricsProvider
	})
}
