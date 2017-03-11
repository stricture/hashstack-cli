package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func uploadFile(path, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		debug(fmt.Sprintf("File: Open local file %s", filename))
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit("There was an error while trying to open the file provided")
	}
	name := filepath.Base(file.Name())
	defer file.Close()
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	part, err := form.CreateFormFile("file", name)
	if err != nil {
		debug("Local HTTP: There was an error preparing the multipart form")
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit("There was an error while preparing the request")
	}
	if _, err := io.Copy(part, file); err != nil {
		debug(fmt.Sprintf("File: There was an error reading the file %s", filename))
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit("There was an error reading the file provided")
	}
	form.Close()

	// Proxy request through progress bar.
	bar := pb.New(body.Len()).SetUnits(pb.U_BYTES)
	bar.Start()
	proxy := bar.NewProxyReader(&body)
	if _, err = postMultipart(path, form.FormDataContentType(), proxy); err != nil {
		writeStdErrAndExit(err.Error())
	}
	bar.Finish()
	fmt.Println("")
	time.Sleep(5 * time.Second)
	f := hashstack.File{Filename: name}
	switch path {
	case "/api/wordlists":
		displayWordlist(f)
	case "/api/rules":
		displayRule(f)
	case "/api/hcstat":
		displayHCStat(f)
	}
}
