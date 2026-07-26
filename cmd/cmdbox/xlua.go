package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrewchambers/goline"
	llib "github.com/arnodel/golua/lib"
	rt "github.com/arnodel/golua/runtime"
	"github.com/terefang/gommons/pkg/ctokener"
	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xlua"
)

func init() {
	subcmd.Register(&XluaCommand{})
}

type XluaCommand struct {
	script     string
	scriptFile string
	version    bool
	extlib     bool
	repl       bool
	dataFile   string
}

func (r *XluaCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.script, "e", "", "statement to execute")
	f.StringVar(&r.dataFile, "D", "", "data-file to load into '_data'")
	f.StringVar(&r.scriptFile, "f", "", "script-file to execute")
	f.BoolVar(&r.version, "v", false, "print version and exit")
	f.BoolVar(&r.repl, "i", false, "interactive (repl) mode")
	f.BoolVar(&r.extlib, "extlib", true, "register extlib")
}

func (r XluaCommand) Info() (string, string) {
	return "xlua", `extended lua environment`
}

func (r XluaCommand) Execute(args []string) int {

	if r.scriptFile != "" {
		_t, _err := os.ReadFile(r.scriptFile)
		if _err != nil {
			panic(_err)
		}
		r.script = string(_t)
	} else {
		r.scriptFile = "_main"
	}

	if r.script == "" && !r.repl && len(args) > 0 {
		// the first non-option arg is the scriptfile
		_t, _err := os.ReadFile(args[0])
		if _err != nil {
			panic(_err)
		}
		r.script = string(_t)
		r.scriptFile = args[0]
		args = args[1:]
	}

	// discard initial shebang and shell comments
	for len(r.script) > 0 && r.script[0] == '#' {
		_idx := strings.IndexByte(r.script, '\n')
		if _idx == -1 {
			break
		}
		r.script = r.script[_idx+1:]
	}

	_lua := rt.New(os.Stdout)
	llib.LoadAll(_lua)

	//xlua.SetOsArgs(_lua)
	xlua.SetArgs(_lua, args)
	if r.extlib {
		xlua.RegisterGlobalFunctions(_lua)
	}

	if r.version {
		_v := _lua.GlobalEnv().Get(rt.StringValue("_VERSION")).AsString()
		fmt.Println(_v)
		return 0
	}

	if r.dataFile != "" {
		if strings.HasSuffix(r.dataFile, ".ldata") || strings.HasSuffix(r.dataFile, ".json") {
			_d, _err := ctokener.LdataParseFile(r.dataFile)
			if _err != nil {
				panic(_err)
			}
			_lval, _err := xlua.ToLuaValueWithLevel(*_lua, _d, -1)
			if _err != nil {
				panic(_err)
			}
			_lua.GlobalEnv().Set(rt.AsValue("_data"), _lval)
		} else {
			panic("Data-file is unknown format: " + filepath.Base(r.dataFile))
		}
	}

	_ret := -1
	if r.script != "" {
		// Compile the chunk. Note that compiling doesn't require a runtime.
		chunk, _err := _lua.CompileAndLoadLuaChunk(r.scriptFile, []byte(r.script), rt.TableValue(_lua.GlobalEnv()))
		if _err != nil {
			panic(_err)
		}
		res, _err := rt.Call1(_lua.MainThread(), rt.FunctionValue(chunk))
		if _err != nil {
			panic(_err)
		}

		if res.IsNil() {
			_ret = 0
		}
	}

	if r.repl {
		r.Repl(_lua)
	}

	return _ret
}

const xluaHistroyFile = ".xlua_history"

func (r XluaCommand) Repl(_lua *rt.Runtime) int {
	goline.SetCompletionCallback(xluaReplComplete)
	goline.SetHintsCallback(xluaReplHint)
	fmt.Println(">>> entering repl mode, use '\\q' to quit.")
	if err := goline.HistoryLoad(xluaHistroyFile); err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, ">>> load history: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, ">>> read history from %s\n", xluaHistroyFile)
	}
	for {
		_line, _err := goline.ReadLine("xlua> ")
		if _err != nil {
			panic(_err)
		}

		if _line == "\\quit" || _line == "\\q" {
			goline.HistorySave(xluaHistroyFile)
			return 0
		}
		goline.HistoryAdd(_line)
		// Compile the chunk.
		chunk, _err := _lua.CompileAndLoadLuaChunk("_line", []byte(_line), rt.TableValue(_lua.GlobalEnv()))
		if _err != nil {
			fmt.Fprintf(os.Stderr, ">>> Parse-Error: %v\n", _err)
			continue
		}
		_, _err = rt.Call1(_lua.MainThread(), rt.FunctionValue(chunk))
		if _err != nil {
			fmt.Fprintf(os.Stderr, ">>> Run-Error: %v\n", _err)
			continue
		}
	}
}

func xluaReplComplete(line string) []string {
	commands := []string{
		"\\quit",
	}

	var matches []string
	for _, command := range commands {
		if strings.HasPrefix(command, line) {
			matches = append(matches, command)
		}
	}
	return matches
}

func xluaReplHint(line string) *goline.Hint {
	switch line {
	case "\\quit":
		return &goline.Hint{Text: " exit", Color: 31}
	default:
		return nil
	}
}
