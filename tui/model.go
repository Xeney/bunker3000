package tui

import (
	"bunker3000/game"
	"bunker3000/player"
	"bunker3000/save"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenMenu screen = iota
	screenDifficulty
	screenClassSelect
	screenGame
	screenGameOver
	screenAchievements
	screenGuide
)

type gamePhase int

const (
	phaseEvent gamePhase = iota
	phaseResult
)

type model struct {
	screen screen

	menuChoice    int
	diffChoice    int
	classChoice   int
	achieveChoice int

	gs        *game.GameState
	gamePhase gamePhase
	resultMsg string
	diff      player.Difficulty

	ready  bool
	width  int
	height int

	scrollOffset int
}

func NewModel() *model {
	return &model{
		screen:      screenMenu,
		menuChoice:  0,
		diffChoice:  1,
		classChoice: 0,
		ready:       false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "pgup":
			m.scrollOffset -= m.height / 3
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
			return m, nil
		case "pgdn":
			m.scrollOffset += m.height / 3
			return m, nil
		case "esc":
			if m.screen == screenGame && m.gamePhase == phaseResult {
				m.gamePhase = phaseEvent
				m.scrollOffset = 0
				return m, nil
			}
			if m.screen != screenMenu {
				m.screen = screenMenu
				m.menuChoice = 0
				m.scrollOffset = 0
				save.DeleteSave()
				return m, nil
			}
		}

		switch m.screen {
		case screenMenu:
			return m.updateMenu(msg)
		case screenDifficulty:
			return m.updateDifficulty(msg)
		case screenClassSelect:
			return m.updateClassSelect(msg)
		case screenGame:
			return m.updateGame(msg)
		case screenGameOver:
			return m.updateGameOver(msg)
		case screenAchievements:
			return m.updateAchievements(msg)
		case screenGuide:
			return m.updateGuide(msg)
		}
	}
	return m, nil
}

func (m model) View() string {
	var content string
	switch m.screen {
	case screenMenu:
		content = m.viewMenu()
	case screenDifficulty:
		content = m.viewDifficulty()
	case screenClassSelect:
		content = m.viewClassSelect()
	case screenGame:
		content = m.viewGame()
	case screenGameOver:
		content = m.viewGameOver()
	case screenAchievements:
		content = m.viewAchievements()
	case screenGuide:
		content = m.viewGuide()
	}

	if m.ready && m.height > 0 {
		lines := strings.Split(content, "\n")
		visible := m.height
		if visible < 1 {
			visible = 1
		}
		total := len(lines)

		if m.scrollOffset > total-visible {
			m.scrollOffset = max(0, total-visible)
		}

		if total > visible {
			start := m.scrollOffset
			end := start + visible
			if end > total {
				end = total
			}
			lines = lines[start:end]

			if m.scrollOffset > 0 || end < total {
				lines = append(lines, "", HelpStyle.Render("PgUp/PgDn — прокрутка"))
			}
		}

		content = strings.Join(lines, "\n")
	}

	return content
}
