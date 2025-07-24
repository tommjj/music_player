package playmanager

import (
	"errors"
	"math/rand"

	"github.com/tommjj/music_player/internal/player"
)

type Song struct {
	Title  string
	Artist string
	Album  string
	Path   string
}

var (
	PlayModeNormal  = "normal"
	PlayModeRepeat  = "repeat"
	PlayModeShuffle = "shuffle"
)

var (
	ErrPlaylistEmpty  = errors.New("playlist is empty")
	ErrSongNotFound   = errors.New("song not found in playlist")
	ErrInvalidIndex   = errors.New("invalid index")
	ErrPlayerNotReady = errors.New("player is not ready")
)

type PlayManager struct {
	playlist []*Song
	// currentIndex is the index of the currently playing song in the playlist.
	currentIndex int
	shuffleList  []int

	playMode string // "normal", "repeat", "shuffle"
	AutoPlay bool

	// Event handlers
	OnCompleted   func(song *Song)
	OnPlay        func(song *Song)
	OnListChanged func(playlist []*Song)

	Player *player.Player
}

func NewPlayManager() *PlayManager {

	player := player.NewPlayer()

	manager := &PlayManager{
		playlist:     []*Song{},
		currentIndex: 0,
		playMode:     PlayModeNormal,
		Player:       player,
	}

	return manager
}

func (pm *PlayManager) SetSongs(songs []*Song) {
	pm.playlist = songs

	if pm.OnListChanged != nil {
		pm.OnListChanged(pm.playlist)
	}
}

func (pm *PlayManager) AddSongs(songs ...*Song) {
	pm.playlist = append(pm.playlist, songs...)

	if pm.OnListChanged != nil {
		pm.OnListChanged(pm.playlist)
	}
}

func (pm *PlayManager) RemoveSong(song *Song) {
	for i, s := range pm.playlist {
		if s.Path == song.Path {
			pm.playlist = append(pm.playlist[:i], pm.playlist[i+1:]...)
			if pm.currentIndex >= i {
				pm.currentIndex--
			}
			break
		}
	}

	if pm.OnListChanged != nil {
		pm.OnListChanged(pm.playlist)
	}
}

func (pm *PlayManager) RemoveSongByIndex(index int) {
	if index < 0 || index >= len(pm.playlist) {
		return
	}
	pm.playlist = append(pm.playlist[:index], pm.playlist[index+1:]...)
	if pm.currentIndex >= index {
		pm.currentIndex--
	}

	if pm.OnListChanged != nil {
		pm.OnListChanged(pm.playlist)
	}
}

func (pm *PlayManager) GetCurrentSong() (*Song, error) {
	if pm.currentIndex < 0 || pm.currentIndex >= len(pm.playlist) {
		return nil, ErrInvalidIndex
	}
	switch pm.playMode {
	case PlayModeShuffle:
		if len(pm.shuffleList) != len(pm.playlist) {
			pm.initShuffleList()
		}
		return pm.playlist[pm.shuffleList[pm.currentIndex]], nil
	default:
		return pm.playlist[pm.currentIndex], nil
	}
}

func (pm *PlayManager) initOnCompleteEventHandler(song *Song) {
	pm.Player.SetOnComplete(func() {
		if pm.OnCompleted != nil {
			pm.OnCompleted(song)
		}

		// Auto Play mode
		if pm.AutoPlay {
			pm.PlayNext()
		}
	})
}

func (pm *PlayManager) play(song *Song) error {
	if err := pm.Player.Play(song.Path); err != nil {
		return err
	}
	// Notify the onPlay callback if set
	if pm.OnPlay != nil {
		pm.OnPlay(song)
	}

	// Add on
	pm.initOnCompleteEventHandler(song)

	return nil
}

func (pm *PlayManager) PlayCurrent() error {
	if pm.Player == nil {
		return ErrPlayerNotReady
	}

	currentSong, err := pm.GetCurrentSong()
	if err != nil {
		return err
	}

	if err := pm.play(currentSong); err != nil {
		return err
	}

	return nil
}

func (pm *PlayManager) PlayNext() error {
	if len(pm.playlist) == 0 {
		return ErrPlaylistEmpty
	}

	switch pm.playMode {
	case PlayModeNormal, "", PlayModeShuffle:
		pm.currentIndex++
		if pm.currentIndex >= len(pm.playlist) {
			pm.currentIndex = 0 // Loop back to the start
		}
	case PlayModeRepeat:
		// Do nothing, stay on the current song
	}

	return pm.PlayCurrent()
}

func (pm *PlayManager) PlayPrevious() error {
	if len(pm.playlist) == 0 {
		return ErrPlaylistEmpty
	}

	switch pm.playMode {
	case PlayModeNormal, "", PlayModeShuffle:
		pm.currentIndex--
		if pm.currentIndex < 0 {
			pm.currentIndex = len(pm.playlist) - 1 // Loop back to the end
		}
	case PlayModeRepeat:
		// Do nothing, stay on the current song
	}

	return pm.PlayCurrent()
}

func (pm *PlayManager) PlaySong(song *Song) error {
	if pm.Player == nil {
		return ErrPlayerNotReady
	}

	songPlaylistIndex := -1
	for i, s := range pm.playlist {
		if s.Path == song.Path {
			songPlaylistIndex = i
			break
		}
	}

	if songPlaylistIndex > 0 {
		pm.currentIndex = songPlaylistIndex
	}

	if err := pm.play(song); err != nil {
		return err
	}

	return nil
}

func (pm *PlayManager) PlaySongByIndex(index int) error {
	if index < 0 || index >= len(pm.playlist) {
		return ErrInvalidIndex
	}
	pm.currentIndex = index
	return pm.PlayCurrent()
}

func (pm *PlayManager) PlayMode() string {
	return pm.playMode
}

func (pm *PlayManager) SetPlayMode(mode string) error {
	if mode == pm.playMode {
		return nil // No change needed
	}

	switch mode {
	case PlayModeNormal, "":
		pm.currentIndex = 0
		pm.playMode = mode
		return nil
	case PlayModeRepeat:
		pm.playMode = mode
		return nil
	case PlayModeShuffle:
		pm.playMode = mode
		pm.initShuffleList()
		pm.currentIndex = 0
		return nil
	default:
		return errors.New("invalid play mode")
	}
}

func (pm *PlayManager) initShuffleList() {
	pm.shuffleList = make([]int, len(pm.playlist))
	for i := range pm.shuffleList {
		pm.shuffleList[i] = i
	}

	shuffle(pm.shuffleList)
}

func shuffle(slice []int) []int {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func (pm *PlayManager) ResetShuffleList() {
	if pm.playMode != PlayModeShuffle {
		return
	}
	pm.initShuffleList()
	pm.currentIndex = 0
}

func (pm *PlayManager) PlayList() []*Song {
	if pm.playMode == PlayModeShuffle {
		if len(pm.shuffleList) != len(pm.playlist) {
			pm.initShuffleList()
		}
		shuffledPlaylist := make([]*Song, len(pm.playlist))
		for i, idx := range pm.shuffleList {
			shuffledPlaylist[i] = pm.playlist[idx]
		}
		return shuffledPlaylist
	}

	return pm.playlist
}
