package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	playmanager "github.com/tommjj/music_player/internal/play_manager"
)

var (
	SongsPath = `D:\Workspace\go\music_player\music`
)

func main() {
	dirEntris, err := os.ReadDir(SongsPath)
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
			Path:   filepath.Join(SongsPath, entry.Name()),
		}

		songs = append(songs, song)
	}

	playManager := playmanager.NewPlayManager()
	playManager.AddSongs(songs...)

	playManager.SetPlayMode(playmanager.PlayModeNormal)
	playManager.Player.SetOnComplete(func() {
		println("Playback completed")

		if err := playManager.PlayNext(); err != nil {
			println("Error playing next song:", err.Error())
		}
	})

	scanner := bufio.NewScanner(os.Stdin)
	for {
		println("Enter command (play, pause/resume, next, prev, exit, set [mode], seek [seconds], ls, go [index], info):")
		scanner.Scan()
		command := scanner.Text()
		command = strings.TrimSpace(command)
		command, value, _ := strings.Cut(command, " ")

		switch command {
		case "play":
			if err := playManager.PlayCurrent(); err != nil {
				println("Error playing current song:", err.Error())
			}
		case "pause":
			playManager.Player.Pause()
			println("Playback paused")
		case "resume":
			playManager.Player.Resume()
			println("Playback resumed")
		case "next":
			if err := playManager.PlayNext(); err != nil {
				println("Error playing next song:", err.Error())
			}
		case "prev":
			if err := playManager.PlayPrevious(); err != nil {
				println("Error playing previous song:", err.Error())
			}
		case "set":
			if value == "" {
				println("Usage: set [mode]")
				continue
			}
			if err := playManager.SetPlayMode(value); err != nil {
				println("Error setting play mode:", err.Error())
			}
		case "seek":
			if value == "" {
				println("Usage: seek [seconds]")
				continue
			}

			seconds, err := time.ParseDuration(value + "s")
			if err != nil {
				println("Invalid duration:", err.Error())
				continue
			}
			if err := playManager.Player.ToPositionByOffset(seconds); err != nil {
				println("Error seeking:", err.Error())
			} else {
				println("Seeked to", seconds)
			}
		case "ls":
			println("Current Playlist:")
			for i, song := range playManager.PlayList() {
				println(i+1, "-", song.Title)
			}
		case "go":
			if value == "" {
				println("Usage: go [index]")
				continue
			}
			index, err := strconv.Atoi(value)
			if err != nil {
				println("Invalid index:", err.Error())
				continue
			}
			index-- // Convert to zero-based index
			if err := playManager.PlaySongByIndex(index); err != nil {
				println("Error playing song at index", index, ":", err.Error())
			} else {
				println("Playing song at index", index)
			}
		case "info":

			info := playManager.Player.Info()
			fmt.Println("Current Song Info:")
			fmt.Println("Filepath:", info.Filepath)
			fmt.Printf("Position: %v / %v\n", info.Current.Round(time.Second), info.Length.Round(time.Second))
			fmt.Println("Volume:", info.Volume)
			fmt.Println("Speed:", info.Speed)
			fmt.Println("Paused:", info.Paused)
		case "exit":
			return
		default:
			println("Unknown command")
		}
	}
}
