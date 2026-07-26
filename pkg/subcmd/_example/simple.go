package _example

import (
	"flag"
	"fmt"

	"github.com/terefang/gommons/pkg/subcmd"
)

// This is an example which demonstrates using the subcommand-library.
//
// Usage is something like this:
//
//  _example help
//  _example help run
//  _example help version
//  _example run 1 2 3
//  _example version
//

//
// We're going to write an application with two sub-commands:
//
//  version -> Show our version
//  run     -> "Do stuff"
//
// To do that we'll defined a pair of objects, and implement
// a series of methods on each one.
//

//
// RUN
//

// runCommand is a subcommand that takes a single optional flag.
type RunCommand struct {
	verbose bool
}

// Arguments adds per-command args to the object.
func (r RunCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&r.verbose, "verbose", false, "Should we be verbose")

}

// Info returns the name of this subcommand, along with a one-line
// description
func (r RunCommand) Info() (string, string) {
	return "run", "Runs some magic, and dumps its arguments."
}

// Execute is invoked if the user specifies `run` as the subcommand.
func (r RunCommand) Execute(args []string) int {

	fmt.Printf("I am a running application!\n")

	fmt.Printf("Verbose flag is %t\n", r.verbose)

	for i, s := range args {
		fmt.Printf("Argument %d is %s\n", i, s)
	}

	return 0
}

//
// VERSION
//

// versionCommand is a subcommand that takes no flags.
type VersionCommand struct {
	subcmd.NoFlags
}

// Name returns the name of this subcommand.
func (h VersionCommand) Info() (string, string) {
	return "version", "Show the application version."
}

// Execute is invoked if the user specifies `version` as the subcommand.
func (h VersionCommand) Execute(args []string) int {
	fmt.Printf("I am application version 1.0\n")
	return 0
}

//
// MAIN
//

// Register the subcommands, and run the one the user chose.
func main() {

	//
	// Register each of our subcommands.
	//
	subcmd.Register(&RunCommand{})
	subcmd.Register(&VersionCommand{})

	//
	// Execute the one the user chose.
	//
	subcmd.Execute()
}
