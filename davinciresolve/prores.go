package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var outputDir = flag.String("output_dir", "", "")

func main() {
	flag.Parse()
	if *outputDir == "" {
		log.Fatal("Please specify --output_dir")
	}
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatal(err)
	}
	for _, f := range flag.Args() {
		base := filepath.Base(f)
		ext := filepath.Ext(base)
		outputPath := filepath.Join(*outputDir, base[:len(base)-len(ext)]+".MOV")
		cmd := exec.Command("ffmpeg", "-hwaccel", "cuda", "-i", f, "-vcodec", "prores", "-profile:v", "2", "-acodec", "pcm_s16le", outputPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
