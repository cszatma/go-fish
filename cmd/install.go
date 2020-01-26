package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cszatma/go-fish/config"
	"github.com/cszatma/go-fish/fatal"
	"github.com/cszatma/go-fish/git"
	"github.com/cszatma/go-fish/hooks"
	"github.com/cszatma/go-fish/util"
	p "github.com/cszatma/printer"
	"github.com/spf13/cobra"
)

var (
	force bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the git hooks",
	Long: `Install will generate missing git hooks and recreate git hooks that were created by go-fish.
By default it will not replace existing git hooks not created by go-fish.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(p.Cyan("Installing git hooks..."))
		util.VerbosePrintln("Finding root directory of git repo")
		rootDir, gitDir, err := git.RootDir()
		if err != nil {
			fatal.ExitErr(err, "Unable to find git directory")
		}

		goFishPath, err := os.Executable()
		if err != nil {
			fatal.ExitErr(err, "Unable to find path of go-fish")
		}

		if util.IsCI() && config.All().SkipCI {
			fmt.Println("CI detected, skipping Git hooks installation")
			return
		}

		hooksPath := filepath.Join(gitDir, "hooks")
		if !util.FileOrDirExists(hooksPath) {
			util.VerbosePrintln(".git/hooks does not exist, creating...")
			err = os.Mkdir(hooksPath, 0755)
			if err != nil {
				fatal.ExitErr(err, "Failed to create .git/hooks directory")
			}
		}

		util.VerbosePrintln("Rendering git hook script")
		script, err := hooks.RenderScript(goFishPath, rootDir, version)
		if err != nil {
			fatal.ExitErr(err, "Failed to generate git hook script")
		}

		util.VerbosePrintf("Creating git hooks")
		err = hooks.CreateHooks(hooksPath, script, force)
		if err != nil {
			fatal.ExitErr(err, "Failed to create git hooks")
		}

		fmt.Println(p.Green("Successfully installed Git hooks! Enjoy! 🎣"))
	},
}

func init() {
	installCmd.Flags().BoolVarP(&force, "force", "f", false, "replaces any existing hooks not generated by go-fish")
	rootCmd.AddCommand(installCmd)
}
