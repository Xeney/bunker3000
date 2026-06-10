package tui

import (
	"bunker3000/achievements"
	"bunker3000/save"
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateAchievements(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenMenu
		m.menuChoice = 2
	}
	return m, nil
}

func (m *model) viewAchievements() string {
	ach := achievements.NewTracker()

	saved, err := save.LoadAchievements()
	if err == nil {
		for _, sa := range saved {
			for _, a := range ach.Achievements {
				if a.ID == sa.ID && sa.Unlocked {
					a.Unlocked = true
				}
			}
		}
	}

	unlocked := 0
	for _, a := range ach.Achievements {
		if a.Unlocked {
			unlocked++
		}
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Получено: %d / %d", unlocked, len(ach.Achievements)))
	lines = append(lines, "")

	for _, a := range ach.Achievements {
		if a.Unlocked {
			lines = append(lines, LabelStyle.Render(" ★ ")+ValueStyle.Render(a.Name))
			lines = append(lines, "    "+DimStyle.Render(a.Desc))
		} else {
			lines = append(lines, DimStyle.Render(" ☆ "+a.Name))
			lines = append(lines, "    "+DimStyle.Render(a.Desc))
		}
		lines = append(lines, "")
	}

	lines = append(lines, HelpStyle.Render("Esc — назад"))

	body := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return RenderInBox("ДОСТИЖЕНИЯ", body)
}
