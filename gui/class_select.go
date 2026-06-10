package gui

import (
	"bunker3000/game"
	"bunker3000/player"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowClassSelect(w fyne.Window, diff player.Difficulty) {
	bg := canvas.NewRectangle(color.RGBA{18, 18, 22, 255})

	title := canvas.NewText("ВЫБОР КЛАССА", color.RGBA{230, 170, 60, 255})
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	classes := []struct {
		class player.ClassType
		name  string
		bonus string
		desc  string
		icon  string
	}{
		{
			player.ClassSurvivor, "Выживальщик", "Нет бонусов",
			"Сбалансированный персонаж. Без особых преимуществ,\nно и без недостатков. Классический выбор.",
			"●",
		},
		{
			player.ClassMedic, "Медик", "+20 HP, x2 лечение",
			"Максимальное здоровье: 120%.\nЛечение восстанавливает вдвое больше.\nОтличный выбор для тех, кто часто рискует.",
			"✚",
		},
		{
			player.ClassScout, "Разведчик", "x2 ресурсы",
			"Двойная добыча еды и воды из событий.\nПозволяет накапливать ресурсы быстрее.\nЛучший выбор для исследования пустоши.",
			"◎",
		},
		{
			player.ClassBrawler, "Боец", "50% урона",
			"Получает вдвое меньше урона от всех источников.\nНезаменим в бою, но не помогает с ресурсами.\nИдеален для агрессивного стиля игры.",
			"⚔",
		},
	}

	selected := 0
	descLabel := widget.NewLabel(classes[0].desc)
	descLabel.Wrapping = fyne.TextWrapWord

	cards := make([]*fyne.Container, 4)
	cardBorders := make([]*canvas.Rectangle, 4)

	classColors := []color.Color{
		color.RGBA{180, 180, 170, 255},
		color.RGBA{80, 200, 80, 255},
		color.RGBA{80, 180, 220, 255},
		color.RGBA{200, 60, 50, 255},
	}

	for i, c := range classes {
		idx := i

		border := canvas.NewRectangle(color.RGBA{50, 45, 38, 255})
		border.CornerRadius = 10

		icon := canvas.NewText(c.icon, classColors[i])
		icon.TextSize = 28

		nameLabel := widget.NewLabelWithStyle(c.name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		bonusLabel := widget.NewLabel(c.bonus)

		inner := container.NewVBox(
			container.NewCenter(icon),
			container.NewCenter(nameLabel),
			container.NewCenter(bonusLabel),
		)
		inner = container.NewPadded(inner)

		card := container.NewStack(border, inner)
		cardBorders[i] = border

		clickBtn := widget.NewButton("", func() {
			selected = idx
			descLabel.SetText(classes[idx].desc)
			for j := 0; j < 4; j++ {
				if j == idx {
					cardBorders[j].FillColor = classColors[j]
				} else {
					cardBorders[j].FillColor = color.RGBA{50, 45, 38, 255}
				}
				cardBorders[j].Refresh()
			}
		})
		clickBtn.Importance = widget.LowImportance

		cards[i] = container.NewStack(card, clickBtn)
	}

	cardBorders[0].FillColor = classColors[0]
	cardBorders[0].Refresh()

	cardRow := container.NewGridWithColumns(4,
		container.NewPadded(cards[0]),
		container.NewPadded(cards[1]),
		container.NewPadded(cards[2]),
		container.NewPadded(cards[3]),
	)

	backBtn := widget.NewButtonWithIcon("  Назад", theme.NavigateBackIcon(), func() {
		ShowDifficulty(w)
	})
	startBtn := widget.NewButtonWithIcon("НАЧАТЬ ИГРУ  ", theme.ConfirmIcon(), func() {
		gs := game.NewGame(diff, classes[selected].class)
		ShowGameScreen(w, gs)
	})
	startBtn.Importance = widget.HighImportance

	nav := container.NewBorder(nil, nil, backBtn, startBtn)

	content := container.NewBorder(
		container.NewPadded(title),
		container.NewPadded(nav),
		nil, nil,
		container.NewVBox(
			cardRow,
			widget.NewSeparator(),
			container.NewPadded(descLabel),
		),
	)

	w.SetContent(container.NewStack(bg, container.NewPadded(content)))
}

func classColor(i int) color.Color {
	colors := []color.Color{
		color.RGBA{180, 180, 170, 255},
		color.RGBA{80, 200, 80, 255},
		color.RGBA{80, 180, 220, 255},
		color.RGBA{200, 60, 50, 255},
	}
	if i < len(colors) {
		return colors[i]
	}
	return colors[0]
}
