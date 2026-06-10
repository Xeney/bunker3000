package gui

import (
	"bunker3000/player"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewStatusPanel(p *player.Player) *fyne.Container {
	cfg := p.GetCfg()

	dayLabel := widget.NewLabelWithStyle(
		fmt.Sprintf("День %d из %d", p.ThisDay, cfg.MaxDays),
		fyne.TextAlignTrailing,
		fyne.TextStyle{Bold: true},
	)
	infoLabel := widget.NewLabel(
		fmt.Sprintf("%s • %s", p.Difficulty.String(), p.Class.String()),
	)

	healthBar := newColoredProgress(float64(p.Health)/float64(player.MaxHealth(p.Class)),
		color.RGBA{200, 60, 50, 255}, "Здоровье", fmt.Sprintf("%d%%", p.Health))

	eatBar := newColoredProgress(float64(p.Eat)/float64(cfg.MaxResource),
		color.RGBA{210, 150, 40, 255}, "Еда", fmt.Sprintf("%d / %d", p.Eat, cfg.MaxResource))

	waterBar := newColoredProgress(float64(p.Water)/float64(cfg.MaxResource),
		color.RGBA{50, 150, 200, 255}, "Вода", fmt.Sprintf("%d / %d", p.Water, cfg.MaxResource))

	return container.NewVBox(
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("СТАТУС", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			dayLabel,
		),
		healthBar,
		eatBar,
		waterBar,
		container.NewHBox(
			widget.NewSeparator(),
			infoLabel,
			widget.NewSeparator(),
		),
	)
}

func newColoredProgress(pct float64, barColor color.Color, label, value string) *fyne.Container {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}

	bg := canvas.NewRectangle(color.RGBA{40, 40, 45, 255})
	bg.CornerRadius = 4

	fg := canvas.NewRectangle(barColor)
	fg.CornerRadius = 4

	bar := container.NewStack(bg, fg)
	bar.Layout = &clipLayout{pct: pct}

	labelW := widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	valueW := widget.NewLabelWithStyle(value, fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true})

	return container.NewBorder(nil, nil, labelW, valueW, bar)
}

type clipLayout struct {
	pct float64
}

func (c *clipLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	bg := objects[0]
	fg := objects[1]
	bg.Resize(size)
	bg.Move(fyne.NewPos(0, 0))

	w := float32(c.pct) * size.Width
	if w < 0 {
		w = 0
	}
	if w > size.Width {
		w = size.Width
	}
	fg.Resize(fyne.NewSize(w, size.Height))
	fg.Move(fyne.NewPos(0, 0))
}

func (c *clipLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.Size{}
	}
	return objects[0].MinSize()
}

func newSectionHeader(text string) *fyne.Container {
	line := canvas.NewRectangle(color.RGBA{230, 170, 60, 255})
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return container.NewBorder(nil, nil, nil, nil,
		container.NewHBox(
			label,
			container.NewPadded(line),
		),
	)
}

func newCategoryLabel(cat string) *canvas.Rectangle {
	var bgColor color.Color
	switch cat {
	case "resource":
		bgColor = color.RGBA{50, 130, 50, 255}
	case "combat":
		bgColor = color.RGBA{160, 40, 40, 255}
	case "helper":
		bgColor = color.RGBA{40, 80, 150, 255}
	case "weather":
		bgColor = color.RGBA{40, 120, 140, 255}
	case "hazard":
		bgColor = color.RGBA{160, 100, 30, 255}
	case "find":
		bgColor = color.RGBA{120, 50, 140, 255}
	default:
		bgColor = color.RGBA{80, 80, 80, 255}
	}
	rect := canvas.NewRectangle(bgColor)
	rect.CornerRadius = 4
	return rect
}

func choiceButton(text string, tapped func()) *widget.Button {
	btn := widget.NewButton(text, tapped)
	return btn
}

func newResultCard(msg string, isGood bool) *fyne.Container {
	borderColor := color.RGBA{80, 180, 80, 255}
	if !isGood {
		borderColor = color.RGBA{200, 60, 50, 255}
	}

	label := widget.NewLabelWithStyle(msg, fyne.TextAlignCenter, fyne.TextStyle{Bold: false})
	label.Wrapping = fyne.TextWrapWord

	border := canvas.NewRectangle(borderColor)
	border.CornerRadius = 8

	content := container.NewBorder(nil, nil, nil, nil,
		container.NewPadded(label),
	)
	content = container.NewStack(border, content)
	content = container.NewPadded(content)
	return content
}

func resourceIcon(char rune) *canvas.Text {
	icon := canvas.NewText(string(char), color.RGBA{230, 170, 60, 255})
	icon.TextSize = 24
	icon.TextStyle = fyne.TextStyle{Bold: true}
	return icon
}
