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
	// OUI一覧
	OUI_LIST_URL = "https://standards-oui.ieee.org/oui/oui.txt"
)

type OUI struct {
	Code    string
	Company string
	Country string
}

type OUIDb struct {
	ouimp map[string]OUI
}

// oui.txtを開く
func openOuiTxt(ouiFilePath string) (io.ReadCloser, error) {
	var err error
	// oui.txtがなければダウンロードする
	if _, err = os.Stat(ouiFilePath); os.IsNotExist(err) {
		err := fetchAndSaveOuiTxt(ouiFilePath)
		if err != nil {
			return nil, err
		}
	}
	// oui.txtを開く
	f, err := os.Open(ouiFilePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// oui.txtをダウンロードして保存
func fetchAndSaveOuiTxt(ouiFilePath string) error {
	// UI表示
	fmt.Println("Downloading oui.txt...")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	// oui.txtをダウンロード
	resp, err := http.Get(OUI_LIST_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 保存
	// ディレクトリがなければ作成
	err = os.MkdirAll(filepath.Dir(ouiFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	// ファイル作成
	file, err := os.Create(ouiFilePath)
	if err != nil {
		return err
	}
	// ファイルへのデータ書き込み
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	s.Stop()

	return nil
}

// OUIDbを作成
func NewOUIDb(ouiFilePath string) (*OUIDb, error) {
	db := OUIDb{
		ouimp: make(map[string]OUI),
	}

	// oui.txtを開く
	r, err := openOuiTxt(ouiFilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// データを読み込む
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

// MACアドレスからOUIを検索
func (db *OUIDb) Lookup(mac net.HardwareAddr) *OUI {
	ret, ok := db.ouimp[strings.ReplaceAll(mac.String(), ":", "")[:6]]
	if !ok {
		return nil
	}
	return &ret
}
