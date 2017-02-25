package invoker

import (
	"fmt"
	"os/exec"
)

func Invoke(dir, script string, args ...string) (string, error) {
	run := fmt.Sprintf("%s/%s/%s.py", dir, script, script)
	if len(run) == 0 {
		return "", fmt.Errorf("No script directory")
	}

	out, err := exec.Command(run, args...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
