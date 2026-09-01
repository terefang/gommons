# CMDBOX EDIT

## NAME

cmdbox edit - lightweight vi-ish terminal text editor

## SYNOPSIS

**edit** [*FILE*]

## DESCRIPTION

**edit** is a Modal, Vi-inspired terminal text editor written in Go. 
It features modal editing, custom viewport scrolling, basic undo capabilities, 
command history, and custom ANSI color palettes.

If *FILE* is specified, **edit** loads the file into the editing buffer. 
If the file does not exist or no argument is given, it initializes an empty buffer.

## MODES

**edit** operates using four primary modes:

* **NORMAL**: Default mode for navigation, text manipulation, and executing modal commands.
* **INSERT**: Mode for entering text directly into the buffer.
* **COMMAND**: Mode for executing file-level operations and exiting the editor via the `:` prompt.
* **HELP**: Modal overlay listing available keybindings.

## KEYBINDINGS

### Navigation (NORMAL Mode)

* **h**, **Left Arrow**: Move cursor left.
* **j**, **Down Arrow**: Move cursor down.
* **k**, **Up Arrow**: Move cursor up.
* **l**, **Right Arrow**: Move cursor right.
* **0**, **Home**: Move cursor to the beginning of the line.
* **$**, **End**: Move cursor to the end of the line.
* **w**, **Shift+Right Arrow**: Advance forward to the start of the next word.
* **b**, **Shift+Left Arrow**: Move backward to the start of the previous word.
* **PageDown**, **Ctrl+f**: Scroll down by one page.
* **PageUp**, **Ctrl+b**: Scroll up by one page.

### Text Insertion (NORMAL -> INSERT Mode)

* **i**: Enter INSERT mode at current cursor position.
* **a**: Enter INSERT mode after current cursor position.
* **o**: Create a new line below the current row and enter INSERT mode.
* **O**: Create a new line above the current row and enter INSERT mode.

### Editing & Manipulation (NORMAL Mode)

* **x**, **Delete**: Delete the character under the cursor.
* **Backspace**: Delete the character before the cursor (or merge with the previous line if at column 0).
* **dd**: Delete (cut) current line to internal clipboard buffer.
* **yy**: Copy (yank) current line to internal clipboard buffer.
* **p**: Paste line from clipboard buffer below current line.
* **J**: Join current line with the next line, placing a single space between them.
* **u**: Undo the previous action (up to 50 history steps).

### Command & Help Modes

* **:**: Enter COMMAND mode.
* **?**, **F1**: Toggle HELP overlay.
* **Esc**: Return to NORMAL mode (from INSERT, COMMAND, or HELP modes).

### INSERT Mode Controls

* **Esc**: Exit to NORMAL mode.
* **Enter**: Split line at cursor position.
* **Tab**: Insert literal tab character.
* **Backspace** / **Delete**: Remove characters.
* **Arrow Keys** / **Home** / **End** / **PgUp** / **PgDn**: Navigate within the buffer.

## COMMAND MODE COMMANDS

Enter COMMAND mode by pressing **:** in NORMAL mode. Use **Up** and **Down** arrow 
keys to cycle through past command history.

* **:w** [*FILE*]
  Write buffer to disk. If *FILE* is supplied, writes to *FILE* (Save As).
* **:r** *FILE*
  Read file into buffer. Fails if unsaved changes exist in current buffer.
* **:r!** *FILE*
  Force-read file into buffer, discarding unsaved changes. If *FILE* is a directory, presents a TUI file picker selection.
* **:q**
  Quit editor. Fails if buffer contains unsaved changes (**E37**).
* **:q!**
  Quit editor without saving changes.
* **:wq**
  Write current file to disk and quit editor.

## EXIT STATUS

* **0**: Successful execution and clean exit.
* **1**: Program failure or initialization error.

## BUGS / LIMITATIONS

* Single line yank/paste operation only (`yy`, `dd`, `p`).
* Undo history depth is capped at 50 snapshot states.
* Tab characters are visually expanded to 4 spaces, but stored as raw tab characters (`\t`).