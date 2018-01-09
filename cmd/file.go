package cmd

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
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
	defer file.Close()

	filename = filepath.Base(file.Name())
	filestat, err := file.Stat()
	if err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit("There was an error getting stats for the file provided")
	}
	filesize := filestat.Size()

	pipeOut, pipeIn := io.Pipe()
	writer := multipart.NewWriter(pipeIn)

	bar := pb.New64(filesize).SetUnits(pb.U_BYTES)
	bar.SetWidth(80)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if _, err = postMultipart(path, writer.FormDataContentType(), pipeOut); err != nil {
			writeStdErrAndExit(err.Error())
		}
		wg.Done()
	}()

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		debug("Local HTTP: There was an error preparing the multipart form")
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit("There was an error while preparing the request")
	}

	out := io.MultiWriter(part, bar)

	bar.Start()

	if _, err = io.Copy(out, file); err != nil {
		debug(fmt.Sprintf("File: There was an error reading the file %s", filename))
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit("There was an error reading the file provided")
	}

	writer.Close()
	pipeIn.Close()
	wg.Wait()

	bar.Finish()
	fmt.Println("")
	time.Sleep(5 * time.Second)

	f := hashstack.File{Filename: filename}
	switch path {
	case "/api/wordlists":
		displayWordlist(f)
	case "/api/rules":
		displayRule(f)
	case "/api/hcstat":
		displayHCStat(f)
	}
}
