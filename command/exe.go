package command

import (
	"os/exec"
)


func ExeBash(cmd string) ([]byte, error) {
	c := exec.Command("sh", "-c", cmd)
	output, err := c.CombinedOutput()
	return output, err
}
