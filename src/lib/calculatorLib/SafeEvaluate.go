package calculatorlib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/maja42/goval"
)

var allowedExpr = regexp.MustCompile(`^[0-9+\-*/().\s]+$`)
var intTok = regexp.MustCompile(`\b(\d+)\b`)

func SafeEvaluate(ev *goval.Evaluator, expr string) (string, error) {
	trimmed := strings.TrimSpace(expr)
	if trimmed == "" {
		return "", fmt.Errorf("empty expression")
	}
	if !allowedExpr.MatchString(trimmed) {
		return "", fmt.Errorf("invalid characters")
	}

	toEval := trimmed

	if strings.Contains(toEval, "/") {
		toEval = intTok.ReplaceAllString(toEval, "$1.0")
	}

	// Evaluate.
	res, err := ev.Evaluate(toEval, nil, nil)
	if err != nil {
		return "", err
	}

	if strings.Contains(trimmed, "/") {
		// Prefer float64 formatting path.
		switch v := res.(type) {
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64), nil
		case int:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
		case int64:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
		default:
			return fmt.Sprintf("%v", res), nil
		}
	}

	switch v := res.(type) {
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		return s, nil
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", res), nil
	}
}
