// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"fmt"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	seedmanagementv1alpha1 "github.com/gardener/gardener/pkg/apis/seedmanagement/v1alpha1"

	constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"

	"github.com/prometheus/client_golang/prometheus"
)

const unknown = "unknown"

var (
	// ScrapeFailures is a metric, which counts the amount scrape issues grouped by kind.
	ScrapeFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "garden_scrape_failure_total",
		Help: "Total count of scraping failures, grouped by kind/group of metric(s)",
	}, []string{"kind"})
)

func mapConditionStatus(status gardencorev1beta1.ConditionStatus) float64 {
	switch status {
	case gardencorev1beta1.ConditionTrue:
		return 1
	case gardencorev1beta1.ConditionFalse:
		return 0
	case gardencorev1beta1.ConditionProgressing:
		return 2
	default:
		return -1
	}
}

func usedAsSeed(shoot *gardencorev1beta1.Shoot, managedSeeds []*seedmanagementv1alpha1.ManagedSeed) bool {
	if shoot.Namespace != constants.GardenNamespace {
		return false
	}
	for _, ms := range managedSeeds {
		if ms.Spec.Shoot.Name == shoot.Name && ms.Namespace == shoot.Namespace {
			return true
		}
	}

	return false
}

func findProject(projects []*gardencorev1beta1.Project, match string) (*string, error) {
	var projectName string
	for _, project := range projects {
		if project.Spec.Namespace != nil && *project.Spec.Namespace == match {
			projectName = project.Name
			break
		}
	}
	if projectName == "" {
		return nil, fmt.Errorf("no project found for shoot %s", match)
	}
	return &projectName, nil
}
