package cmd

import (
	"bufio"
	"bytes"
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
	liststat := fmt.Sprintf("%d/%d (%d%%)", list.RecoveredCount, list.DigestCount, list.RecoveredCount/list.DigestCount)
	fmt.Printf("ID..............: %d\n", list.ID)
	fmt.Printf("Name............: %s\n", list.Name)
	fmt.Printf("Hash Mode.......: %d\n", list.HashMode)
	fmt.Printf("Recoverd........: %s\n", liststat)
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
	Short: "Displays a list of all lists associated with the provided project (-h or --help for subcommands)",
	Long: `
Displays a list of all lists associated with the provided project. If list_name|list_id is provied, details will be displayed for 
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
	if !hashMode.IsBinary && !hashMode.IsScrapable {
		file, err := os.Open(filename)
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error opening the provided file.")
		}
		defer file.Close()
		var body bytes.Buffer
		form := multipart.NewWriter(&body)
		_, name := filepath.Split(file.Name())
		part, err := form.CreateFormFile("file", file.Name())
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit(new(requestCreateError).Error())
		}
		if _, err := io.Copy(part, file); err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error reading the provided file.")
		}
		form.WriteField("hash_mode", strconv.Itoa(hashMode.HashMode))
		form.WriteField("name", name)
		isHexSaltStr := "false"
		if flIsHexSalt {
			isHexSaltStr = "true"
		}
		form.WriteField("is_hex_salt", isHexSaltStr)
		form.Close()
		// Proxy request through progress bar.
		bar := pb.New(body.Len()).SetUnits(pb.U_BYTES)
		bar.SetWidth(80)
		bar.Start()
		proxy := bar.NewProxyReader(&body)
		resp, err = postMultipart(fmt.Sprintf("/api/projects/%d/lists/nonbinary", pid), form.FormDataContentType(), proxy)
		if err != nil {
			fmt.Println("")
			fmt.Println("")
			fmt.Println("This error likely occurred because you did not have any valid hashes.")
			writeStdErrAndExit(err.Error())
		}
		bar.Finish()
		fmt.Println("")
	}
	if hashMode.IsBinary && hashMode.IsScrapable {
		_, name := filepath.Split(filename)
		file, err := os.Open(filename)
		if err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error opening the provided file.")
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
				continue
			}
			hashes = append(hashes, binaryScrapableItem{
				Filename: parts[0],
				Hash:     parts[1],
			})
		}
		if err := scanner.Err(); err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit("There was an error reading the file.")
		}
		req := binaryScrapableRequest{
			ProjectID: pid,
			HashMode:  hashMode.HashMode,
			Hashes:    hashes,
			Name:      name,
		}
		resp, err = postJSON(fmt.Sprintf("/api/projects/%d/lists/binaryscrapable", pid), req)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
	}

	if hashMode.IsBinary {
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
		resp, err = postJSON(fmt.Sprintf("/api/projects/%d/lists/binary", pid), req)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
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
	Short: "Add a new list to a project",
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
	Short: "Delete a list from a project",
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

var plainsListCmd = &cobra.Command{
	Use:   "plains <project_name|project_id> <list_name|list_id>",
	Short: "Download plains for a list",
	Long:  "Download plains for a list",
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

var hashesListCmd = &cobra.Command{
	Use:   "hashes <project_name|project_id> <list_name|list_id>",
	Short: "Download hashes for a list",
	Long:  "Download hashes for a list",
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
	listCmd.AddCommand(plainsListCmd)
	listCmd.AddCommand(hashesListCmd)
	RootCmd.AddCommand(listCmd)
}
