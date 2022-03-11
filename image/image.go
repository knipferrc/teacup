// Package image provides an image bubble which can render
// images as strings.
package image

import (
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
)

type convertImageToStringMsg string
type errorMsg error

// Constants used throughout.
const (
	Padding = 1
)

// ToString converts an image to a string representation of an image.
func ToString(width int, img image.Image) string {
	img = imaging.Resize(img, width, 0, imaging.Lanczos)
	b := img.Bounds()
	imageWidth := b.Max.X
	h := b.Max.Y
	str := strings.Builder{}

	for heightCounter := 0; heightCounter < h; heightCounter += 2 {
		for x := imageWidth; x < width; x += 2 {
			str.WriteString(" ")
		}

		for x := 0; x < imageWidth; x++ {
			c1, _ := colorful.MakeColor(img.At(x, heightCounter))
			color1 := lipgloss.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, heightCounter+1))
			color2 := lipgloss.Color(c2.Hex())
			str.WriteString(lipgloss.NewStyle().Foreground(color1).
				Background(color2).Render("▀"))
		}

		str.WriteString("\n")
	}

	return str.String()
}

// convertImageToStringCmd redraws the image based on the width provided.
func convertImageToStringCmd(width int, filename string) tea.Cmd {
	return func() tea.Msg {
		imageContent, err := os.Open(filepath.Clean(filename))
		if err != nil {
			return errorMsg(err)
		}

		img, _, err := image.Decode(imageContent)
		if err != nil {
			return errorMsg(err)
		}

		imageString := ToString(width, img)

		return convertImageToStringMsg(imageString)
	}
}

// Bubble represents the properties of a code bubble.
type Bubble struct {
	Viewport    viewport.Model
	BorderColor lipgloss.AdaptiveColor
	Borderless  bool
	FileName    string
	ImageString string
}

// New creates a new instance of code.
func New(borderless bool, borderColor lipgloss.AdaptiveColor) Bubble {
	viewPort := viewport.New(0, 0)
	border := lipgloss.NormalBorder()

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	viewPort.Style = lipgloss.NewStyle().
		PaddingLeft(Padding).
		PaddingRight(Padding).
		Border(border).
		BorderForeground(borderColor)

	return Bubble{
		Viewport:    viewPort,
		Borderless:  borderless,
		BorderColor: borderColor,
	}
}

// Init initializes the code bubble.
func (b Bubble) Init() tea.Cmd {
	return nil
}

// SetFileName sets current file to highlight, this
// returns a cmd which will highlight the text.
func (b *Bubble) SetFileName(filename string) tea.Cmd {
	b.FileName = filename

	return convertImageToStringCmd(b.Viewport.Width, filename)
}

// SetBorderColor sets the current color of the border.
func (b *Bubble) SetBorderColor(color lipgloss.AdaptiveColor) {
	b.BorderColor = color
}

// SetSize sets the size of the bubble.
func (b *Bubble) SetSize(w, h int) {
	b.Viewport.Width = w - b.Viewport.Style.GetHorizontalFrameSize()
	b.Viewport.Height = h - b.Viewport.Style.GetVerticalFrameSize()

	b.Viewport.SetContent(lipgloss.NewStyle().
		Width(b.Viewport.Width).
		Height(b.Viewport.Height).
		Render(b.ImageString))
}

// Update handles updating the UI of a code bubble.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case convertImageToStringMsg:
		b.FileName = ""
		b.ImageString = lipgloss.NewStyle().
			Width(b.Viewport.Width).
			Height(b.Viewport.Height).
			Render(string(msg))

		b.Viewport.SetContent(b.ImageString)

		return b, nil
	case errorMsg:
		b.FileName = ""
		b.ImageString = lipgloss.NewStyle().
			Width(b.Viewport.Width).
			Height(b.Viewport.Height).
			Render("Error: " + msg.Error())

		b.Viewport.SetContent(b.ImageString)

		return b, nil
	}

	b.Viewport, cmd = b.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the code bubble.
func (b Bubble) View() string {
	border := lipgloss.NormalBorder()

	if b.Borderless {
		border = lipgloss.HiddenBorder()
	}

	b.Viewport.Style = lipgloss.NewStyle().
		PaddingLeft(Padding).
		PaddingRight(Padding).
		Border(border).
		BorderForeground(b.BorderColor)

	return b.Viewport.View()
}