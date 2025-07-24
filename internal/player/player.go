// Package player provides an audio player that can play, pause, and control audio playback.
// this is a wrapper around the beep library for audio playback.
package player

import (
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var (
	DefaultAudioQuality = 4 // Resampling quality
)

// Player represents an audio player that can play, pause, and control audio playback.
type Player struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeekCloser
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	filepath   string

	onComplete func()

	quality     int     // quality is the resampling quality for audio playback.
	volumeValue float64 // volumeValue is the current volume level.
	radioValue  float64 // radioValue is the current radio volume level.
}

func NewPlayer() *Player {
	return &Player{
		quality:     DefaultAudioQuality,
		volumeValue: 0,   // Default volume level
		radioValue:  1.0, // Default radio volume level
	}
}

func (p *Player) SetOnComplete(callback func()) {
	// SetOnComplete sets a callback function to be called when playback completes.
	p.onComplete = callback
}

// Play starts playing the audio file specified by filename.
// It initializes the audio system, opens the file, and starts playback.
func (p *Player) Play(filename string) error {
	if p.streamer != nil {
		p.Close() // Close any existing stream before playing a new one
	}

	streamer, format, err := p.loadStreamer(filename)
	if err != nil {
		return err
	}

	p.filepath = filename
	p.sampleRate = format.SampleRate
	p.streamer = streamer

	//
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/5))

	p.ctrl = &beep.Ctrl{Streamer: beep.Seq(streamer, beep.Callback(func() {
		if p.onComplete != nil {
			go p.onComplete()
		}
	}))}
	p.resampler = beep.ResampleRatio(p.quality, p.radioValue, p.ctrl)
	p.volume = &effects.Volume{Streamer: p.resampler, Base: 2, Volume: p.volumeValue}

	speaker.Play(p.volume)
	return nil
}

func (p *Player) Pause() {
	if p.ctrl == nil {
		return
	}

	speaker.Lock()
	defer speaker.Unlock()

	p.ctrl.Paused = true
}

func (p *Player) Resume() {
	if p.ctrl == nil {
		return
	}

	speaker.Lock()
	defer speaker.Unlock()
	p.ctrl.Paused = false
}

func (p *Player) IsPaused() bool {
	if p.ctrl == nil {
		return true // If no control, consider it paused
	}
	return p.ctrl.Paused
}

func (p *Player) VolumeUp() {
	p.volumeValue += 0.1

	if p.volume == nil {
		return
	}

	speaker.Lock()
	defer speaker.Unlock()

	p.volume.Volume = p.volumeValue
}

func (p *Player) VolumeDown() {
	p.volumeValue -= 0.1

	if p.volume == nil {
		return
	}
	speaker.Lock()
	defer speaker.Unlock()

	p.volume.Volume = p.volumeValue
}

func (p *Player) SetVolume(volume float64) {
	p.volumeValue = volume

	if p.volume == nil {
		return
	}
	speaker.Lock()
	defer speaker.Unlock()
	p.volume.Volume = p.volumeValue

}

func (p *Player) ToPosition(pos time.Duration) error {
	if p.streamer == nil {
		return os.ErrInvalid // No file loaded
	}

	speaker.Lock()
	defer speaker.Unlock()

	length := p.sampleRate.D(p.streamer.Len())
	if pos < 0 || pos >= length {
		return os.ErrInvalid // Position out of bounds
	}

	if err := p.streamer.Seek(p.sampleRate.N(pos)); err != nil {
		return err
	}
	return nil
}

// ToPositionByOffset moves the playback position by a specified offset.
func (p *Player) ToPositionByOffset(offset time.Duration) error {
	if p.streamer == nil {
		return os.ErrInvalid // No file loaded
	}

	speaker.Lock()
	defer speaker.Unlock()

	newPos := p.streamer.Position()
	newPos += p.sampleRate.N(offset)

	if newPos < 0 {
		newPos = 0
	} else if newPos >= p.streamer.Len() {
		newPos = p.streamer.Len() - 1
	}

	if err := p.streamer.Seek(newPos); err != nil {
		return err
	}

	return nil
}

// Replay currently loaded audio file.
func (p *Player) Replay() error {
	if p.filepath == "" {
		return os.ErrInvalid // No file loaded
	}

	p.Play(p.filepath)

	return nil
}

// Close stops playback and releases resources.
// It closes the speaker and the streamer, and resets the player state.
func (p *Player) Close() {
	speaker.Lock()
	defer speaker.Unlock()

	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}
	if p.ctrl != nil {
		p.ctrl.Paused = true
	}

	p.ctrl = nil
	p.resampler = nil
	p.volume = nil
	p.sampleRate = 0
	p.filepath = ""
}

type Info struct {
	Filepath string
	Current  time.Duration
	Length   time.Duration
	Volume   float64
	Speed    float64
	Paused   bool
}

// Info returns the current playback information.
// It includes the file path, current position, total length, volume, speed, and paused state.
func (p *Player) Info() *Info {
	speaker.Lock()
	defer speaker.Unlock()
	if p.streamer == nil {
		return &Info{
			Filepath: "",
			Current:  0,
			Length:   0,
			Volume:   p.volumeValue,
			Speed:    p.radioValue,
			Paused:   true,
		}
	}
	return &Info{
		Filepath: p.filepath,
		Current:  p.sampleRate.D(p.streamer.Position()),
		Length:   p.sampleRate.D(p.streamer.Len()),
		Volume:   p.volumeValue,
		Speed:    p.radioValue,
		Paused:   p.ctrl.Paused,
	}
}

// Auto loads the audio file by file format
// just supports mp3 and wav formats
func (p *Player) loadStreamer(filename string) (beep.StreamSeekCloser, beep.Format, error) {
	ext := filepath.Ext(filename)
	if ext != ".mp3" && ext != ".wav" {
		return nil, beep.Format{}, os.ErrInvalid
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, beep.Format{}, err
	}

	var streamer beep.StreamSeekCloser
	var format beep.Format

	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	}

	if err != nil {
		f.Close() // đóng nếu decode thất bại
		return nil, beep.Format{}, err
	}

	return streamer, format, nil
}
