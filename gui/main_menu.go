package gui

import (
	"bunker3000/game"
	"bunker3000/save"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowMainMenu(w fyne.Window) {
	bg := canvas.NewRectangle(color.RGBA{18, 18, 22, 255})

	title := canvas.NewText("БУНКЕР-3000", color.RGBA{230, 170, 60, 255})
	title.TextSize = 36
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Симулятор выживания в пустоши", color.RGBA{140, 135, 120, 255})
	subtitle.TextSize = 16

	ver := canvas.NewText("Версия 2.0.0", color.RGBA{100, 95, 85, 255})
	ver.TextSize = 12

	newBtn := widget.NewButton("  НОВАЯ ИГРА", func() {
		ShowDifficulty(w)
	})
	newBtn.Importance = widget.HighImportance

	loadBtn := widget.NewButton("  ЗАГРУЗИТЬ", func() {
		if !save.HasSave() {
			showMsgPopup(w, "Сохранение не найдено.")
			return
		}
		sd, err := save.Load()
		if err != nil {
			showMsgPopup(w, "Ошибка загрузки сохранения.")
			return
		}
		if sd.Player.Lock {
			showMsgPopup(w, "Эта игра уже завершена.")
			return
		}
		gs := game.LoadFromSave(sd)
		ShowGameScreen(w, gs)
	})
	if !save.HasSave() {
		loadBtn.Disable()
	}

	achBtn := widget.NewButton("  ДОСТИЖЕНИЯ", func() {
		ShowAchievements(w)
	})

	exitBtn := widget.NewButton("  ВЫХОД", func() {
		w.Close()
	})

	btnBox := container.NewVBox(
		container.NewPadded(newBtn),
		container.NewPadded(loadBtn),
		container.NewPadded(achBtn),
		container.NewPadded(exitBtn),
	)

	titleBox := container.NewVBox(
		container.NewPadded(title),
		subtitle,
	)

	center := container.NewCenter(
		container.NewBorder(
			titleBox, ver, nil, nil,
			btnBox,
		),
	)

	content := container.NewStack(bg, center)
	w.SetContent(content)
}

func showMsgPopup(w fyne.Window, msg string) {
	var pop *widget.PopUp
	pop = widget.NewModalPopUp(
		container.NewVBox(
			container.NewPadded(widget.NewLabel(msg)),
			widget.NewButton("OK", func() { pop.Hide() }),
		),
		w.Canvas(),
	)
	pop.Show()
}
