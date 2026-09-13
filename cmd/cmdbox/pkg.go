package main

import (
	"archive/tar"
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xar"
	"github.com/terefang/gommons/pkg/xcompress"
	"github.com/terefang/gommons/pkg/xdg"
	"github.com/terefang/gommons/pkg/xfile"
	"github.com/terefang/gommons/pkg/xfmt"
	"github.com/terefang/gommons/pkg/xnss"
	"github.com/terefang/gommons/pkg/xstrings"
	"github.com/therootcompany/xz"
)

const OPKG_PKG_EXTENSION = ".opk"
const IPKG_PKG_EXTENSION = ".ipk"
const ZPKG_PKG_EXTENSION = ".pkz"
const DEB_EXTENSION = ".deb"
const TGZ_EXTENSION = ".tgz"
const TXZ_EXTENSION = ".txz"
const XPKG_PKG_INDEX = "Packages"
const JPKG_PKG_INDEX = "index.json"
const ZPKG_PKG_INDEX = "index.lst"

func init() {
	subcmd.Register(&PkgCommand{})
}

type PkgCommand struct {
	doInstall        bool
	doList           bool
	doSearch         bool
	doSearchProvides bool
	noDeps           bool
	dryRun           bool
	repoBase         string
	rootDir          string
	config           string
	installDb        string
	useIndexJson     bool
	usePackagesGz    bool
	useIndexList     bool
	useIndex         string
	noRepo           bool
}

func (r *PkgCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&r.doInstall, "i", false, "install packages")
	f.BoolVar(&r.doList, "l", false, "list package info")
	f.BoolVar(&r.doSearch, "s", true, "search for packages")
	f.BoolVar(&r.doSearchProvides, "p", false, "search for provides packages")
	f.BoolVar(&r.noRepo, "norepo", false, "dont update pkg index")
	f.BoolVar(&r.noDeps, "nodeps", false, "dont resolve and install dependencies")
	f.BoolVar(&r.dryRun, "dryrun", false, "dry-run installation")
	f.BoolVar(&r.useIndexJson, "json", false, "index.json repository")
	f.BoolVar(&r.useIndexList, "lst", false, "index.lst repository")
	f.BoolVar(&r.usePackagesGz, "opkg", false, "Packages.gz repository")
	f.StringVar(&r.useIndex, "index", "", "type repository")
	f.StringVar(&r.repoBase, "repo", "/pkg/var/pkg/packages", "repository base")
	f.StringVar(&r.rootDir, "root", "/pkg", "installation root directory")
	f.StringVar(&r.installDb, "idb", "/pkg/var/pkg/db", "installation db directory")
	f.StringVar(&r.config, "config", "/pkg/etc/pkg.repo", "repository configuration file")
}

func (r PkgCommand) Info() (string, string) {
	return "pkg", `Installs/Lists/Searches (i/o)pkg-style packages/repositories`
}

// Package Metadata structure matching Opkg/Debian control fields
type Package struct {
	Name         string
	Version      string
	Architecture string
	Section      string
	Filename     string
	Depends      []string
	Provides     []string
	Size         uint64
	Description  string
}

