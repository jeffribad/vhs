package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// VHS is the main struct that holds the state of the tape player.
type VHS struct {
	Options    *Options
	Page       *rod.Page
	Browser    *rod.Browser
	Terminal   *Terminal
	Writer     *Writer
	Errors     []error
}

// Options holds the configuration options for VHS.
type Options struct {
	FontFamily   string
	FontSize     float64
	LineHeight   float64
	LetterSpacing float64
	Width        int
	Height       int
	Padding      string
	Framerate    float64
	PlaybackSpeed float64
	Theme        Theme
	Output       Output
}

// Output holds the output file paths for the different output formats.
type Output struct {
	GIF  string
	MP4  string
	WebM string
	Frames string
}

// DefaultOptions returns the default options for VHS.
func DefaultOptions() *Options {
	return &Options{
		FontFamily:    defaultFontFamily,
		FontSize:      defaultFontSize,
		LineHeight:    defaultLineHeight,
		LetterSpacing: defaultLetterSpacing,
		Width:         defaultWidth,
		Height:        defaultHeight,
		Padding:       defaultPadding,
		Framerate:     defaultFramerate,
		PlaybackSpeed: defaultPlaybackSpeed,
		Theme:         DefaultTheme,
	}
}

// New creates a new VHS instance with the given options.
func New(opts *Options) *VHS {
	v := &VHS{
		Options: opts,
		Errors:  []error{},
	}
	return v
}

// Setup initializes the browser and terminal for recording.
func (v *VHS) Setup() error {
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()

	v.Browser = rod.New().ControlURL(u).MustConnect()
	v.Page = v.Browser.MustPage("")

	return nil
}

// Cleanup tears down the browser and cleans up temporary files.
func (v *VHS) Cleanup() {
	if v.Browser != nil {
		v.Browser.MustClose()
	}
}

// Render processes the tape file and renders the output.
func (v *VHS) Render(tape string) error {
	cmds, err := Parse(tape)
	if err != nil {
		return fmt.Errorf("failed to parse tape: %w", err)
	}

	for _, cmd := range cmds {
		if err := cmd.Execute(v); err != nil {
			v.Errors = append(v.Errors, err)
		}
	}

	return nil
}

// resolveOutputPath resolves the output path relative to the tape file.
func resolveOutputPath(tapePath, outputPath string) string {
	if filepath.IsAbs(outputPath) {
		return outputPath
	}
	if strings.HasPrefix(outputPath, "./") || strings.HasPrefix(outputPath, "../") {
		return filepath.Join(filepath.Dir(tapePath), outputPath)
	}
	return outputPath
}

// ensureDir creates the directory for the given file path if it doesn't exist.
func ensureDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}
