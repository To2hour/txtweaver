package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cliTest",
	Short: "Hugo is a very fast static site generator",
	Long:  `命令行的提示`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run hugo...")
		fmt.Println(args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
