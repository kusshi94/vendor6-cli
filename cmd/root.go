package cmd

import (
	"fmt"
	"os"

	"github.com/kusshi94/vendor6-cli/pkg/infra"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// フラグ引数用変数
var rootOpts = struct {
	ouiFilePath string
}{}

// コマンド
var rootCmd = &cobra.Command{
	Use:   "vendor6-cli",
	Short: "Start vendor6-cli",
	Long:  `Start vendor6-cli, an Interactive CLI tool to identify vendors by IPv6 address
Interactively entering an IPv6 address returns the vendor name for that address.
You can exit by typing "exit"
OUI information is downloaded from https://standards-oui.ieee.org/oui/oui.txt

IPv6アドレスを入力すると、そのアドレスのベンダー名を返します。
"exit"と入力すると終了します。
OUI情報は https://standards-oui.ieee.org/oui/oui.txt からダウンロードされます。
`,
	Example: `$ vendor6-cli
>: 2001:db8::0a00:7ff:fe12:3456
Apple, Inc.
>: 2001:db8::6666:b3ff:fe11:1111
TP-LINK TECHNOLOGIES CO.,LTD.
>: exit
$
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// DB初期化
		db, err := infra.NewOUIDb(rootOpts.ouiFilePath)
		if err != nil {
			return err
		}

		prompt := promptui.Prompt{
			Label:    ">",
			Validate: func(s string) error { return nil },
		}

		// welcome message
		fmt.Println("Enter an IPv6 address to get the vendor name. If you want to exit, type \"exit\".")

		// メインループ
		for {
			// データ入力
			input, err := prompt.Run()
			if err != nil {
				return err
			}

			// 終了
			if input == "exit" {
				return nil
			}

			// 入力を処理
			fmt.Println(ipToVendor(input, db))
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&rootOpts.ouiFilePath, "oui-file", "f", "./oui.txt", "OUI file path")
}