// PrettyPrint displays a formatted single-package detailed view
func (p Package) PrettyPrint() {
	deps := "None"
	if len(p.Depends) > 0 {
		deps = strings.Join(p.Depends, ", ")
	}
	prov := p.Name
	if len(p.Provides) > 0 {
		prov = strings.Join(p.Provides, ", ")
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Package      : %s\n", p.Name)
	fmt.Printf("Version      : %s\n", p.Version)
	fmt.Printf("Architecture : %s\n", p.Architecture)
	fmt.Printf("Section      : %s\n", p.Section)
	fmt.Printf("Filename     : %s\n", p.Filename)
	fmt.Printf("Size         : %s (%d bytes)\n", xfmt.ByteCountIEC(int64(p.Size)), p.Size)
	fmt.Printf("Depends      : %s\n", deps)
	fmt.Printf("Provides     : %s\n", prov)
	fmt.Printf("Description  : %s\n", p.Description)
	fmt.Println("--------------------------------------------------")
}

// OpkgManager handles repository operations and local installations
type OpkgManager struct {
	RootDir   string
	InstallDb string
	IndexName string
	RepoURL   string
	Packages  map[string]Package
	Provides  map[string][]string
	NoDeps    bool
	DryRun    bool
	Installed []string
	NoRepo    bool
}

func NewOpkgManager(rootDir, repoURL string) *OpkgManager {
	return &OpkgManager{
		RootDir:   rootDir,
		RepoURL:   repoURL,
		Packages:  make(map[string]Package),
		Provides:  make(map[string][]string),
		Installed: make([]string, 0),
	}
}

// Find matches packages by name, description, or wildcard pattern.
// Pattern supports standard shell wildcards (e.g., "lib*", "*zlib*", "*pdf").
func (m *OpkgManager) Find(pattern string) []Package {
	var matches []Package
	pattern = strings.ToLower(pattern)

	for _, pkg := range m.Packages {
		name := strings.ToLower(pkg.Name)
		matched, err := path.Match(pattern, name)
		if err == nil && matched {
			matches = append(matches, pkg)
		}
	}

	if len(matches) > 0 {
		return matches
	}

	// If no wildcard characters are supplied, wrap in wildcards for substring matching
	if !strings.ContainsAny(pattern, "*?[]") {
		pattern = "*" + pattern + "*"
	}

	for _, pkg := range m.Packages {
		name := strings.ToLower(pkg.Name)
		matched, err := path.Match(pattern, name)
		if err == nil && matched {
			matches = append(matches, pkg)
		}
	}

	return matches
}

func (m *OpkgManager) FindProvides(pattern string) map[string][]string {
	var matches map[string][]string = make(map[string][]string)
	pattern = strings.ToLower(pattern)

	// If no wildcard characters are supplied, wrap in wildcards for substring matching
	if !strings.ContainsAny(pattern, "*?[]") {
		pattern = "*" + pattern + "*"
	}

	for prov, pkg := range m.Provides {
		matched, err := path.Match(pattern, prov)
		if err == nil && matched {
			matches[prov] = append(matches[prov], pkg...)
		}
	}

	return matches
}

func (m *OpkgManager) Search(pattern string) {
	list := m.Find(pattern)
	PrettyPrintTable(list)
	llen := len(list)
	if llen == 0 {
		m.ListProvides([]string{pattern})
	} else {
		fmt.Printf("\n%d packages found.\n", llen)
	}
}

// PrettyPrintTable displays a list of packages in a column-aligned table
func PrettyPrintTable(pkgs []Package) {
	if len(pkgs) == 0 {
		fmt.Println("\nNo packages found.")
		return
	}

	// Dynamic column width calculation
	nameW, verW, archW, sectW := 9, 9, 5, 6 // 7, 7, 4, 5 // Minimal header lengths

	for _, p := range pkgs {
		if len(p.Name) > nameW {
			nameW = len(p.Name)
		}
		if len(p.Version) > verW {
			verW = len(p.Version)
		}
		if len(p.Architecture) > archW {
			archW = len(p.Architecture)
		}
		if len(p.Section) > archW {
			archW = len(p.Section)
		}
	}

	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%-%ds  %%10s\n", nameW, verW, archW, sectW)

	// Print Header
	fmt.Printf(format, "PACKAGE", "VERSION", "ARCH", "GROUP", "SIZE")
	fmt.Println(strings.Repeat("-", nameW+verW+archW+sectW+18))

	// Print Rows
	for _, p := range pkgs {
		fmt.Printf(format, p.Name, p.Version, p.Architecture, p.Section, xfmt.ByteCountIEC(int64(p.Size)))
	}
}

// Update fetches and parses the Packages index file from the repository
func (m *OpkgManager) Update() error {
	if m.IndexName == "" {
		return m.UpdateIndex(XPKG_PKG_INDEX)
	}
	return m.UpdateIndex(m.IndexName)
}

func (m *OpkgManager) UpdateIndex(indexName string) error {
	fmt.Println("Updating package database ... ")
	indexURL := fmt.Sprintf("%s/%s", m.RepoURL, indexName)
	if strings.HasPrefix(m.RepoURL, "deb+") {
		tempRepoUrl := strings.SplitN(m.RepoURL, "|", 2)
		indexURL = fmt.Sprintf("%s/%s/%s", tempRepoUrl[0][4:], tempRepoUrl[1], indexName)
	}
	fmt.Println(indexURL)

	var body []byte
	if strings.HasPrefix(indexURL, "https://") || strings.HasPrefix(indexURL, "http://") {

		resp, err := http.Get(indexURL)
		if err != nil {
			return fmt.Errorf("failed to fetch index: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.New(resp.Status)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading http content: %v", err)
		}

	} else if strings.HasPrefix(indexURL, "file://") {

		_url, err := url.Parse(indexURL)
		if err != nil {
			return fmt.Errorf("error parsing %s: %v", indexURL, err)
		}

		body, err = os.ReadFile(_url.Path)
		if err != nil {
			return fmt.Errorf("failed to fetch index: %w", err)
		}

	} else if strings.HasPrefix(indexURL, "/") || strings.HasPrefix(m.RepoURL, "./") {

		var err error
		body, err = os.ReadFile(indexURL)
		if err != nil {
			return fmt.Errorf("failed to fetch index: %w", err)
		}

	} else {
		return errors.New("invalid Repo URL " + m.RepoURL)
	}

	err := m.parsePackagesIndexContainer(indexName, body)
	if err != nil {
		return fmt.Errorf("error parsing %s: %v", indexName, err)
	}

	fmt.Printf("Updated index: %d packages available.\n", len(m.Packages))
	return nil
}

func (m *OpkgManager) parsePackagesIndexContainer(indexName string, data []byte) error {
	if strings.HasSuffix(indexName, ".gz") {
		_data, err := xcompress.GzipDecompressData(data)
		if err != nil {
			return fmt.Errorf("failed to decompress %s: %w", indexName, err)
		}
		data = _data
	} else if strings.HasSuffix(indexName, ".zst") {
		_data, err := xcompress.ZstdDecompressData(data)
		if err != nil {
			return fmt.Errorf("failed to decompress %s: %w", indexName, err)
		}
		data = _data
	}
	m.parsePackagesIndex(indexName, string(data))
	return nil
}

func renderProgressBar(prefix string, current, unit, total int, width int) {
	if unit == -1 || (current%unit) == 0 {
		percent := float64(current) / float64(total)
		filledLength := int(percent * float64(width))

		// Construct bar strings
		bar := strings.Repeat("=", filledLength) + strings.Repeat("-", width-filledLength)

		// Print using carriage return \r to overwrite the line
		fmt.Printf("\r%s [%s] %3.0f%% (%d/%d)", prefix, bar, percent*100, current, total)

		if unit > 0 {
			st := time.Second * 5
			st *= time.Duration(unit)
			st /= time.Duration(total)
			if st > time.Millisecond*500 {
				st = time.Millisecond * 500
			}
			time.Sleep(st)
		}
	}
}

func (m *OpkgManager) parsePackageEntry(block string) Package {
	var pkg Package
	lines := strings.Split(block, "\n")
	for _, line := range lines {
		//fmt.Println(line)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "Package":
			pkg.Name = val
		case "Version":
			pkg.Version = val
		case "Architecture":
			pkg.Architecture = val
		case "Section":
			pkg.Section = val
		case "Filename":
			pkg.Filename = val
		case "Description":
			pkg.Description = val
		case "Size":
			pkg.Size, _ = strconv.ParseUint(val, 10, 63)
		case "Depends":
			deps := strings.Split(val, ",")
			for _, dep := range deps {
				pkg.Depends = append(pkg.Depends, strings.TrimSpace(dep))
			}
		case "Provides":
			deps := strings.Split(val, ",")
			for _, dep := range deps {
				pkg.Provides = append(pkg.Provides, strings.TrimSpace(dep))
			}
		}
	}
	return pkg
}

