package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	playmanager "github.com/tommjj/music_player/internal/play_manager"
	"github.com/tommjj/music_player/internal/tui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("please set path")
		return
	}
	songsPath := os.Args[1]

	dirEntris, err := os.ReadDir(songsPath)
	if err != nil {
		panic(err)
	}

	songs := []*playmanager.Song{}

	for _, entry := range dirEntris {
		if entry.IsDir() {
			continue
		}

		song := &playmanager.Song{
			Title:  entry.Name(),
			Artist: "Unknown Artist",
			Album:  "Unknown Album",
			Path:   filepath.Join(songsPath, entry.Name()),
		}

		songs = append(songs, song)
	}

	playManager := playmanager.NewPlayManager()
	playManager.AddSongs(songs...)
	playManager.AutoPlay = true

	model := tui.NewModel(playManager)
	app := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := app.Run(); err != nil {
		println("Error starting TUI:", err.Error())
		os.Exit(1)
	}
}
