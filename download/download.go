package download

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/sfragata/download-artifact/utils"
)

// Options options do be passed to GetFile
type Options struct {
	URL           url.URL
	FolderPath    string
	Filename      string
	FileExtension string
}

// WriteCounter counts the number of bytes written to it. By implementing the Write method,
// it is of the io.Writer interface and we can pass this into io.TeeReader()
// Every write to this writer, will print the progress of the file write.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints the progress of a file write
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// Return again and print current status of download
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// GetFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
// We pass an io.TeeReader into Copy() to report progress on the download.
func GetFile(options Options) error {

	// Get the data
	resp, err := http.Get(options.URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filename := options.Filename

	if utils.IsEmpty(filename) {
		_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if err != nil {
			return err
		}
		filename = params["filename"]
	} else if !utils.IsEmpty(options.FileExtension) {
		filename += "." + options.FileExtension
	}

	targetFolder := options.FolderPath

	if !strings.HasSuffix(targetFolder, "/") {
		targetFolder += "/"
	}

	file := fmt.Sprintf("%s%s", targetFolder, filename)

	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	out, err := os.Create(file + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Create our bytes counter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Printf("\nfile %s", file)

	// Rename the tmp file back to the original file
	err = os.Rename(file+".tmp", file)
	if err != nil {
		return err
	}

	return nil
}
