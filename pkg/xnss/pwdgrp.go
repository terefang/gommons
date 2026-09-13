package xnss

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Passwd = struct {
	Name   string
	Passwd string
	Uid    int
	Gid    int
	Gecos  string
	Dir    string
	Shell  string
}

type Group = struct {
	Name   string
	Passwd string
	Gid    int
	Mem    []string
}

// Cache structures
type pwdCacheEntry struct {
	modTime time.Time
	byName  map[string]*Passwd
	byUid   map[string]*Passwd
}

type grpCacheEntry struct {
	modTime time.Time
	byName  map[string]*Group
	byGid   map[string]*Group
}

var (
	pwdMu    sync.RWMutex
	pwdCache = make(map[string]*pwdCacheEntry)

	grpMu    sync.RWMutex
	grpCache = make(map[string]*grpCacheEntry)
)

func GetPwd(user string) (*Passwd, error) {
	return GetPwdFile(user, "/etc/passwd")
}

func GetPwdFile(user string, pwdfile string) (*Passwd, error) {
	fi, err := os.Stat(pwdfile)
	if err != nil {
		return nil, err
	}
	modTime := fi.ModTime()

	// 1. Try Read Lock (Cache Hit)
	pwdMu.RLock()
	entry, exists := pwdCache[pwdfile]
	if exists && entry.modTime.Equal(modTime) {
		p := findPwd(entry, user)
		pwdMu.RUnlock()
		if p != nil {
			return p, nil
		}
		return nil, errors.Errorf("%s - entry not found for: %s", pwdfile, user)
	}
	pwdMu.RUnlock()

	// 2. Cache Miss or File Modified -> Acquire Write Lock
	pwdMu.Lock()
	defer pwdMu.Unlock()

	// Double-check condition after acquiring write lock
	entry, exists = pwdCache[pwdfile]
	if !exists || !entry.modTime.Equal(modTime) {
		newEntry, err := parsePwdFile(pwdfile, modTime)
		if err != nil {
			return nil, err
		}
		pwdCache[pwdfile] = newEntry
		entry = newEntry
	}

	p := findPwd(entry, user)
	if p != nil {
		return p, nil
	}

	return nil, errors.Errorf("%s - entry not found for: %s", pwdfile, user)
}

func findPwd(entry *pwdCacheEntry, key string) *Passwd {
	lowerKey := strings.ToLower(key)
	if p, ok := entry.byName[lowerKey]; ok {
		return p
	}
	if p, ok := entry.byUid[key]; ok {
		return p
	}
	return nil
}

func parsePwdFile(pwdfile string, modTime time.Time) (*pwdCacheEntry, error) {
	f, err := os.Open(pwdfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entry := &pwdCacheEntry{
		modTime: modTime,
		byName:  make(map[string]*Passwd),
		byUid:   make(map[string]*Passwd),
	}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}

		a := strings.Split(s, ":")
		if len(a) < 7 {
			return nil, errors.Errorf("%s - invalid line: %s", pwdfile, s)
		}

		uid, err := strconv.Atoi(a[2])
		if err != nil {
			return nil, err
		}

		gid, err := strconv.Atoi(a[3])
		if err != nil {
			return nil, err
		}

		v := &Passwd{
			Name:   a[0],
			Passwd: a[1],
			Uid:    uid,
			Gid:    gid,
			Gecos:  a[4],
			Dir:    a[5],
			Shell:  a[6],
		}

		entry.byName[strings.ToLower(a[0])] = v
		entry.byUid[a[2]] = v
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return entry, nil
}

func GetGroup(group string) (*Group, error) {
	return GetGroupFile(group, "/etc/group")
}

func GetGroupFile(group string, grpfile string) (*Group, error) {
	fi, err := os.Stat(grpfile)
	if err != nil {
		return nil, err
	}
	modTime := fi.ModTime()

	// 1. Try Read Lock (Cache Hit)
	grpMu.RLock()
	entry, exists := grpCache[grpfile]
	if exists && entry.modTime.Equal(modTime) {
		g := findGroup(entry, group)
		grpMu.RUnlock()
		if g != nil {
			return g, nil
		}
		return nil, errors.Errorf("%s - entry not found for: %s", grpfile, group)
	}
	grpMu.RUnlock()

	// 2. Cache Miss or File Modified -> Acquire Write Lock
	grpMu.Lock()
	defer grpMu.Unlock()

	// Double-check condition after acquiring write lock
	entry, exists = grpCache[grpfile]
	if !exists || !entry.modTime.Equal(modTime) {
		newEntry, err := parseGroupFile(grpfile, modTime)
		if err != nil {
			return nil, err
		}
		grpCache[grpfile] = newEntry
		entry = newEntry
	}

	g := findGroup(entry, group)
	if g != nil {
		return g, nil
	}

	return nil, errors.Errorf("%s - entry not found for: %s", grpfile, group)
}

func findGroup(entry *grpCacheEntry, key string) *Group {
	lowerKey := strings.ToLower(key)
	if g, ok := entry.byName[lowerKey]; ok {
		return g
	}
	if g, ok := entry.byGid[key]; ok {
		return g
	}
	return nil
}

func parseGroupFile(grpfile string, modTime time.Time) (*grpCacheEntry, error) {
	f, err := os.Open(grpfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entry := &grpCacheEntry{
		modTime: modTime,
		byName:  make(map[string]*Group),
		byGid:   make(map[string]*Group),
	}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}

		a := strings.Split(s, ":")
		if len(a) < 4 {
			return nil, errors.Errorf("%s - invalid line: %s", grpfile, s)
		}

		gid, err := strconv.Atoi(a[2])
		if err != nil {
			return nil, err
		}

		var mem []string
		if a[3] != "" {
			mem = strings.Split(a[3], ",")
		}

		v := &Group{
			Name:   a[0],
			Passwd: a[1],
			Gid:    gid,
			Mem:    mem,
		}

		entry.byName[strings.ToLower(a[0])] = v
		entry.byGid[a[2]] = v
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return entry, nil
}
