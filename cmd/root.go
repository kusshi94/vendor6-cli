package cmd

import (
	"fmt"
	"os"

	"github.com/kusshi94/vendor6-cli/pkg/infra"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// variables for command line options
var rootOpts = struct {
	ouiFilePath string
	printAllInfo bool
}{}

// main command
var rootCmd = &cobra.Command{
	Use:   "vendor6-cli [IPv6 address...]",
	Short: "Start vendor6-cli",
	Long:  `Start vendor6-cli, an Interactive CLI tool to identify vendors by IPv6 address

If IPv6 address is given as arguments, the vendor name for that address is returned.

If no arguments are given, the interactive mode is started.
In interactive mode, enter an IPv6 address to get the vendor name.
If you want to exit interactive mode, type "exit".

OUI information is automatically downloaded from https://standards-oui.ieee.org/oui/oui.txt to the current directory.
If you want to use a different file, use the -f option to specify the file path.`,
	DisableFlagsInUseLine: true,
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
		// Initialize OUI database
		db, err := infra.NewOUIDb(rootOpts.ouiFilePath)
		if err != nil {
			return err
		}

		// If arguments are given, process them and exit
		if len(args) > 0 {
			for _, arg := range args {
				fmt.Println(ipToVendor(arg, db, rootOpts.printAllInfo))
			}
			return nil
		}

		// -- Interactive mode --

		// Initialize prompt
		prompt := promptui.Prompt{
			Label:    ">",
			Validate: func(s string) error { return nil },
		}

		// Print usage message
		fmt.Println("Enter an IPv6 address to get the vendor name. If you want to exit, type \"exit\".")

		// Start interactive mode
		for {
			// Get input from user interactively
			input, err := prompt.Run()
			if err != nil {
				return err
			}

			// Exit if input is "exit"
			if input == "exit" {
				return nil
			}

			// Print vendor name
			fmt.Println(ipToVendor(input, db, rootOpts.printAllInfo))
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
	rootCmd.Flags().BoolVarP(&rootOpts.printAllInfo, "all", "a", false, "Print all information of OUI")
}
