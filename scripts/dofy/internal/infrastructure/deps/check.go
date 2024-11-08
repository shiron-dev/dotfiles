package deps

import "os/exec"

func checkInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
