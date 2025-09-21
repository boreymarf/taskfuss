package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/models"
)

func ValidateValueByDataType(value string, dataType string) error {
	switch dataType {
	case "int":
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("expected integer value, got %s", value)
		}
	case "float":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("expected float value, got %s", value)
		}
	case "bool":
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("expected boolean value (true/false), got %s", value)
		}
	case "duration":
		// Parse duration in format like "1h30m20s" or "1:30:20"
		if _, err := time.ParseDuration(normalizeDuration(value)); err != nil {
			return fmt.Errorf("expected duration value, got %s", value)
		}
	case "string":
		// All values are valid for string type
		return nil
	default:
		return fmt.Errorf("unknown data type: %s", dataType)
	}
	return nil
}

// Helper function to normalize duration formats
// If duration is in format "1:30:20", convert to "1h30m20s"
func normalizeDuration(durationStr string) string {
	if strings.Contains(durationStr, ":") {
		parts := strings.Split(durationStr, ":")
		if len(parts) == 3 {
			// HH:MM:SS format
			return parts[0] + "h" + parts[1] + "m" + parts[2] + "s"
		} else if len(parts) == 2 {
			// MM:SS format
			return parts[0] + "m" + parts[1] + "s"
		}
	}
	return durationStr
}

func EvaluateCondition(condition string, args []any) (bool, error) {
	switch strings.ToUpper(condition) {
	case "==":
		if len(args) < 2 {
			return false, fmt.Errorf("insufficient arguments for equality comparison, need at least 2, got %d", len(args))
		}

		// Check if all arguments are equal
		for i := 1; i < len(args); i++ {
			if args[i] != args[i-1] {
				return false, nil
			}
		}
		return true, nil

	case ">=", ">", "<", "<=":
		if len(args) != 2 {
			return false, fmt.Errorf("%s requires exactly 2 arguments, got %d", condition, len(args))
		}

		switch a := args[0].(type) {
		case int:
			b, ok := args[1].(int)
			if !ok {
				return false, fmt.Errorf("both arguments must be of the same type (int) for %s", condition)
			}
			switch condition {
			case ">=":
				return a >= b, nil
			case ">":
				return a > b, nil
			case "<":
				return a < b, nil
			case "<=":
				return a <= b, nil
			}

		case time.Duration:
			b, ok := args[1].(time.Duration)
			if !ok {
				return false, fmt.Errorf("both arguments must be of the same type (time.Duration) for %s", condition)
			}
			switch condition {
			case ">=":
				return a >= b, nil
			case ">":
				return a > b, nil
			case "<":
				return a < b, nil
			case "<=":
				return a <= b, nil
			}

		default:
			return false, fmt.Errorf("unsupported type for %s: %T, must be int or time.Duration", condition, a)
		}

	case "AND", "OR", "NAND", "NOR":
		if len(args) < 1 {
			return false, fmt.Errorf("%s requires at least 1 argument, got %d", condition, len(args))
		}

		// Check all arguments are bool
		for i, arg := range args {
			if _, ok := arg.(bool); !ok {
				return false, fmt.Errorf("all arguments must be bool for %s, got %T at position %d", condition, arg, i)
			}
		}

		switch strings.ToUpper(condition) {
		case "AND":

			for _, arg := range args {
				if !arg.(bool) {
					return false, nil
				}
			}
			return true, nil

		case "OR":
			for _, arg := range args {
				if arg.(bool) {
					return true, nil
				}
			}
			return false, nil

		case "NAND":
			for _, arg := range args {
				if !arg.(bool) {
					return true, nil
				}
			}
			return false, nil

		case "NOR":
			for _, arg := range args {
				if arg.(bool) {
					return false, nil
				}
			}
			return true, nil
		}

	default:
		return false, fmt.Errorf("unsupported condition: %s", condition)
	}

	return false, fmt.Errorf("unreachable code: %s", condition)
}

func EvaluateAtomRequirement(snapshot *models.RequirementSnapshot, entry *models.RequirementEntry) (bool, error) {
	if snapshot.Type != "atom" {
		return false, fmt.Errorf("expected atom, got %s", snapshot.Type)
	}

	actualValue, err := parseValue(entry.Value, snapshot.DataType.String)
	if err != nil {
		return false, err
	}

	targetValue, err := parseValue(snapshot.TargetValue, snapshot.DataType.String)
	if err != nil {
		return false, err
	}

	return EvaluateCondition(snapshot.Operator.String, []any{actualValue, targetValue})
}

// Helper function to parse target value based on data type
func parseValue(valueStr, dataType string) (any, error) {
	switch dataType {
	case "bool":
		return strconv.ParseBool(valueStr)
	case "int":
		return strconv.Atoi(valueStr)
	case "time":
		return time.ParseDuration(valueStr)
	case "string":
		return valueStr, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType)
	}
}
