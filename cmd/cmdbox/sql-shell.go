package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/andrewchambers/goline"
	// _ "github.com/duckdb/duckdb-go/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xstrings"
	_ "modernc.org/sqlite"
)

func init() {
	subcmd.Register(&SqlShellCommand{})
}

type SqlShellCommand struct {
	connect string
	DB      *sql.DB
}

func (r *SqlShellCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.connect, "connect", "", "connect to database url")
}

func (r SqlShellCommand) Info() (string, string) {
	return "sqlshell", `a simple sql shell similar to sqlite.`
}

func (r SqlShellCommand) Execute(args []string) int {
	return r.Repl()
}

const cmdSqlShellHistroyFile = ".cmd_sqlshell_history"

func (r *SqlShellCommand) Repl() int {
	goline.SetCompletionCallback(sqlShellReplComplete)
	goline.SetHintsCallback(sqlShellReplHint)
	fmt.Println(">>> entering repl mode, use '.q' to quit.")
	// if err := goline.HistoryLoad(xdg.UserHomePath(cmdSqlShellHistroyFile)); err != nil && !errors.Is(err, os.ErrNotExist) {
	//		fmt.Fprintf(os.Stderr, ">>> load history: %v\n", err)
	//	} else {
	//		fmt.Fprintf(os.Stderr, ">>> read history from %s\n", cmdSqlShellHistroyFile)
	//	}
	for {
		_line, _err := goline.ReadLine("sqlshell> ")
		if _err != nil {
			panic(_err)
		}

		_line = strings.TrimSpace(_line)
		goline.HistoryAdd(_line)
		if _line[0] == '.' {
			if len(_line) > 1 {
				_args := xstrings.SplitWsSq(_line)
				if _line[1] == 'q' {
					//goline.HistorySave(xdg.UserHomePath(cmdSqlShellHistroyFile))
					return 0
				} else if _line[1] == 'h' || _line[1] == '?' {
					r.ReplHelp()
					continue
				} else if _line[1] == 'o' {
					if r.DB != nil {
						fmt.Println("database open, close active first.")
						continue
					}
					_db, _err := sql.Open(_args[1], _args[2])
					if _err != nil {
						fmt.Println("Error opening database: ", _err)
						continue
					}
					r.DB = _db
				} else if _line[1] == 'c' {
					if r.DB == nil {
						fmt.Println("database not open.")
						continue
					}
					_err := r.DB.Close()
					if _err != nil {
						panic(_err)
					}
					r.DB = nil
				} else if _line[1] == 'd' {
					for _, driver := range sql.Drivers() {
						fmt.Println(driver)
					}
					continue
				} else {
					fmt.Println("unknown command:", _line)
				}
			} else {
				fmt.Println("unknown command:", _line)
			}
		} else {
			for !sqlShellStatementComplete(_line) {
				_nline, _err := goline.ReadLine("...> ")
				if _err != nil {
					panic(_err)
				}

				_nline = strings.TrimSpace(_nline)
				_line = strings.Join([]string{_line, _nline}, " ")
			}
			//_args := xstrings.SplitWsSq(_line)
			//fmt.Println(xstrings.StringifyNoErrorWithLevel(_args, -1))

			if r.DB == nil {
				fmt.Println("database not open")
				continue
			}
			_rows, _err := r.DB.Query(_line)
			if _err != nil {
				fmt.Println("Error query: ", _err)
				continue
			}
			r.PrintRows(_rows)
		}
	}
}

func (r *SqlShellCommand) ReplHelp() {
	fmt.Println(`Dot commands (single-line only):
  .help                 Show this message
  .open /driver/ /uri/  open database		    
  .close                close database		    
  .quit                 Exit the shell

SQL statements are terminated with a semicolon (;).
Multi-line statements are supported.`)
}

func (r *SqlShellCommand) PrintRows(rows *sql.Rows) {
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	var resultRows [][]string
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		row := make([]string, len(cols))
		for i, v := range vals {
			row[i] = sqlShellFormatValue(v)
		}
		resultRows = append(resultRows, row)
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if len(cols) > 0 {
		sqlShellPrintResult(cols, resultRows)
	}
}

func sqlShellFormatValue(v any) string {
	if v == nil {
		return "NULL"
	}
	switch val := v.(type) {
	case []byte:
		return string(val)
	case float64:
		// Trim trailing zeros for cleaner output.
		s := fmt.Sprintf("%g", val)
		return s
	case int64:
		s := fmt.Sprintf("%d", val)
		return s
	default:
		return fmt.Sprintf("%v", val)
	}
}

func sqlShellPrintResult(cols []string, rows [][]string) {
	w := os.Stdout
	if len(cols) == 0 {
		return
	}

	// Compute column widths: max of header and all cell values.
	widths := make([]int, len(cols))
	for i, c := range cols {
		widths[i] = utf8.RuneCountInString(c)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				widths[i] = max(widths[i], utf8.RuneCountInString(cell))
			}
		}
	}

	// Header.
	sqlShellPrintTableRow(cols, widths)

	// Separator.
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("-", w)
	}
	fmt.Fprintln(w, strings.Join(parts, "  "))

	// Data rows.
	for _, row := range rows {
		sqlShellPrintTableRow(row, widths)
	}
}

func sqlShellPrintTableRow(cells []string, widths []int) {
	w := os.Stdout
	parts := make([]string, len(widths))
	for i, width := range widths {
		var cell string
		if i < len(cells) {
			cell = cells[i]
		}
		pad := max(0, width-utf8.RuneCountInString(cell))
		parts[i] = cell + strings.Repeat(" ", pad)
	}
	fmt.Fprintln(w, strings.Join(parts, "  "))
}

func sqlShellReplComplete(line string) []string {
	commands := []string{
		".quit",
	}

	var matches []string
	for _, command := range commands {
		if strings.HasPrefix(command, line) {
			matches = append(matches, command)
		}
	}
	return matches
}

func sqlShellReplHint(line string) *goline.Hint {
	switch line {
	case ".quit":
		return &goline.Hint{Text: " exit", Color: 31}
	default:
		return nil
	}
}

// sqlShellStatementComplete returns true when buf contains at least one full SQL
// statement ending with ';' (ignoring quotes and comments at a basic level).
func sqlShellStatementComplete(buf string) bool {
	inSingle := false
	inDouble := false
	for i := 0; i < len(buf); i++ {
		c := buf[i]
		switch {
		case c == '\'' && !inDouble:
			// Check for escaped quote (doubled).
			if inSingle && i+1 < len(buf) && buf[i+1] == '\'' {
				i += 1
			} else {
				inSingle = !inSingle
			}
		case c == '"' && !inSingle:
			inDouble = !inDouble
		case c == ';' && !inSingle && !inDouble:
			return true
		}
	}
	return false
}
