package gui

import (
	"bunker3000/player"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowDifficulty(w fyne.Window) {
	bg := canvas.NewRectangle(color.RGBA{18, 18, 22, 255})

	title := canvas.NewText("ВЫБОР СЛОЖНОСТИ", color.RGBA{230, 170, 60, 255})
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	difficulties := []struct {
		diff    player.Difficulty
		name    string
		desc    string
		details string
		color   color.Color
	}{
		{player.DifficultyEasy, "Лёгкая", "10 дней, щадящий режим",
			"Ресурсов с запасом, урон от голода/жажды небольшой.\nИдеально для первого знакомства с игрой.",
			color.RGBA{80, 180, 80, 255}},
		{player.DifficultyNormal, "Обычная", "7 дней, стандартные условия",
			"Сбалансированный режим. Требует внимания к ресурсам\nи правильного выбора в событиях.",
			color.RGBA{210, 150, 40, 255}},
		{player.DifficultyHard, "Сложная", "5 дней, жёсткий вызов",
			"Ресурсов в обрез, каждый день на счету.\nДля опытных игроков, кто ищет испытаний.",
			color.RGBA{200, 60, 50, 255}},
		{player.DifficultyEndless, "Бесконечный", "∞ дней, растущая сложность",
			"Продержитесь максимально долго! Расход ресурсов растёт каждые 5 дней.\nПроверьте, насколько вас хватит.",
			color.RGBA{180, 80, 180, 255}},
	}

	selected := 1

	detailsLabel := widget.NewLabel(difficulties[1].details)
	detailsLabel.Wrapping = fyne.TextWrapWord

	cards := make([]*fyne.Container, 4)
	cardBorders := make([]*canvas.Rectangle, 4)

	for i, d := range difficulties {
		idx := i

		border := canvas.NewRectangle(color.RGBA{50, 45, 38, 255})
		border.CornerRadius = 10

		diffDot := canvas.NewCircle(d.color)
		nameLabel := widget.NewLabelWithStyle(d.name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		descLabel := widget.NewLabel(d.desc)

		inner := container.NewVBox(
			container.NewHBox(
				container.NewPadded(diffDot),
				nameLabel,
			),
			descLabel,
		)
		inner = container.NewPadded(inner)

		card := container.NewStack(border, inner)
		cardBorders[i] = border

		clickBox := widget.NewButton("", func() {
			selected = idx
			detailsLabel.SetText(difficulties[idx].details)
			for j := 0; j < 4; j++ {
				if j == idx {
					cardBorders[j].FillColor = difficulties[j].color
				} else {
					cardBorders[j].FillColor = color.RGBA{50, 45, 38, 255}
				}
				cardBorders[j].Refresh()
			}
		})
		clickBox.Importance = widget.LowImportance

		cards[i] = container.NewStack(card, clickBox)
	}

	cardBorders[1].FillColor = difficulties[1].color
	cardBorders[1].Refresh()

	cardRow := container.NewGridWithColumns(4,
		container.NewPadded(cards[0]),
		container.NewPadded(cards[1]),
		container.NewPadded(cards[2]),
		container.NewPadded(cards[3]),
	)

	backBtn := widget.NewButtonWithIcon("  Назад", theme.NavigateBackIcon(), func() {
		ShowMainMenu(w)
	})
	nextBtn := widget.NewButtonWithIcon("Далее  ", theme.NavigateNextIcon(), func() {
		ShowClassSelect(w, difficulties[selected].diff)
	})
	nextBtn.Importance = widget.HighImportance

	nav := container.NewBorder(nil, nil, backBtn, nextBtn, container.NewCenter(nextBtn))

	content := container.NewBorder(
		container.NewPadded(title),
		container.NewPadded(nav),
		nil, nil,
		container.NewVBox(
			cardRow,
			widget.NewSeparator(),
			container.NewPadded(detailsLabel),
		),
	)

	w.SetContent(container.NewStack(bg, container.NewPadded(content)))
}
