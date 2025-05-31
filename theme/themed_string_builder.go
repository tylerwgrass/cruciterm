package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

type ThemedStringBuilder struct {
	sb    strings.Builder
	theme lipgloss.Style
}

func NewThemedStringBuilder(theme lipgloss.Style) ThemedStringBuilder {
	var sb strings.Builder
	return ThemedStringBuilder{
		sb:    sb,
		theme: theme,
	}
}

func (tsb *ThemedStringBuilder) WriteString(s string) {
	tsb.sb.WriteString(tsb.theme.Render(s))
}

func (tsb *ThemedStringBuilder) WriteStyledString(s string, style lipgloss.Style) {
	tsb.sb.WriteString(style.Render(s))
}

func (tsb ThemedStringBuilder) String() string {
	return tsb.theme.Render(tsb.sb.String())
}
