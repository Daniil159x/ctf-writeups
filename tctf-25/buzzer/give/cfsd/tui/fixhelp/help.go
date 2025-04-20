package fixhelp

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

var (
	NextWidgetKey key.Binding = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next widget"),
	)
	ExitKey key.Binding = key.NewBinding(
		key.WithKeys("ctrl+c", "q", "next widget"),
		key.WithHelp("q", "exit"),
	)
)

type WrapKeyMap struct {
	binds []key.Binding
}

func FromBindings(binds []key.Binding) WrapKeyMap {
	binds = append(binds, NextWidgetKey, ExitKey)
	return WrapKeyMap{binds}
}

func FromViewport(keys viewport.KeyMap) WrapKeyMap {
	return FromBindings([]key.Binding{
		keys.Up,
		keys.Down,
	})
}

// FullHelp implements help.KeyMap.
func (w WrapKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{w.ShortHelp()}
}

// ShortHelp implements help.KeyMap.
func (w WrapKeyMap) ShortHelp() []key.Binding {
	return w.binds
}
