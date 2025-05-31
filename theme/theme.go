package theme

import (
	"image/color"

	"github.com/charmbracelet/lipgloss/v2"
	tint "github.com/lrstanley/bubbletint"
)

func Default() string {
	return "rose_pine"
}

var theme lipgloss.Style

func Init() lipgloss.Style {
	tint.NewDefaultRegistry()
	tint.SetTintID(Default())
	theme = lipgloss.NewStyle().
		Foreground(tint.Fg()).
		BorderForeground(tint.Fg())
	return theme
}

func Get() lipgloss.Style {
	return theme
}

func SetWidth(w int) {
	theme = theme.Width(w)
}

func SetHeight(h int) {
	theme = theme.Height(h)
}

func Foreground() color.Color {
	return tint.Fg()
}

func Background() color.Color {
	return tint.Bg()
}

func Primary() color.Color {
	return tint.Yellow()
}

func Secondary() color.Color {
	return tint.Cyan()
}

func Red() color.Color {
	return tint.Red()
}

func Green() color.Color {
	return tint.Green()
}

func Apply(input string) string {
	return theme.Render(input)
}
