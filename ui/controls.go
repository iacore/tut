package ui

import (
	"github.com/gdamore/tcell/v2"

	"github.com/RasmusLindroth/tut/config"
)

type Control struct {
	key   config.InputAction
	Label string
	Len   int
}

func NewControl(c *config.Config, k config.InputAction, first bool) Control {
	label, length := config.ColorFromKey(c, k, first)
	return Control{
		key:   k,
		Label: label,
		Len:   length,
	}
}

func (c Control) Click() *tcell.EventKey {
	for _, k := range c.key.Keys {
		return tcell.NewEventKey(k, 0, tcell.ModNone)
	}
	for _, r := range c.key.Runes {
		return tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone)
	}
	return tcell.NewEventKey(tcell.KeyRune, 0, tcell.ModNone)
}
