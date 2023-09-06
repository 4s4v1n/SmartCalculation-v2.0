package view

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
	"strconv"
	"strings"
	"telvina/APG2_SmartCalc/pkg/configurator"
	"telvina/APG2_SmartCalc/pkg/presenter"
)

type CalculatorView struct {
	application      fyne.App
	calculatorWindow fyne.Window
	helpWindow       fyne.Window
	plotWindow       fyne.Window
	creditWindow     fyne.Window

	helpWindowOpened   bool
	creditWindowOpened bool
	plotWindowOpened   bool

	presenter *presenter.Presenter
}

const (
	termTypeContainerPosition       = 2
	paymentTypeContainerPosition    = 4
	creditSumContainerPosition      = 8
	creditTermContainerPosition     = 9
	percentContainerPosition        = 12
	monthlyPaymentContainerPosition = 15
	overPaymentContainerPosition    = 16
	fullPaymentContainerPosition    = 17

	undoContainerPosition   = 20
	clearContainerPosition  = 27
	equalContainerPosition  = 33
	inputContainerPosition  = 34
	xInputContainerPosition = 36
	xBeginContainerPosition = 39
	xEndContainerPosition   = 40
)

const helpMsg = "\t\t\t\t\tSmartCalc v2 by Anton Savin\n" +
	"This project implements desktop calculator with some mathematics functions.\n\n" +
	"Usage:\n" +
	"\tYou can input expression in line edit from keyboard or use calculator keys.\n" +
	"\tIf you want to clear line edit, use (AC) button.\n" +
	"\tIf you want to delete last character, use (<-) button.\n" +
	"\tTo get expression value you should press (=) equal button.\n" +
	"\tIf expression has non valid symbols, you will get error message.\n" +
	"\tIf expression valid, you will get the answer.\n" +
	"\tIf expression contains x character, you can set x value, or plot graph from toolbar button.\n" +
	"\tAlso from toolbar you can open creditor calculator, get previous expression, clear history, open help.\n" +
	"\tExpressions keeper saves between program runs in file\n\n" +
	"Configure:\n" +
	"\tYou can configure visual part of calculator (labels, buttons, entries) colors\n" +
	"\tAlso there is ability to set history set file and logs location form config file\n\n" +
	"Technical part:\n" +
	"\tView part were write in GO using fyne. Core part were made in C++. Core uses as shared library."

