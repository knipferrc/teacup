// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/teacup/dirfs"
)

// Init initializes the filetree with files from the current directory.
func (b Bubble) Init() tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	cmd = getDirectoryListingCmd(dirfs.CurrentDirectory, true)
	cmds = append(cmds, cmd, textinput.Blink)

	return tea.Batch(cmds...)
}