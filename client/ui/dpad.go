package ui

import (
	"gioui.org/io/key"
)

// RemoteKey mendefinisikan tombol remote TV yang dinormalisasi
type RemoteKey int

const (
	KeyNone RemoteKey = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeySelect // OK / DpadCenter / Enter
	KeyBack   // Back / Escape
	KeyMenu   // Menu / Context
	KeyZoom   // Shortcut tombol kaca pembesar (misal: tombol 'Z', Media Play, atau Dpad Up saat di mode baca)
)

// MapKeyEvent mengonversi event hardware key dari Android TV / Keyboard ke RemoteKey
func MapKeyEvent(e key.Event) RemoteKey {
	if e.State != key.Press {
		return KeyNone
	}

	switch e.Name {
	case key.NameUpArrow:
		return KeyUp
	case key.NameDownArrow:
		return KeyDown
	case key.NameLeftArrow:
		return KeyLeft
	case key.NameRightArrow:
		return KeyRight
	case key.NameReturn, key.NameEnter, key.NameSpace, "DpadCenter", "Select", "ButtonA", "Center":
		return KeySelect
	case key.NameEscape, key.NameBack, "ButtonB":
		return KeyBack
	case "Menu", "M", "Context":
		return KeyMenu
	case "Z", "MediaPlayPause", "MediaPlay":
		return KeyZoom
	default:
		return KeyNone
	}
}
