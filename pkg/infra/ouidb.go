package infra

import (
	"bufio"
	"fmt"
	"io"
	"log"
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
	// OUI一覧
	OUI_LIST_URL = "https://standards-oui.ieee.org/oui/oui.txt"
)

var OUIFILEPATH = "./oui.txt"

func init() {
	// oui.txtがなければダウンロード
	_, err := os.Stat(OUIFILEPATH)
	if os.IsNotExist(err) {
		err = fetchAndSaveOuiList()
		if err != nil {
			log.Fatal(err)
		}
	}
}

type OUI struct {
	Code    string
	Company string
	Country string
}

type OUIDb struct {
	ouimp map[string]OUI
}

func getOuiList() (io.Reader, error) {
	var err error
	if _, err = os.Stat(OUIFILEPATH); !os.IsNotExist(err) {
		f, err := os.Open(OUIFILEPATH)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	return nil, err
}

func fetchOuiList() (io.ReadCloser, error) {
	resp, err := http.Get(OUI_LIST_URL)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func fetchAndSaveOuiList() error {
	fmt.Println("Downloading oui.txt...")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	rb, err := fetchOuiList()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(OUIFILEPATH), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(OUIFILEPATH)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, rb)
	if err != nil {
		return err
	}
	s.Stop()
	return nil
}

func NewOUIDb() (*OUIDb, error) {
	db := OUIDb{
		ouimp: make(map[string]OUI),
	}

	r, err := getOuiList()
	if os.IsNotExist(err) {
		err = fetchAndSaveOuiList()
		if err != nil {
			return nil, err
		}
		r, err = getOuiList()
		if err != nil {
			return nil, err
		}
	}

	pat_hex := regexp.MustCompile(`^([0-9A-F][0-9A-F]-[0-9A-F][0-9A-F]-[0-9A-F][0-9A-F]) +\(hex\)\t\t(.*)$`)
	pat_base16 := regexp.MustCompile(`^([0-9A-F]{6})     \(base 16\)\t\t(.*)$`)
	pat_country := regexp.MustCompile(`^\t\t\t\t[A-Z][A-Z]$`)
	pat_tab := regexp.MustCompile(`^\t\t\t\t`)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	// 最初の4行は捨てる
	for i := 0; i < 4; i++ {
		scanner.Scan()
	}

	newoui := OUI{}
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case pat_hex.MatchString(line):
		case pat_base16.MatchString(line):
			strs := pat_base16.FindStringSubmatch(line)
			newoui.Code = strings.ToLower(strs[1])
			newoui.Company = strs[2]
		case pat_country.MatchString(line):
			newoui.Country = strings.TrimSpace(line)
		case pat_tab.MatchString(line):
		case line == "":
			db.ouimp[newoui.Code] = newoui
			newoui = OUI{}
		}
	}

	return &db, nil
}

func (db *OUIDb) Lookup(mac net.HardwareAddr) *OUI {
	ret, ok := db.ouimp[strings.ReplaceAll(mac.String(), ":", "")[:6]]
	if !ok {
		return nil
	}
	return &ret
}
