package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/terefang/gommons/pkg/stemplate"
	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xtui/chooseui"
)

func init() {
	subcmd.Register(&chooseFileCommand{})
}

// Structure for our options and state.
type chooseFileCommand struct {

	// Command to execute
	exec     string
	template string
	// Filenames we'll let the user choose between
	files []string
}

// Arguments adds per-command args to the object.
func (cf *chooseFileCommand) Arguments(f *flag.FlagSet) {
	if cf != nil {
		f.StringVar(&cf.exec, "execute", "", "Command to execute once a selection has been made")
		f.StringVar(&cf.template, "template", "", "Template to format once a selection has been made")
	}
}

// Info returns the name of this subcommand.
func (cf *chooseFileCommand) Info() (string, string) {
	return "choose-file", `Choose a file, interactively.

Details:

This command presents a directory view, showing you all the files beneath
the named directory.  You can navigate with the keyboard, and press RETURN
to select a file.

Optionally you can press TAB to filter the list via an input field.

Uses:

This is ideal for choosing videos, roms, etc.  For example launch a
video file, interactively:

   $ xine "$(cmdbox choose-file ~/Videos)"
   $ cmdbox choose-file -execute="xine {}" ~/Videos
   $ cmdbox choose-file -template="xine {} && something else" ~/Videos |sh

See also 'cmdbox help choose-stdin'.`
}

// Execute is invoked if the user specifies `choose-file` as the subcommand.
func (cf *chooseFileCommand) Execute(args []string) int {

	//
	// Get our starting directory
	//
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	//
	// Find files
	//
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {

			// Null info?  That probably means that the
			// destination we're trying to walk doesn't exist.
			if info == nil {
				return nil
			}

			// We'll add anything that isn't a directory
			if !info.IsDir() {
				if !strings.Contains(path, "/.") && !strings.HasPrefix(path, ".") {
					cf.files = append(cf.files, path)
				}
			}
			return nil
		})

	if err != nil {
		fmt.Printf("error walking %s: %s\n", dir, err.Error())
		return 1
	}
	if len(cf.files) < 1 {
		fmt.Printf("Failed to find any files beneath %s\n", dir)
		return 1
	}

	//
	// Launch the UI
	//
	chooser := chooseui.New(cf.files)
	choice := chooser.Choose()

	//
	// Did something get chosen?  If not terminate
	//
	if choice == "" {

		return 1
	}

	//
	// Are we executing?
	//
	if cf.exec != "" {

		//
		// Split into command and arguments
		//
		run := stemplate.Expand(cf.exec, choice, "")

		//
		// Run it.
		//
		cmd := exec.Command(run[0], run[1:]...)
		out, errr := cmd.CombinedOutput()
		if errr != nil {
			fmt.Printf("Error running '%v': %s\n", run, errr.Error())
			return 1
		}

		//
		// And we're done
		//
		fmt.Printf("%s\n", out)
		return 0

	} else if cf.template != "" {

		out := stemplate.ExpandString(cf.template, choice)
		fmt.Printf("%s\n", out)
		return 0
	}

	//
	// We're not executing, so show the user's choice
	//
	fmt.Printf("%s\n", choice)
	return 0
}
