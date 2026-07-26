package xdg

import (
	"os"
	"path/filepath"
)

// XDG user directories environment variables.
const (
	PathDesktopDir     = "~/Desktop"
	PathDocumentsDir   = "~/Documents"
	PathDownloadDir    = "~/Downloads"
	PathTemplatesDir   = "~/Templates"
	PathPublicShareDir = "~/Public"
	PathProjectsDir    = "~/Projects"
	PathDataDir        = "~/.local/share"

	EnvDesktopDir     = "XDG_DESKTOP_DIR"
	EnvDocumentsDir   = "XDG_DOCUMENTS_DIR"
	EnvDownloadDir    = "XDG_DOWNLOAD_DIR"
	EnvTemplatesDir   = "XDG_TEMPLATES_DIR"
	EnvPublicShareDir = "XDG_PUBLICSHARE_DIR"
	EnvProjectsDir    = "XDG_PROJECTS_DIR"
	EnvDataDir        = "XDG_DATA_HOME"

	EnvRuntimeDir  = "XDG_RUNTIME_DIR"
	EnvMusicDir    = "XDG_MUSIC_DIR"
	EnvPicturesDir = "XDG_PICTURES_DIR"
	EnvVideosDir   = "XDG_VIDEOS_DIR"
)

func DocumentsDir() string {
	_path := os.Getenv(EnvDocumentsDir)
	if _path == "" {
		return ExpandHome(PathDocumentsDir)
	}
	return _path
}

func DesktopDir() string {
	_path := os.Getenv(EnvDesktopDir)
	if _path == "" {
		return ExpandHome(PathDesktopDir)
	}
	return _path
}

func DownloadDir() string {
	_path := os.Getenv(EnvDownloadDir)
	if _path == "" {
		return ExpandHome(PathDownloadDir)
	}
	return _path
}

func TemplatesDir() string {
	_path := os.Getenv(EnvTemplatesDir)
	if _path == "" {
		return ExpandHome(PathTemplatesDir)
	}
	return _path
}

func PublicDir() string {
	_path := os.Getenv(EnvPublicShareDir)
	if _path == "" {
		return ExpandHome(PathPublicShareDir)
	}
	return _path
}

func ProjectsDir() string {
	_path := os.Getenv(EnvProjectsDir)
	if _path == "" {
		return ExpandHome(PathProjectsDir)
	}
	return _path
}

// -------------------------------------------------------------------

func DocumentsPath(path string) string {
	return filepath.Join(DocumentsDir(), path)
}

func DesktopPath(path string) string {
	return filepath.Join(DesktopDir(), path)
}

func DownloadPath(path string) string {
	return filepath.Join(DownloadDir(), path)
}

func TemplatesPath(path string) string {
	return filepath.Join(TemplatesDir(), path)
}

func PublicPath(path string) string {
	return filepath.Join(PublicDir(), path)
}

func ProjectsPath(path string) string {
	return filepath.Join(ProjectsDir(), path)
}
