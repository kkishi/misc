package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run() error {
	// Import the destination if it is not imported yet.
	if _, err := os.Stat("/tank-mirror/photos"); os.IsNotExist(err) {
		fmt.Println("Importing the destination")
		cmd := exec.Command("sudo", "zpool", "import", "tank-mirror")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Find the previous snapshots of the source.
	var buf bytes.Buffer
	{
		cmd := exec.Command("zfs", "list", "-t", "snapshot", "tank/photos", "-o", "name", "-H")
		cmd.Stdout = &buf
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	snapshots := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(snapshots) == 0 {
		return errors.New("no snapshot found")
	}
	prev := snapshots[len(snapshots)-1]
	fmt.Printf("Found a previous snapshot: %s\n", prev)

	// Determine the name of the next snapshot.
	next := fmt.Sprintf("tank/photos@%s", time.Now().Format("2006-01-02-03:04"))
	if prev == next {
		return fmt.Errorf("Previous and next snapshots ")
	}
	{
		cmd := exec.Command("sudo", "zfs", "snapshot", next)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Printf("Created a new snapshot: %s\n", next)
	}

	// Send.
	{
		cmdStr := fmt.Sprintf("sudo zfs send -i %s %s | sudo zfs receive -F tank-mirror/photos", prev, next)
		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Println("zfs send complete")
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
