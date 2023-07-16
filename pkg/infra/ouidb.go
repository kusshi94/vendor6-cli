package infra

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

const (
	// OUI_LIST_URL is the URL of the IEEE OUI list.
	OUI_LIST_URL = "https://standards-oui.ieee.org/oui/oui.txt"
)

var (
	// regular expressions for parsing oui.txt
	pat_hex = regexp.MustCompile(`^([0-9A-F][0-9A-F]-[0-9A-F][0-9A-F]-[0-9A-F][0-9A-F]) +\(hex\)\t\t(.*)$`)
	pat_base16 = regexp.MustCompile(`^([0-9A-F]{6})     \(base 16\)\t\t(.*)$`)
	pat_country = regexp.MustCompile(`^\t\t\t\t[A-Z][A-Z]$`)
	pat_tab = regexp.MustCompile(`^\t\t\t\t`)
)

// OUI represents a single entry in the OUI list.
type OUI struct {
	Code    string // OUI code (e.g. "38-9C-B2")
	Company string // Company name (e.g. "Apple, Inc.")
	Country string // Country code (e.g. "US")
	Address string // Address (e.g. "1 Infinite Loop Cupertino CA 95014 US")
}

// String returns a string representation of the OUI.
func (oui OUI) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s", oui.Code, oui.Company, oui.Country, oui.Address)
}

// OUI database
type OUIDb struct {
	ouimp map[string]OUI
}

// Create a new OUI database from the given file.
// If the file does not exist, it will be downloaded from the IEEE website.
func NewOUIDb(ouiFilePath string) (*OUIDb, error) {
	db := OUIDb{
		ouimp: make(map[string]OUI),
	}

	// open oui.txt
	// if the file does not exist, download it from the IEEE website
	r, err := openOuiTxt(ouiFilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// open scanner
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	// skip the first 4 lines
	for i := 0; i < 4; i++ {
		scanner.Scan()
	}

	// parse oui.txt line by line
	newoui := OUI{}
	for scanner.Scan() {
		// get a line
		line := scanner.Text()
		switch {
		case pat_hex.MatchString(line):
			// hex format
			// do nothing
		case pat_base16.MatchString(line):
			// base16 format
			// e.g. "389CB2     (base 16)\t\tApple, Inc."
			for i, s := range pat_base16.FindStringSubmatch(line) {
				if i == 1 {
					newoui.Code = strings.ToLower(s)
				} else if i == 2 {
					newoui.Company = s
				}
			}
		case pat_country.MatchString(line):
			// country code
			newoui.Country = strings.TrimSpace(line)
		case pat_tab.MatchString(line):
			// address
			if newoui.Address != "" {
				newoui.Address += " "
			}
			newoui.Address += strings.TrimSpace(line)
		case line == "":
			// blank line (end of entry)
			// add the entry to the database
			if newoui.Code != "" {
				db.ouimp[newoui.Code] = newoui
				newoui = OUI{}
			}
		}
	}

	return &db, nil
}

// Lookup the given MAC address in the database and if found, return the OUI.
func (db *OUIDb) Lookup(mac net.HardwareAddr) *OUI {
	ret, ok := db.ouimp[strings.ReplaceAll(mac.String(), ":", "")[:6]]
	if !ok {
		return nil
	}
	return &ret
}

// Open oui.txt and return the file handle.
// If the file does not exist, download it from the IEEE website.
func openOuiTxt(ouiFilePath string) (io.ReadCloser, error) {
	var err error
	// If oui.txt does not exist, download it from the IEEE website.
	if _, err = os.Stat(ouiFilePath); os.IsNotExist(err) {
		err := fetchAndSaveOuiTxt(ouiFilePath)
		if err != nil {
			return nil, err
		}
	}
	// open oui.txt
	f, err := os.Open(ouiFilePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Download oui.txt from the IEEE website and save it to the given file path.
func fetchAndSaveOuiTxt(ouiFilePath string) error {
	// Display a spinner while downloading
	fmt.Println("Downloading oui.txt...")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	// -- Download --
	resp, err := http.Get(OUI_LIST_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// -- Save to file --
	// If the directory does not exist, create it.
	err = os.MkdirAll(filepath.Dir(ouiFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	// Create a file
	file, err := os.Create(ouiFilePath)
	if err != nil {
		return err
	}
	// Save
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
