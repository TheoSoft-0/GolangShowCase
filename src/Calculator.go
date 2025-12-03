package main

// Modern, safer Fyne calculator with sanitized evaluation and correct button actions.
import (
	calculatorlib "01Hello/src/lib/calculatorLib"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app" // Added import for direct color changes if needed later
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

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Modern Fyne Calculator")

	myWindow.Resize(fyne.NewSize(175, 200))
	myWindow.SetFixedSize(true)
	resultLabel := widget.NewLabel("0")
	resultLabel.Alignment = fyne.TextAlignTrailing
	resultLabel.TextStyle.Bold = true

	displayContainer := container.NewPadded(resultLabel)

	evaluator := goval.NewEvaluator()

	buttonLabels := []string{
		"7", "8", "9", "/",
		"4", "5", "6", "*",
		"1", "2", "3", "-",
		"C", "0", "=", "+",
	}

	type actionFunc func()
	specialActions := map[string]actionFunc{
		"C": func() {
			resultLabel.SetText("")
		},
		"=": func() {
			expr := resultLabel.Text
			res, err := calculatorlib.SafeEvaluate(evaluator, expr)
			if err != nil {
				resultLabel.SetText("err")

			} else {
				resultLabel.SetText(res)
			}
		},
	}

	buttonGrid := container.New(layout.NewGridLayout(4))

	for _, lbl := range buttonLabels {
		label := lbl
		var onClick func()

		if act, ok := specialActions[label]; ok {
			onClick = act
		} else {
			onClick = func() { resultLabel.SetText(resultLabel.Text + label) }
		}

		btn := widget.NewButton(label, onClick)
		if label == "=" {
			btn.Importance = widget.HighImportance
		}
		if label == "C" {
			btn.Importance = widget.WarningImportance
		}
		buttonGrid.Add(btn)
	}

	calculatorUI := container.New(layout.NewVBoxLayout(),
		displayContainer, // Use the display container wrapping the label
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