func (m *OpkgManager) parsePackagesIndex(indexName string, data string) {
	if strings.HasPrefix(indexName, XPKG_PKG_INDEX) {
		m.parsePackages(data)
	} else if strings.HasPrefix(indexName, JPKG_PKG_INDEX) {
		m.parseIndexJson(data)
	} else if strings.HasPrefix(indexName, ZPKG_PKG_INDEX) {
		m.parseIndexList(data)
	}
	fmt.Println("")
}

// index.json:
// {
//	  "version": 2,
//	  "architecture": "x86_64",
//	  "packages": {
//	    "464xlat": "13",
//	    "6in4": "29",
//	    "6rd": "13",
//	    "6to4": "13",
//	    "adb": "5.0.2~6fe92d1a-r4",
//	    "adb-enablemodem": "2017.03.05-r1",
//	    "aeonsemi-as21xxx-firmware": "20260221-r1",
//	    "agetty": "2.41.5-r1"
//	  }
//	}
// <K>-<V>.<EXT>

type PackageJsonIndex struct {
	Version      int               `json:"version"`
	Architecture string            `json:"architecture"`
	Packages     map[string]string `json:"packages"`
}

func (m *OpkgManager) parseIndexJson(data string) {
	var index PackageJsonIndex
	err := json.Unmarshal([]byte(data), &index)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	blen := len(index.Packages)
	i := 0
	_section := filepath.Base(m.RepoURL)
	for pName, pVersion := range index.Packages {
		renderProgressBar("Indexing", i, 100, blen, 50)

		var pkg Package
		pkg.Name = pName
		pkg.Version = pVersion
		pkg.Architecture = index.Architecture
		pkg.Provides = []string{pName}
		pkg.Depends = make([]string, 0)
		pkg.Filename = fmt.Sprintf("%s-%s%s", pName, pVersion, ZPKG_PKG_EXTENSION)
		pkg.Section = _section
		xName := fmt.Sprintf("%s/%s", pkg.Name, pkg.Version)
		m.Packages[xName] = pkg
		if pkg.Provides != nil && len(pkg.Provides) > 0 {
			for _, prov := range pkg.Provides {
				m.Provides[prov] = append(m.Provides[prov], xName)
			}
		} else {
			m.Provides[pkg.Name] = append(m.Provides[pkg.Name], xName)
		}
		i++
	}
}

type PackageListIndex struct {
	Packages map[string]string
	Files    map[string]string
}

