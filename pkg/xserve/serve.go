package xserve

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

type HttpServerContext struct {
	Context    context.Context
	Muxer      *http.ServeMux
	Svc        *http.Server
	PathPrefix string
	RootPath   string
	Header     []string
	HttpFs     fs.FS
}

func (s *HttpServerContext) SetPathPrefix(n string) {
	s.PathPrefix = n
}

func NewHttpServerContext() *HttpServerContext {
	c := &HttpServerContext{}
	return c
}

func (s *HttpServerContext) LogS(p ...string) {
	log.Println(p)
}

func (s *HttpServerContext) LogR(r *http.Request, p string) {
	log.Printf("%s %s - %s", r.Method, r.URL, p)
}

func (s *HttpServerContext) LogFmtR(r *http.Request, _fmt string, _args ...interface{}) {
	log.Printf("%s %s - "+_fmt, r.Method, r.URL, _args)
}

func (s *HttpServerContext) Init() {
	//s.NotifyContext, _ = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s.Muxer = http.NewServeMux()
	s.SetPathPrefix("")
	s.Context = context.Background()
}

func (s *HttpServerContext) SetHandler(path string, f func(*HttpServerContext, http.ResponseWriter, *http.Request)) {
	s.Muxer.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		f(s, w, r)
	})
}

func (s *HttpServerContext) RunAndServe(hostAddr string) {
	s.Svc = &http.Server{
		Addr:    hostAddr,
		Handler: s.Muxer,
	}

	if err := s.Svc.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		s.LogS("server closed")
	} else {
		s.LogS(err.Error())
	}
}

func (s *HttpServerContext) RegisterHandlers() {
	s.HttpFs = os.DirFS(s.RootPath)

	// setup fileserver
	_fs := http.FileServerFS(s.HttpFs)
	s.Muxer.Handle(s.PathPrefix+"/", http.StripPrefix(s.PathPrefix+"/", neuter(_fs)))
}

func (s *HttpServerContext) SetRootPath(p string) {
	s.RootPath = p
}

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		fmt.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
