package cmd

import (
	"fmt"

	"github.com/TouchBistro/goutils/color"
	"github.com/TouchBistro/goutils/fatal"
	"github.com/cszatmary/go-fish/git"
	"github.com/cszatmary/go-fish/hooks"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type uninstallOptions struct {
	removeAll bool
}

var uninstallOpts uninstallOptions

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Args:  cobra.NoArgs,
	Short: "Uninstall the git hooks",
	Long: `Uninstall will remove all the git hooks generated by go-fish.
By default it will not touch hooks not created by go-fish.`,
	Run: func(cmd *cobra.Command, args []string) {
		gitDir, err := git.GitDir()
		if err != nil {
			fatal.ExitErr(err, "Failed to find .git directory")
		}
		log.Debugf("Found git directory at: %s", gitDir)

		fmt.Println("Uninstalling git hooks...")
		err = hooks.RemoveHooks(gitDir, uninstallOpts.removeAll)
		if err != nil {
			fatal.ExitErr(err, "Failed to remove git hooks")
		}

		fmt.Println(color.Green("Successfully uninstalled Git hooks! 🎣"))
	},
}

func init() {
	uninstallCmd.Flags().BoolVar(&uninstallOpts.removeAll, "all", false, "deletes all hooks including ones not generated by go-fish")
	rootCmd.AddCommand(uninstallCmd)
}
