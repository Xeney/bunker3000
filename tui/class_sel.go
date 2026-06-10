package tui

import (
	"bunker3000/game"
	"bunker3000/player"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type classEntry struct {
	class player.ClassType
	name  string
	bonus string
	desc  string
}

var classes = []classEntry{
	{player.ClassSurvivor, "Выживальщик", "Нет бонусов",
		"Сбалансированный персонаж. Без преимуществ и недостатков."},
	{player.ClassMedic, "Медик", "+20 HP, x2 лечение",
		"Максимальное здоровье: 120%. Лечение вдвое эффективнее."},
	{player.ClassScout, "Разведчик", "x2 ресурсы",
		"Двойная добыча еды и воды из событий."},
	{player.ClassBrawler, "Боец", "50% урона",
		"Получает вдвое меньше урона от всех источников."},
}

func (m *model) updateClassSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.classChoice > 0 {
			m.classChoice--
		}
	case "down", "j":
		if m.classChoice < len(classes)-1 {
			m.classChoice++
		}
	case "enter", " ":
		m.gs = game.NewGame(m.diff, classes[m.classChoice].class)
		m.gs.StartDay()
		m.screen = screenGame
		m.gamePhase = phaseEvent
		return m, nil
	case "esc":
		m.screen = screenDifficulty
	}
	return m, nil
}

func (m *model) viewClassSelect() string {
	var lines []string
	lines = append(lines, TitleStyle.Render("ВЫБОР КЛАССА"))
	lines = append(lines, "")

	for i, c := range classes {
		if i == m.classChoice {
			lines = append(lines, SelectedStyle.Render("> "+c.name+" <"))
			lines = append(lines, "  "+LabelStyle.Render("Бонус:")+" "+ValueStyle.Render(c.bonus))
			lines = append(lines, "  "+DimStyle.Render(c.desc))
		} else {
			lines = append(lines, UnselectedStyle.Render("  "+c.name))
			lines = append(lines, "   "+DimStyle.Render(c.bonus))
		}
		lines = append(lines, "")
	}

	lines = append(lines, HelpStyle.Render("↑/↓ выбор · Enter начать · Esc назад"))

	body := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return RenderInBox("ВЫБОР КЛАССА", body)
}
