package gui

import (
	"bunker3000/game"
	"bunker3000/player"
	"bunker3000/save"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type gameScreen struct {
	w   fyne.Window
	gs  *game.GameState

	dayLabel    *widget.Label
	healthLabel *widget.Label
	eatLabel    *widget.Label
	waterLabel  *widget.Label
	infoLabel   *widget.Label

	healthClip *clipLayout
	eatClip    *clipLayout
	waterClip  *clipLayout
	healthFG   *canvas.Rectangle
	eatFG      *canvas.Rectangle
	waterFG    *canvas.Rectangle

	eventHeader *canvas.Rectangle
	eventTitle  *widget.Label
	eventDesc   *widget.Label
	choiceBtn1  *widget.Button
	choiceBtn2  *widget.Button
	choices     *fyne.Container

	resultBorder *canvas.Rectangle
	resultLabel  *widget.Label
	resultBox    *fyne.Container
	continueBtn  *widget.Button

	root *fyne.Container
}

func ShowGameScreen(w fyne.Window, gs *game.GameState) {
	s := &gameScreen{w: w, gs: gs}
	s.buildUI()
	s.startDay()
	w.SetContent(container.NewPadded(s.root))
}

func (s *gameScreen) buildUI() {

	dayLbl := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	infoLbl := widget.NewLabel("")
	healthLbl := widget.NewLabel("")
	eatLbl := widget.NewLabel("")
	waterLbl := widget.NewLabel("")

	healthFG := canvas.NewRectangle(color.RGBA{190, 50, 40, 255})
	healthFG.CornerRadius = 3
	eatFG := canvas.NewRectangle(color.RGBA{200, 140, 30, 255})
	eatFG.CornerRadius = 3
	waterFG := canvas.NewRectangle(color.RGBA{40, 140, 190, 255})
	waterFG.CornerRadius = 3

	healthClip := &clipLayout{pct: 1}
	eatClip := &clipLayout{pct: 1}
	waterClip := &clipLayout{pct: 1}

	healthBar := newClipBar(healthFG, healthClip)
	eatBar := newClipBar(eatFG, eatClip)
	waterBar := newClipBar(waterFG, waterClip)

	statusBox := container.NewVBox(
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("СТАТУС", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			dayLbl,
		),
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("Здоровье", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			healthLbl, healthBar),
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("Еда", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			eatLbl, eatBar),
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("Вода", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			waterLbl, waterBar),
		container.NewHBox(widget.NewSeparator(), infoLbl, widget.NewSeparator()),
	)

	s.dayLabel = dayLbl
	s.infoLabel = infoLbl
	s.healthLabel = healthLbl
	s.eatLabel = eatLbl
	s.waterLabel = waterLbl
	s.healthClip = healthClip
	s.eatClip = eatClip
	s.waterClip = waterClip
	s.healthFG = healthFG
	s.eatFG = eatFG
	s.waterFG = waterFG

	hdr := canvas.NewRectangle(color.RGBA{80, 80, 80, 255})
	hdr.CornerRadius = 6
	s.eventHeader = hdr

	evTitle := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	evTitle.Wrapping = fyne.TextWrapWord
	s.eventTitle = evTitle

	evDesc := widget.NewLabel("")
	evDesc.Wrapping = fyne.TextWrapWord
	s.eventDesc = evDesc

	choice1 := widget.NewButton("", nil)
	choice2 := widget.NewButton("", nil)
	choice1.Importance = widget.HighImportance
	choice2.Importance = widget.HighImportance
	s.choiceBtn1 = choice1
	s.choiceBtn2 = choice2

	s.choices = container.NewVBox(
		s.choiceBtn1,
		container.NewPadded(s.choiceBtn2),
	)

	resBorder := canvas.NewRectangle(color.RGBA{60, 140, 60, 255})
	resBorder.CornerRadius = 8
	s.resultBorder = resBorder

	resLbl := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	resLbl.Wrapping = fyne.TextWrapWord
	s.resultLabel = resLbl

	resContent := container.NewStack(resBorder, container.NewPadded(resLbl))
	s.resultBox = container.NewPadded(resContent)
	s.resultBox.Hide()

	contBtn := widget.NewButton("Продолжить →", nil)
	contBtn.Importance = widget.HighImportance
	contBtn.Hide()
	s.continueBtn = contBtn

	eventArea := container.NewVBox(
		statusBox,
		widget.NewSeparator(),
		container.NewStack(hdr, container.NewPadded(evTitle)),
		container.NewPadded(evDesc),
		s.choices,
		s.resultBox,
		container.NewPadded(contBtn),
	)

	s.root = container.NewStack(
		canvas.NewRectangle(color.RGBA{18, 18, 22, 255}),
		container.NewPadded(eventArea),
	)
}

func (s *gameScreen) refreshStatus() {
	p := &s.gs.Player
	cfg := p.GetCfg()
	maxHP := player.MaxHealth(p.Class)

	if s.gs.Endless {
		s.dayLabel.SetText(fmt.Sprintf("День %d (∞)", p.ThisDay))
	} else {
		s.dayLabel.SetText(fmt.Sprintf("День %d из %d", p.ThisDay, cfg.MaxDays))
	}
	s.infoLabel.SetText(fmt.Sprintf("%s • %s", p.Difficulty.String(), p.Class.String()))
	s.healthLabel.SetText(fmt.Sprintf("%d%%", p.Health))
	s.eatLabel.SetText(fmt.Sprintf("%d / %d", p.Eat, cfg.MaxResource))
	s.waterLabel.SetText(fmt.Sprintf("%d / %d", p.Water, cfg.MaxResource))

	s.healthClip.pct = float64(p.Health) / float64(maxHP)
	if s.healthClip.pct < 0 {
		s.healthClip.pct = 0
	}
	s.eatClip.pct = float64(p.Eat) / float64(cfg.MaxResource)
	if s.eatClip.pct < 0 {
		s.eatClip.pct = 0
	}
	s.waterClip.pct = float64(p.Water) / float64(cfg.MaxResource)
	if s.waterClip.pct < 0 {
		s.waterClip.pct = 0
	}

	s.healthFG.Refresh()
	s.eatFG.Refresh()
	s.waterFG.Refresh()
}

func (s *gameScreen) switchToChoice() {
	s.choices.Show()
	s.resultBox.Hide()
	s.continueBtn.Hide()
}

func (s *gameScreen) switchToResult(msg string, isGood bool) {
	s.choices.Hide()
	if isGood {
		s.resultBorder.FillColor = color.RGBA{60, 140, 60, 255}
	} else {
		s.resultBorder.FillColor = color.RGBA{160, 50, 40, 255}
	}
	s.resultBorder.Refresh()
	s.resultLabel.SetText(msg)
	s.resultBox.Show()
	s.continueBtn.Show()
}

func (s *gameScreen) goToGameOver() {
	s.gs.Save()
	save.DeleteSave()
	ShowGameOverScreen(s.w, s.gs)
}

func (s *gameScreen) startDay() {
	s.gs.StartDay()

	if s.gs.DayError != nil {
		s.switchToResult(s.gs.DayError.Error(), false)
		s.continueBtn.SetText("К статистике")
		s.continueBtn.OnTapped = s.goToGameOver
		return
	}

	if s.gs.IsGameOver() {
		s.goToGameOver()
		return
	}

	s.refreshStatus()

	ev := s.gs.CurrentEvent
	if ev == nil {
		s.switchToResult("Ошибка: событие не получено", false)
		s.continueBtn.SetText("К статистике")
		s.continueBtn.OnTapped = s.goToGameOver
		return
	}

	s.eventHeader.FillColor = catColor(ev.Category)
	s.eventHeader.Refresh()
	s.eventTitle.SetText(eventCategoryLabel(ev.Category))
	s.eventDesc.SetText(ev.Message)

	s.choiceBtn1.SetText("1. " + ev.Variants[0])
	s.choiceBtn2.SetText("2. " + ev.Variants[1])

	s.choiceBtn1.OnTapped = func() { s.executeChoice(1) }
	s.choiceBtn2.OnTapped = func() { s.executeChoice(2) }

	s.continueBtn.SetText("Продолжить →")
	s.switchToChoice()
}

func (s *gameScreen) executeChoice(choice int) {
	resultMsg, _ := s.gs.ExecuteChoice(choice)
	s.refreshStatus()

	if s.gs.IsGameOver() {
		s.switchToResult(resultMsg, false)
		s.continueBtn.SetText("К статистике")
		s.continueBtn.OnTapped = s.goToGameOver
		return
	}

	s.switchToResult(resultMsg, isGoodResult(resultMsg))
	s.continueBtn.OnTapped = func() {
		s.gs.NextDay()
		s.startDay()
	}
}

func newClipBar(fg *canvas.Rectangle, cl *clipLayout) *fyne.Container {
	bg := canvas.NewRectangle(color.RGBA{40, 40, 45, 255})
	bg.CornerRadius = 3
	c := container.NewStack(bg, fg)
	c.Layout = &stackClipLayout{bg: bg, fg: fg, cl: cl}
	return c
}

type stackClipLayout struct {
	bg *canvas.Rectangle
	fg *canvas.Rectangle
	cl *clipLayout
}

func (l *stackClipLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	l.bg.Resize(size)
	l.bg.Move(fyne.NewPos(0, 0))

	w := float32(l.cl.pct) * size.Width
	if w < 0 {
		w = 0
	}
	if w > size.Width {
		w = size.Width
	}
	l.fg.Resize(fyne.NewSize(w, size.Height))
	l.fg.Move(fyne.NewPos(0, 0))
}

func (l *stackClipLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return l.bg.MinSize()
}

func eventCategoryLabel(cat string) string {
	switch cat {
	case "resource":
		return "[РЕСУРСЫ]"
	case "combat":
		return "[БИТВА]"
	case "helper":
		return "[ПОМОЩЬ]"
	case "weather":
		return "[ПОГОДА]"
	case "hazard":
		return "[ОПАСНОСТЬ]"
	case "find":
		return "[НАХОДКА]"
	default:
		return "[СОБЫТИЕ]"
	}
}

func isGoodResult(msg string) bool {
	return !strings.Contains(msg, "[-]") &&
		!strings.Contains(msg, "СМЕРТЬ") &&
		!strings.Contains(msg, "ТРАВМА") &&
		!strings.Contains(msg, "ОПАСНОСТЬ") &&
		!strings.Contains(msg, "ОТРАВЛЕНИЕ") &&
		!strings.Contains(msg, "ПАДЕНИЕ") &&
		!strings.Contains(msg, "БЕГСТВО") &&
		!strings.Contains(msg, "НЕУДАЧА")
}

func catColor(cat string) color.Color {
	switch cat {
	case "resource":
		return color.RGBA{50, 130, 50, 255}
	case "combat":
		return color.RGBA{160, 40, 40, 255}
	case "helper":
		return color.RGBA{40, 80, 150, 255}
	case "weather":
		return color.RGBA{40, 120, 140, 255}
	case "hazard":
		return color.RGBA{160, 100, 30, 255}
	case "find":
		return color.RGBA{120, 50, 140, 255}
	default:
		return color.RGBA{80, 80, 80, 255}
	}
}
