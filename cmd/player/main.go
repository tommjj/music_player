package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tommjj/music_player/internal/player"
)

func main() {
	control := make(chan rune)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
			control <- r
		}
	}()

	ctx, cal := context.WithCancel(context.Background())
	player := player.NewPlayer()

	player.SetOnComplete(func() {
		fmt.Println("Playback completed")
		cal()
	})

	filepath.Walk(`D:\Workspace\go\music_player\music`, func(path string, d os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error walking the path %q: %v\n", path, err)
			return nil
		}
		switch filepath.Ext(path) {
		case ".mp3", ".wav":
		default:
			return nil // Skip non-mp3 files
		}

		ctx, cal = context.WithCancel(context.Background())
		fmt.Printf("Play file: %s\n", filepath.Base(path))

		player.SetOnComplete(func() {
			fmt.Printf("Playback of %s completed\n", filepath.Base(path))
			cal()
		})

		err = player.Play(path)
		if err != nil {
			fmt.Printf("Player.Play(%s) failed: %v\n", path, err)
			return nil
		}

		go func() {
			for {
				select {
				case r := <-control:
					switch r {
					case 'q', 'Q':
						fmt.Println("Quitting...")
						os.Exit(0)
					case 'n', 'N': // Next file
						fmt.Println("Skipping to next file...")
						cal() // Cancel current playback
						return
					case ' ':
						if player.IsPaused() {
							fmt.Println("Resuming playback...")
							player.Resume()
						} else {
							fmt.Println("Pausing playback...")
							player.Pause()
						}
					case 'a', 'A': // seek
						fmt.Println("Priv 1 second")
						player.ToPositionByOffset(-1 * time.Second)
					case 's', 'S': // seek
						fmt.Println("Next 1 second")
						player.ToPositionByOffset(time.Second)
					}

				case <-ctx.Done():
					return
				}
			}
		}()

		<-ctx.Done()

		return nil
	})

}
