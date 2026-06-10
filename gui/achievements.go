package gui

import (
	"bunker3000/achievements"
	"bunker3000/save"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowAchievements(w fyne.Window) {
	title := widget.NewLabelWithStyle("ДОСТИЖЕНИЯ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	ach := achievements.NewTracker()

	sd, err := save.Load()
	if err == nil {
		for _, sa := range sd.Achievements {
			for _, a := range ach.Achievements {
				if a.ID == sa.ID && sa.Unlocked {
					a.Unlocked = true
				}
			}
		}
	}

	var contentText string
	unlocked := 0
	for _, a := range ach.Achievements {
		if a.Unlocked {
			unlocked++
		}
	}

	if unlocked == 0 {
		contentText = "Достижения не найдены. Сыграйте игру, чтобы получить их."
	} else {
		contentText = fmt.Sprintf("Получено: %d / %d\n", unlocked, len(ach.Achievements))
		for _, a := range ach.Achievements {
			if a.Unlocked {
				contentText += fmt.Sprintf("\n[X] %s\n    %s", a.Name, a.Desc)
			} else {
				contentText += fmt.Sprintf("\n[ ] %s\n    %s", a.Name, a.Desc)
			}
		}
	}

	achLabel := widget.NewLabel(contentText)
	achLabel.Wrapping = fyne.TextWrapWord

	backBtn := widget.NewButton("Назад", func() {
		ShowMainMenu(w)
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		achLabel,
		widget.NewSeparator(),
		backBtn,
	)

	scroll := container.NewVScroll(content)
	w.SetContent(scroll)
}
