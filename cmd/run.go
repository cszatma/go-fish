package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TouchBistro/goutils/color"
	"github.com/TouchBistro/goutils/command"
	"github.com/TouchBistro/goutils/fatal"
	"github.com/cszatmary/go-fish/config"
	"github.com/cszatmary/go-fish/git"
	"github.com/cszatmary/go-fish/hooks"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the specified git hook",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hookName := args[0]
		log.WithFields(log.Fields{
			"hook": hookName,
		}).Debug("Running hook")

		hook, exists := config.GetHook(hookName)
		if !exists {
			log.WithFields(log.Fields{
				"hook": hookName,
			}).Debug("No action defined for hook in config, skipping")
			return
		}

		rootDir, err := git.RootDir()
		if err != nil {
			fatal.ExitErr(err, "Unable to find root directory of git repo")
		}

		fmt.Printf(color.Cyan("🎣 go-fish > %s\n"), hookName)

		// If pre-commit, set STAGED_FILES env var
		if hookName == hooks.PreCommit {
			log.WithFields(log.Fields{
				"hook": hookName,
			}).Debug("Getting staged files")
			stagedFiles, err := git.StagedFiles()
			if err != nil {
				fatal.ExitErr(err, "Failed to get staged files")
			}

			log.WithFields(log.Fields{
				"hook":  hookName,
				"files": stagedFiles,
			}).Debug("Found staged files")
			os.Setenv("STAGED_FILES", strings.Join(stagedFiles, "\n"))
		}

		log.Debugf("Running: %q", hook.Run)
		err = command.Exec("sh", []string{"-c", hook.Run}, hookName, func(c *exec.Cmd) {
			c.Dir = rootDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
		})
		if err != nil {
			fatal.ExitErrf(err, "Hook failed: %s", hookName)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
