package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func getListByID(list *hashstack.List) {
	path := fmt.Sprintf("/api/projects/%d/lists/%d", list.ProjectID, list.ID)
	if err := getJSON(path, list); err != nil {
		writeStdErrAndExit(err.Error())
	}
}
func getListByName(list *hashstack.List) {
	path := fmt.Sprintf("/api/projects/%d/lists?name=%s", list.ProjectID, list.Name)
	if err := getJSON(path, list); err != nil {
		writeStdErrAndExit(err.Error())
	}
}
func getList(projectID int64, arg string) hashstack.List {
	list := hashstack.List{
		ProjectID: projectID,
	}
	i, err := strconv.Atoi(arg)
	if err != nil {
		list.Name = arg
		getListByName(&list)
	} else {
		list.ID = int64(i)
		getListByID(&list)
	}
	return list
}

func displayList(list hashstack.List) {
	liststat := fmt.Sprintf("%d/%d (%0.2f%%) hashes", list.RecoveredCount, list.DigestCount, percentOf(int(list.RecoveredCount), int(list.DigestCount)))
	fmt.Printf("ID..............: %d\n", list.ID)
	fmt.Printf("Name............: %s\n", list.Name)
	fmt.Printf("Hash Mode.......: %d\n", list.HashMode)
	fmt.Printf("Cracked.........: %s\n", liststat)
	fmt.Println()
}

func displayLists(arg string) {
	project := getProject(arg)
	count, err := getListCount(project.ID)
	if err != nil {
		writeStdErrAndExit(err.Error())
	}
	if count < 1 {
		writeStdErrAndExit("You have not created any lists for this project.")
	}
	path := fmt.Sprintf("/api/projects/%d/lists", project.ID)
	var lists []hashstack.List
	if err := getRangeJSON(path, &lists); err != nil {
		writeStdErrAndExit(err.Error())
	}
	for _, l := range lists {
		displayList(l)
	}
}

var listCmd = &cobra.Command{
	Use:   "lists <project_name|project_id> [list_name|list_id]",
	Short: "Displays a list of all lists associated with the provided project (-h or --help for subcommands).",
	Long: `
Displays a list of all lists associated with the provided project. If list_name|list_id is provided, details will be displayed for
that specific list. Additional subcommands are available.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			writeStdErrAndExit("project_name or project_id required. -h --help for subcommands.")
		case 1:
			displayLists(args[0])
		case 2:
			project := getProject(args[0])
			displayList(getList(project.ID, args[1]))
		default:
			cmd.Usage()
		}
	},
}

type binaryRequest struct {
	ProjectID   int64  `json:"project_id"`
	HashMode    int    `json:"hash_mode"`
	EncodedHash string `json:"encoded_hash"`
	Filename    string `json:"filename"`
	Name        string `json:"name"`
}

type multiTrackRequest struct {
	ProjectID int64            `json:"project_id"`
	HashMode  int              `json:"hash_mode"`
	Hashes    []multiTrackItem `json:"hashes"`
	Name      string           `json:"name"`
}

type multiTrackItem struct {
	Hash     string `json:"hash"`
	Filename string `json:"filename"`
}

var (
	flIsHexSalt bool
)

func uploadList(pid int64, mode int, filename string) {
	var (
		hashMode hashstack.HashMode
		resp     []byte
	)
	if err := getJSON(fmt.Sprintf("/api/hash_modes?mode=%d", mode), &hashMode); err != nil {
		writeStdErrAndExit("The selected mode is not supported by the server.")
	}
	if !hashMode.IsSupported {
		writeStdErrAndExit("The selected mode is not supported by the server.")
	}
	if !hashMode.IsBinary && hashMode.Upload == "" {
		file, err := os.Open(filename)
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error opening the provided file.")
		}
		defer file.Close()

		_, filenamesplit := filepath.Split(file.Name())
		filestat, err := file.Stat()
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error getting stats for the file provided")
		}
		filesize := filestat.Size()

		if filesize > int64(15*1024*1024) {
			writeStdErrAndExit("This list exceeds the maximum size supported by the server (15 MB).")
		}

		pipeOut, pipeIn := io.Pipe()
		writer := multipart.NewWriter(pipeIn)

		bar := pb.New64(filesize).SetUnits(pb.U_BYTES)
		bar.SetWidth(80)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			resp, err = postMultipart(fmt.Sprintf("/api/projects/%d/lists/multi", pid), writer.FormDataContentType(), pipeOut)
			if err != nil {
				fmt.Println("")
				fmt.Println("")
				fmt.Println("This error likely occurred because you did not have any valid hashes.")
				writeStdErrAndExit(err.Error())
			}
			wg.Done()
		}()

		part, err := writer.CreateFormFile("file", file.Name())
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit(new(requestCreateError).Error())
		}

		out := io.MultiWriter(part, bar)

		bar.Start()

		if _, err := io.Copy(out, file); err != nil {
			bar.Finish()
			debug(fmt.Sprintf("Error: %s", err.Error()))
			if err.Error() == "io: read/write on closed pipe" {
				writeStdErrAndExit("The list exceeded the maxmimum size supported by the server (15 MB).")
			}
			writeStdErrAndExit("There was an error reading the provided file.")
		}
		writer.WriteField("hash_mode", strconv.Itoa(hashMode.HashMode))
		writer.WriteField("name", filenamesplit)
		isHexSaltStr := "false"
		if flIsHexSalt {
			isHexSaltStr = "true"
		}
		writer.WriteField("is_hex_salt", isHexSaltStr)

		writer.Close()
		pipeIn.Close()
		wg.Wait()

		bar.Finish()
		fmt.Println("")
	} else if hashMode.Upload != "" {
		_, name := filepath.Split(filename)
		file, err := os.Open(filename)
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error opening the provided file.")
		}
		defer file.Close()
		var hashes []multiTrackItem
		var lineNum int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				writeStdErrAndExit(fmt.Sprintf("Line number %d is not in the format %s!\n\nThis is required for this format.", lineNum, hashMode.Upload))
			}
			hashes = append(hashes, multiTrackItem{
				Filename: parts[0],
				Hash:     parts[1],
			})
		}
		if len(hashes) < 1 {
			writeStdErrAndExit("There were no parsable hashes in the provided file")
		}
		if err := scanner.Err(); err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error reading the file.")
		}
		req := multiTrackRequest{
			ProjectID: pid,
			HashMode:  hashMode.HashMode,
			Hashes:    hashes,
			Name:      name,
		}
		fmt.Printf("Uploading %d hashes from %s...\n", len(hashes), name)
		resp, err = postJSON(fmt.Sprintf("/api/projects/%d/lists/multitrack", pid), req)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Printf("Upload complete...\n\n")
	} else if hashMode.IsBinary {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error reading the provided file.")
		}
		_, name := filepath.Split(filename)
		req := binaryRequest{
			ProjectID:   pid,
			HashMode:    hashMode.HashMode,
			EncodedHash: base64.StdEncoding.EncodeToString(data),
			Filename:    name,
			Name:        name,
		}

		fmt.Printf("Uploading binary hash from %s...\n", name)
		resp, err = postJSON(fmt.Sprintf("/api/projects/%d/lists/binary", pid), req)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Printf("Upload complete...\n\n")
	} else {
		writeStdErrAndExit("Malformed hash_mode. Please contact support!")
	}

	var list hashstack.List
	if err := json.Unmarshal(resp, &list); err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit(new(jsonServerError).Error())
	}
	displayList(list)
}

var addListCmd = &cobra.Command{
	Use:   "add <project_name|project_id> <mode> <file>",
	Short: "Add a new list to a project.",
	Long: `
Add a new file containing one or more hashes to a project by project_name or project_id. Modes can be viewed
using the "modes" subcommand. The file name must be unique across projects.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			writeStdErrAndExit("project_name|project_id, mode, and file are required.")
		}
		var (
			pidStr   = args[0]
			modeStr  = args[1]
			filename = args[2]
		)
		mode, err := strconv.Atoi(modeStr)
		if err != nil {
			writeStdErrAndExit("mode is invalid")
		}
		project := getProject(pidStr)
		uploadList(project.ID, mode, filename)
	},
}

