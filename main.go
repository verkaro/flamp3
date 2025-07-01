// flamp3: convert FLAC files to MP3s using flac and lame
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	outDir    string
	recursive bool
	quality   string
)

func init() {
	flag.StringVar(&outDir, "out", "out", "Output directory")
	flag.BoolVar(&recursive, "recursive", false, "Recurse into directories")
	flag.StringVar(&quality, "quality", "6", "LAME VBR quality level (0=best,9=worst)")
}

func main() {
	flag.Parse()
	inputs := flag.Args()
	if len(inputs) == 0 {
		log.Fatal("Usage: flamp3 [options] <file or dir>...")
	}

	// Verify required tools
	for _, tool := range []string{"flac", "lame", "sh"} {
		if _, err := exec.LookPath(tool); err != nil {
			log.Fatalf("Required tool '%s' not found in PATH", tool)
		}
	}

	// create base output directory
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	for _, in := range inputs {
		info, err := os.Stat(in)
		if err != nil {
			log.Printf("Skipping %s: %v", in, err)
			continue
		}
		if info.IsDir() {
			if recursive {
				root := filepath.Base(in)
				filepath.WalkDir(in, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						log.Printf("Error reading %s: %v", path, err)
						return nil
					}
					if d.IsDir() {
						return nil
					}
					if filepath.Ext(path) != ".flac" {
						return nil
					}
					rel, _ := filepath.Rel(in, path)
					rel = filepath.Join(root, rel)
					processFlac(path, rel)
					return nil
				})
			} else {
				log.Printf("Skipping directory %s (use -recursive)", in)
			}
		} else {
			if filepath.Ext(in) == ".flac" {
				processFlac(in, filepath.Base(in))
			} else {
				log.Printf("Skipping non-FLAC file %s", in)
			}
		}
	}
}

// processFlac converts a single FLAC to MP3, preserving relPath under outDir
func processFlac(inputPath, relPath string) {
	base := relPath[:len(relPath)-len(filepath.Ext(relPath))]
	outRel := base + ".mp3"
	outPath := filepath.Join(outDir, outRel)

	// skip if output exists
	if _, err := os.Stat(outPath); err == nil {
		log.Printf("Skipping %s: output already exists at %s", inputPath, outPath)
		return
	} else if !os.IsNotExist(err) {
		log.Printf("Error checking %s: %v", outPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		log.Printf("Failed to create dir for %s: %v", outPath, err)
		return
	}

	// build shell pipeline: flac decode to stdout | lame vbr -> outPath
	cmdStr := fmt.Sprintf("flac -c -d %q | lame --vbr-new -V %s - %q", inputPath, quality, outPath)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Conversion failed for %s: %v", inputPath, err)
		return
	}

	fmt.Printf("Converted %s -> %s\n", inputPath, outPath)
}

