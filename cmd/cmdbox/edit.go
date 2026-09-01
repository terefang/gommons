package main

/*
================================================================================
  APPLICATION NOTICE
================================================================================

  This editor is primarily "vibe-coded" built in collaboration with an AI
  assistant (Gemini). While structural patterns, state transitions,
  and core logic follow clean software design principles (Model-View-Controller /
  Elm architecture via Charm's Bubble Tea framework), much of the feature set,
  layout choices, keybinding flows, and utility behaviors were organically
  iterated on the fly through real-time AI prompting and conversational code
  generation.

  LICENSE:
  Distributed under the MIT License.

  (c) 2026 - Created via AI Collaboration (Gemini) & the Open Source Community.
================================================================================
*/

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xstrings"
	"github.com/terefang/gommons/pkg/xtui/chooseui"
)

type EditCommand struct {
	subcmd.NoFlags
}

func init() {
	subcmd.Register(&EditCommand{})
}

func (r EditCommand) Info() (string, string) {
	return "edit", "edit a textfile vi-like"
}

func (r EditCommand) Execute(args []string) int {
	filename := ""
	if len(args) > 0 {
		filename = args[0]
	}
	m := initialModel(filename)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}
	return 0
}

const ANSI_Black_Dark = "0"
const ANSI_Red_Dark = "1"
const ANSI_Green_Dark = "2"
const ANSI_Yellow_Dark = "3"
const ANSI_Blue_Dark = "4"
const ANSI_Magenta_Dark = "5"
const ANSI_Cyan_Dark = "6"
const ANSI_White_Dark = "7"

const ANSI_Black = "8"
const ANSI_Red = "9"
const ANSI_Green = "10"
const ANSI_Yellow = "11"
const ANSI_Blue = "12"
const ANSI_Magenta = "13"
const ANSI_Cyan = "14"
const ANSI_White = "15"

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
	ModeHelp
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeCommand:
		return "COMMAND"
	case ModeHelp:
		return "HELP"
	default:
		return "UNKNOWN"
	}
}

type Position struct {
	Row int
	Col int
}

type model struct {
	mode           Mode
	lines          []string
	cursor         Position
	commandInput   string
	commandHistory []string
	cmdHistoryIdx  int
	errMsg         string
	warnMsg        string
	clipboard      string
	filename       string
	dirty          bool
	width, height  int
	pendingOp      string
	history        [][]string
	previousMode   Mode
	xOffset        int
	yOffset        int
	ready          bool
}

func initialModel(filename string) model {
	lines := []string{""}
	if filename != "" {
		data, err := os.ReadFile(filename)
		if err == nil {
			lines = strings.Split(string(data), "\n")
		}
	}
	if len(lines) == 0 {
		lines = []string{""}
	}

	m := model{
		mode:           ModeNormal,
		lines:          lines,
		cursor:         Position{Row: 0, Col: 0},
		filename:       filename,
		history:        make([][]string, 0),
		commandHistory: make([]string, 0),
	}
	m.saveHistory()
	return m
}

func (m *model) saveHistory() {
	snapshot := make([]string, len(m.lines))
	copy(snapshot, m.lines)
	m.history = append(m.history, snapshot)
	if len(m.history) > 50 {
		m.history = m.history[1:]
	}
}

func (m *model) undo() {
	if len(m.history) > 1 {
		m.history = m.history[:len(m.history)-1]
		last := m.history[len(m.history)-1]
		m.lines = make([]string, len(last))
		copy(m.lines, last)
		m.clampCursor()
	}
}

func (m *model) clampCursor() {
	if len(m.lines) == 0 {
		m.lines = []string{""}
	}

	if m.cursor.Row < 0 {
		m.cursor.Row = 0
	}
	if m.cursor.Row >= len(m.lines) {
		m.cursor.Row = len(m.lines) - 1
	}

	runes := []rune(m.lines[m.cursor.Row])
	runeLen := len(runes)

	if m.mode == ModeInsert {
		if m.cursor.Col > runeLen {
			m.cursor.Col = runeLen
		}
	} else {
		if runeLen == 0 {
			m.cursor.Col = 0
		} else if m.cursor.Col >= runeLen {
			m.cursor.Col = runeLen - 1
		}
	}

	if m.cursor.Col < 0 {
		m.cursor.Col = 0
	}

	m.adjustViewport()
}

