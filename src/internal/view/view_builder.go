package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"telvina/APG2_SmartCalc/pkg/configurator"
)

func coloredButton(x, y, width, height float32, text string, rgba configurator.RGBA, action func()) *fyne.Container {
	items := container.New(
		layout.NewMaxLayout(),
		widget.NewButton(text, action),
		canvas.NewRectangle(color.RGBA{
			R: rgba.R,
			G: rgba.G,
			B: rgba.B,
			A: rgba.A,
		}),
	)

	items.Move(fyne.NewPos(x, y))
	items.Resize(fyne.NewSize(width, height))

	return items
}

func coloredLabel(x, y, width, height float32, text string, rgba configurator.RGBA) *fyne.Container {
	items := container.New(
		layout.NewMaxLayout(),
		widget.NewLabel(text),
		canvas.NewRectangle(color.RGBA{
			R: rgba.R,
			G: rgba.G,
			B: rgba.B,
			A: rgba.A,
		}),
	)

	items.Move(fyne.NewPos(x, y))
	items.Resize(fyne.NewSize(width, height))

	return items
}

func coloredEntry(x, y, width, height float32, text string, rgba configurator.RGBA) *fyne.Container {
	items := container.New(
		layout.NewMaxLayout(),
		widget.NewEntry(),
		canvas.NewRectangle(color.RGBA{
			R: rgba.R,
			G: rgba.G,
			B: rgba.B,
			A: rgba.A,
		}),
	)

	items.Objects[0].(*widget.Entry).SetPlaceHolder(text)
	items.Objects[0].(*widget.Entry).TextStyle = fyne.TextStyle{
		Monospace: true,
	}

	items.Move(fyne.NewPos(x, y))
	items.Resize(fyne.NewSize(width, height))

	return items
}

func disableAll(objects []fyne.CanvasObject) {
	for _, item := range objects {
		item.Hide()
	}
}

func enableAll(objects []fyne.CanvasObject) {
	for _, item := range objects {
		item.Show()
	}
}
