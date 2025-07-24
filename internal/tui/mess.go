package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	playmanager "github.com/tommjj/music_player/internal/play_manager"
)

type songChangedMsg struct {
	Song *playmanager.Song
}

func newSongChangedMsg(song *playmanager.Song) tea.Cmd {
	return func() tea.Msg {
		return songChangedMsg{Song: song}
	}
}

type songCompletedMsg struct {
	Song *playmanager.Song
}

func newSongCompletedMsg(song *playmanager.Song) tea.Cmd {
	return func() tea.Msg {
		return songCompletedMsg{Song: song}
	}
}

type playModeChangedMsg struct {
	Mode string // e.g., "normal", "repeat", "shuffle"
}

type playErrorMsg struct {
	Error error
}

func tickEverySecond() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type TickMsg time.Time
