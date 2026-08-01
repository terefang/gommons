package xlua

import (
	"database/sql"
	"errors"
	"fmt"
)

type sqlValue struct {
	db     *sql.DB
	tx     *sql.Tx
	driver string
	uri    string
	rows   *sql.Rows
	header []string
}

func (s *sqlValue) Begin() error {
	if s.tx == nil {
		_tx, _err := s.db.Begin()
		if _err != nil {
			return _err
		}
		s.tx = _tx
	} else {
		return errors.New("transaction already started")
	}
	return nil
}

func (s *sqlValue) Commit() error {
	if s.tx != nil {
		_err := s.tx.Commit()
		if _err != nil {
			return _err
		}
		s.tx = nil
	} else {
		return errors.New("transaction not started")
	}
	return nil
}

func (s *sqlValue) Rollback() error {
	if s.tx != nil {
		_err := s.tx.Rollback()
		if _err != nil {
			return _err
		}
		s.tx = nil
	} else {
		return errors.New("transaction not started")
	}
	return nil
}

func (s *sqlValue) Exec(query string, args ...any) (int64, error) {
	var _res sql.Result
	var _err error
	if s.tx != nil {
		_res, _err = s.tx.Exec(query, args...)
	} else {
		_res, _err = s.db.Exec(query, args...)
	}
	if _err != nil {
		return 0, _err
	}
	// TODO process _res result
	return _res.RowsAffected()
}

func (s *sqlValue) Query(query string, args ...any) error {
	var _res *sql.Rows
	var _err error
	if s.tx != nil {
		_res, _err = s.tx.Query(query, args...)
	} else {
		_res, _err = s.db.Query(query, args...)
	}
	if _err != nil {
		return _err
	}
	s.rows = _res
	s.header, _ = s.rows.Columns()
	// TODO process _res result
	return nil
}

func (s *sqlValue) Next() bool {
	if s.rows != nil {
		return s.rows.Next()
	}
	return false
}

func (s *sqlValue) Header() []string {
	if s.rows != nil {
		s.header, _ = s.rows.Columns()
		return s.header
	}
	return make([]string, 0)
}

func (s *sqlValue) Row() ([]any, error) {
	if s.rows != nil {
		_lh := len(s.header)
		vals := make([]any, _lh)
		ptrs := make([]any, _lh)
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := s.rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		for i, v := range vals {
			vals[i] = sqlConvertValue(v)
		}
		return vals, nil
	}
	return nil, errors.New("no rows available")
}

func (s *sqlValue) CloseQuery() error {
	if s.rows == nil {
		return errors.New("no rows available")
	}
	_err := s.rows.Close()
	if _err != nil {
		return _err
	}
	s.header = nil
	s.rows = nil
	return nil
}

func sqlConvertValue(v any) any {
	if v == nil {
		return "NULL"
	}
	switch val := v.(type) {
	case []byte:
		return string(val)
	case float64:
		return val
	case int64:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

// _db.Ping()
// _db.QueryRow()
func (s *sqlValue) Close() error {
	if s.tx != nil {
		s.tx.Rollback()
		s.tx = nil
	}
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			return err
		}
		s.driver = ""
		s.uri = ""
		return nil
	}
	return errors.New("no database available")
}
