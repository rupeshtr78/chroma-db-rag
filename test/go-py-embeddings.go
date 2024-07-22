package test

import (
	"fmt"
	"os/exec"
	"strings"
)

func EmbedMain() {
	cmd := exec.Command("python3", "embedding.py", "your sentence here")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strings.TrimSpace(string(output)))
}
