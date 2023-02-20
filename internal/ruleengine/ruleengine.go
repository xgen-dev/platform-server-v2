package ruleengine

import "xgen.ai/utils"

type Condition struct {
	Operator   string       `json:"operator"`
	Value      interface{}  `json:"value"`
	Variable   string       `json:"variable"`
	Property   string       `json:"property"`
	Logic      string       `json:"logic"`
	Conditions []*Condition `json:"conditions"`
}

type Compare func(a, b interface{}) bool

func satisfiesCondition(c *Condition, obj interface{}, context map[string]interface{}) bool {
	o := operatorMap[c.Operator]
	value := c.Value
	if c.Variable != "" && context != nil {
		value = context[c.Variable]
	}
	property := utils.GetPathValue(obj, c.Property)
	return o.match(property, value)
}

func satisfiesRule(c *Condition, obj interface{}, context map[string]interface{}) bool {
	if c.Conditions == nil || c.Logic == "" {
		return false
	}

	var comp Compare
	switch c.Logic {
	case "and":
		comp = func(a, b interface{}) bool {
			return a.(bool) && b.(bool)
		}
	case "or":
		comp = func(a, b interface{}) bool {
			return a.(bool) || b.(bool)
		}
	default:
		return false
	}

	result := comp(true, true)
	for _, condition := range c.Conditions {
		if condition.Conditions != nil && condition.Logic != "" {
			subResult := satisfiesRule(condition, obj, context)
			result = comp(result, subResult)
		} else {
			subResult := satisfiesCondition(condition, obj, context)
			result = comp(result, subResult)
		}
	}

	return result
}

func Trigger(c *Condition, context map[string]interface{}) bool {
	return satisfiesRule(c, context, context)
}

func Filter(c *Condition, data []interface{}, context map[string]interface{}) []interface{} {
	var result []interface{}
	for _, item := range data {
		if satisfiesRule(c, item, context) {
			result = append(result, item)
		}
	}
	return result
}
