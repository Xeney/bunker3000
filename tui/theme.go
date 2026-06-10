package tui

import "github.com/charmbracelet/lipgloss"

var (
	Green   = lipgloss.Color("10")
	DarkGray  = lipgloss.Color("8")
	LightGray = lipgloss.Color("7")
	White     = lipgloss.Color("15")
	Black     = lipgloss.Color("0")

	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(Green).
		Foreground(Green).
		Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true).
			Align(lipgloss.Center).
			Width(50)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(Black).
			Background(Green).
			Padding(0, 1)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(Green).
			Padding(0, 1)

	DimStyle = lipgloss.NewStyle().
		Foreground(DarkGray).
		Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
		Foreground(DarkGray).
		Align(lipgloss.Center).
		Width(50)

	LabelStyle = lipgloss.NewStyle().
		Foreground(LightGray).
		Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(Green)

	ResultGoodStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	ResultBadStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	ProgressBg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("236"))

	ProgressFg = lipgloss.NewStyle().
			Foreground(Green)
)

func ProgressBar(pct float64, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	filled := int(float64(width) * pct)
	if filled > width {
		filled = width
	}
	empty := width - filled
	bar := "[" + ProgressFg.Render(renderBar(filled, '█')) +
		ProgressBg.Render(renderBar(empty, '░')) + "]"
	return bar
}

func renderBar(n int, ch rune) string {
	if n <= 0 {
		return ""
	}
	s := make([]rune, n)
	for i := range s {
		s[i] = ch
	}
	return string(s)
}

func RenderInBox(title string, body string) string {
	return BoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			TitleStyle.Render(title),
			"",
			body,
		),
	)
}
