package conf

import (
	"os/exec"
	"runtime"
	"strings"
)

type EnvInfo struct {
	os        string
	osVersion string
	arch      string
}

func ScanEnvInfo() *EnvInfo {
	osVersion, _ := exec.Command("uname", "-r").Output()
	arch, _ := exec.Command("uname", "-p").Output()
	return &EnvInfo{
		os:        runtime.GOOS,
		osVersion: strings.TrimSpace(string(osVersion)),
		arch:      strings.TrimSpace(string(arch)),
	}
}
