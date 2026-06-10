package tui

import (
	"bunker3000/game"
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateGameOver(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ", "esc":
		m.screen = screenMenu
		m.menuChoice = 0
	}
	return m, nil
}

func (m *model) viewGameOver() string {
	p := &m.gs.Player
	cfg := p.GetCfg()

	isVictory := p.Health > 0 && p.Lock && !m.gs.Endless

	var titleStr string
	if isVictory {
		titleStr = "ПОБЕДА!"
	} else {
		titleStr = "ПОРАЖЕНИЕ"
	}

	var statusLine string
	if isVictory {
		statusLine = fmt.Sprintf("Выжил %d из %d дней", p.ThisDay, cfg.MaxDays)
	} else if m.gs.Endless {
		statusLine = fmt.Sprintf("Продержался %d дней (∞)", p.ThisDay)
	} else if p.Health == 0 {
		statusLine = fmt.Sprintf("Погиб на %d дне", p.ThisDay)
	} else {
		statusLine = fmt.Sprintf("Игра завершена на %d дне", p.ThisDay)
	}

	stats := fmt.Sprintf(
		"Статус:    %s\nЗдоровье:   %d%%\nЕда:        %d / %d\nВода:       %d / %d\nСложность:  %s\nКласс:      %s",
		statusLine,
		p.Health,
		p.Eat, cfg.MaxResource,
		p.Water, cfg.MaxResource,
		p.Difficulty.String(),
		p.Class.String(),
	)

	epilogueText := game.GenerateEpilogue(m.gs)

	achText := ""
	unlocked := 0
	total := len(m.gs.Achievements.Achievements)
	for _, a := range m.gs.Achievements.Achievements {
		if a.Unlocked {
			unlocked++
		}
	}
	if unlocked == 0 {
		achText = "Достижений не получено"
	} else {
		achText = fmt.Sprintf("Достижения: %d / %d", unlocked, total)
		for _, a := range m.gs.Achievements.Achievements {
			if a.Unlocked {
				achText += "\n  ★ " + a.Name + " — " + a.Desc
			}
		}
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		LabelStyle.Render("Статистика:"),
		ValueStyle.Render(stats),
		"",
		LabelStyle.Render("Эпилог:"),
		DimStyle.Render(epilogueText),
		"",
		LabelStyle.Render("Достижения:"),
		ValueStyle.Render(achText),
		"",
		HelpStyle.Render("Enter — в главное меню"),
	)

	return RenderInBox(titleStr, body)
}
