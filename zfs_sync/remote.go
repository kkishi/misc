package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var addr = flag.String("addr", "", "")

func run() error {
	// Find the previous snapshots of the source.
	var buf bytes.Buffer
	{
		cmd := exec.Command("zfs", "list", "-t", "snapshot", "tank/photos", "-o", "name", "-H")
		cmd.Stdout = &buf
		cmd.Stderr = os.Stderr
		log.Printf("command: %s\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
		log.Printf("existing snapshots:\n%s", buf.String())
	}
	snapshots := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(snapshots) == 0 {
		return errors.New("no snapshot found")
	}
	prev := snapshots[len(snapshots)-1]
	log.Printf("Found a previous snapshot: %s\n", prev)

	// Determine the name of the next snapshot.
	next := fmt.Sprintf("tank/photos@%s", time.Now().Format("2006-01-02-15:04:05"))
	if prev == next {
		return fmt.Errorf("previous and next snapshots have the same name (%s)", next)
	}
	{
		cmd := exec.Command("sudo", "zfs", "snapshot", next)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Printf("command: %s\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
		log.Printf("Created a new snapshot: %s\n", next)
	}
	// The new snapshot is not useful unless the subsequent send/receive succeeds.
	// If something goes wrong, it should be destroyed.
	undo := true // This variable is set to false when this function succeeds.
	defer func() error {
		if !undo {
			return nil
		}
		cmd := exec.Command("sudo", "zfs", "destroy", next)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Printf("command: %s\n", cmd)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error destroying the snapshot %s: %v", next, err)
		}
		log.Printf("Destroyed the snapshot %s\n", next)
		return nil
	}()

	// Set up receive.
	receiveCMD := exec.Command("ssh", *addr, "zfs receive -F -v tank-mirror/photos")
	receiveCMD.Stdout = os.Stdout
	receiveCMD.Stderr = os.Stderr
	receiveStdin, err := receiveCMD.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdinpipe for receiveCMD: %w", err)
	}

	// Send.
	sendCMD := exec.Command("zfs", "send", "-v", "-i", prev, next)
	sendCMD.Stdout = receiveStdin
	sendCMD.Stderr = os.Stderr
	log.Printf("command: %s\n", sendCMD)
	go func() {
		defer receiveStdin.Close()
		err := sendCMD.Run()
		if err != nil {
			log.Printf("sendCMD error: %v", err)
		} else {
			log.Println("sendCMD success")
		}
	}()
	log.Printf("command: %s\n", receiveCMD)
	if err := receiveCMD.Run(); err != nil {
		return fmt.Errorf("receiveCMD error: %v", err)
	}

	log.Println("zfs send complete")
	undo = false

	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
