package evaluation

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Operator string

const (
	OpEquals      Operator = "equals"
	OpNotEquals   Operator = "not_equals"
	OpContains    Operator = "contains"
	OpNotContains Operator = "not_contains"
	OpStartsWith  Operator = "starts_with"
	OpEndsWith    Operator = "ends_with"
	OpIn          Operator = "in"
	OpNotIn       Operator = "not_in"
	OpGT          Operator = "gt"
	OpGTE         Operator = "gte"
	OpLT          Operator = "lt"
	OpLTE         Operator = "lte"
	OpExists      Operator = "exists"
	OpNotExists   Operator = "not_exists"
)

// EvalOperator evaluates whether an attribute value matches the condition value
// using the given operator. The attributeExists flag indicates whether the
// attribute was present in the context at all.
func EvalOperator(op Operator, attributeValue string, conditionValue string, attributeExists bool) bool {
	switch op {
	case OpExists:
		return attributeExists
	case OpNotExists:
		return !attributeExists
	case OpEquals:
		return attributeExists && attributeValue == conditionValue
	case OpNotEquals:
		return !attributeExists || attributeValue != conditionValue
	case OpContains:
		return attributeExists && strings.Contains(attributeValue, conditionValue)
	case OpNotContains:
		return !attributeExists || !strings.Contains(attributeValue, conditionValue)
	case OpStartsWith:
		return attributeExists && strings.HasPrefix(attributeValue, conditionValue)
	case OpEndsWith:
		return attributeExists && strings.HasSuffix(attributeValue, conditionValue)
	case OpIn:
		return attributeExists && inList(attributeValue, conditionValue)
	case OpNotIn:
		return !attributeExists || !inList(attributeValue, conditionValue)
	case OpGT:
		return attributeExists && compareNumeric(attributeValue, conditionValue) > 0
	case OpGTE:
		return attributeExists && compareNumeric(attributeValue, conditionValue) >= 0
	case OpLT:
		return attributeExists && compareNumeric(attributeValue, conditionValue) < 0
	case OpLTE:
		return attributeExists && compareNumeric(attributeValue, conditionValue) <= 0
	default:
		return false
	}
}

func inList(value, jsonArray string) bool {
	var list []string
	if err := json.Unmarshal([]byte(jsonArray), &list); err != nil {
		return false
	}
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func compareNumeric(a, b string) int {
	fa, errA := strconv.ParseFloat(a, 64)
	fb, errB := strconv.ParseFloat(b, 64)
	if errA != nil || errB != nil {
		return strings.Compare(a, b)
	}
	switch {
	case fa > fb:
		return 1
	case fa < fb:
		return -1
	default:
		return 0
	}
}
