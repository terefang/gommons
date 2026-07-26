package xdg

import (
	"os"
	"path/filepath"
	"strings"
)

func ExpandHome(path string) string {
	_home := UserHomeDir()
	if path[0] == '~' {
		return filepath.Join(_home, path[1:])
	}
	if strings.HasPrefix(path, "$HOME") {
		return filepath.Join(_home, path[5:])
	}
	if strings.HasPrefix(path, "$(HOME)") {
		return filepath.Join(_home, path[7:])
	}
	if strings.HasPrefix(path, "${HOME}") {
		return filepath.Join(_home, path[7:])
	}
	return path
}

func UserHomeDir() string {
	_home, _err := os.UserHomeDir()
	if _err != nil {
		_user := os.Getenv("USER")
		if _user == "" {
			return "/home/user"
		}
		return "/home/" + _user
	}
	return _home
}

func UserCacheDir() string {
	_home, _err := os.UserCacheDir()
	if _err != nil {
		return ExpandHome("~/.cache")
	}
	return _home
}

func UserConfigDir() string {
	_home, _err := os.UserConfigDir()
	if _err != nil {
		return ExpandHome("~/.config")
	}
	return _home
}

func UserDataDir() string {
	_path := os.Getenv(EnvDataDir)
	if _path == "" {
		return ExpandHome(PathDataDir)
	}
	return _path
}

// -------------------------------------------------------------------

func UserHomePath(path string) string {
	return filepath.Join(UserHomeDir(), path)
}

func UserCachePath(path string) string {
	return filepath.Join(UserCacheDir(), path)
}

func UserConfigPath(path string) string {
	return filepath.Join(UserConfigDir(), path)
}

func UserDataPath(path string) string {
	return filepath.Join(UserDataDir(), path)
}
