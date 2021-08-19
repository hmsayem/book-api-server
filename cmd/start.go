package cmd

import (
	"github.com/spf13/cobra"
	"go-rest-api/api"
)

// startCmd represents the start command
var port string
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		api.Run(port)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVar(&port, "port", "8000", "Set port number")
}
