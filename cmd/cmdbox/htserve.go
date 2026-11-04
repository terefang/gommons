package main

import (
	"flag"
	"log"
	"os"

	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xserve"
)

func init() {
	subcmd.Register(&SvcCommand{})
}

type SvcCommand struct {
	hostAddr string
	prefix   string
	rootDir  string
	headers  []string
}

func (r *SvcCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.hostAddr, "hostAddr", ":5555", "listen address and port")
	f.StringVar(&r.prefix, "prefix", "/other", "path prefix")
	f.StringVar(&r.rootDir, "root", ".", "use directory as root")
}

func (r SvcCommand) Info() (string, string) {
	return "serve", `runs the http server`
}

func (r *SvcCommand) Execute(args []string) int {
	_hsc := xserve.NewHttpServerContext()
	//basic init
	_hsc.Init()
	//set config
	_hsc.SetPathPrefix(r.prefix)
	var _fi, _err = os.Stat(r.rootDir)
	if (_err == nil) && _fi.IsDir() {
		log.Printf("Adding root from %s", r.rootDir)
		_hsc.SetRootPath(r.rootDir)
	}
	//set handlers
	_hsc.RegisterHandlers()
	//serve
	log.Printf("Starting services at %s", r.hostAddr)
	_hsc.RunAndServe(r.hostAddr)
	return 0
}
