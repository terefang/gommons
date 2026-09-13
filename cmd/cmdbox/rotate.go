package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xcompress"
	"github.com/terefang/gommons/pkg/xfile"
)

func init() {
	subcmd.Register(&RotateCommand{NumFiles: 10})
}

// Option flags matching shtool rotate parameters
type RotateCommand struct {
	Verbose    bool
	Force      bool
	NumFiles   int
	SizeStr    string
	Copy       bool
	Remove     bool
	ArchiveDir string
	Compress   bool
	Background bool
	Delay      bool
	Pad        int
	Owner      string
	Group      string
	Mode       string
	// Evaluated properties
	SizeBytes int64
}

const msgPrefix = "rotate"

func (cfg *RotateCommand) Info() (string, string) {
	return "rotate", `rotate files with optional compression`
}

func (cfg *RotateCommand) Execute(files []string) int {

	if err := cfg.validateAndPrepareConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "%s:Error: %v\n", msgPrefix, err)
		return 1
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "%s:Error: no log files specified\n", msgPrefix)
		return 1
	}

	for _, file := range files {
		if err := cfg.processLogFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "%s:Error: %v\n", msgPrefix, err)
			return 1
		}
	}
	return 0
}

func (cfg *RotateCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&cfg.Verbose, "v", false, "Verbose output")
	f.BoolVar(&cfg.Verbose, "verbose", false, "Verbose output")
	f.BoolVar(&cfg.Force, "f", false, "Force archive dir creation / ignore missing files")
	f.BoolVar(&cfg.Force, "force", false, "Force archive dir creation / ignore missing files")

	f.IntVar(&cfg.NumFiles, "n", 10, "Number of rotated files to keep")
	f.IntVar(&cfg.NumFiles, "num-files", 10, "Number of rotated files to keep")

	f.StringVar(&cfg.SizeStr, "s", "", "Min size to trigger rotation")
	f.StringVar(&cfg.SizeStr, "size", "", "Min size to trigger rotation")

	f.BoolVar(&cfg.Copy, "c", false, "Copy and truncate instead of move")
	f.BoolVar(&cfg.Copy, "copy", false, "Copy and truncate instead of move")

	f.BoolVar(&cfg.Remove, "r", false, "Remove instead of truncating/re-creating")
	f.BoolVar(&cfg.Remove, "remove", false, "Remove instead of truncating/re-creating")

	f.StringVar(&cfg.ArchiveDir, "a", "", "Archive directory path")
	f.StringVar(&cfg.ArchiveDir, "archive-dir", "", "Archive directory path")

	f.BoolVar(&cfg.Compress, "z", false, "Compression")
	f.BoolVar(&cfg.Compress, "compress", false, "Compression")

	f.BoolVar(&cfg.Background, "b", false, "Background compression")
	f.BoolVar(&cfg.Background, "background", false, "Background compression")

	f.BoolVar(&cfg.Delay, "d", false, "Delay compression on log.0")
	f.BoolVar(&cfg.Delay, "delay", false, "Delay compression on log.0")

	f.IntVar(&cfg.Pad, "p", 1, "Number padding length")
	f.IntVar(&cfg.Pad, "pad", 1, "Number padding length")

	f.StringVar(&cfg.Owner, "o", "", "File owner")
	f.StringVar(&cfg.Owner, "owner", "", "File owner")

	f.StringVar(&cfg.Group, "g", "", "File group")
	f.StringVar(&cfg.Group, "group", "", "File group")

	f.StringVar(&cfg.Mode, "m", "", "File mode (e.g. 0755)")
	f.StringVar(&cfg.Mode, "mode", "", "File mode (e.g. 0755)")

}

func (cfg *RotateCommand) validateAndPrepareConfig() error {
	if cfg.NumFiles <= 0 {
		return fmt.Errorf("invalid argument `%d' to option -n", cfg.NumFiles)
	}

	cfg.SizeBytes = -1
	// Size parsing
	if cfg.SizeStr != "" {
		re := regexp.MustCompile(`^([0-9]+)([KkMmGg])?$`)
		matches := re.FindStringSubmatch(cfg.SizeStr)
		if len(matches) == 0 {
			return fmt.Errorf("invalid argument `%s' to option -s", cfg.SizeStr)
		}
		val, _ := strconv.ParseInt(matches[1], 10, 64)
		unit := strings.ToUpper(matches[2])
		switch unit {
		case "K":
			val *= 1024
		case "M":
			val *= 1024 * 1024
		case "G":
			val *= 1024 * 1024 * 1024
		}
		cfg.SizeBytes = val
	}

	// Consistency checks
	if cfg.Delay && !cfg.Compress {
		return fmt.Errorf("option -d requires option -z")
	}

	// Ensure target directory write permissions
	if cfg.ArchiveDir != "" {
		info, err := os.Stat(cfg.ArchiveDir)
		if err != nil {
			if os.IsNotExist(err) {
				if !cfg.Force {
					return fmt.Errorf("archive directory `%s' does not exist", cfg.ArchiveDir)
				}
				if err := os.MkdirAll(cfg.ArchiveDir, 0755); err != nil {
					return err
				}
			} else {
				return err
			}
		} else if !info.IsDir() {
			return fmt.Errorf("archive directory `%s' is not a directory", cfg.ArchiveDir)
		}
	}

	return nil
}

