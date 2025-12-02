package main

// Modern, safer Fyne calculator with sanitized evaluation and correct button actions.
import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/maja42/goval"
)

func makeUI() (*widget.Label, *widget.Entry) {
	out := widget.NewLabel("Hello world!")
	in := widget.NewEntry()

	in.OnChanged = func(content string) {
		out.SetText("Hello " + content + "!")
	}
	return out, in
}

// allowedExpr is a regexp that permits only digits, operators, parentheses, decimals and spaces.
// This prevents injection of names, function calls, or other unexpected tokens.
var allowedExpr = regexp.MustCompile(`^[0-9+\-*/().\s]+$`)
var intTok = regexp.MustCompile(`\b(\d+)\b`)

// safeEvaluate validates the expression then evaluates it using goval.
// It returns a user-friendly string or an error.
func safeEvaluate(ev *goval.Evaluator, expr string) (string, error) {
	trimmed := strings.TrimSpace(expr)
	if trimmed == "" {
		return "", fmt.Errorf("empty expression")
	}

	// Validate characters before any transformation.
	if !allowedExpr.MatchString(trimmed) {
		return "", fmt.Errorf("invalid characters")
	}

	toEval := trimmed
	// If the expression contains division, convert integer tokens to floats.
	// This ensures divisions produce floating-point results.
	if strings.Contains(toEval, "/") {
		toEval = intTok.ReplaceAllString(toEval, "$1.0")
	}

	// Evaluate.
	res, err := ev.Evaluate(toEval, nil, nil)
	if err != nil {
		return "", err
	}

	// Format numeric results. If original contained '/', prefer float formatting.
	if strings.Contains(trimmed, "/") {
		// Prefer float64 formatting path.
		switch v := res.(type) {
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64), nil
		case int:
			// Shouldn't normally happen after conversion, but handle anyway.
			return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
		case int64:
			return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
		default:
			return fmt.Sprintf("%v", res), nil
		}
	}

	// No division. Return integer-like results without unnecessary decimals.
	switch v := res.(type) {
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		// Remove trailing ".0" if the value is integral.
		s := strconv.FormatFloat(v, 'f', -1, 64)
		return s, nil
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", res), nil
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Modern Fyne Calculator")

	// Slightly larger fixed size for comfortable buttons.
	myWindow.Resize(fyne.NewSize(150, 200))
	myWindow.SetFixedSize(true)

	// Single entry for both input and output.
	resultEntry := widget.NewEntry()
	resultEntry.SetPlaceHolder("0")
	resultEntry.Disable()

	// Create a single evaluator and reuse it.
	evaluator := goval.NewEvaluator()

	// Buttons in reading order.
	buttonLabels := []string{
		"7", "8", "9", "/",
		"4", "5", "6", "*",
		"1", "2", "3", "-",
		"C", "0", "=", "+",
	}

	// Map special actions to named handlers.
	type actionFunc func()
	specialActions := map[string]actionFunc{
		"C": func() {
			resultEntry.SetText("")
		},
		"=": func() {
			expr := resultEntry.Text
			res, err := safeEvaluate(evaluator, expr)
			if err != nil {
				// Show concise error on invalid input.
				resultEntry.SetText("err")
				return
			}
			resultEntry.SetText(res)
		},
	}

	buttonGrid := container.New(layout.NewGridLayout(4))

	// Build buttons. Use local variable capture to avoid closure bugs.
	for _, lbl := range buttonLabels {
		label := lbl // local copy for closure safety
		var onClick func()

		if act, ok := specialActions[label]; ok {
			onClick = act
		} else {
			// Append text safely. Use closure capturing label variable.
			onClick = func() {
				// Keep input simple. No extra validation here.
				resultEntry.SetText(resultEntry.Text + label)
			}
		}

		btn := widget.NewButton(label, onClick)
		buttonGrid.Add(btn)
	}

	// Layout
	calculatorUI := container.New(layout.NewVBoxLayout(),
		resultEntry,
		widget.NewSeparator(),
		buttonGrid,
	)

	finalContent := container.New(layout.NewStackLayout(),
		container.New(layout.NewVBoxLayout(),
			layout.NewSpacer(),
			calculatorUI,
		),
	)

	myWindow.SetContent(finalContent)
	myWindow.ShowAndRun()
}
