package tui

import (
	"bunker3000/player"
	"bunker3000/save"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateGame(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.gamePhase {
	case phaseEvent:
		switch msg.String() {
		case "1":
			r, _ := m.gs.ExecuteChoice(1)
			m.resultMsg = r
			if m.gs.IsGameOver() {
				m.gs.Save()
				save.DeleteSave()
				m.screen = screenGameOver
				return m, nil
			}
			m.gamePhase = phaseResult
		case "2":
			r, _ := m.gs.ExecuteChoice(2)
			m.resultMsg = r
			if m.gs.IsGameOver() {
				m.gs.Save()
				save.DeleteSave()
				m.screen = screenGameOver
				return m, nil
			}
			m.gamePhase = phaseResult
		}

	case phaseResult:
		switch msg.String() {
		case "enter", " ":
			m.gs.NextDay()
			m.gs.StartDay()
			if m.gs.IsGameOver() || m.gs.DayError != nil {
				m.gs.Save()
				save.DeleteSave()
				m.screen = screenGameOver
				return m, nil
			}
			m.gamePhase = phaseEvent
		}
	}

	return m, nil
}

func (m *model) viewGame() string {
	p := &m.gs.Player
	cfg := p.GetCfg()
	maxHP := player.MaxHealth(p.Class)

	dayStr := fmt.Sprintf("День %d", p.ThisDay)
	if m.gs.Endless {
		dayStr += " (∞)"
	} else {
		dayStr += fmt.Sprintf(" из %d", cfg.MaxDays)
	}

	healthPct := float64(p.Health) / float64(maxHP)
	eatPct := float64(p.Eat) / float64(cfg.MaxResource)
	waterPct := float64(p.Water) / float64(cfg.MaxResource)

	status := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center,
			LabelStyle.Render(" День: "),
			ValueStyle.Render(dayStr),
		),
		LabelStyle.Render(" Здоровье: ")+ProgressBar(healthPct, 20)+fmt.Sprintf(" %d%%", p.Health),
		LabelStyle.Render(" Еда:      ")+ProgressBar(eatPct, 20)+fmt.Sprintf(" %d/%d", p.Eat, cfg.MaxResource),
		LabelStyle.Render(" Вода:     ")+ProgressBar(waterPct, 20)+fmt.Sprintf(" %d/%d", p.Water, cfg.MaxResource),
		DimStyle.Render(fmt.Sprintf(" %s · %s", p.Difficulty.String(), p.Class.String())),
	)

	switch m.gamePhase {
	case phaseEvent:
		ev := m.gs.CurrentEvent
		if ev == nil {
			return RenderInBox("БУНКЕР-3000",
				lipgloss.JoinVertical(lipgloss.Left,
					status,
					"",
					ResultBadStyle.Render(" Ошибка: событие не получено"),
					"",
					HelpStyle.Render("Esc — в главное меню"),
				),
			)
		}

		catLabel := eventCategoryLabel(ev.Category)

		eventBlock := lipgloss.JoinVertical(lipgloss.Left,
			"",
			LabelStyle.Render(" ["+catLabel+"]"),
			"",
			ValueStyle.Render(" "+ev.Message),
			"",
			ValueStyle.Render(" 1. "+ev.Variants[0]),
			ValueStyle.Render(" 2. "+ev.Variants[1]),
			"",
			HelpStyle.Render("1/2 — сделать выбор · Esc — выйти"),
		)

		return RenderInBox("БУНКЕР-3000",
			lipgloss.JoinVertical(lipgloss.Left,
				status,
				eventBlock,
			),
		)

	case phaseResult:
		isGood := isGoodResultTUI(m.resultMsg)
		resultStyle := ResultGoodStyle
		if !isGood {
			resultStyle = ResultBadStyle
		}

		resultBlock := lipgloss.JoinVertical(lipgloss.Left,
			"",
			resultStyle.Render(m.resultMsg),
			"",
			HelpStyle.Render("Enter — продолжить · Esc — выйти"),
		)

		return RenderInBox("БУНКЕР-3000",
			lipgloss.JoinVertical(lipgloss.Left,
				status,
				resultBlock,
			),
		)
	}

	return ""
}

func eventCategoryLabel(cat string) string {
	switch cat {
	case "resource":
		return "РЕСУРСЫ"
	case "combat":
		return "БИТВА"
	case "helper":
		return "ПОМОЩЬ"
	case "weather":
		return "ПОГОДА"
	case "hazard":
		return "ОПАСНОСТЬ"
	case "find":
		return "НАХОДКА"
	default:
		return "СОБЫТИЕ"
	}
}

func isGoodResultTUI(msg string) bool {
	return !strings.Contains(msg, "[-]") &&
		!strings.Contains(msg, "СМЕРТЬ") &&
		!strings.Contains(msg, "ТРАВМА") &&
		!strings.Contains(msg, "ОПАСНОСТЬ") &&
		!strings.Contains(msg, "ОТРАВЛЕНИЕ") &&
		!strings.Contains(msg, "ПАДЕНИЕ") &&
		!strings.Contains(msg, "БЕГСТВО") &&
		!strings.Contains(msg, "НЕУДАЧА")
}