func (m *model) adjustViewport() {
	if !m.ready {
		return
	}

	maxVisibleRows := m.viewportHeight()
	if maxVisibleRows < 1 {
		maxVisibleRows = 1
	}

	if m.cursor.Row < m.yOffset {
		m.yOffset = m.cursor.Row
	} else if m.cursor.Row >= m.yOffset+maxVisibleRows {
		m.yOffset = m.cursor.Row - maxVisibleRows + 1
	}

	if m.yOffset < 0 {
		m.yOffset = 0
	}

	maxVisibleCols := m.width - 5
	if maxVisibleCols < 1 {
		maxVisibleCols = 1
	}

	tabWidth := 4
	rawRunes := []rune(m.lines[m.cursor.Row])
	colIdx := m.cursor.Col
	if colIdx > len(rawRunes) {
		colIdx = len(rawRunes)
	}
	visualCol := len([]rune(expandTabs(string(rawRunes[:colIdx]), tabWidth)))

	if visualCol < m.xOffset {
		m.xOffset = visualCol
	} else if visualCol >= m.xOffset+maxVisibleCols {
		m.xOffset = visualCol - maxVisibleCols + 1
	}

	if m.xOffset < 0 {
		m.xOffset = 0
	}
}

func (m model) viewportHeight() int {
	h := m.height - 2
	if h < 1 {
		return 1
	}
	return h
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.adjustViewport()

	case tea.KeyMsg:
		m.errMsg = ""
		m.warnMsg = ""

		if msg.String() == "f1" {
			if m.mode == ModeHelp {
				m.mode = m.previousMode
			} else {
				m.previousMode = m.mode
				m.mode = ModeHelp
			}
			return m, nil
		}

		switch m.mode {
		case ModeNormal:
			return m.updateNormal(msg)
		case ModeInsert:
			return m.updateInsert(msg)
		case ModeCommand:
			return m.updateCommand(msg)
		case ModeHelp:
			return m.updateHelp(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) pageSize() int {
	return m.viewportHeight()
}

func (m model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "?":
		m.mode = m.previousMode
	}
	return m, nil
}

func (m model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.pendingOp != "" {
		if m.pendingOp == "d" && key == "d" {
			m.saveHistory()
			if m.cursor.Row >= 0 && m.cursor.Row < len(m.lines) {
				m.clipboard = m.lines[m.cursor.Row]
				if len(m.lines) > 1 {
					m.lines = append(m.lines[:m.cursor.Row], m.lines[m.cursor.Row+1:]...)
				} else {
					m.lines[0] = ""
				}
			}
			m.dirty = true
			m.clampCursor()
		} else if m.pendingOp == "y" && key == "y" {
			if m.cursor.Row >= 0 && m.cursor.Row < len(m.lines) {
				m.clipboard = m.lines[m.cursor.Row]
			}
		}
		m.pendingOp = ""
		return m, nil
	}

	switch key {
	case "?":
		m.previousMode = m.mode
		m.mode = ModeHelp
	case "i":
		m.mode = ModeInsert
	case "a":
		m.mode = ModeInsert
		m.cursor.Col++
		m.clampCursor()
	case "o":
		m.saveHistory()
		insertIdx := m.cursor.Row + 1
		if insertIdx > len(m.lines) {
			insertIdx = len(m.lines)
		}
		m.lines = append(m.lines[:insertIdx], append([]string{""}, m.lines[insertIdx:]...)...)
		m.cursor.Row = insertIdx
		m.cursor.Col = 0
		m.mode = ModeInsert
		m.dirty = true
	case "O":
		m.saveHistory()
		m.lines = append(m.lines[:m.cursor.Row], append([]string{""}, m.lines[m.cursor.Row:]...)...)
		m.cursor.Col = 0
		m.mode = ModeInsert
		m.dirty = true
	case "J":
		if m.cursor.Row < len(m.lines)-1 {
			m.saveHistory()
			currentLine := m.lines[m.cursor.Row]
			nextLine := strings.TrimLeft(m.lines[m.cursor.Row+1], " \t")

			joinCol := len([]rune(currentLine))
			if currentLine != "" && !strings.HasSuffix(currentLine, " ") && nextLine != "" {
				currentLine += " "
			}

			m.lines[m.cursor.Row] = currentLine + nextLine
			m.lines = append(m.lines[:m.cursor.Row+1], m.lines[m.cursor.Row+2:]...)
			m.cursor.Col = joinCol
			m.dirty = true
		}
	case ":":
		m.mode = ModeCommand
		m.commandInput = ""
		m.cmdHistoryIdx = len(m.commandHistory)

	case "pgdown", "ctrl+f":
		m.cursor.Row += m.pageSize()
	case "pgup", "ctrl+b":
		m.cursor.Row -= m.pageSize()

	case "h", "left":
		m.cursor.Col--
	case "j", "down":
		m.cursor.Row++
	case "k", "up":
		m.cursor.Row--
	case "l", "right":
		m.cursor.Col++
	case "0", "home":
		m.cursor.Col = 0
	case "$", "end":
		r := []rune(m.lines[m.cursor.Row])
		if len(r) > 0 {
			m.cursor.Col = len(r) - 1
		} else {
			m.cursor.Col = 0
		}
	case "w", "shift+right":
		runes := []rune(m.lines[m.cursor.Row])
		for m.cursor.Col < len(runes) && runes[m.cursor.Col] != ' ' {
			m.cursor.Col++
		}
		for m.cursor.Col < len(runes) && runes[m.cursor.Col] == ' ' {
			m.cursor.Col++
		}
	case "b", "shift+left":
		runes := []rune(m.lines[m.cursor.Row])
		if m.cursor.Col > 0 {
			m.cursor.Col--
			for m.cursor.Col > 0 && runes[m.cursor.Col] == ' ' {
				m.cursor.Col--
			}
			for m.cursor.Col > 0 && runes[m.cursor.Col-1] != ' ' {
				m.cursor.Col--
			}
		}
	case "x", "delete":
		runes := []rune(m.lines[m.cursor.Row])
		if len(runes) > 0 && m.cursor.Col < len(runes) {
			m.saveHistory()
			newRunes := append(runes[:m.cursor.Col], runes[m.cursor.Col+1:]...)
			m.lines[m.cursor.Row] = string(newRunes)
			m.dirty = true
		}
	case "backspace":
		runes := []rune(m.lines[m.cursor.Row])
		if m.cursor.Col > 0 {
			m.saveHistory()
			safeCol := m.cursor.Col
			if safeCol > len(runes) {
				safeCol = len(runes)
			}
			newRunes := append(runes[:safeCol-1], runes[safeCol:]...)
			m.lines[m.cursor.Row] = string(newRunes)
			m.cursor.Col--
			m.dirty = true
		} else if m.cursor.Row > 0 {
			m.saveHistory()
			prevRunes := []rune(m.lines[m.cursor.Row-1])
			prevLen := len(prevRunes)
			m.lines[m.cursor.Row-1] = string(append(prevRunes, runes...))
			m.lines = append(m.lines[:m.cursor.Row], m.lines[m.cursor.Row+1:]...)
			m.cursor.Row--
			m.cursor.Col = prevLen
			m.dirty = true
		}
	case "p":
		if m.clipboard != "" {
			m.saveHistory()
			insertIdx := m.cursor.Row + 1
			if insertIdx > len(m.lines) {
				insertIdx = len(m.lines)
			}
			m.lines = append(m.lines[:insertIdx], append([]string{m.clipboard}, m.lines[insertIdx:]...)...)
			m.cursor.Row = insertIdx
			m.dirty = true
		}
	case "u":
		m.undo()
	case "d", "y":
		m.pendingOp = key
	}

	m.clampCursor()
	return m, nil
}

func (m model) updateInsert(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	runes := []rune(m.lines[m.cursor.Row])

	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		m.cursor.Col--
		m.clampCursor()

	case "pgdown", "ctrl+f":
		m.cursor.Row += m.pageSize()
		m.clampCursor()
	case "pgup", "ctrl+b":
		m.cursor.Row -= m.pageSize()
		m.clampCursor()

	case "left":
		m.cursor.Col--
		m.clampCursor()
	case "right":
		m.cursor.Col++
		m.clampCursor()
	case "up":
		m.cursor.Row--
		m.clampCursor()
	case "down":
		m.cursor.Row++
		m.clampCursor()
	case "home":
		m.cursor.Col = 0
	case "end":
		m.cursor.Col = len(runes)

	case "tab":
		m.saveHistory()
		safeCol := m.cursor.Col
		if safeCol > len(runes) {
			safeCol = len(runes)
		}
		newRunes := append(runes[:safeCol], append([]rune{'\t'}, runes[safeCol:]...)...)
		m.lines[m.cursor.Row] = string(newRunes)
		m.cursor.Col++
		m.dirty = true

	case "backspace":
		if m.cursor.Col > 0 {
			m.saveHistory()
			safeCol := m.cursor.Col
			if safeCol > len(runes) {
				safeCol = len(runes)
			}
			newRunes := append(runes[:safeCol-1], runes[safeCol:]...)
			m.lines[m.cursor.Row] = string(newRunes)
			m.cursor.Col--
			m.dirty = true
		} else if m.cursor.Row > 0 {
			m.saveHistory()
			prevRunes := []rune(m.lines[m.cursor.Row-1])
			prevLen := len(prevRunes)
			m.lines[m.cursor.Row-1] = string(append(prevRunes, runes...))
			m.lines = append(m.lines[:m.cursor.Row], m.lines[m.cursor.Row+1:]...)
			m.cursor.Row--
			m.cursor.Col = prevLen
			m.dirty = true
		}
	case "delete":
		if len(runes) > 0 && m.cursor.Col < len(runes) {
			m.saveHistory()
			newRunes := append(runes[:m.cursor.Col], runes[m.cursor.Col+1:]...)
			m.lines[m.cursor.Row] = string(newRunes)
			m.dirty = true
		}
	case "enter":
		m.saveHistory()
		safeCol := m.cursor.Col
		if safeCol > len(runes) {
			safeCol = len(runes)
		}
		lineHead := string(runes[:safeCol])
		lineTail := string(runes[safeCol:])
		m.lines[m.cursor.Row] = lineHead
		m.lines = append(m.lines[:m.cursor.Row+1], append([]string{lineTail}, m.lines[m.cursor.Row+1:]...)...)
		m.cursor.Row++
		m.cursor.Col = 0
		m.dirty = true
	default:
		// Accept multi-byte runes or key strings representing single characters
		runesInput := []rune(msg.String())
		if len(runesInput) > 0 && (len(msg.String()) == 1 || utf8.RuneCountInString(msg.String()) == 1) {
			m.saveHistory()
			safeCol := m.cursor.Col
			if safeCol > len(runes) {
				safeCol = len(runes)
			}
			newRunes := append(runes[:safeCol], append(runesInput, runes[safeCol:]...)...)
			m.lines[m.cursor.Row] = string(newRunes)
			m.cursor.Col += len(runesInput)
			m.dirty = true
		} else {
			m.warnMsg = fmt.Sprintf("WARN: unhandled event `%s`", msg.String())
		}
	}
	m.clampCursor()
	return m, nil
}

