package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	TICK time.Duration = 30 * time.Millisecond

	PREF_MONEY_H string = "MONEY_PER_HOUR"
)

var (
	start          bool
	startTime      time.Time
	moneyPerSecond float64

	timeText    *canvas.Text
	moneyText   *canvas.Text
	moneySlider *widget.Slider
	btnStart    *widget.Button
	btnStop     *widget.Button

	lastSize fyne.Size

	a fyne.App
	w fyne.Window
)

func main() {

	a = app.NewWithID("com.github.sebastianrau.money-clock")

	w = a.NewWindow("Hello World")
	w.Resize(fyne.NewSize(300, 300))
	w.SetMaster()
	w.SetTitle("Money Clock")

	timeText = canvas.NewText("00:00:00.000", theme.Color(theme.ColorNamePressed))
	timeText.Alignment = fyne.TextAlignCenter
	timeText.TextSize = 32

	moneyText = canvas.NewText("0,00 €", theme.Color(theme.ColorNameForeground))
	moneyText.Alignment = fyne.TextAlignCenter
	moneyText.TextSize = 128

	btnStart = widget.NewButton("Start", func() {
		start = true
		startTime = time.Now()

		btnStart.Disable()
		btnStop.Enable()
		moneySlider.Disable()
	})

	btnStop = widget.NewButton("Stop", func() {
		start = false
		UpdateGui()

		btnStart.Enable()
		btnStop.Disable()
		moneySlider.Enable()
	})
	btnStop.Disable()

	sliderValue := a.Preferences().FloatWithFallback(PREF_MONEY_H, 80.0)
	moneyPerSecond = sliderValue / 3600.0

	sliderLabel := widget.NewLabel(fmt.Sprintf("%.0f €/h", sliderValue))
	sliderLabel.Alignment = fyne.TextAlignTrailing

	moneySlider = widget.NewSlider(10, 200)
	moneySlider.Step = 5
	moneySlider.SetValue(sliderValue)
	moneySlider.OnChanged = func(f float64) {
		sliderLabel.SetText(fmt.Sprintf("%d €/h", int(f)))
	}

	moneySlider.OnChangeEnded = func(f float64) {
		moneyPerSecond = f / 3600.0
		a.Preferences().SetFloat(PREF_MONEY_H, f)
	}

	containerBot := container.NewGridWithColumns(2,
		container.NewGridWithColumns(2, btnStart, btnStop),
		container.New(layout.NewFormLayout(), sliderLabel, moneySlider),
	)

	content := container.NewBorder(timeText, containerBot, nil, nil, moneyText)

	w.SetContent(content)
	UpdateGui()

	go func() {
		t := time.NewTicker(TICK)
		for range t.C {
			if start {
				UpdateGui()
			}
		}
	}()

	w.ShowAndRun()
}

func UpdateGui() {
	if start {
		since := time.Since(startTime)
		m := fmt.Sprintf("%.2f €", float64(since/time.Second)*moneyPerSecond)
		t := fmtDuration(since)
		timeText.Text = t
		moneyText.Text = m
	}

	actualSize := w.Canvas().Size()
	if lastSize.Width != actualSize.Width {
		sizeM := calculateMaxFontSize("XXXXX.XX €", actualSize.Width, moneyText.TextStyle)
		sizeT := sizeM * 0.75
		timeText.TextSize = sizeT
		moneyText.TextSize = sizeM

		lastSize = actualSize
	}

	timeText.Refresh()
	moneyText.Refresh()

}

func fmtDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := d.Milliseconds() % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func calculateMaxFontSize(text string, maxWidth float32, textStyle fyne.TextStyle) float32 {
	fontSize := float32(12) // Startgröße
	for {
		size := fyne.MeasureText(text, fontSize, textStyle) // Größe des Textes messen
		if size.Width > maxWidth {
			break // Stoppen, wenn die maximale Breite überschritten wird
		}
		fontSize++
	}
	return fontSize - 5 // Letzte passende Größe zurückgeben
}