func (m *OpkgManager) parseIndexList(data string) {
	var index PackageListIndex
	index.Packages = make(map[string]string)
	index.Files = make(map[string]string)

	blocks := strings.Split(data, "\n")
	blen := len(blocks)
	fmt.Println("")
	for i, block := range blocks {
		renderProgressBar("Indexing", i, 100, blen, 50)

		if strings.TrimSpace(block) == "" {
			continue
		}

		// name version arch section size filename
		tokens := xstrings.SplitWsWithExtendedQuotes(block)
		if len(tokens) >= 4 {
			var pkg Package
			pkg.Name = tokens[0]
			pkg.Version = tokens[1]
			pkg.Architecture = tokens[2]
			pkg.Section = tokens[3]
			if len(tokens) >= 6 {
				pkg.Size, _ = strconv.ParseUint(tokens[4], 10, 64)
				pkg.Filename = tokens[5]
			} else {
				pkg.Filename = fmt.Sprintf("%s-%s%s", pkg.Name, pkg.Version, ZPKG_PKG_EXTENSION)
			}
			xName := fmt.Sprintf("%s/%s", pkg.Name, pkg.Version)
			m.Packages[xName] = pkg
			if pkg.Provides != nil && len(pkg.Provides) > 0 {
				for _, prov := range pkg.Provides {
					m.Provides[prov] = append(m.Provides[prov], xName)
				}
			} else {
				m.Provides[pkg.Name] = append(m.Provides[pkg.Name], xName)
			}
		}
	}
}

// Simple RFC 822 / Control file parser
func (m *OpkgManager) parsePackages(data string) {
	blocks := strings.Split(data, "\n\n")
	blen := len(blocks)
	fmt.Println("")
	for i, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		renderProgressBar("Indexing", i, 100, blen, 50)
		//fmt.Println("-----")

		var pkg Package = m.parsePackageEntry(block)
		if pkg.Name != "" {
			pName := fmt.Sprintf("%s/%s", pkg.Name, pkg.Version)
			m.Packages[pName] = pkg
			if pkg.Provides != nil && len(pkg.Provides) > 0 {
				for _, prov := range pkg.Provides {
					m.Provides[prov] = append(m.Provides[prov], pName)
				}
			} else {
				m.Provides[pkg.Name] = append(m.Provides[pkg.Name], pName)
			}
		}
	}
}

func (m *OpkgManager) List(args []string) {
	found := false
	count := 0
	for _, n := range args {
		list := m.Find(n)
		for _, pkg := range list {
			found = true
			pkg.PrettyPrint()
			count++
		}
	}
	fmt.Printf("%d packages found.\n", count)

	if !found {
		m.ListProvides(args)
	}
}

func (m *OpkgManager) ListProvides(args []string) {
	for _, n := range args {
		for prov, pkgs := range m.FindProvides(n) {
			fmt.Printf("\n%s provided by %s\n", prov, strings.Join(pkgs, ", "))
		}
	}
}

func (m *OpkgManager) IsInstalled(name string) bool {
	for _, n := range m.Installed {
		if n == name {
			return true
		}
	}
	fD := filepath.Clean(filepath.Join(m.InstallDb, name))
	fM := filepath.Clean(filepath.Join(fD, "MANIFEST"))
	return xfile.FileExists(fM)
}

func (m *OpkgManager) WriteManifest(manifest []string, pkgName string) error {
	if pkgName == "" {
		return errors.New("no package name provided")
	}
	fD := filepath.Clean(filepath.Join(m.InstallDb, pkgName))
	fM := filepath.Clean(filepath.Join(fD, "MANIFEST"))
	if err := xfile.CreatePathIfNotExists(fD); err != nil {
		return fmt.Errorf("error creating %s: %v", fD, err)
	}
	fmt.Printf("writing manifest %s\n", fM)
	return xfile.WriteStringToFile(fM, strings.Join(manifest, "\n"))
}

// Install resolves dependencies and unpacks the target package
func (m *OpkgManager) Install(pkgName string, asDep bool) error {
	if m.IsInstalled(pkgName) {
		if asDep {
			fmt.Printf("package %s already installed\n", pkgName)
			return nil
		}
		fmt.Printf("package %s already installed\n", pkgName)
		//return fmt.Errorf("package %s already installed", pkgName)
		return nil
	}

	pName, exists := m.Provides[pkgName]
	if !exists {
		return fmt.Errorf("package %s not found in feed", pkgName)
	}

	if len(pName) == 1 {
		pkg, exists := m.Packages[pName[0]]
		if !exists {
			return fmt.Errorf("package %s not found in feed", pkgName)
		}

		if (!m.NoDeps) && (!asDep) {
			// Resolve direct dependencies
			for _, dep := range pkg.Depends {
				if dep != "" {
					fmt.Printf("Resolving dependency: %s\n", dep)
					dep0 := strings.SplitN(dep, " ", 2)
					if !m.IsInstalled(dep0[0]) {
						provs, _ok := m.Provides[dep0[0]]
						if _ok {
							if len(provs) == 1 {
								provs0 := strings.SplitN(provs[0], "/", 2)
								if err := m.Install(provs0[0], true); err != nil {
									return fmt.Errorf("failed dependency %s: %w", provs0[0], err)
								}
							} else {
								return fmt.Errorf("package %s provided by multiple deps: %v", dep, strings.Join(provs, ", "))
							}
						} else if err := m.Install(dep0[0], true); err != nil {
							return fmt.Errorf("failed dependency %s: %w", dep, err)
						}
						m.Installed = append(m.Installed, dep0[0])
					}
				}
			}
		}
		m.Installed = append(m.Installed, pkg.Name)
		return m.InstallPkg(pkg)
	} else {
		return fmt.Errorf("package %s provided by multiple deps: %v", pkgName, strings.Join(pName, ", "))
	}
}

