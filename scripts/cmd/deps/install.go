package deps

import "os/exec"

func installHomebrew() {
	exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)")
}
