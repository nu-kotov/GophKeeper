package client

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Gophkeeper",
	Short: "CLI клиент для сохранения самых секретных данных",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(sendTextCmd)
}
