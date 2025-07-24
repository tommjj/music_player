package tui

import (
	playmanager "github.com/tommjj/music_player/internal/play_manager"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type Item struct {
	Song *playmanager.Song
}

func (i Item) Title() string       { return i.Song.Title }
func (i Item) Description() string { return i.Song.Artist + " - " + i.Song.Album }
func (i Item) FilterValue() string { return i.Song.Title }

type delegateKeyMap struct {
	choose     key.Binding
	togglePlay key.Binding
	play       key.Binding
	next       key.Binding
	priv       key.Binding
	changeMode key.Binding
	next10s    key.Binding
	priv10s    key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.togglePlay,
		d.play,
		d.next,
		d.priv,
		d.next10s,
		d.priv10s,
		d.changeMode,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.togglePlay,
			d.play,
			d.next,
			d.priv,
			d.next10s,
			d.priv10s,
			d.changeMode,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "paly song"),
		),
		togglePlay: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "toggle play/pause"),
		),
		play: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "play"),
		),
		next: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "next"),
		),
		priv: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "priv"),
		),
		changeMode: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "change mode"),
		),
		next10s: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "next 10s"),
		),
		priv10s: key.NewBinding(
			key.WithKeys("z"),
			key.WithHelp("z", "priv 10s"),
		),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	help := []key.Binding{
		keys.choose,
		keys.togglePlay,
		keys.play,
		keys.next,
		keys.priv,
		keys.next10s,
		keys.priv10s,
		keys.changeMode,
	}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}
