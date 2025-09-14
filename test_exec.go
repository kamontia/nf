package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Testing os/exec by running 'sleep 1'...")
	cmd := exec.Command("sleep", "1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command returned an error: %v\n", err)
	}
	fmt.Println("Test program finished successfully.")
}
