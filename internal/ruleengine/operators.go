package ruleengine

import (
	"regexp"
	"strings"
)

type operator struct {
	slugs       []string
	description string
	match       func(interface{}, interface{}) bool
}

func listContains(list []string, item string) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}
	return false
}

func textIncludes(text, item string) bool {
	return strings.Contains(text, item)
}

var operators = []operator{
	{
		description: "Equals",
		slugs:       []string{"=", "===", "equal_to", "equals", "is", "is_equal_to"},
		match: func(property, value interface{}) bool {
			return property == value
		},
	},
	{
		slugs: []string{"!=", "!==", "doesnt_equal", "doesnt_equal_to", "is_not", "is_not_equal_to", "not", "not_equal", "not_equal_to", "!equals"},
		match: func(property, value interface{}) bool {
			return property != value
		},
	},
	{
		slugs: []string{"contains", "has", "includes"},
		match: func(property, value interface{}) bool {
			return textIncludes(property.(string), value.(string))
		},
	},
	{
		slugs: []string{"!contains", "!has", "!includes", "doesnt_contain", "doesnt_have", "not_contains"},
		match: func(property, value interface{}) bool {
			return !textIncludes(property.(string), value.(string))
		},
	},
	{
		slugs: []string{"greater_than", ">"},
		match: func(property, value interface{}) bool {
			switch property.(type) {
			case int:
				return property.(int) > value.(int)
			case float64:
				return property.(float64) > value.(float64)
			case string:
				return property.(string) > value.(string)
			default:
				return false
			}
		},
	},
	{
		slugs: []string{"greater_than_or_equal_to", ">=", "gte"},
		match: func(property, value interface{}) bool {
			switch property.(type) {
			case int:
				return property.(int) >= value.(int)
			case float64:
				return property.(float64) >= value.(float64)
			case string:
				return property.(string) >= value.(string)
			default:
				return false
			}
		},
	},
	{
		slugs: []string{"less_than", "<"},
		match: func(property, value interface{}) bool {
			switch property.(type) {
			case int:
				return property.(int) < value.(int)
			case float64:
				return property.(float64) < value.(float64)
			case string:
				return property.(string) < value.(string)
			default:
				return false
			}
		},
	},
	{
		slugs: []string{"less_than_or_equal_to", "<=", "lte"},
		match: func(property, value interface{}) bool {
			switch property.(type) {
			case int:
				return property.(int) <= value.(int)
			case float64:
				return property.(float64) <= value.(float64)
			case string:
				return property.(string) <= value.(string)
			default:
				return false
			}
		},
	},
	{
		slugs: []string{"matches_regex", "regex"},
		match: func(property, value interface{}) bool {
			regex, err := regexp.Compile(value.(string))
			if err != nil {
				return false
			}
			return regex.MatchString(property.(string))
		},
	},
	{
		slugs: []string{"any", "some", "in"},
		match: func(property, value interface{}) bool {

			if arr, ok := value.([]string); ok {
				return listContains(arr, property.(string))
			}

			if str, ok := value.(string); ok {
				return textIncludes(str, property.(string))
			}

			return false
		},
	},
	{
		slugs: []string{
			"in_fuzzy",
			"any_fuzzy",
			"any_fuzzy_match",
			"any_loose_match",
			"some_fuzzy",
			"some_loose",
			"some_fuzzy_match",
			"some_loose_match",
		},
		match: func(property, value interface{}) bool {
			if arr, ok := value.([]interface{}); ok {
				for _, item := range arr {
					if str, ok := item.(string); ok && textIncludes(property.(string), str) {
						return true
					}
				}
			}

			if str, ok := value.(string); ok {
				return textIncludes(property.(string), str)
			}
			return false
		},
	},
	{
		slugs: []string{"all", "every"},
		match: func(property, value interface{}) bool {
			if arr, ok := value.([]interface{}); ok {
				for _, item := range arr {
					if item != property {
						return false
					}
				}
				return true
			} else if str, ok := value.(string); ok {
				return str == property
			}
			return false
		},
	},
}

func expandOperators() [][]interface{} {
	var arr [][]interface{}
	for _, item := range operators {
		for _, slug := range item.slugs {
			arr = append(arr, []interface{}{slug, item})
		}
	}
	return arr
}

var operatorMap = func() map[string]operator {
	allOperators := expandOperators()
	m := make(map[string]operator)
	for _, item := range allOperators {
		m[item[0].(string)] = item[1].(operator)
	}
	return m
}()
