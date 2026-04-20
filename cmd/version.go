package cmd

import (
	"fmt"
	"txtweaver/internal"

	"github.com/spf13/cobra"
)

var testFlag string
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "短的short",
	Long:  `长的long`,
	Run: func(cmd *cobra.Command, args []string) {
		test := internal.Test
		fmt.Println("从其他地方获取的值" + test)
		fmt.Println("version命令")
		if testFlag != "" {
			fmt.Println(testFlag)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&testFlag, "TTT", "T", "value", "usage")
}
