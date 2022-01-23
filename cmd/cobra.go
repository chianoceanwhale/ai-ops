package cmd

import (
	"errors"

	au "ai-ops/cmd/admin"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "ai-ops",
	Short:             "-v",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `ai-ops`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		usageStr := `ai-ops 1.0.0 欢迎使用，可以是用 -h 查看命令`
		log.Printf("%s\n", usageStr)
	},
}

func init() {
	rootCmd.AddCommand(au.StartCmd)

}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
