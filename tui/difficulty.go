package tui

import (
	"bunker3000/player"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type diffEntry struct {
	diff    player.Difficulty
	name    string
	desc    string
	details string
}

var difficulties = []diffEntry{
	{player.DifficultyEasy, "Лёгкая", "10 дней, щадящий режим",
		"Ресурсов с запасом, урон небольшой. Идеально для первого знакомства."},
	{player.DifficultyNormal, "Обычная", "7 дней, стандартные условия",
		"Сбалансированный режим. Требует внимания к ресурсам."},
	{player.DifficultyHard, "Сложная", "5 дней, жёсткий вызов",
		"Ресурсов в обрез, каждый день на счету. Для опытных игроков."},
	{player.DifficultyEndless, "Бесконечный", "∞ дней, растущая сложность",
		"Продержитесь максимально долго! Расход ресурсов растёт каждые 5 дней."},
}

func (m *model) updateDifficulty(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.diffChoice > 0 {
			m.diffChoice--
		}
	case "down", "j":
		if m.diffChoice < len(difficulties)-1 {
			m.diffChoice++
		}
	case "enter", " ":
		m.diff = difficulties[m.diffChoice].diff
		m.screen = screenClassSelect
		m.classChoice = 0
	case "esc":
		m.screen = screenMenu
		m.menuChoice = 0
	}
	return m, nil
}

func (m *model) viewDifficulty() string {
	var lines []string
	lines = append(lines, TitleStyle.Render("ВЫБОР СЛОЖНОСТИ"))
	lines = append(lines, "")

	for i, d := range difficulties {
		if i == m.diffChoice {
			lines = append(lines, SelectedStyle.Render("> "+d.name+" <")+"  "+
				ValueStyle.Render(d.desc))
		} else {
			lines = append(lines, UnselectedStyle.Render("  "+d.name)+"  "+DimStyle.Render(d.desc))
		}
	}

	lines = append(lines, "")
	lines = append(lines, DimStyle.Render(difficulties[m.diffChoice].details))
	lines = append(lines, "")
	lines = append(lines, HelpStyle.Render("↑/↓ выбор · Enter выбрать · Esc назад"))

	body := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return RenderInBox("ВЫБОР СЛОЖНОСТИ", body)
}
