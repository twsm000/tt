package translator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TransCLI struct{}

func (t TransCLI) Translate(word string) (string, error) {
	fpath, _ := os.Executable()
	fpath = filepath.Join(filepath.Dir(fpath), "trans")
	cmd := exec.Command(fpath, "-brief", word)
	data, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\n"), nil
}
