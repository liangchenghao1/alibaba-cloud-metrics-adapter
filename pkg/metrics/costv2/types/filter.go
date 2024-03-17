package costv2

import (
	"fmt"
	"strings"
)

type Filter struct {
	Namespace      []string
	ControllerName []string
	ControllerKind []string
	Pod            []string
	Label          map[string][]string
}

func (f *Filter) GetLabelSelectorStr() string {
	if f.GetKubePodLabelStr() == "" {
		return fmt.Sprintf(`namespace=%s,created_by_kind=%s,created_by_name=%s,pod=%s`,
			f.getNamespaceStr(), f.getControllerKindStr(), f.getControllerNameStr(), f.getPodStr())
	}
	return fmt.Sprintf(`namespace=%s,created_by_kind=%s,created_by_name=%s,pod=%s,%s`,
		f.getNamespaceStr(), f.getControllerKindStr(), f.getControllerNameStr(), f.getPodStr(), f.GetKubePodLabelStr())
}

func (f *Filter) GetKubePodInfoStr() string {
	return fmt.Sprintf(`namespace=~"%s",created_by_kind=~"%s",created_by_name=~"%s",pod=~"%s"`,
		f.getNamespaceStr(), f.getControllerKindStr(), f.getControllerNameStr(), f.getPodStr())
}

func (f *Filter) GetKubePodLabelStr() string {
	// only support single label currently
	// todo: check promql special symbol conversion, eg. "label_a/b" -> "label_a_b"
	res := ""
	if f.Label != nil {
		for key, value := range f.Label {
			//res += fmt.Sprintf(`label_%s="%s"`, key, value[0])
			res += fmt.Sprintf(`label_%s=%s`, key, value[0])
		}
	}
	return res
}

func (f *Filter) getNamespaceStr() string {
	if f.Namespace != nil {
		return strings.Join(f.Namespace, "|")
	}
	return ".*"
}

func (f *Filter) getControllerNameStr() string {
	if f.ControllerName != nil {
		return strings.Join(f.ControllerName, "|")
	}
	return ".*"
}

func (f *Filter) getControllerKindStr() string {
	if f.ControllerKind != nil {
		return strings.Join(f.ControllerKind, "|")
	}
	return ".*"
}

func (f *Filter) getPodStr() string {
	if f.Pod != nil {
		return strings.Join(f.Pod, "|")
	}
	return ".*"
}

// parseFilterParts Split the filter string
func parseFilterParts(filterStr string) []string {
	prefixes := []string{"namespace:", "controllerName:", "controllerKind:", "pod:", "label["}

	// for each prefix, the prefix in the filter string is replaced with a special symbol("\x1f") plus prefix for segmentation.
	for _, prefix := range prefixes {
		filterStr = strings.Replace(filterStr, "+"+prefix, "\x1f"+prefix, -1)
	}

	return strings.Split(filterStr, "\x1f")
}

// ParseFilter Parses the given string to *Filter
func ParseFilter(filterStr string) (*Filter, error) {
	filter := &Filter{}

	filterParts := parseFilterParts(filterStr)

	for _, part := range filterParts {
		if part == "" {
			continue
		}

		// handles the contents of fields inside ""
		kv := strings.SplitN(part, `:"`, 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid filter format: %s", part)
		}

		key := strings.Trim(kv[0], `"`)
		value := strings.Trim(kv[1], `"`)

		switch {
		case strings.HasPrefix(key, "namespace"):
			filter.Namespace = strings.Split(value, ",")
		case strings.HasPrefix(key, "controllerName"):
			filter.ControllerName = strings.Split(value, ",")
		case strings.HasPrefix(key, "controllerKind"):
			filter.ControllerKind = strings.Split(value, ",")
		case strings.HasPrefix(key, "pod"):
			filter.Pod = strings.Split(value, ",")
		case strings.HasPrefix(key, "label["):
			filter.Label = make(map[string][]string)
			labelKey := strings.TrimPrefix(key, "label[")
			labelKey = strings.TrimSuffix(labelKey, "]")
			filter.Label[labelKey] = []string{value}
		default:
			return nil, fmt.Errorf("unsupported filter key: %s", key)
		}
	}

	return filter, nil
}
