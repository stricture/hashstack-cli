package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"

	"io"

	"path/filepath"

	"bufio"

	"strings"

	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var listCmd = &cobra.Command{
	Use:    "lists",
	Short:  "Subcommands can be used to interact with hashstack lists",
	Long:   "Subcommands can be used to interact with hashstack lists",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("try hashstack-cli lists -h")
	},
}

type binaryRequest struct {
	ProjectID   int64  `json:"project_id"`
	HashMode    int    `json:"hash_mode"`
	EncodedHash string `json:"encoded_hash"`
	Filename    string `json:"filename"`
	Name        string `json:"name"`
}

type binaryScrapableRequest struct {
	ProjectID int64                 `json:"project_id"`
	HashMode  int                   `json:"hash_mode"`
	Hashes    []binaryScrapableItem `json:"hashes"`
	Name      string                `json:"name"`
}

type binaryScrapableItem struct {
	Hash     string `json:"hash"`
	Filename string `json:"filename"`
}

var newListCmd = &cobra.Command{
	Use:    "new [project_id] [mode] [file]",
	Short:  "Upload a new hash or hash list to hashstack",
	Long:   "Upload a new hash or hash list to hashstack",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			pid      = args[0]
			mode     = args[1]
			filename = args[2]
			hashMode hashstack.HashMode
			resp     []byte
		)
		projectID, err := strconv.Atoi(pid)
		if err != nil {
			writeStdErrAndExit("[project_id] is not valid")
		}
		if err := getJSON(fmt.Sprintf("/api/hash_modes?mode=%s", mode), &hashMode); err != nil {
			writeStdErrAndExit(err.Error())
		}
		if !hashMode.IsSupported {
			writeStdErrAndExit("the selected mode is not supported by the server")
		}
		if !hashMode.IsBinary && !hashMode.IsScrapable {
			file, err := os.Open(filename)
			if err != nil {
				debug(err.Error())
				writeStdErrAndExit("there was an error opening the provided file")
			}
			defer file.Close()
			var body bytes.Buffer
			form := multipart.NewWriter(&body)
			_, name := filepath.Split(file.Name())
			part, err := form.CreateFormFile("file", file.Name())
			if err != nil {
				debug(err.Error())
				writeStdErrAndExit("there was an error generating the request")
			}
			if _, err := io.Copy(part, file); err != nil {
				debug(err.Error())
				writeStdErrAndExit("there was an error reading the provided file")
			}
			form.WriteField("hash_mode", strconv.Itoa(hashMode.HashMode))
			form.WriteField("name", name)
			form.Close()
			resp, err = postMultipart(fmt.Sprintf("/api/projects/%s/lists/nonbinary", pid), form.FormDataContentType(), &body)
			if err != nil {
				writeStdErrAndExit(err.Error())
			}

		}
		if hashMode.IsBinary && hashMode.IsScrapable {
			_, name := filepath.Split(filename)
			file, err := os.Open(filename)
			if err != nil {
				debug(err.Error())
				writeStdErrAndExit("there was an error opening the provided file")
			}
			defer file.Close()
			var hashes []binaryScrapableItem
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) != 2 {
					// TODO: Should we warn here?
					continue
				}
				hashes = append(hashes, binaryScrapableItem{
					Filename: parts[0],
					Hash:     parts[1],
				})
			}
			if err := scanner.Err(); err != nil {
				debug(err.Error())
				writeStdErrAndExit("there was an error reading the file")
			}
			req := binaryScrapableRequest{
				ProjectID: int64(projectID),
				HashMode:  hashMode.HashMode,
				Hashes:    hashes,
				Name:      name,
			}
			resp, err = postJSON(fmt.Sprintf("/api/projects/%s/lists/binaryscrapable", pid), req)
			if err != nil {
				writeStdErrAndExit(err.Error())
			}
		}

		if hashMode.IsBinary {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				debug(err.Error())
				writeStdErrAndExit("there was an error reading the provided file")
			}
			_, name := filepath.Split(filename)
			req := binaryRequest{
				ProjectID:   int64(projectID),
				HashMode:    hashMode.HashMode,
				EncodedHash: base64.StdEncoding.EncodeToString(data),
				Filename:    name,
				Name:        name,
			}
			resp, err = postJSON(fmt.Sprintf("/api/projects/%s/lists/binary", pid), req)
			if err != nil {
				writeStdErrAndExit(err.Error())
			}
		}

		var list hashstack.List
		if err := json.Unmarshal(resp, &list); err != nil {
			debug(err.Error())
			writeStdErrAndExit("error decoding json returned form server")
		}
		fmt.Printf("JSON LIST\n")
		fmt.Printf("%+v\n", list)
	},
}

func init() {
	listCmd.AddCommand(newListCmd)
	RootCmd.AddCommand(listCmd)
}