func (m model) updateCommand(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal

	case "up":
		if len(m.commandHistory) > 0 && m.cmdHistoryIdx > 0 {
			m.cmdHistoryIdx--
			if m.cmdHistoryIdx >= 0 && m.cmdHistoryIdx < len(m.commandHistory) {
				m.commandInput = m.commandHistory[m.cmdHistoryIdx]
			}
		}
	case "down":
		if len(m.commandHistory) > 0 && m.cmdHistoryIdx < len(m.commandHistory)-1 {
			m.cmdHistoryIdx++
			if m.cmdHistoryIdx < len(m.commandHistory) {
				m.commandInput = m.commandHistory[m.cmdHistoryIdx]
			}
		} else {
			m.cmdHistoryIdx = len(m.commandHistory)
			m.commandInput = ""
		}

	case "enter":
		cmd := strings.TrimSpace(m.commandInput)
		if cmd != "" {
			m.commandHistory = append(m.commandHistory, cmd)
			_cq := xstrings.SplitWsDq(cmd)
			if len(_cq) > 0 {
				m.mode = ModeNormal
				switch _cq[0] {
				case "r":
					if m.dirty {
						m.errMsg = "E37: No write since last change (add ! to override)"
						return m, nil
					}
					fallthrough
				case "r!":
					if len(_cq) >= 2 {
						_info, err := os.Stat(_cq[1])
						if err != nil {
							m.errMsg = fmt.Sprintf("E212: Can't open file for reading: %v", err)
							return m, nil
						}
						if _info.IsDir() {
							_fn, err := chooseui.SelectFileFromPath(strings.TrimSpace(_cq[1]), false)
							if err != nil {
								m.errMsg = fmt.Sprintf("E212: Can't choose file: %v", err)
							}
							m.filename = _fn
						} else {
							m.filename = _cq[1]
						}
					} else {
						m.errMsg = "E32: No file/directory name"
						return m, nil
					}
					_buf, err := os.ReadFile(m.filename)
					if err != nil {
						m.errMsg = fmt.Sprintf("E212: Can't open file for reading: %v", err)
						return m, nil
					}
					m.saveHistory()
					m.lines = strings.Split(string(_buf), "\n")
					if len(m.lines) == 0 {
						m.lines = []string{""}
					}
					m.dirty = false
					m.clampCursor()
					return m, nil
				case "q":
					if m.dirty {
						m.errMsg = "E37: No write since last change (add ! to override)"
						return m, nil
					}
					return m, tea.Quit
				case "q!":
					return m, tea.Quit
				case "w":
					targetFile := m.filename
					if len(_cq) >= 2 {
						targetFile = strings.TrimSpace(_cq[1])
					}
					if targetFile == "" {
						m.errMsg = "E32: No file name"
						return m, nil
					}
					err := os.WriteFile(targetFile, []byte(strings.Join(m.lines, "\n")), 0644)
					if err != nil {
						m.errMsg = fmt.Sprintf("E212: Can't open file for writing: %v", err)
						return m, nil
					}
					m.filename = targetFile
					m.dirty = false
				case "wq":
					if m.filename == "" {
						m.errMsg = "E32: No file name"
						return m, nil
					}
					err := os.WriteFile(m.filename, []byte(strings.Join(m.lines, "\n")), 0644)
					if err != nil {
						m.errMsg = fmt.Sprintf("E212: Can't open file for writing: %v", err)
						return m, nil
					}
					return m, tea.Quit
				default:
					m.errMsg = fmt.Sprintf("E492: Not an editor command: %s", cmd)
				}
			}
		} else {
			m.mode = ModeNormal
		}
	case "backspace":
		if len(m.commandInput) > 0 {
			m.commandInput = m.commandInput[:len(m.commandInput)-1]
		} else {
			m.mode = ModeNormal
		}
	default:
		if len(msg.String()) == 1 {
			m.commandInput += msg.String()
		}
	}
	return m, nil
}

