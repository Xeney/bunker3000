package gui

import (
	"bunker3000/game"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowGameOverScreen(w fyne.Window, gs *game.GameState) {
	p := &gs.Player
	cfg := p.GetCfg()

	var titleText string
	if p.Health > 0 && p.Lock && !gs.Endless {
		titleText = "ПОБЕДА!"
	} else {
		titleText = "ПОРАЖЕНИЕ"
	}

	title := widget.NewLabelWithStyle(titleText, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	var statusLine string
	if p.Health > 0 && p.Lock && !gs.Endless {
		statusLine = fmt.Sprintf("ПОБЕДА! %d из %d дней", p.ThisDay, cfg.MaxDays)
	} else if gs.Endless {
		statusLine = fmt.Sprintf("ПОРАЖЕНИЕ на %d дне (∞)", p.ThisDay)
	} else if p.Health == 0 {
		statusLine = fmt.Sprintf("ПОРАЖЕНИЕ на %d дне", p.ThisDay)
	} else {
		statusLine = fmt.Sprintf("Игра завершена на %d дне", p.ThisDay)
	}

	stats := widget.NewLabel(fmt.Sprintf(
		"Статус:    %s\nЗдоровье:   %d%%\nЕда:        %d / %d\nВода:       %d / %d\nСложность:  %s\nКласс:      %s",
		statusLine,
		p.Health,
		p.Eat, cfg.MaxResource,
		p.Water, cfg.MaxResource,
		p.Difficulty.String(),
		p.Class.String(),
	))

	// Epilogue
	epilogueText := game.GenerateEpilogue(gs)
	epilogueLabel := widget.NewLabel(epilogueText)
	epilogueLabel.Wrapping = fyne.TextWrapWord

	achLabel := widget.NewLabel("")
	var achText string
	unlocked := 0
	total := len(gs.Achievements.Achievements)
	for _, a := range gs.Achievements.Achievements {
		if a.Unlocked {
			unlocked++
		}
	}
	achText = fmt.Sprintf("Достижения: %d / %d\n", unlocked, total)
	for _, a := range gs.Achievements.Achievements {
		if a.Unlocked {
			achText += fmt.Sprintf("[X] %s - %s\n", a.Name, a.Desc)
		}
	}
	achLabel.SetText(achText)

	menuBtn := widget.NewButton("В главное меню", func() {
		ShowMainMenu(w)
	})
	exitBtn := widget.NewButton("Выход", func() {
		w.Close()
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("ИТОГОВАЯ СТАТИСТИКА", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		stats,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("ЭПИЛОГ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		epilogueLabel,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("ДОСТИЖЕНИЯ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		achLabel,
		widget.NewSeparator(),
		container.NewHBox(menuBtn, exitBtn),
	)

	scroll := container.NewVScroll(content)
	w.SetContent(scroll)
}
