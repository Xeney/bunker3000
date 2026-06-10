package tui

import (
	"bunker3000/game"
	"bunker3000/save"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuItem struct {
	label string
	action func(m *model) tea.Cmd
}

func (m *model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	items := menuItems(m)

	switch msg.String() {
	case "up", "k":
		if m.menuChoice > 0 {
			m.menuChoice--
		}
	case "down", "j":
		if m.menuChoice < len(items)-1 {
			m.menuChoice++
		}
	case "enter", " ":
		cmd := items[m.menuChoice].action(m)
		return m, cmd
	}
	return m, nil
}

func menuItems(m *model) []menuItem {
	items := []menuItem{
		{label: "Новая игра", action: func(m *model) tea.Cmd {
			m.screen = screenDifficulty
			m.diffChoice = 1
			m.scrollOffset = 0
			return nil
		}},
		{label: "Загрузить", action: func(m *model) tea.Cmd {
			if !save.HasSave() {
				return nil
			}
			sd, err := save.Load()
			if err != nil {
				return nil
			}
			if sd.Player.Lock {
				return nil
			}
			m.gs = game.LoadFromSave(sd)
			m.gs.StartDay()
			m.screen = screenGame
			m.gamePhase = phaseEvent
			m.scrollOffset = 0
			return nil
		}},
		{label: "Достижения", action: func(m *model) tea.Cmd {
			m.screen = screenAchievements
			m.achieveChoice = 0
			m.scrollOffset = 0
			return nil
		}},
		{label: "Руководство", action: func(m *model) tea.Cmd {
			m.screen = screenGuide
			m.achieveChoice = 0
			m.scrollOffset = 0
			return nil
		}},
		{label: "Выход", action: func(m *model) tea.Cmd {
			return tea.Quit
		}},
	}
	return items
}

func (m *model) viewMenu() string {
	title := " БУНКЕР-3000 "
	subtitle := "Симулятор выживания в пустоши"
	ver := "Версия 2.5.0"

	var lines []string
	items := menuItems(m)
	for i, item := range items {
		if i == m.menuChoice {
			lines = append(lines, SelectedStyle.Render("> "+item.label+" <"))
		} else {
			lines = append(lines, UnselectedStyle.Render("  "+item.label))
		}
	}

	body := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		subtitle,
		"",
		lipgloss.JoinVertical(lipgloss.Left, lines...),
		"",
		HelpStyle.Render("↑/↓ выбор · Enter подтвердить · q выход"),
		ver,
	)

	return RenderInBox("БУНКЕР-3000", body)
}
