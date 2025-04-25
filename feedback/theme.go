package feedback

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type emojiButton struct {
	fyne.Theme
}

func (b *emojiButton) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if n == theme.ColorNameButton {
		return color.Transparent
	}

	return b.Theme.Color(n, v)
}

func (b *emojiButton) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameText {
		return 32
	}

	return b.Theme.Size(n)
}