func New(config configurator.Config) *CalculatorView {
	c := &CalculatorView{
		application: app.New(),
		presenter:   presenter.New(config),
	}

	c.calculatorWindow = c.application.NewWindow("SmartCalc v2")
	c.calculatorWindow.Resize(fyne.NewSize(700, 625))
	c.calculatorWindow.SetFixedSize(true)
	c.calculatorWindow.SetMaster()

	var itemsContainer []fyne.CanvasObject
	var x float32 = 0
	var y float32 = 200

	buttonsText := [34]string{
		"7", "8", "9", "+", "-", "*", "/",
		"4", "5", "6", "^", "%", "(", ")",
		"1", "2", "3", "sin", "cos", "tan", "<-",
		"0", ".", "x", "asin", "acos", "atan", "AC",
		"pi", "e", "sqrt", "ln", "lg", "="}

	creditLabelsText := [8]string{
		"Credit sum", "Credit term", "TermType", "Percent",
		"Monthly payment type", "Monthly payment", "Overpayment", "Full payment"}

	for i := 0; i < 5; i++ {
		for j := 0; j < 7; j++ {
			if i*7+j == 34 {
				break
			}

			buttonText := buttonsText[i*7+j]
			itemsContainer = append(itemsContainer, coloredButton(x, y, 100, 80,
				buttonText, config.ButtonsColor, func() {
					input := itemsContainer[inputContainerPosition].(*fyne.Container).Objects[0].(*widget.Entry)
					input.SetText(input.Text + buttonText)
				}))

			x += 100
		}

		x = 0
		y += 80
	}

	itemsContainer[equalContainerPosition].(*fyne.Container).Resize(fyne.Size{
		Width:  200,
		Height: 80,
	})

	itemsContainer = append(itemsContainer, coloredEntry(0, 0, 600, 100,
		"Input expression:", config.EntriesColor))
	inputExpression := itemsContainer[inputContainerPosition].(*fyne.Container).Objects[0].(*widget.Entry)

	itemsContainer = append(itemsContainer, coloredLabel(600, 0, 100, 100,
		"Input x value", config.LabelsColor))

	itemsContainer = append(itemsContainer, coloredEntry(600, 100, 100, 100,
		"", config.EntriesColor))
	xInput := itemsContainer[xInputContainerPosition].(*fyne.Container).Objects[0].(*widget.Entry)

	itemsContainer = append(itemsContainer, coloredLabel(0, 100, 150, 100,
		"Plot begin:", config.LabelsColor))
	itemsContainer = append(itemsContainer, coloredLabel(300, 100, 150, 100,
		"Plot end:", config.LabelsColor))

	itemsContainer = append(itemsContainer, coloredEntry(150, 100, 150, 100,
		"", config.EntriesColor))
	xBeginInput := itemsContainer[xBeginContainerPosition].(*fyne.Container).Objects[0].(*widget.Entry)

	itemsContainer = append(itemsContainer, coloredEntry(450, 100, 150, 100,
		"", config.EntriesColor))
	xEndInput := itemsContainer[xEndContainerPosition].(*fyne.Container).Objects[0].(*widget.Entry)

	itemsContainer[clearContainerPosition].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
		inputExpression.SetText("")
	}

	itemsContainer[undoContainerPosition].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
		if len(inputExpression.Text) > 0 {
			inputExpression.SetText(inputExpression.Text[:len(inputExpression.Text)-1])
		}
	}

	itemsContainer[equalContainerPosition].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
		if err := c.presenter.InvokeCalculatorValidate(inputExpression.Text); err != nil {
			inputExpression.SetText(err.Error())
			return
		}

		var xValue float64
		var err error

		if strings.Contains(inputExpression.Text, "x") {
			xValue, err = strconv.ParseFloat(xInput.Text, 64)
			if err != nil {
				inputExpression.SetText("invalid input")
				return
			}
		}

		c.presenter.StoreExpression(inputExpression.Text)
		logrus.Info(inputExpression.Text)
		inputExpression.SetText(fmt.Sprintf("%f", c.presenter.InvokeCalculator(inputExpression.Text, xValue)))
	}

	menu := fyne.NewMenu("Actions",
		fyne.NewMenuItem("Plot", func() {
			begin, err := strconv.ParseFloat(xBeginInput.Text, 64)
			if err != nil {
				xBeginInput.SetText("invalid input")
				return
			}

			end, err := strconv.ParseFloat(xEndInput.Text, 64)
			if err != nil {
				xEndInput.SetText("invalid input")
				return
			}

			if begin >= end {
				xBeginInput.SetText("invalid input")
				xEndInput.SetText("invalid input")
				return
			}

			if err = c.presenter.InvokeCalculatorValidate(inputExpression.Text); err != nil {
				inputExpression.SetText(err.Error())
				return
			}

			if !strings.Contains(inputExpression.Text, "x") {
				inputExpression.SetText("expression doesn't contain x")
				return
			}

			if !c.plotWindowOpened {
				c.plotWindow = c.application.NewWindow(inputExpression.Text)
				c.plotWindowOpened = true

				c.plotWindow.Resize(fyne.NewSize(600, 600))
				c.plotWindow.SetFixedSize(true)
				c.plotWindow.SetPadded(false)
				c.plotWindow.SetOnClosed(func() {
					c.plotWindowOpened = false

					enableAll(itemsContainer)
				})

				disableAll(itemsContainer)

				abscissa, ordinate := c.presenter.InvokeFixAbscissaOrdinate(c.presenter.InvokeAbscissa(begin, end),
					c.presenter.InvokeOrdinate(inputExpression.Text, begin, end))

				graph := chart.Chart{
					Width:  600,
					Height: 600,
					XAxis: chart.XAxis{
						Name: "x",
					},
					YAxis: chart.YAxis{
						Name: inputExpression.Text,
					},
					Series: []chart.Series{
						chart.ContinuousSeries{
							Name:    inputExpression.Text,
							XValues: abscissa,
							YValues: ordinate,
						},
					},
				}

				buffer := bytes.NewBuffer([]byte{})
				if err = graph.Render(chart.PNG, buffer); err != nil {
					logrus.Info(err)
					return
				}

				image := canvas.NewImageFromReader(buffer, "PLOT")
				image.Resize(fyne.NewSize(600, 600))
				image.Move(fyne.NewPos(0, 0))

				c.presenter.StoreExpression(inputExpression.Text)
				logrus.Info(inputExpression.Text)

				c.plotWindow.SetContent(container.NewWithoutLayout(image))
				c.plotWindow.Show()
			}
		}),
		fyne.NewMenuItem("Credit", func() {
			if !c.creditWindowOpened {
				c.creditWindow = c.application.NewWindow("SmartCalc v3 credit")
				c.creditWindowOpened = true

				c.creditWindow.Resize(fyne.NewSize(400, 450))
				c.creditWindow.SetFixedSize(true)
				c.creditWindow.SetPadded(false)
				c.creditWindow.SetOnClosed(func() {
					c.creditWindowOpened = false

					enableAll(itemsContainer)
				})

				disableAll(itemsContainer)

				var creditItemsContainer []fyne.CanvasObject
				for i, text := range creditLabelsText {
					creditItemsContainer = append(creditItemsContainer,
						coloredLabel(0, float32(i*50), 200, 50, text, config.LabelsColor))
				}

				creditItemsContainer = append(creditItemsContainer, coloredEntry(200, 0, 200, 50,
					"", config.EntriesColor))
				creditSumInput := creditItemsContainer[creditSumContainerPosition].(*fyne.Container).Objects[0].(*widget.Entry)

				creditItemsContainer = append(creditItemsContainer, coloredEntry(200, 50, 200, 50,
					"", config.EntriesColor))
				creditTermInput := creditItemsContainer[creditTermContainerPosition].(*fyne.Container).Objects[0].(*widget.Entry)

				termTypeLabel := creditItemsContainer[termTypeContainerPosition].(*fyne.Container).Objects[0].(*widget.Label)
				paymentTypeLabel := creditItemsContainer[paymentTypeContainerPosition].(*fyne.Container).Objects[0].(*widget.Label)

				creditItemsContainer = append(creditItemsContainer, coloredButton(200, 100, 100, 50,
					"Years", config.ButtonsColor, func() {
						termTypeLabel.SetText("Years")
					}))
				creditItemsContainer = append(creditItemsContainer, coloredButton(300, 100, 100, 50,
					"Months", config.ButtonsColor, func() {
						termTypeLabel.SetText("Months")
					}))

				creditItemsContainer = append(creditItemsContainer, coloredEntry(200, 150, 200, 50,
					"", config.EntriesColor))
				percentInput := creditItemsContainer[percentContainerPosition].(*fyne.Container).Objects[0].(*widget.Entry)

				creditItemsContainer = append(creditItemsContainer, coloredButton(200, 200, 100, 50,
					"Annuity", config.ButtonsColor, func() {
						paymentTypeLabel.SetText("Annuity")
					}))
				creditItemsContainer = append(creditItemsContainer, coloredButton(300, 200, 100, 50,
					"Different", config.ButtonsColor, func() {
						paymentTypeLabel.SetText("Different")
					}))

				for i := 0; i < 3; i++ {
					creditItemsContainer = append(creditItemsContainer, coloredLabel(200, 250+float32(i*50),
						200, 50, "None", config.LabelsColor))
				}

				monthlyPaymentLabel := creditItemsContainer[monthlyPaymentContainerPosition].(*fyne.Container).Objects[0].(*widget.Label)
				fullPaymentLabel := creditItemsContainer[fullPaymentContainerPosition].(*fyne.Container).Objects[0].(*widget.Label)
				overPaymentLabel := creditItemsContainer[overPaymentContainerPosition].(*fyne.Container).Objects[0].(*widget.Label)

				creditItemsContainer = append(creditItemsContainer, coloredButton(0, 400, 400, 50,
					"Calculate", config.ButtonsColor, func() {
						creditSum, err := strconv.ParseFloat(creditSumInput.Text, 64)
						if err != nil {
							return
						}

						creditTerm, err := strconv.ParseFloat(creditTermInput.Text, 64)
						if err != nil {
							return
						}

						percent, err := strconv.ParseFloat(percentInput.Text, 64)
						if err != nil {
							return
						}

						if !c.presenter.InvokeCreditorValidate(paymentTypeLabel.Text, termTypeLabel.Text,
							creditTerm, percent, creditSum) {
							monthlyPaymentLabel.SetText("invalid input")
							fullPaymentLabel.SetText("invalid input")
							overPaymentLabel.SetText("invalid input")
							return
						}

						result := c.presenter.InvokeCreditorCalculate(paymentTypeLabel.Text, termTypeLabel.Text,
							creditTerm, percent, creditSum)

						overPaymentLabel.SetText(strconv.FormatFloat(result.OverPay, 'f', 2, 64))
						fullPaymentLabel.SetText(strconv.FormatFloat(result.FullPay, 'f', 2, 64))

						if result.LastPay == 0 {
							monthlyPaymentLabel.SetText(strconv.FormatFloat(
								result.MonthPay, 'f', 2, 64))
						} else {
							monthlyPaymentLabel.SetText(strconv.FormatFloat(
								result.MonthPay, 'f', 2, 64) +
								"..." + strconv.FormatFloat(result.LastPay, 'f', 2, 64))
						}
					}))

				c.creditWindow.SetContent(container.NewWithoutLayout(creditItemsContainer...))
				c.creditWindow.Show()
			}
		}),
		fyne.NewMenuItem("Previous expression", func() {
			if expression := c.presenter.GetExpression(); expression != "" {
				inputExpression.SetText(expression)
			}
		}),
		fyne.NewMenuItem("Clear history", func() {
			c.presenter.ClearExpressions()
		}),
		fyne.NewMenuItem("Help", func() {
			if !c.helpWindowOpened {
				c.helpWindow = c.application.NewWindow("SmartCalc v3 help")
				c.helpWindowOpened = true

				c.helpWindow.Resize(fyne.NewSize(600, 400))
				c.helpWindow.SetFixedSize(true)
				c.helpWindow.SetContent(widget.NewLabel(helpMsg))
				c.helpWindow.SetOnClosed(func() {
					c.helpWindowOpened = false

					enableAll(itemsContainer)
				})

				disableAll(itemsContainer)

				c.helpWindow.Show()
			}
		}),
	)
	mainMenu := fyne.NewMainMenu(menu)

	c.calculatorWindow.SetMainMenu(mainMenu)
	c.calculatorWindow.SetPadded(false)
	c.calculatorWindow.SetContent(container.NewWithoutLayout(itemsContainer...))
	c.calculatorWindow.SetOnClosed(func() {
		c.presenter.SaveExpressions()
	})

	return c
}

func (c *CalculatorView) Run() {
	c.calculatorWindow.ShowAndRun()
	defer func() {
		c.presenter.ReleaseModel()
		c.presenter.ReleaseLogger()
	}()
}