func expandTabs(s string, tabWidth int) string {
	if tabWidth < 1 {
		tabWidth = 4
	}
	var b strings.Builder
	column := 0
	for _, r := range s {
		if r == '\t' {
			spaces := tabWidth - (column % tabWidth)
			b.WriteString(strings.Repeat(" ", spaces))
			column += spaces
		} else {
			b.WriteRune(r)
			column++
		}
	}
	return b.String()
}

func (m model) renderBuffer() string {
	var b strings.Builder
	lineNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Width(4)
	cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color(ANSI_White)).Foreground(lipgloss.Color(ANSI_Black_Dark))

	maxLineLen := m.width - 5
	if maxLineLen < 1 {
		maxLineLen = 1
	}

	maxRows := m.viewportHeight()

	if m.yOffset < 0 {
		m.yOffset = 0
	}
	endRow := m.yOffset + maxRows
	if endRow > len(m.lines) {
		endRow = len(m.lines)
	}

	tabWidth := 4

	for i := m.yOffset; i < endRow; i++ {
		if i < 0 || i >= len(m.lines) {
			continue
		}
		rawLine := m.lines[i]
		rawRunes := []rune(rawLine)
		lineNum := lineNumStyle.Render(fmt.Sprintf("%3d ", i+1))
		b.WriteString(lineNum)

		line := expandTabs(rawLine, tabWidth)

		cursorVisualCol := -1
		if m.cursor.Row == i {
			safeCol := m.cursor.Col
			if safeCol < 0 {
				safeCol = 0
			}
			if safeCol > len(rawRunes) {
				safeCol = len(rawRunes)
			}
			cursorVisualCol = len([]rune(expandTabs(string(rawRunes[:safeCol]), tabWidth)))
		}
		runesline := []rune(line)
		runesLen := len(runesline)

		for j := m.xOffset; j < m.xOffset+maxLineLen; j++ {
			ch := " "
			if j >= 0 && j < runesLen {
				ch = string(runesline[j])
			}

			if m.cursor.Row == i && j == cursorVisualCol {
				b.WriteString(cursorStyle.Render(ch))
			} else {
				b.WriteString(ch)
			}
		}

		if i < endRow-1 {
			b.WriteString("\n")
		}
	}

	renderedRows := endRow - m.yOffset
	if renderedRows < 0 {
		renderedRows = 0
	}
	for i := renderedRows; i < maxRows; i++ {
		b.WriteString("\n~")
	}

	return b.String()
}

