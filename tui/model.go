package tui

import (
	"bunker3000/game"
	"bunker3000/player"
	"bunker3000/save"

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

	ready bool
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.screen == screenGame && m.gamePhase == phaseResult {
				m.gamePhase = phaseEvent
				return m, nil
			}
			if m.screen != screenMenu {
				m.screen = screenMenu
				m.menuChoice = 0
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
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case screenMenu:
		return m.viewMenu()
	case screenDifficulty:
		return m.viewDifficulty()
	case screenClassSelect:
		return m.viewClassSelect()
	case screenGame:
		return m.viewGame()
	case screenGameOver:
		return m.viewGameOver()
	case screenAchievements:
		return m.viewAchievements()
	}
	return ""
}
