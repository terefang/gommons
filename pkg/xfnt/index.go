package xfnt

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/terefang/gommons/pkg/xstrings"
	"golang.org/x/image/font/sfnt"
)

type FontMetaInfo struct {
	Path       string
	RelPath    string
	Name       string
	Family     string
	SubFamily  string
	PostScript string
	Style      string
	Weight     int
	Stretch    int
}

// Print outputs the single font entry with tree-like formatting.
func (f FontMetaInfo) Print() {
	family := f.Family
	if family == "" {
		family = "Unknown Family"
	}

	style := f.Style
	if style == "" {
		style = "Normal"
	}

	stretch := f.Stretch
	if stretch == 0 {
		stretch = 100 // Fallback to normal width
	}

	fmt.Println(f.Path)
	fmt.Println(xstrings.Slugify(fmt.Sprintf("%s %s %d %d", f.Family, style, f.Weight, stretch)))
	fmt.Printf("p=`%s`\nn=`%s`\nf=`%s` s=`%s`\n", f.PostScript, f.Name, f.Family, f.SubFamily)
	fmt.Printf(" - Style: `%s`, Weight: `%d`, Stretch: `%d`\n\n", style, f.Weight, stretch)
}

type FontJob struct {
	List []*FontMetaInfo
	Base string
}

func ScanFontDirectory(searchDir string) (*FontJob, error) {
	job := new(FontJob)
	job.List = make([]*FontMetaInfo, 0)
	job.Base = searchDir
	err := filepath.WalkDir(searchDir, job.Run)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (fj *FontJob) Run(path string, d fs.DirEntry, err error) error {
	if err != nil || d.IsDir() {
		return nil
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".ttf" && ext != ".otf" {
		return nil
	}

	fnt, err := ExtractFontMetaInfo(path)
	if err == nil {
		fnt.RelPath, err = filepath.Rel(fj.Base, path)
		fnt.RelPath = "./" + fnt.RelPath
		fj.List = append(fj.List, fnt)
	}
	return nil
}

func ExtractFontMetaInfo(path string) (*FontMetaInfo, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	font, err := sfnt.Parse(data)
	if err == nil {
		var buffer sfnt.Buffer

		// Extract Family Name (ID 1)
		familyName, err := font.Name(&buffer, sfnt.NameIDFamily)
		if err != nil || familyName == "" {
			familyName = "Unknown Family"
		}

		// Extract Style (Subfamily ID 2)
		subfamily, _ := font.Name(&buffer, sfnt.NameIDSubfamily)
		style := parseStyle(subfamily)

		weight, stretch := parseOS2Table(data)

		fontName, err := font.Name(&buffer, sfnt.NameIDCompatibleFull)
		if err != nil || fontName == "" {
			fontName, err = font.Name(&buffer, sfnt.NameIDFull)
			if err != nil || fontName == "" {
				fontName = familyName + "/" + subfamily + "/" + style
			}
		}

		postName, _ := font.Name(&buffer, sfnt.NameIDPostScript)
		// Clean up file path separators to use forward slashes
		cleanPath := filepath.ToSlash(path)

		return &FontMetaInfo{
			Path:       cleanPath,
			Family:     familyName,
			SubFamily:  subfamily,
			PostScript: postName,
			Name:       fontName,
			Style:      style,
			Weight:     weight,
			Stretch:    stretch,
		}, nil
	}

	return nil, err
}

// parseOS2Table manually extracts usWeightClass and usWidthClass from raw TTF/OTF bytes
func parseOS2Table(data []byte) (weight int, stretch int) {
	weight = 400
	stretch = 1000

	if len(data) < 12 {
		return
	}

	// Read table directory entry count
	numTables := int(binary.BigEndian.Uint16(data[4:6]))
	offset := 12

	// Search directory entries (16 bytes each) for tag "OS/2"
	for i := 0; i < numTables; i++ {
		if offset+16 > len(data) {
			break
		}
		tag := string(data[offset : offset+4])
		tableOffset := int(binary.BigEndian.Uint32(data[offset+8 : offset+12]))

		if tag == "OS/2" {
			// OS/2 table offsets: usWeightClass is at +4 (uint16), usWidthClass is at +6 (uint16)
			if tableOffset+8 <= len(data) {
				w := binary.BigEndian.Uint16(data[tableOffset+4 : tableOffset+6])
				widthClass := binary.BigEndian.Uint16(data[tableOffset+6 : tableOffset+8])

				if w > 0 {
					weight = int(w)
				}
				stretch = mapWidthClass(widthClass)
			}
			break
		}
		offset += 16
	}

	return weight, stretch
}

func mapWidthClass(widthClass uint16) int {
	switch widthClass {
	case 1:
		return 500
	case 2:
		return 625
	case 3:
		return 750
	case 4:
		return 875
	case 5:
		return 1000
	case 6:
		return 1125
	case 7:
		return 1250
	case 8:
		return 1500
	case 9:
		return 2000
	default:
		return 1000
	}
}

func parseStyle(subfamily string) string {
	subLower := strings.ToLower(subfamily)
	if strings.Contains(subLower, "italic") {
		return "Italic"
	} else if strings.Contains(subLower, "oblique") {
		return "Oblique"
	}
	return "Normal"
}