func deleteList(projectID int64, listID int64) {
	path := fmt.Sprintf("/api/projects/%d/lists/%d", projectID, listID)
	if err := deleteHTTP(path); err != nil {
		writeStdErrAndExit(err.Error())
	}
	fmt.Println("The list was deleted successfully.")
}

var delListCmd = &cobra.Command{
	Use:   "delete <project_name|project_id> <list_name|list_id>",
	Short: "Delete a list from a project.",
	Long: `
Delete a list from a project by project_name or project_id. Deleting a list also deletes
the associated plains.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_name|project_id and list_name|list_id are required.")
		}
		project := getProject(args[0])
		list := getList(project.ID, args[1])
		if ok := promptDelete("this list"); !ok {
			writeStdErrAndExit("Not deleting list.")
		}
		deleteList(project.ID, list.ID)
	},
}

var crackedListCmd = &cobra.Command{
	Use:   "cracked <project_name|project_id> <list_name|list_id>",
	Short: "Download cracked hashes for a list.",
	Long:  "Download cracked hashes for a list.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_name|project_id and list_id is required.")
		}
		project := getProject(args[0])
		list := getList(project.ID, args[1])
		body, err := getReader(fmt.Sprintf("/api/projects/%d/lists/%d/plains", project.ID, list.ID))
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		io.Copy(os.Stdout, body)
	},
}

var uncrackedListCmd = &cobra.Command{
	Use:   "uncracked <project_name|project_id> <list_name|list_id>",
	Short: "Download uncracked hashes for a list.",
	Long:  "Download uncracked hashes for a list.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_name|project_id and list_id is required.")
		}
		project := getProject(args[0])
		list := getList(project.ID, args[1])
		body, err := getReader(fmt.Sprintf("/api/projects/%d/lists/%d/hashes", project.ID, list.ID))
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		io.Copy(os.Stdout, body)
	},
}

func init() {
	addListCmd.PersistentFlags().BoolVar(&flIsHexSalt, "hex-salt", false, "Assume is given in hex")
	listCmd.AddCommand(addListCmd)
	listCmd.AddCommand(delListCmd)
	listCmd.AddCommand(crackedListCmd)
	listCmd.AddCommand(uncrackedListCmd)
	RootCmd.AddCommand(listCmd)
}
