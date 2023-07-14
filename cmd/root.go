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
	downloadOuiTxt bool
}{}

// コマンド
var rootCmd = &cobra.Command{
	Use:   "vendor6-cli",
	Short: "Start vendor6-cli",
	Long:  `Start vendor6-cli, an Interactive CLI tool to identify vendors by IPv6 address`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// DB初期化
		db, err := infra.NewOUIDb()
		if err != nil {
			return err
		}

		prompt := promptui.Prompt{
			Label:    ">",
			Validate: func(s string) error { return nil },
		}

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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vendor6-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
