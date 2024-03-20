package metrics

import (
	"context"

	operatorv1alpha1 "github.com/gardener/gardener/pkg/apis/operator/v1alpha1"

	"github.com/prometheus/client_golang/prometheus"
)

// collectGardenMetrics collects metrics about the Garden resource(s).
func (c gardenMetricsCollector) collectGardenMetrics(ch chan<- prometheus.Metric) {
	c.client.WaitForCacheSync(c.ctx)
	gardenList := &operatorv1alpha1.GardenList{}
	if err := c.client.Client().List(context.Background(), gardenList); err != nil {
		ScrapeFailures.With(prometheus.Labels{"kind": "garden"}).Inc()
		return
	}

	for _, garden := range gardenList.Items {
		for _, condition := range garden.Status.Conditions {
			if condition.Type == "" {
				continue
			}

			metric, err := prometheus.NewConstMetric(
				c.descs[metricGardenGardenInfo],
				prometheus.GaugeValue,
				mapConditionStatus(condition.Status),
				[]string{
					garden.Name,
					string(condition.Type),
					string(garden.Status.LastOperation.Type),
				}...,
			)
			if err != nil {
				ScrapeFailures.With(prometheus.Labels{"kind": "garden"}).Inc()
				continue
			}
			ch <- metric
		}
	}
}
