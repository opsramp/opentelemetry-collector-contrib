// SPDX-License-Identifier: Apache-2.0

package filter // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/matcher/internal/filter"
import (
	"fmt"
	"go.uber.org/multierr"
	"regexp"
)

type filterRegex struct {
	IncludeRegex              []string
	ExcludeRegex              []string
	CapturePathSubstringRegex string
}

func (fr filterRegex) apply(items []*item) ([]*item, error) {
	filteredItems := make([]*item, 0, len(items))
	var errs error
	var capturedPathSubstring string

	// Compile the regular expression
	re := regexp.MustCompile(fr.CapturePathSubstringRegex)

	// Find a match and capture the group
	for _, item := range items {
		match := re.FindStringSubmatch(item.value)

		if len(match) > 1 {
			capturedPathSubstring = match[1]
		} else {
			continue
		}

		if len(fr.IncludeRegex) > 0 && len(fr.ExcludeRegex) == 0 {

			_, err := callIncludeFunc(capturedPathSubstring, fr.IncludeRegex)
			if err != nil {
				errs = multierr.Append(errs, err)
				continue
			} else {
				filteredItems = append(filteredItems, item)
			}

		} else if len(fr.ExcludeRegex) > 0 && len(fr.IncludeRegex) == 0 {

			_, err := callExcludeFunc(capturedPathSubstring, fr.ExcludeRegex)
			if err != nil {
				errs = multierr.Append(errs, err)
				continue
			} else {
				filteredItems = append(filteredItems, item)
			}

		} else {
			for _, includeregexVar := range fr.IncludeRegex {
				if includeregexVar != "" {
					re := regexp.MustCompile(includeregexVar)
					is_match := re.MatchString(capturedPathSubstring)
					if is_match {
						_, err := callExcludeFunc(capturedPathSubstring, fr.ExcludeRegex)
						if err != nil {
							errs = multierr.Append(errs, err)
							continue
						} else {
							filteredItems = append(filteredItems, item)
						}
					}
				}
			}
		}
	}
	return filteredItems, errs
}

// FilterRegex excludes or includes files whose include_regex or exclude_regex matches
func FilterRegex(capture_path_SubstringRegex string, include_regex []string, exclude_regex []string) Option {
	return filterRegex{CapturePathSubstringRegex: capture_path_SubstringRegex, IncludeRegex: include_regex, ExcludeRegex: exclude_regex}
}

func callIncludeFunc(value string, includeregexlist []string) (string, error) {

	for _, includeregex := range includeregexlist {
		if includeregex != "" {
			re := regexp.MustCompile(includeregex)
			is_match := re.MatchString(value)
			if is_match {
				return value, nil
			}
		}
	}

	return "", fmt.Errorf("'%s' does not match with whole include regex list hence not considering this namsepace\n", value)
}

func callExcludeFunc(value string, excluderegexlist []string) (string, error) {

	for _, excluderegex := range excluderegexlist {
		if excluderegex != "" {
			re := regexp.MustCompile(excluderegex)
			is_match := re.MatchString(value)
			if is_match {
				return "", fmt.Errorf("'%s' matched with exclude regex hence not considering this namsepace\n", value)
			}
		}
	}

	return value, nil
}
