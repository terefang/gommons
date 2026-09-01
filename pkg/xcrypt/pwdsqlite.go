package xcrypt

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/terefang/gommons/pkg/xstrings"
	_ "modernc.org/sqlite"
	"modernc.org/sqlite/vfs"
)

// ReadGroupsFromSqlite reads group membership information from the SQLite
// database at f.
//
// The database is opened read-only and is expected to contain an accountgroups
// table with groupname and username columns. Multiple rows for the same group
// are accumulated into a single group entry. Values containing multiple group
// // members may use the same separators as the text-based group files:
// // space, comma, semicolon, colon, slash, or pipe.
//
// The returned map is keyed by group name, with each value containing the
// usernames assigned to that group.
//
// An error is returned if the database cannot be opened or read, or if the
// expected table or columns are unavailable.
func ReadGroupsFromSqlite(f string) (map[string][]string, error) {
	var _ga = make(map[string][]string)
	_fdir := filepath.Dir(f)

	name, fsvfs, err := vfs.New(os.DirFS(_fdir))
	if err == nil {
		defer fsvfs.Close()
		_fname := filepath.Base(f)
		db, err := sql.Open("sqlite", "file:"+_fname+"?vfs="+name)
		if err == nil {
			defer db.Close()
			rows, err := db.Query("SELECT groupname,username FROM accountgroups;")
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var gname string
					var uname string
					err = rows.Scan(&gname, &uname)
					if err == nil {
						_ga[strings.ToUpper(gname)] = append(_ga[strings.ToUpper(gname)], xstrings.SplitByDefaultSet(uname)...)
					} else {
						return _ga, err
					}
				}
			}
			return _ga, err
		}
		return _ga, err
	}
	return _ga, err
}

// ReadAccountsFromSqlite reads account credentials and role assignments from
// the SQLite database at f.
//
// The database is opened read-only and is expected to contain the accounts and
// accountroles tables.
//
// The accounts table must contain username and password columns. Each username
// is mapped to its stored password hash in the returned accounts map.
//
// The accountroles table must contain username and rolename columns. Multiple
// role assignments for the same username are accumulated into a single string
// in the returned accountroles map. Role values may contain multiple roles
// separated by spaces, commas, semicolons, colons, slashes, or pipes, using the
// same format as the corresponding text-based role file.
//
// An error is returned if the database cannot be opened or read, or if the
// expected tables or columns are unavailable.
func ReadAccountsFromSqlite(f string) (map[string]string, map[string]string, error) {
	var _ac = make(map[string]string)
	var _ar = make(map[string]string)

	_fdir := filepath.Dir(f)

	name, fsvfs, err := vfs.New(os.DirFS(_fdir))
	if err == nil {
		defer fsvfs.Close()
		_fname := filepath.Base(f)
		db, err := sql.Open("sqlite", "file:"+_fname+"?vfs="+name)
		if err == nil {
			defer db.Close()
			rows, err := db.Query("SELECT username,password FROM accounts;")
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var uname string
					var pwd string
					err = rows.Scan(&uname, &pwd)
					if err == nil {
						_ac[strings.ToLower(uname)] = pwd
					} else {
						return _ac, _ar, err
					}
				}
				rows2, err := db.Query("SELECT username,rolename FROM accountroles;")
				if err == nil {
					defer rows2.Close()
					for rows2.Next() {
						var uname string
						var rname string
						err = rows2.Scan(&uname, &rname)
						if err == nil {
							_r, _ok := _ar[strings.ToLower(uname)]
							if !_ok {
								_ar[strings.ToLower(uname)] = strings.ToUpper(rname)
							} else {
								_ar[strings.ToLower(uname)] = _r + "," + strings.ToUpper(rname)
							}
						} else {
							return _ac, _ar, err
						}
					}
				}
			}
			return _ac, _ar, err
		}
		return _ac, _ar, err
	}
	return _ac, _ar, err
}
