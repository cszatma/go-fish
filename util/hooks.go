package util

import (
	"bufio"
	"os"
	"strings"
)

var hookList = [...]string{
	"applypatch-msg",
	"pre-applypatch",
	"post-applypatch",
	"pre-commit",
	"prepare-commit-msg",
	"commit-msg",
	"post-commit",
	"pre-rebase",
	"post-checkout",
	"post-merge",
	"pre-push",
	"pre-receive",
	"update",
	"post-receive",
	"post-update",
	"push-to-checkout",
	"pre-auto-gc",
	"post-rewrite",
	"sendemail-validate",
}

func isGoFish(path string) (bool, error) {
	const id = "# Hook created by go-fish"

	f, err := os.Open(path)

	if err != nil {
		return false, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	// Look for `id` in file to see if it was generated by go-fish
	for s.Scan() {
		if strings.Contains(s.Text(), id) {
			return true, nil
		}
	}

	return false, nil
}

func writeHook(path, script string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(script)
	f.Sync()

	return err
}

func deleteHook(path string) error {
	err := os.Remove(path)
	return err
}

func createHook(name, path, script string, force bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		VerbosePrintf("Creating hook: %s\n", name)
		return writeHook(path, script)
	}

	if is, err := isGoFish(path); err != nil {
		return err
	} else if is || force {
		VerbosePrintf("Updating existing hook: %s\n", name)
		return writeHook(path, script)
	}

	VerbosePrintf("Skipping existing user hook: %s\n", name)
	return nil
}

func removeHook(name, path string, force bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		VerbosePrintf("Hook does not exists, skipping: %s\n", name)
		return nil
	}

	if is, err := isGoFish(path); err != nil {
		return err
	} else if is || force {
		VerbosePrintf("Removing hook: %s\n", name)
		return deleteHook(path)
	}

	VerbosePrintf("Skipping user hook: %s\n", name)
	return nil
}

// CreateHooks creates each git hook.
func CreateHooks(path, script string, force bool) error {
	for _, v := range hookList {
		hookPath := path + "/" + v
		err := createHook(v, hookPath, script, force)

		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveHooks removes each git hook.
func RemoveHooks(path string, force bool) error {
	for _, v := range hookList {
		hookPath := path + "/" + v
		err := removeHook(v, hookPath, force)

		if err != nil {
			return err
		}
	}

	return nil
}
