package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

var (
	addr    = flag.String("addr", "", "")
	dataset = flag.String("dataset", "", "")
)

func run(ctx context.Context) error {
	// Find the previous snapshots of the source.
	var snapshots []string
	{
		var buf bytes.Buffer
		cmd := exec.CommandContext(ctx, "zfs", "list", "-t", "snapshot", *dataset, "-o", "name", "-H")
		cmd.Stdout = &buf
		cmd.Stderr = os.Stderr
		log.Printf("command: %s\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
		snapshots = strings.Fields(buf.String())
	}
	var prev string
	if len(snapshots) == 0 {
		log.Println("No previous snapshot found")
	} else {
		prev = snapshots[len(snapshots)-1]
		log.Printf("Found a previous snapshot: %s\n", prev)
	}

	// Determine the name of the next snapshot.
	next := *dataset + "@" + time.Now().Format("2006-01-02-15:04:05")
	if prev == next {
		return fmt.Errorf("previous and next snapshots have the same name (%s)", next)
	}
	{
		cmd := exec.CommandContext(ctx, "sudo", "zfs", "snapshot", next)
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
	defer func() {
		if undo {
			cmd := exec.CommandContext(ctx, "sudo", "zfs", "destroy", next)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			log.Printf("command: %s\n", cmd)
			if err := cmd.Run(); err != nil {
				log.Printf("error destroying the snapshot %s: %v", next, err)
				return
			}
			log.Printf("Destroyed the snapshot %s\n", next)
		}
	}()

	// Set up receive.
	receiveCMD := exec.CommandContext(ctx, "ssh", *addr, "zfs receive -F -v "+strings.Replace(*dataset, "tank", "tank-mirror", 1))
	receiveCMD.Stdout = os.Stdout
	receiveCMD.Stderr = os.Stderr
	receiveStdin, err := receiveCMD.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdinpipe for receiveCMD: %w", err)
	}

	// Send.
	var sendCMD *exec.Cmd
	if prev == "" {
		sendCMD = exec.CommandContext(ctx, "zfs", "send", "-v", next)
	} else {
		sendCMD = exec.CommandContext(ctx, "zfs", "send", "-v", "-i", prev, next)
	}
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

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		s := <-c
		signal.Reset(s)
		fmt.Println(s)
		cancel()
	}()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
