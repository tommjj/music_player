// Package tui provides a terminal user interface for the music player.
package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	playmanager "github.com/tommjj/music_player/internal/play_manager"
)

var docStyle = lipgloss.NewStyle().Margin(0, 1, 0, 1)

var ()

type Model struct {
	playmanager *playmanager.PlayManager

	list list.Model

	progress progress.Model
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tickEverySecond()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(Item); ok {
				if err := m.playmanager.PlaySong(item.Song); err != nil {
					return m, nil
				}
				return m, newSongChangedMsg(item.Song)
			}
		case " ":
			if m.playmanager.Player.IsPaused() {
				m.playmanager.Player.Resume()
			} else {
				m.playmanager.Player.Pause()
			}
		case "z":
			m.playmanager.Player.ToPositionByOffset(time.Second * -10)
		case "x":
			m.playmanager.Player.ToPositionByOffset(time.Second * 10)
		case "a":
			m.playmanager.PlayPrevious()
		case "s":
			m.playmanager.PlayNext()
		case "p":
			m.playmanager.PlayCurrent()
		case "m":
			var cmd tea.Cmd
			switch m.playmanager.PlayMode() {
			case playmanager.PlayModeNormal:
				cmd = m.setPlayMode(playmanager.PlayModeShuffle)
			case playmanager.PlayModeShuffle:
				cmd = m.setPlayMode(playmanager.PlayModeRepeat)
			case playmanager.PlayModeRepeat:
				cmd = m.setPlayMode(playmanager.PlayModeNormal)
			}
			return m, cmd
		case "?":
			m.list.Help.ShowAll = true
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}

	case songChangedMsg:
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-4)
		m.progress.Width = msg.Width - h
	case TickMsg:
		return m, tickEverySecond()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) setPlayMode(playMode string) tea.Cmd {
	m.playmanager.SetPlayMode(playMode)

	playlist := m.playmanager.PlayList()
	items := make([]list.Item, len(playlist))
	for i, song := range playlist {
		items[i] = Item{Song: song}
	}

	return m.list.SetItems(items)
}

func (m Model) View() string {
	info := m.playmanager.Player.Info()
	title := "no song play"
	currentSong, _ := m.playmanager.GetCurrentSong()
	if currentSong != nil {
		title = currentSong.Title
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		docStyle.Render(m.list.View()),
		docStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				"",
				fmt.Sprintf("[%v] %v", m.playmanager.PlayMode(), title),
				"",
				m.progress.ViewAs(float64(info.Current)/float64(info.Length)),
			),
		),
	)
}

func NewModel(pm *playmanager.PlayManager) Model {
	playlist := pm.PlayList()

	items := make([]list.Item, len(playlist))
	for i, song := range playlist {
		items[i] = Item{Song: song}
	}

	list := list.New(items, newItemDelegate(newDelegateKeyMap()), 0, 0)

	list.SetShowStatusBar(false)
	list.SetShowTitle(false)
	list.SetShowHelp(false)

	prs := progress.New(progress.WithScaledGradient("#2f0bfdff", "#2f0bfdff"))
	prs.ShowPercentage = false

	return Model{
		playmanager: pm,
		list:        list,
		progress:    prs,
	}
}
