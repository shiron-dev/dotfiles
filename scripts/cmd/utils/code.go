package utils

import "os/exec"

func OpenWithCode(path ...string) {
	args := []string{"-n", "-w"}
	args = append(args, path...)
	exec.Command("code", args...).Run()
}
