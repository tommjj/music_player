package player

import (
	"os"
	"testing"
	"time"
)

func TestPlayer_Play(t *testing.T) {
	player := NewPlayer()
	fileName := `D:\Workspace\go\music_player\music\(君の名は  Kimi no Na wa) Nandemonaiya  Kamishiraishi Mone (Maxone Remix) ♪.mp3`

	err := player.Play(fileName)

	if err != nil {
		t.Errorf("Player.Play(%s) failed: %v", fileName, err)
	} else {
		t.Logf("Player.Play(%s) succeeded", fileName)
	}

	time.Sleep(10 * time.Second) // Allow some time for playback
}

func TestReadDirTime(t *testing.T) {
	dir := `D:\Workspace\go\music_player\music`

	start := time.Now()
	_, error := os.ReadDir(dir)
	t.Logf("ReadDir(%s) took %v", dir, time.Since(start))
	if error != nil {
		t.Errorf("ReadDir(%s) failed: %v", dir, error)
		return
	}
}
