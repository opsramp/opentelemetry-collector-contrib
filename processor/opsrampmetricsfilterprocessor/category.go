// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampmetricsfilterprocessor"

import (
	"fmt"
	"strings"
)

// categorySet is a bitmask of the categories a metric belongs to. A metric
// referenced by both a k8s_pod alert and a non-pod alert carries both bits.
type categorySet uint8

const (
	categoryPod categorySet = 1 << iota
	categoryCluster

	categoryAll = categoryPod | categoryCluster
)

// Config values accepted by metric_categories, compared case-insensitively.
const (
	podMetricCategory     = "podmetric"
	clusterMetricCategory = "clustermetric"
)

// podResourceType is the only resourceType mapped to the pod category; every
// other resourceType (including an absent one) is treated as cluster scope.
const podResourceType = "k8s_pod"

func (c categorySet) String() string {
	var names []string
	if c&categoryPod != 0 {
		names = append(names, "podMetric")
	}
	if c&categoryCluster != 0 {
		names = append(names, "clusterMetric")
	}
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ",")
}

// categoryForResourceType maps an alert definition resourceType to its category.
func categoryForResourceType(resourceType string) categorySet {
	if normalizeResourceType(resourceType) == podResourceType {
		return categoryPod
	}
	return categoryCluster
}

// normalizeResourceType tolerates the "Pod"/"k8s-pod"/"k8s_pod" spellings that
// appear across OpsRamp alert definition payloads.
func normalizeResourceType(resourceType string) string {
	normalized := strings.ToLower(strings.TrimSpace(resourceType))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	if normalized == "pod" {
		return podResourceType
	}
	return normalized
}

// parseCategoryMask converts the configured category names into a mask. An
// empty list selects every category, preserving the pre-existing behaviour.
func parseCategoryMask(categories []string) (categorySet, error) {
	if len(categories) == 0 {
		return categoryAll, nil
	}

	var mask categorySet
	for _, category := range categories {
		switch strings.ToLower(strings.TrimSpace(category)) {
		case podMetricCategory:
			mask |= categoryPod
		case clusterMetricCategory:
			mask |= categoryCluster
		default:
			return 0, fmt.Errorf("invalid metric_categories value %q (valid values: podMetric, clusterMetric)", category)
		}
	}
	return mask, nil
}