func (m *OpkgManager) InstallPkg(pkg Package) error {
	tName := filepath.Base(pkg.Filename)
	tmpIPK, err := os.CreateTemp(xdg.UserCacheDir(), "*."+tName)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", tmpIPK.Name(), err)
	}
	defer os.Remove(tmpIPK.Name())
	defer tmpIPK.Close()

	err = m.FetchToFile(tmpIPK, pkg)
	if err != nil {
		return fmt.Errorf("error fetching %s: %v", pkg.Name, err)
	}

	if strings.HasSuffix(pkg.Filename, TGZ_EXTENSION) {
		return m.extractTGZ(tmpIPK.Name(), pkg.Name)
	} else if strings.HasSuffix(pkg.Filename, TXZ_EXTENSION) {
		return m.extractTXZ(tmpIPK.Name(), pkg.Name)
	} else if strings.HasSuffix(pkg.Filename, ZPKG_PKG_EXTENSION) {
		return m.extractZPK(tmpIPK.Name(), pkg.Name)
	} else if strings.HasSuffix(pkg.Filename, DEB_EXTENSION) {
		return m.extractDeb(tmpIPK.Name(), pkg.Name)
	}
	return m.extractIPK(tmpIPK.Name(), pkg.Name)
}

func (m *OpkgManager) InstallUrl(url string, pkgName string) error {
	tName := filepath.Base(url)
	tmpIPK, err := os.CreateTemp(xdg.UserCacheDir(), "*."+tName)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", tmpIPK.Name(), err)
	}
	defer os.Remove(tmpIPK.Name())
	defer tmpIPK.Close()

	err = m.FetchUrlToFile(tmpIPK, filepath.Base(url), url)
	if err != nil {
		return fmt.Errorf("error fetching %s: %v", pkgName, err)
	}

	if strings.HasSuffix(tName, TGZ_EXTENSION) {
		return m.extractTGZ(tmpIPK.Name(), pkgName)
	} else if strings.HasSuffix(tName, TXZ_EXTENSION) {
		return m.extractTXZ(tmpIPK.Name(), pkgName)
	} else if strings.HasSuffix(url, ZPKG_PKG_EXTENSION) {
		return m.extractZPK(tmpIPK.Name(), pkgName)
	}
	return m.extractIPK(tmpIPK.Name(), pkgName)
}

func (m *OpkgManager) FetchToFile(tmpIPK *os.File, pkg Package) error {
	ipkURL := fmt.Sprintf("%s/%s", m.RepoURL, pkg.Filename)
	if strings.HasPrefix(m.RepoURL, "deb+") {
		tmpRepoUrl := strings.SplitN(m.RepoURL, "|", 2)
		ipkURL = fmt.Sprintf("%s/%s", tmpRepoUrl[0][4:], pkg.Filename)
	}
	return m.FetchUrlToFile(tmpIPK, pkg.Name, ipkURL)
}