func (cfg *RotateCommand) processLogFile(targetPath string) error {
	fi, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) && cfg.Force {
			return nil
		}
		return fmt.Errorf("logfile `%s' not found", targetPath)
	}

	ldir := filepath.Dir(targetPath)
	file := filepath.Base(targetPath)

	adir := ldir
	if cfg.ArchiveDir != "" {
		if filepath.IsAbs(cfg.ArchiveDir) || strings.HasPrefix(cfg.ArchiveDir, "./") {
			adir = cfg.ArchiveDir
		} else {
			adir = filepath.Join(ldir, cfg.ArchiveDir)
		}
	}

	// Size threshold check
	if cfg.SizeBytes > 0 && fi.Size() < cfg.SizeBytes {
		if cfg.Verbose {
			fmt.Printf("%s: still too small in size -- skipping\n", filepath.Join(ldir, file))
		}
		return nil
	}

	if cfg.Verbose {
		fmt.Printf("rotating %s\n", filepath.Join(ldir, file))
	}

	// Calculate highest rotated filename index
	outNum := formatIndex(cfg.NumFiles-1, cfg.Pad)
	outPath := filepath.Join(adir, fmt.Sprintf("%s.%s", file, outNum))
	if cfg.Compress {
		outPath = fmt.Sprintf("%s%s", outPath, ".gz")
	}

	if _, err := os.Stat(outPath); err == nil {
		// Delete last out-rotated archive
		_ = os.Remove(outPath)
	}

	// Shift archive rotation sequence
	for i := cfg.NumFiles - 1; i > 0; i-- {
		m := formatIndex(i, cfg.Pad)
		n := formatIndex(i-1, cfg.Pad)

		if i-1 == 0 && cfg.Delay {
			uncompressed0 := filepath.Join(adir, fmt.Sprintf("%s.%s", file, n))
			if _, err := os.Stat(uncompressed0); err != nil {
				if cfg.Verbose {
					fmt.Printf("%s: error: %v\n", uncompressed0, err)
				}
				continue
			}

			targetCompressed := filepath.Join(adir, fmt.Sprintf("%s.%s%s", file, m, ".gz"))

			if cfg.Background {
				tempMoved := filepath.Join(adir, fmt.Sprintf("%s.%s", file, m))
				if err := os.Rename(uncompressed0, tempMoved); err != nil {
					return err
				}

				go func(src, dst string) {
					_ = xcompress.GzipCompressFile(dst, src)
					_ = os.Remove(src)
				}(tempMoved, targetCompressed)

			} else {
				if err := xcompress.GzipCompressFile(targetCompressed, uncompressed0); err != nil {
					return err
				}

				_ = os.Remove(uncompressed0)
			}

			_ = xfile.ApplyFilePermissionString(targetCompressed, cfg.Owner, cfg.Group, cfg.Mode)
			if cfg.Verbose {
				fmt.Println("rotated", uncompressed0, "to", targetCompressed)
			}

		} else {
			src := filepath.Join(adir, fmt.Sprintf("%s.%s", file, n))
			if cfg.Compress {
				src = src + ".gz"
			}
			dst := filepath.Join(adir, fmt.Sprintf("%s.%s", file, m))
			if cfg.Compress {
				dst = dst + ".gz"
			}

			if _, err := os.Stat(src); err != nil {
				continue
			}

			if err := os.Rename(src, dst); err != nil {
				return err
			}
			if cfg.Verbose {
				fmt.Println("rotated", src, "to", dst)
			}
		}
	}

	// Move / Copy log sequence
	targetIdx0 := filepath.Join(adir, fmt.Sprintf("%s.%s", file, formatIndex(0, cfg.Pad)))
	srcFile := filepath.Join(ldir, file)

	if cfg.Copy {
		if err := xfile.CopyFilePreserve(srcFile, targetIdx0); err != nil {
			return err
		}

		if !cfg.Remove {
			if err := os.Truncate(srcFile, 0); err != nil {
				return err
			}
		}
	} else {
		if err := os.Rename(srcFile, targetIdx0); err != nil {
			return err
		}

		if !cfg.Remove {
			f, err := os.OpenFile(srcFile, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			_ = f.Close()

			_ = xfile.ApplyFilePermissionString(srcFile, cfg.Owner, cfg.Group, cfg.Mode)
		}
	}

	// Immediate rotation compression
	if cfg.Compress && !cfg.Delay {
		targetCompressed := targetIdx0 + ".gz"

		if cfg.Background {
			go func(src, dst string) {
				_ = xcompress.GzipCompressFile(src, dst)
				_ = os.Remove(src)
			}(targetIdx0, targetCompressed)
		} else {
			if err := xcompress.GzipCompressFile(targetCompressed, targetIdx0); err != nil {
				return err
			}
			_ = os.Remove(targetIdx0)
		}

		_ = xfile.ApplyFilePermissionString(targetCompressed, cfg.Owner, cfg.Group, cfg.Mode)
		if cfg.Verbose {
			fmt.Println("rotated", targetIdx0, "to", targetCompressed)
		}
	}

	return nil
}

// Utility functions

func formatIndex(idx, pad int) string {
	return fmt.Sprintf(fmt.Sprintf("%%0%dd", pad), idx)
}
