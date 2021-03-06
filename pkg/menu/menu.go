package menu

import (
	"fmt"
	"pokered/pkg/audio"
	"pokered/pkg/joypad"
	"pokered/pkg/screen"
	"pokered/pkg/store"
	"sort"
)

// Cancel menu cancel
const Cancel = "CANCEL"

var downArrowBlinkCnt = 6 * 10

// MaxZIndex get max z index
func MaxZIndex() uint {
	selectZ := uint(screen.MenuMin)
	for _, s := range CurSelectMenus {
		if s.z > selectZ {
			selectZ = s.z
		}
	}
	if CurListMenu.z > selectZ {
		return CurListMenu.z
	}
	return selectZ
}

// VBlank script executed in VBlank
func VBlank() {
	listZ, done := CurListMenu.z, false
	sort.Sort(CurSelectMenus)

	newCurSelectMenus := []*SelectMenu{}
	for _, m := range CurSelectMenus {
		if m.z == 0 {
			continue
		}

		if listZ > 0 && listZ < m.z {
			screen.AddLayer("listmenu", int(listZ), CurListMenu.image, 0, 0)
			done = true
		}
		screen.AddLayer(fmt.Sprintf("selectmenu%d", m.z), int(m.z), m.image, 0, 0)
		newCurSelectMenus = append(newCurSelectMenus, m)
	}
	if !done && CurListMenu.z > 0 {
		screen.AddLayer("listmenu", int(CurListMenu.z), CurListMenu.image, 0, 0)
	}
	CurSelectMenus = newCurSelectMenus
}

func HandleMenuInput(current, maxItem uint, wrap bool) uint {
	switch {
	case joypad.Joy5.Up:
		if current > 0 {
			return current - 1
		} else if wrap {
			return maxItem
		}
	case joypad.Joy5.Down:
		if current < maxItem {
			return current + 1
		} else if wrap {
			return 0
		}
	}

	if joypad.Joy5.A || joypad.Joy5.B {
		if !store.Flag.CD60.DontPlayMenuSound {
			audio.PlaySound(audio.SFX_PRESS_AB)
		}
	}
	return current
}
