package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/schollz/progressbar/v3"
)

func main() {
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "No API key has been configured. Please export ASSEMBLYAI_API_KEY and run the command again.")
		os.Exit(1)
	}

	args := os.Args[1:]

	var filePath string

	if len(args) == 1 {
		filePath = args[0]
	} else {
		fmt.Fprintf(os.Stderr, "Error: Expected 1 argument, but got %d.\n", len(args))
		os.Exit(1)
	}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to open %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer f.Close()

	ctx := context.Background()

	fi, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to get file information for %s: %v\n", filePath, err)
		os.Exit(1)
	}

	progressReader := io.TeeReader(f, newProgressReader(fi.Size()))

	client := assemblyai.NewClient(apiKey)

	url, err := client.Upload(ctx, progressReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to upload file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully uploaded %s to AssemblyAI.\n\n", filePath)
	fmt.Printf("Use the following URL to transcribe the file:\n\n%s\n\n", url)
	fmt.Println("The URL is only accessible from AssemblyAI's servers.")
}

// progressReader implements the io.Writer interface and updates the progress
// bar when data is written.
type progressReader struct {
	progress *progressbar.ProgressBar
}

func newProgressReader(maxBytes int64) *progressReader {
	return &progressReader{
		progress: progressbar.DefaultBytes(maxBytes),
	}
}

func (p *progressReader) Write(b []byte) (int, error) {
	if err := p.progress.Add(len(b)); err != nil {
		return 0, err
	}
	return len(b), nil
}
