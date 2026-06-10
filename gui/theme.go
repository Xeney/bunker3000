package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type bunkerTheme struct {
	base fyne.Theme
}

func NewTheme() fyne.Theme {
	return &bunkerTheme{base: theme.DefaultTheme()}
}

func (t *bunkerTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{18, 18, 22, 255}
	case theme.ColorNameForeground:
		return color.RGBA{220, 220, 210, 255}
	case theme.ColorNamePrimary:
		return color.RGBA{230, 170, 60, 255}
	case theme.ColorNameInputBackground:
		return color.RGBA{30, 30, 36, 255}
	case theme.ColorNameInputBorder:
		return color.RGBA{60, 55, 45, 255}
	case theme.ColorNameDisabled:
		return color.RGBA{60, 60, 55, 255}
	case theme.ColorNameDisabledButton:
		return color.RGBA{40, 40, 38, 255}
	case theme.ColorNameError:
		return color.RGBA{200, 60, 50, 255}
	case theme.ColorNameSuccess:
		return color.RGBA{80, 180, 80, 255}
	case theme.ColorNameWarning:
		return color.RGBA{210, 150, 40, 255}
	case theme.ColorNameHover:
		return color.RGBA{60, 55, 45, 255}
	case theme.ColorNamePressed:
		return color.RGBA{80, 70, 50, 255}
	case theme.ColorNameFocus:
		return color.RGBA{230, 170, 60, 180}
	case theme.ColorNameHeaderBackground:
		return color.RGBA{28, 28, 32, 255}
	case theme.ColorNameMenuBackground:
		return color.RGBA{24, 24, 30, 255}
	case theme.ColorNameOverlayBackground:
		return color.RGBA{24, 24, 30, 220}
	case theme.ColorNameScrollBar:
		return color.RGBA{60, 55, 45, 200}
	case theme.ColorNameSeparator:
		return color.RGBA{50, 45, 38, 255}
	case theme.ColorNameSelection:
		return color.RGBA{230, 170, 60, 80}
	case theme.ColorNameShadow:
		return color.RGBA{0, 0, 0, 120}
	case theme.ColorNamePlaceHolder:
		return color.RGBA{120, 120, 110, 255}
	default:
		return t.base.Color(name, variant)
	}
}

func (t *bunkerTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t *bunkerTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t *bunkerTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInputRadius:
		return 6
	case theme.SizeNameSelectionRadius:
		return 4
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 18
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInnerPadding:
		return 6
	case theme.SizeNameLineSpacing:
		return 4
	default:
		return t.base.Size(name)
	}
}