func (m model) View() string {
	if !m.ready {
		return "Initializing Editor..."
	}

	var b strings.Builder

	b.WriteString(m.renderBuffer())
	b.WriteString("\n")

	statusBg := ANSI_Blue_Dark
	switch m.mode {
	case ModeInsert:
		statusBg = ANSI_Green_Dark
	case ModeCommand:
		statusBg = ANSI_Red_Dark
	case ModeHelp:
		statusBg = ANSI_Magenta_Dark
	}

	modeStyle := lipgloss.NewStyle().Background(lipgloss.Color(statusBg)).Foreground(lipgloss.Color(ANSI_White)).Bold(true).Padding(0, 1)
	statusStyle := lipgloss.NewStyle().Background(lipgloss.Color(ANSI_Black_Dark)).Foreground(lipgloss.Color(ANSI_White)).Width(m.width)

	modified := ""
	if m.dirty {
		modified = " [+]"
	}

	fname := m.filename
	if fname == "" {
		fname = "[No Name]"
	}

	fileInfo := fmt.Sprintf(" %s%s | %d:%d | Press ? or F1 for Help", fname, modified, m.cursor.Row+1, m.cursor.Col+1)
	modeStr := modeStyle.Render(m.mode.String())

	b.WriteString(statusStyle.Render(modeStr + fileInfo))
	b.WriteString("\n")

	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ANSI_White)).Background(lipgloss.Color(ANSI_Red_Dark)).Width(m.width)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ANSI_Black_Dark)).Background(lipgloss.Color(ANSI_Yellow)).Width(m.width)
	if m.warnMsg != "" {
		b.WriteString(warnStyle.Render(m.warnMsg))
	} else if m.errMsg != "" {
		b.WriteString(errStyle.Render(m.errMsg))
	} else if m.mode == ModeCommand {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(ANSI_White)).Render(":" + m.commandInput))
	}

	if m.mode == ModeHelp {
		return m.renderHelpPopup(b.String())
	}

	return b.String()
}

