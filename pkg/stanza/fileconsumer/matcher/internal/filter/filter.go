// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package filter // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/matcher/internal/filter"

import (
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/multierr"
)

type Option interface {
	// Returned error is for explanatory purposes only.
	// All options will be called regardless of error.
	apply([]*item) ([]*item, error)
}

func Filter(values []string, regex *regexp.Regexp, opts ...Option) ([]string, error) {
	var errs error
	items := make([]*item, 0, len(values))
	for _, value := range values {
		it, err := newItem(value, regex)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		items = append(items, it)
	}
	for _, opt := range opts {
		var applyErr error
		items, applyErr = opt.apply(items)
		errs = multierr.Append(errs, applyErr)
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.value)
	}
	return result, errs
}

func FilterFiles(values []string, includeregex []*regexp.Regexp, excluderegex []*regexp.Regexp) ([]string, error) {

	result := []string{}
	var errs error

	for _, value := range values {

		parts := strings.Split(value, "/")
		parts = strings.SplitN(parts[4], "_", 2)
		namespace := parts[0]

		if len(includeregex) == 0 && len(excluderegex) == 0 {

			result = append(result, value)

		} else if len(includeregex) > 0 && len(excluderegex) == 0 {

			_, err := callIncludeFunc(namespace, includeregex)
			if err != nil {
				errs = multierr.Append(errs, err)
				continue
			} else {
				result = append(result, value)
			}

		} else if len(excluderegex) > 0 && len(includeregex) == 0 {

			_, err := callExcludeFunc(namespace, excluderegex)
			if err != nil {
				errs = multierr.Append(errs, err)
				continue
			} else {
				result = append(result, value)
			}

		} else {
			for _, includeregexVar := range includeregex {
				if includeregexVar != nil {
					match := includeregexVar.MatchString(namespace)
					if match {
						_, err := callExcludeFunc(namespace, excluderegex)
						if err != nil {
							errs = multierr.Append(errs, err)
							continue
						} else {
							result = append(result, value)
						}
					}
				}
			}
		}
	}
	return result, errs
}

type item struct {
	value    string
	captures map[string]string

	// Used when an Option is unable to interpret the value.
	// For example, a numeric sort may fail to parse the value as a number.
	err error
}

func callIncludeFunc(value string, includeregexlist []*regexp.Regexp) (string, error) {

	for _, includeregex := range includeregexlist {
		if includeregex != nil {
			match := includeregex.MatchString(value)
			if match {
				return value, nil
			}
		}
	}

	return "", fmt.Errorf("'%s' does not match with whole include regex list hence not considering this namsepace\n", value)
}

func callExcludeFunc(value string, excluderegexlist []*regexp.Regexp) (string, error) {

	for _, excluderegex := range excluderegexlist { //kubeing=   include:{^kube} exclude:{ing$,}
		if excluderegex != nil {
			match := excluderegex.MatchString(value)
			if match {
				return "", fmt.Errorf("'%s' does match with exclude regex hence not considering this namsepace\n", value)
			}
		}
	}

	return value, nil
}

func newItem(value string, regex *regexp.Regexp) (*item, error) {
	if regex == nil {
		return &item{
			value: value,
		}, nil
	}

	match := regex.FindStringSubmatch(value)
	if match == nil {
		return nil, fmt.Errorf("'%s' does not match regex", value)
	}
	it := &item{
		value:    value,
		captures: make(map[string]string),
	}
	for i, name := range regex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		it.captures[name] = match[i]
	}
	return it, nil
}
