package script

import (
	"pokered/pkg/joypad"
	"pokered/pkg/menu"
	"pokered/pkg/store"
	"pokered/pkg/text"
	"pokered/pkg/util"
	"pokered/pkg/widget"
	"pokered/pkg/world"
)

// general counter used in various functions
var counter uint

// ScriptMap script ID -> script
var scriptMap = newScriptMap()

func newScriptMap() map[uint]func() {
	result := map[uint]func(){}
	result[store.Overworld] = halt
	result[store.WhiteScreen] = whiteScreen
	result[store.ExecText] = execText
	result[store.WidgetStartMenu] = widgetStartMenu
	result[store.WidgetStartMenu2] = widgetStartMenu2
	result[store.WidgetBag] = widgetBag
	result[store.WidgetTrainerCard] = widgetTrainerCard
	result[store.WidgetPlayerNamingScreen] = widgetPlayerNamingScreen
	result[store.WidgetRivalNamingScreen] = widgetRivalNamingScreen
	result[store.WidgetPartyMenu] = widgetPartyMenu
	result[store.WidgetPartyMenuSelect] = widgetPartyMenuSelect
	result[store.TwoOptionMenu] = handleTwoOption
	result[store.WidgetStats] = widgetStats
	result[store.WidgetStats2] = widgetStats2
	result[store.WidgetPokedexPage] = widgetPokedexPage
	result[store.WidgetStarterPokedexPage] = widgetStarterPokedexPage
	result[store.FadeOutToBlack] = fadeOutToBlack
	result[store.FadeOutToWhite] = fadeOutToWhite
	result[store.LoadMapData] = loadMapData
	result[store.TitleCopyright] = titleCopyright
	result[store.TitleBlank] = titleBlank
	result[store.TitleIntroScene] = titleIntroScene
	result[store.TitleWhiteOut] = titleWhiteOut
	result[store.TitlePokemonRed] = titlePokemonRed
	result[store.TitleMenu] = titleMenu
	result[store.TitleMenu2] = titleMenu2
	result[store.OakSpeech0] = oakSpeech0
	result[store.OakSpeech1] = oakSpeech1
	result[store.OakSpeech2] = oakSpeech2
	result[store.IntroducePlayer] = introducePlayer
	result[store.ChoosePlayerName] = choosePlayerName
	result[store.ChoosePlayerName2] = choosePlayerName2
	result[store.CustomPlayerName] = customPlayerName
	result[store.AfterChoosePlayerName] = afterChoosePlayerName
	result[store.AfterCustomPlayerName] = afterCustomPlayerName
	result[store.IntroduceRival] = introduceRival
	result[store.ChooseRivalName] = chooseRivalName
	result[store.ChooseRivalName2] = chooseRivalName2
	result[store.CustomRivalName] = customRivalName
	result[store.AfterChooseRivalName] = afterChooseRivalName
	result[store.AfterCustomRivalName] = afterCustomRivalName
	result[store.LetsGoPlayer] = letsGoPlayer
	result[store.ShrinkPlayer] = shrinkPlayer
	return result
}

// Current return current script
func Current() func() {
	scr := store.Script()

	switch s := scr.(type) {
	case int:
		sc, ok := scriptMap[uint(s)]
		if !ok {
			util.NotRegisteredError("scriptMap", store.ScriptID())
			return halt
		}
		return sc
	case uint:
		sc, ok := scriptMap[s]
		if !ok {
			util.NotRegisteredError("scriptMap", store.ScriptID())
			return halt
		}
		return sc
	case func():
		return func() {
			s()
			nextScript()
		}
	default:
		return scriptMap[0]
	}
}

func nextScript() {
	if store.ScriptLength() > 1 {
		store.PopScript()
		return
	}
	store.SetScriptID(store.Overworld)
}

func halt() {}

func execText() {
	if len([]rune(text.CurText)) == 0 {
		nextScript()
	}

	if text.InScroll {
		text.ScrollTextUpOneLine(text.TextBoxImage)
		return
	}

	if store.FrameCounter > 0 {
		joypad.Joypad()
		if joypad.JoyHeld.A || joypad.JoyHeld.B {
			store.FrameCounter = 0
			return
		}
		store.FrameCounter--
		if store.FrameCounter > 0 {
			store.DelayFrames = 1
			return
		}
		return
	}

	target := text.TextBoxImage
	if widget.DexPageScreen() != nil {
		target = widget.DexPageScreen()
	}
	text.CurText = text.PlaceStringOneByOne(target, text.CurText)
	if len([]rune(text.CurText)) == 0 {
		nextScript()
	}
}

func fadeOutToBlack() {
	if store.FadeCounter <= 0 {
		store.SetScriptID(store.Overworld)
		return
	}

	store.FadeCounter--

	if store.Palette < 1 {
		store.Palette = 1
		return
	}

	store.Palette--
	store.DelayFrames = 8

	if store.FadeCounter <= 0 {
		store.PopScript()
	}
}

func fadeOutToWhite() {
	if store.FadeCounter <= 0 {
		nextScript()
		return
	}

	store.FadeCounter--

	if store.Palette > 8 {
		store.Palette = 8
		return
	}

	store.Palette++
	store.DelayFrames = 8

	if store.FadeCounter <= 0 {
		nextScript()
	}
}

func loadMapData() {
	mapID, warpID := world.WarpTo[0], world.WarpTo[1]
	if mapID < 0 {
		return
	}
	world.LoadWorldData(mapID)

	// ref: LoadDestinationWarpPosition
	if warpID < 0 {
		return
	}
	warpTo := world.CurWorld.Object.WarpTos[warpID]
	p := store.SpriteData[0]
	p.MapXCoord, p.MapYCoord = warpTo.XCoord, warpTo.YCoord

	store.SetScriptID(store.Overworld)
}

// InOakSpeech returns if game mode is OakSpeech
func InOakSpeech() bool {
	scriptID := store.ScriptID()
	inOakSpeechScript := scriptID >= store.WidgetPlayerNamingScreen && scriptID <= store.ShrinkPlayer
	inText := scriptID == store.ExecText

	if inOakSpeechScript {
		return true
	}

	if inText {
		return OakSpeechScreen != nil
	}

	return false
}

// InTitle returns if game mode is title
func InTitle() bool {
	scriptID := store.ScriptID()
	return scriptID >= store.TitleCopyright && scriptID <= store.TitleMenu2
}

func handleTwoOption() {
	m := menu.CurSelectMenu()
	pressed := menu.HandleSelectMenuInput()

	switch {
	case pressed.A:
		m.Close()
		store.TwoOptionResult = m.Index()
		store.PopScript()
	}
}