func (m *OpkgManager) FetchUrlToFile(tmpIPK *os.File, pkgName, ipkURL string) error {
	if strings.HasPrefix(ipkURL, "https://") || strings.HasPrefix(ipkURL, "http://") {
		// Fetch .ipk file (ipk is an ar or tar.gz file containing data.tar.gz)
		fmt.Printf("Downloading %s from %s...\n", pkgName, ipkURL)
		resp, err := http.Get(ipkURL)
		if err != nil {
			return fmt.Errorf("error http-get %s: %v", pkgName, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.New(resp.Status)
		}

		count, err := io.Copy(tmpIPK, resp.Body)
		if err != nil {
			return fmt.Errorf("error reading http-body %s: %v", pkgName, err)
		}
		if count == 0 {
			return fmt.Errorf("error reading http-body %s: length 0", pkgName)
		}
	} else if strings.HasPrefix(ipkURL, "file://") {
		fmt.Printf("Fetching %s...\n", pkgName)
		_url, err := url.Parse(ipkURL)
		if err != nil {
			return fmt.Errorf("error get %s: %v", pkgName, err)
		}
		err = xfile.CopyToFile(tmpIPK, _url.Path)
		if err != nil {
			return fmt.Errorf("error reading body %s: %v", pkgName, err)
		}
	} else if strings.HasPrefix(ipkURL, "/") {
		fmt.Printf("Fetching %s...\n", pkgName)
		err := xfile.CopyToFile(tmpIPK, ipkURL)
		if err != nil {
			return fmt.Errorf("error get %s: %v", pkgName, err)
		}
	} else {
		return errors.New("invalid Repo URL " + m.RepoURL)
	}
	return nil
}

// Extract processes the .deb container to extract data.tar.gz into RootDir
func (m *OpkgManager) extractDeb(filePath string, pkgName string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", filePath, err)
	}
	defer f.Close()

	// Handle outer ar layer of .deb
	tarReader, err := xar.NewReader(f)
	if err != nil {
		return fmt.Errorf("error deb: %v", err)
	}

	okToInstall := false
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error ar: %v", err)
		}

		// we should see a package meta-data file
		// in advance of the installation contents
		// control.tar(.gz)
		if strings.HasPrefix(header.Name, "./control.tar") || strings.HasPrefix(header.Name, "control.tar") {
			okToInstall = true
			continue
		}

		// Locate the data payload data.tar(.gz)
		if strings.HasPrefix(header.Name, "./data.tar") || strings.HasPrefix(header.Name, "data.tar") {
			if !okToInstall {
				//return fmt.Errorf("package missing control data: %s", filePath)
				fmt.Printf("package missing control data: %s\n", filePath)
				continue
			}

			if strings.HasSuffix(header.Name, ".tar") {
				return m.extractArchive(tarReader, pkgName)
			} else if strings.HasSuffix(header.Name, ".gz") {
				return m.extractDataTarGz(tarReader, pkgName)
			} else if strings.HasSuffix(header.Name, ".bz2") {
				return m.extractDataTarBz2(tarReader, pkgName)
			} else if strings.HasSuffix(header.Name, ".xz") {
				return m.extractDataTarXz(tarReader, pkgName)
			} else if strings.HasSuffix(header.Name, ".zst") {
				return m.extractDataTarZstd(tarReader, pkgName)
			}
		}
	}

	return fmt.Errorf("data.tar(.gz) not found inside .deb package")
}

// Extract processes the .ipk container to extract data.tar.gz into RootDir
func (m *OpkgManager) extractIPK(filePath string, pkgName string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", filePath, err)
	}
	defer f.Close()

	// Handle outer gzip layer if IPK is formatted as tar.gz
	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to read outer container: %w", err)
	}
	defer gz.Close()

	tarReader := tar.NewReader(gz)

	okToInstall := false
	// .ipk from opkg/openwrt lack the control meta-data
	if pkgName != "" || strings.HasSuffix(filePath, IPKG_PKG_EXTENSION) {
		okToInstall = true
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error tar: %v", err)
		}

		// we should see a package meta-data file
		// in advance of the installation contents
		// control.tar(.gz) or CONTROL
		if strings.EqualFold(header.Name, "./CONTROL") || strings.EqualFold(header.Name, "CONTROL") {
			okToInstall = true
			if pkgName == "" {
				buf := make([]byte, header.Size)
				tarReader.Read(buf)
				pkg := m.parsePackageEntry(string(buf))
				if pkg.Name != "" {
					pkgName = pkg.Name
				}
				pkg.PrettyPrint()
			}
			continue
		} else if strings.HasPrefix(header.Name, "./control.tar") || strings.HasPrefix(header.Name, "control.tar") {
			okToInstall = true
			continue
		}

		// Locate the data payload data.tar(.gz)
		if strings.HasPrefix(header.Name, "./data.tar") || strings.HasPrefix(header.Name, "data.tar") {
			if !okToInstall {
				//return fmt.Errorf("package missing control data: %s", filePath)
				fmt.Printf("package missing control data: %s\n", filePath)
				continue
			}

			if strings.HasSuffix(header.Name, ".tar") {
				return m.extractArchive(tarReader, pkgName)
			} else if strings.HasSuffix(header.Name, ".gz") {
				return m.extractDataTarGz(tarReader, pkgName)
			} else if strings.HasSuffix(header.Name, ".bz2") {
				return m.extractDataTarBz2(tarReader, pkgName)
			} else if strings.HasSuffix(header.Name, ".xz") {
				return m.extractDataTarXz(tarReader, pkgName)
			} else if strings.HasSuffix(header.Name, ".zst") {
				return m.extractDataTarZstd(tarReader, pkgName)
			}
		}
	}

	return fmt.Errorf("data.tar(.gz) not found inside .ipk package")
}

// Extract processes the .pkz container to extract into RootDir
func (m *OpkgManager) extractZPK(filePath string, pkgName string) error {
	r, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", filePath, err)
	}
	defer r.Close()

	return m.extractDataTarZstd(r, pkgName)
}

func (m *OpkgManager) extractTGZ(filePath string, pkgName string) error {
	r, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", filePath, err)
	}
	defer r.Close()

	return m.extractDataTarGz(r, pkgName)
}

func (m *OpkgManager) extractTXZ(filePath string, pkgName string) error {
	r, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", filePath, err)
	}
	defer r.Close()

	return m.extractDataTarXz(r, pkgName)
}

