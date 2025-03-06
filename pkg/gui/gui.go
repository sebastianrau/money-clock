package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	TICK time.Duration = 30 * time.Millisecond

	PREF_MONEY_H          string  = "MONEY_PER_HOUR"
	PREF_MONEY_H_FALLBACK float64 = 80.0
)

type MainGui struct {
	start     bool
	startTime time.Time
	timeCount time.Duration

	pauseTime  time.Duration
	pauseBegin time.Time

	moneyPerHour float64
	moneyValue   float64

	lastSize fyne.Size

	timeText  *canvas.Text
	moneyText *canvas.Text

	moneySlider *widget.Slider
	sliderLabel *widget.Label

	btnStart *widget.Button
	btnPause *widget.Button
	btnStop  *widget.Button

	w fyne.Window
	a fyne.App
}

func NewMainGui(w fyne.Window, a fyne.App) *MainGui {

	mg := &MainGui{
		w: w,
		a: a,
	}

	mg.moneyPerHour = a.Preferences().FloatWithFallback(PREF_MONEY_H, PREF_MONEY_H_FALLBACK)

	mg.w.SetMaster()
	mg.w.SetTitle("Money Clock")

	mg.timeText = canvas.NewText("00:00:00.000", theme.Color(theme.ColorNamePressed))
	mg.timeText.Alignment = fyne.TextAlignCenter

	mg.moneyText = canvas.NewText("0,00 €", theme.Color(theme.ColorNameForeground))
	mg.moneyText.Alignment = fyne.TextAlignCenter

	mg.btnStart = widget.NewButton("Start", func() {
		mg.start = true
		mg.startTime = time.Now()
		mg.pauseTime = 0

		mg.btnStart.Disable()
		mg.btnPause.Enable()
		mg.btnStop.Enable()
		mg.moneySlider.Disable()
	})

	mg.btnPause = widget.NewButton("Pause", func() {
		mg.start = !mg.start

		if mg.start {
			mg.pauseTime += time.Since(mg.pauseBegin)
			mg.btnPause.Importance = widget.MediumImportance
		} else {
			mg.pauseBegin = time.Now()
			mg.btnPause.Importance = widget.HighImportance
		}
		mg.btnPause.Refresh()
	})
	mg.btnPause.Disable()

	mg.btnStop = widget.NewButton("Stop", func() {
		mg.start = false
		mg.UpdateGui()

		mg.btnStart.Enable()
		mg.btnStop.Disable()
		mg.btnPause.Importance = widget.MediumImportance
		mg.btnPause.Disable()
		mg.moneySlider.Enable()

	})
	mg.btnStop.Disable()

	mg.sliderLabel = widget.NewLabel(fmt.Sprintf("%.0f €/h", mg.moneyPerHour))
	mg.sliderLabel.Alignment = fyne.TextAlignTrailing

	mg.moneySlider = widget.NewSlider(10, 350)
	mg.moneySlider.Step = 5
	mg.moneySlider.SetValue(mg.moneyPerHour)
	mg.moneySlider.OnChanged = func(f float64) {
		mg.sliderLabel.SetText(fmt.Sprintf("%d €/h", int(f)))
	}

	mg.moneySlider.OnChangeEnded = func(f float64) {
		mg.moneyPerHour = f
		a.Preferences().SetFloat(PREF_MONEY_H, f)
	}

	containerBot := container.NewGridWithRows(2,
		container.NewGridWithColumns(3, mg.btnStart, mg.btnPause, mg.btnStop),
		container.New(layout.NewFormLayout(), mg.sliderLabel, mg.moneySlider),
	)

	content := container.NewBorder(mg.timeText, containerBot, nil, nil, mg.moneyText)

	go func() {
		t := time.NewTicker(TICK)
		for range t.C {
			if mg.start {
				mg.timeCount = (time.Since(mg.startTime) - mg.pauseTime)
				mg.moneyValue = float64(mg.timeCount/time.Second) * (mg.moneyPerHour / 3600)
				mg.UpdateGui()
			}
		}
	}()

	mg.w.SetContent(content)
	mg.UpdateGui()

	return mg

}

func (m *MainGui) UpdateGui() {
	if m.start {
		m.moneyText.Text = fmt.Sprintf("%.2f €", m.moneyValue)
		m.timeText.Text = fmtDuration(m.timeCount)
	}

	actualSize := m.w.Canvas().Size()
	if m.lastSize.Width != actualSize.Width {
		sizeM := calculateMaxFontSize("XXXXX.XX €", actualSize, m.moneyText.TextStyle)
		sizeT := sizeM * 0.75
		m.timeText.TextSize = sizeT
		m.moneyText.TextSize = sizeM

		m.lastSize = actualSize
	}

	m.timeText.Refresh()
	m.moneyText.Refresh()
}
