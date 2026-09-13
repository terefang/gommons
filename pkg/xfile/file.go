package xfile

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strconv"
	"syscall"

	"github.com/terefang/gommons/pkg/xnss"
)

func MakePath(path string, perm fs.FileMode) error {
	if FileExists(path) {
		return nil
	}
	return os.MkdirAll(path, perm)
}

func FileExists(filepath string) bool {
	if _fi, err := os.Stat(filepath); err == nil && (!_fi.IsDir()) && (!os.IsNotExist(err)) {
		return true
	}
	return false
}

func PathExists(filepath string) bool {
	if _fi, err := os.Stat(filepath); err == nil && _fi.IsDir() && (!os.IsNotExist(err)) {
		return true
	}
	return false
}

func CreatePath(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func CreatePathIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return CreatePath(dir)
	}
	return errors.New("dir already exists")
}

func IsEmpty(dir string) (bool, error) {
	es, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(es) == 0, nil
}

func DeletePath(dir string) error {
	return os.RemoveAll(dir)
}

func DeletePathIfExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return errors.New("dir does not exist")
	}
	return DeletePath(dir)
}

func DeletePathIfEmpty(dir string) error {
	if b, err := IsEmpty(dir); err != nil {
		return err
	} else if !b {
		return errors.New("dir is not empty")
	}
	return DeletePath(dir)
}

func IsRegularFile(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return fi.Mode().IsRegular()
}

func CopyFilePreserve(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if fi, err := in.Stat(); err == nil {
		_ = out.Chmod(fi.Mode())
	}
	return nil
}

func ApplyFilePermissions(target string, cOwner, cGroup int, cMode int64) error {
	if cOwner != -1 || cGroup != -1 {
		uid, gid := -1, -1

		// Fetch current owner details if updating only one attribute
		if fi, err := os.Stat(target); err == nil {
			if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
				uid = int(stat.Uid)
				gid = int(stat.Gid)
			}
		}

		if cOwner != -1 {
			uid = cOwner
		}
		if cGroup != -1 {
			gid = cGroup
		}

		err := os.Chown(target, uid, gid)
		if err != nil {
			return err
		}
	}

	if cMode != -1 {
		err := os.Chmod(target, os.FileMode(cMode))
		if err != nil {
			return err
		}
	}

	return nil
}

func ApplyFilePermissionString(target string, cOwner, cGroup, cMode string) error {
	nOwner := int(-1)
	nGroup := int(-1)
	nMode := int64(-1)
	if cOwner != "" {
		pOwner, err := xnss.GetPwd(cOwner)
		if err == nil {
			nOwner = pOwner.Uid
		} else {
			nOwner, err = strconv.Atoi(cOwner)
			if err != nil {
				return err
			}
		}
	}
	if cGroup != "" {
		pGroup, err := xnss.GetGroup(cGroup)
		if err == nil {
			nGroup = pGroup.Gid
		} else {
			nGroup, err = strconv.Atoi(cGroup)
			if err != nil {
				return err
			}
		}
	}
	if cMode != "" {
		n, err := strconv.ParseInt(cMode, 8, 32)
		if err != nil {
			return err
		}
		nMode = n
	}
	return ApplyFilePermissions(target, nOwner, nGroup, nMode)
}