// Extracts binary data files directly into the root filesystem destination
func (m *OpkgManager) extractDataTarGz(r io.Reader, pkgName string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("error gunzip %s: %v", pkgName, err)
	}
	defer gz.Close()

	return m.extractArchive(gz, pkgName)
}

func (m *OpkgManager) extractDataTarBz2(r io.Reader, pkgName string) error {
	gz := bzip2.NewReader(r)

	return m.extractArchive(gz, pkgName)
}

func (m *OpkgManager) extractDataTarXz(r io.Reader, pkgName string) error {
	gz, err := xz.NewReader(r, xz.DefaultDictMax)
	if err != nil {
		return fmt.Errorf("error unxz %s: %v", pkgName, err)
	}

	return m.extractArchive(gz, pkgName)
}

func (m *OpkgManager) extractDataTarZstd(r io.Reader, pkgName string) error {
	gz, err := zstd.NewReader(r)
	if err != nil {
		return fmt.Errorf("error unzstd %s: %v", pkgName, err)
	}
	defer gz.Close()

	return m.extractArchive(gz, pkgName)
}

func (m *OpkgManager) extractArchive(gz io.Reader, pkgName string) error {

	if m.DryRun {
		fmt.Println("dry-run mode, not installing anything!")
	}

	manifest := make([]string, 0)
	if !m.NoRepo && !m.DryRun {
		defer m.WriteManifest(manifest, pkgName)
	}

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error tar: %v", err)
		}

		targetPath := filepath.Join(m.RootDir, filepath.Clean(header.Name))
		targetPath = filepath.Clean(targetPath)

		hUname := header.Uname
		if hUname == "" {
			hUname = fmt.Sprintf("%d", header.Uid)
			_u, err := xnss.GetPwd(hUname)
			if err == nil {
				hUname = _u.Name
			}
		}
		hGname := header.Gname
		if hGname == "" {
			hGname = fmt.Sprintf("%d", header.Gid)
			_g, err := xnss.GetGroup(hUname)
			if err == nil {
				hGname = _g.Name
			}
		}

		if strings.EqualFold(header.Name, "./CONTROL") || strings.EqualFold(header.Name, "CONTROL") {
			fmt.Printf("! %-5s %-5s %s %s\n", hUname, hGname, os.FileMode(header.Mode).String(), header.Name)
			if pkgName == "" {
				buf := make([]byte, header.Size)
				tr.Read(buf)
				pkg := m.parsePackageEntry(string(buf))
				pkg.PrettyPrint()
				if pkg.Name != "" {
					pkgName = pkg.Name
				}
			}
			continue
		}

		switch header.Typeflag {
		case tar.TypeLink, tar.TypeSymlink:
			linkTargetPath := header.Linkname
			if strings.HasPrefix(linkTargetPath, "/") {
				linkTargetPath = filepath.Join(m.RootDir, linkTargetPath[1:])
				linkTargetPath, err = filepath.Rel(filepath.Dir(targetPath), linkTargetPath)
				//} else if strings.HasPrefix(linkTargetPath, "./") {
				//    linkTargetPath = filepath.Join(filepath.Dir(targetPath), filepath.Clean(linkTargetPath))
				//} else if strings.HasPrefix(linkTargetPath, "../") {
				//    linkTargetPath = filepath.Join(filepath.Dir(targetPath), linkTargetPath)
				//} else {
				//    linkTargetPath = filepath.Join(filepath.Dir(targetPath), linkTargetPath)
			}
			linkTargetPath = filepath.Clean(linkTargetPath)
			fmt.Printf("S %-5s %-5s %s %s -> %s\n", hUname, hGname, os.FileMode(header.Mode).String(), targetPath, linkTargetPath)
			//if strings.ContainsRune(linkTargetPath, '/') && !strings.HasPrefix(linkTargetPath, m.RootDir) {
			//    return fmt.Errorf("illegal package path %s", linkTargetPath)
			//}
			//fmt.Printf("*-----> %s\n", linkTargetPath)
			if !m.DryRun {
				manifest = append(manifest, fmt.Sprintf("S %s %s %s %s -> %s", hUname, hGname, os.FileMode(header.Mode).String(), targetPath, linkTargetPath))
				if err := os.Symlink(linkTargetPath, targetPath); err != nil {
					return fmt.Errorf("unable to symlink: %v", err)
				}
			}
		case tar.TypeDir:
			fmt.Printf("D %-5s %-5s %s %s\n", hUname, hGname, os.FileMode(header.Mode).String(), targetPath)
			if !m.DryRun {
				manifest = append(manifest, fmt.Sprintf("D %s %s %s %s", hUname, hGname, os.FileMode(header.Mode).String(), targetPath))
				if err := os.MkdirAll(targetPath, 0755); err != nil {
					return fmt.Errorf("unable to mkpath %s: %v", targetPath, err)
				}
				xfile.ApplyFilePermissions(targetPath, -1, -1, header.Mode)
				xfile.ApplyFilePermissionString(targetPath, hUname, hGname, "")
			}
		case tar.TypeReg:
			fmt.Printf("F %-5s %-5s %s %-20s %s\n", hUname, hGname, os.FileMode(header.Mode).String(), targetPath, xfmt.ByteCountIEC(header.Size))
			if !m.DryRun {
				manifest = append(manifest, fmt.Sprintf("F %s %s %s %s %d", hUname, hGname, os.FileMode(header.Mode).String(), targetPath, header.Size))
				dir := filepath.Dir(targetPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("unable to mkpath %s: %v", dir, err)
				}
				out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, header.FileInfo().Mode())
				if err != nil {
					return fmt.Errorf("unable to create %s: %v", targetPath, err)
				}
				if _, err := io.Copy(out, tr); err != nil {
					out.Close()
					return err
				}
				out.Close()
				xfile.ApplyFilePermissions(targetPath, -1, -1, header.Mode)
				xfile.ApplyFilePermissionString(targetPath, hUname, hGname, "")
			}
		case tar.TypeBlock:
			fmt.Printf("!B %-5s %-5s %s %s\n", hUname, hGname, os.FileMode(header.Mode).String(), targetPath)
		case tar.TypeChar:
			fmt.Printf("!C %-5s %-5s %s %s\n", hUname, hGname, os.FileMode(header.Mode).String(), targetPath)
		case tar.TypeFifo:
			fmt.Printf("!P %-5s %-5s %s %s\n", hUname, hGname, os.FileMode(header.Mode).String(), targetPath)
		default:
			fmt.Printf("!X %-5s %-5s %s %s\n", hUname, hGname, os.FileMode(header.Mode).String(), targetPath)
		}
	}
	return nil
}