func (m model) renderHelpPopup(backgroundView string) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ANSI_Yellow)).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ANSI_Cyan))

	helpContent := strings.Join([]string{
		titleStyle.Render("Editor Keybindings"),
		fmt.Sprintf("%s - %s", keyStyle.Render("h/j/k/l or Arrows"), descStyle.Render("Navigate cursor")),
		fmt.Sprintf("%s - %s", keyStyle.Render("PgUp / PgDn / C+f / C+b"), descStyle.Render("Page view scroll")),
		fmt.Sprintf("%s - %s", keyStyle.Render("0 / $"), descStyle.Render("Start / End of line")),
		fmt.Sprintf("%s - %s", keyStyle.Render("w / b"), descStyle.Render("Forward / Backward word")),
		fmt.Sprintf("%s - %s", keyStyle.Render("i / a"), descStyle.Render("Enter Insert mode (before/after)")),
		fmt.Sprintf("%s - %s", keyStyle.Render("o / O"), descStyle.Render("Insert line (below/above)")),
		fmt.Sprintf("%s - %s", keyStyle.Render("x"), descStyle.Render("Delete char")),
		fmt.Sprintf("%s - %s", keyStyle.Render("dd / yy / p"), descStyle.Render("Cut / Copy / Paste line")),
		fmt.Sprintf("%s - %s", keyStyle.Render("J"), descStyle.Render("Join lines")),
		fmt.Sprintf("%s - %s", keyStyle.Render("u"), descStyle.Render("Undo action")),
		fmt.Sprintf("%s - %s", keyStyle.Render(":w / :q / :wq"), descStyle.Render("Save / Quit / Save+Quit")),
		fmt.Sprintf("%s - %s", keyStyle.Render(":w / :r filepath"), descStyle.Render("SaveAs / Load")),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color(ANSI_White_Dark)).Render("Press Esc, q, or ? to close help"),
	}, "\n")

	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ANSI_Cyan_Dark)).
		Padding(1, 2).
		Background(lipgloss.Color(ANSI_Black_Dark))

	dialog := dialogBoxStyle.Render(helpContent)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
		lipgloss.WithWhitespaceChars(" "),
	)
}
