package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
)

func fmtDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := d.Milliseconds() % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func calculateMaxFontSize(text string, maxSize fyne.Size, textStyle fyne.TextStyle) float32 {
	fontSize := float32(12) // Startgröße
	for {
		size := fyne.MeasureText(text, fontSize, textStyle) // Größe des Textes messen
		if size.Width > maxSize.Width || size.Height > (maxSize.Height/2) {
			break // Stoppen, wenn die maximale Breite überschritten wird
		}
		fontSize++
	}
	return fontSize * 0.8 // Letzte passende Größe zurückgeben
}