func (r PkgCommand) Execute(args []string) int {
	// Initialize with a local installation target directory and repository endpoint
	pm := NewOpkgManager(r.rootDir, r.repoBase)
	pm.DryRun = r.dryRun
	pm.NoRepo = r.noRepo

	if !r.noRepo {
		if r.config != "" && xfile.FileExists(r.config) {
			r.readConfig(pm)
		}

		pm.InstallDb = r.installDb

		if r.useIndexJson {
			pm.IndexName = JPKG_PKG_INDEX
		}

		if r.useIndexList {
			pm.IndexName = ZPKG_PKG_INDEX
		}

		if r.usePackagesGz {
			pm.IndexName = XPKG_PKG_INDEX + ".gz"
		}

		if err := pm.Update(); err != nil {
			fmt.Printf("Update failed: %v\n", err)
			//return 1
		}
		fmt.Printf("Defined Packages: %d\n", len(pm.Packages))
	}

	if r.noDeps {
		pm.NoDeps = true
	}
	if r.doSearchProvides {
		if len(args) != 0 {
			pm.ListProvides(args)
		}
	} else if r.doList {
		if len(args) != 0 {
			pm.List(args)
		}
	} else if r.doInstall {
		for _, pkgToInstall := range args {
			fmt.Printf("\nInstalling package: %s ...\n", pkgToInstall)
			if (strings.HasPrefix(pkgToInstall, "https://") ||
				strings.HasPrefix(pkgToInstall, "http://") ||
				strings.HasPrefix(pkgToInstall, "file:///") ||
				strings.HasPrefix(pkgToInstall, "/") ||
				strings.HasPrefix(pkgToInstall, "./")) &&
				(strings.HasSuffix(pkgToInstall, OPKG_PKG_EXTENSION) ||
					strings.HasSuffix(pkgToInstall, IPKG_PKG_EXTENSION) ||
					strings.HasSuffix(pkgToInstall, DEB_EXTENSION) ||
					strings.HasSuffix(pkgToInstall, TGZ_EXTENSION) ||
					strings.HasSuffix(pkgToInstall, TXZ_EXTENSION) ||
					strings.HasSuffix(pkgToInstall, ZPKG_PKG_EXTENSION)) {

				if err := pm.InstallUrl(pkgToInstall, ""); err != nil {
					fmt.Printf("Failed: %v\n", err)
					return 1
				} else {
					fmt.Println("OK.")
				}

			} else {
				if err := pm.Install(pkgToInstall, false); err != nil {
					fmt.Printf("Failed: %v\n", err)
					return 1
				} else {
					fmt.Println("OK.")
				}
			}
		}
	} else if r.doSearch {
		if len(args) == 0 {
			pm.Search("*")
		} else {
			pm.Search(args[0])
		}
	}
	return 0
}

func (r *PkgCommand) readConfig(pm *OpkgManager) {
	f, err := os.Open(r.config)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		tokens := xstrings.SplitWsWithExtendedQuotes(s)
		if len(tokens) == 2 {
			switch tokens[0] {
			case "Index":
				r.useIndexJson = false
				r.usePackagesGz = false
				pm.IndexName = tokens[1]
			case "Repository":
				pm.RepoURL = tokens[1]
			case "Root":
				pm.RootDir = tokens[1]
			case "Manifests":
				pm.InstallDb = tokens[1]
			}
		}
	}
}
