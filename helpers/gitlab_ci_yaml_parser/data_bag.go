package gitlab_ci_yaml_parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/gitlab-org/gitlab-runner/helpers"
)

type DataBag map[string]interface{}

func (m *DataBag) Get(keys ...string) (interface{}, bool) {
	return helpers.GetMapKey(*m, keys...)
}

func (m *DataBag) GetSlice(keys ...string) ([]interface{}, bool) {
	slice, ok := helpers.GetMapKey(*m, keys...)
	if slice != nil {
		return slice.([]interface{}), ok
	}
	return nil, false
}

func (m *DataBag) GetStringSlice(keys ...string) (slice []string, ok bool) {
	rawSlice, ok := m.GetSlice(keys...)
	if !ok {
		return
	}

	for _, rawElement := range rawSlice {
		if element, ok := rawElement.(string); ok {
			slice = append(slice, element)
		}
	}
	return
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) && len(v) > 7 {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (m *DataBag) GetAllJobs() (result []string, ok bool) { // (result []string, ok bool) {
	var keys = []string{}
	for k := range *m {
		keys = append(keys, k)
	}
	for i := range keys {
		value, ok := helpers.GetMapKey(*m, keys[i])
		if ok {
			value, ok = value.(map[string]interface{})
			if ok {
				result = append(result, keys[i])
			}
		}
	}
	var out = []string{"variables", "workflow"}
	result = Filter(result, func(el string) bool { return !contains(out, el) })

	return
}

func (m *DataBag) GetSubOptions(keys ...string) (result DataBag, ok bool) {
	value, ok := helpers.GetMapKey(*m, keys...)
	if ok {
		result, ok = value.(map[string]interface{})
	}
	return
}

func (m *DataBag) GetString(keys ...string) (result string, ok bool) {
	value, ok := helpers.GetMapKey(*m, keys...)
	if ok {
		result, ok = value.(string)
	}
	return
}

func (m *DataBag) Decode(result interface{}, keys ...string) error {
	value, ok := m.Get(keys...)
	if !ok {
		return fmt.Errorf("key not found %v", strings.Join(keys, "."))
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

func convertMapToStringMap(in interface{}) (out interface{}, err error) {
	mapString := make(map[string]interface{})

	switch convMap := in.(type) {
	case map[string]interface{}:
		mapString = convMap
	case map[interface{}]interface{}:
		for k, v := range convMap {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("failed to convert %v to string", k)
			}
			mapString[key] = v
		}
	default:
		return in, nil
	}

	for k, v := range mapString {
		mapString[k], err = convertMapToStringMap(v)
		if err != nil {
			return
		}
	}
	return mapString, nil
}

func (m *DataBag) Sanitize() (err error) {
	n := make(DataBag)
	for k, v := range *m {
		n[k], err = convertMapToStringMap(v)
		if err != nil {
			return
		}
	}
	*m = n
	return
}

func getOptionsMap(optionKey string, primary, secondary DataBag) (value DataBag) {
	value, ok := primary.GetSubOptions(optionKey)
	if !ok {
		value, _ = secondary.GetSubOptions(optionKey)
	}

	return
}

func getOptions(optionKey string, primary, secondary DataBag) (value []interface{}, ok bool) {
	value, ok = primary.GetSlice(optionKey)
	if !ok {
		value, ok = secondary.GetSlice(optionKey)
	}

	return
}
