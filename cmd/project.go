package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var projectCmd = &cobra.Command{
	Use:    "project",
	Short:  "Subcommands can be used to interact with hashstack projects",
	Long:   "Subcommands can be used to interact with hashstack projects",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("try hashstack-cli project -h")
	},
}

var listProjectCmd = &cobra.Command{
	Use:    "list",
	Short:  "Prints a list of all your projects",
	Long:   "Prints a list of all your projects",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		var projects []hashstack.Project
		if err := getRangeJSON("/api/projects", &projects); err != nil {
			writeStdErrAndExit(err.Error())
		}
		tbl := uitable.New()
		tbl.AddRow("ID", "NAME", "DESCRIPTION", "UPDATED")
		for _, p := range projects {
			tm := time.Unix(p.UpdatedAt, 0)
			tbl.AddRow(p.ID, p.Name, p.Description, tm.String())
		}
		fmt.Println(tbl)
	},
}

var getProjectCmd = &cobra.Command{
	Use:    "get [id]",
	Short:  "Get information about a project by id",
	Long:   "Get information about a project by id",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("[id] is required")
		}
		var project hashstack.Project
		if err := getJSON(fmt.Sprintf("/api/projects/%s", args[0]), &project); err != nil {
			writeStdErrAndExit(err.Error())
		}
		tbl := uitable.New()
		tbl.AddRow("NAME", "DESCRIPTION", "UPDATED", "OWNER")
		tbl.AddRow(project.Name, project.Description, time.Unix(project.UpdatedAt, 0).String(), project.Owner.Username)
		fmt.Println(tbl)
	},
}

type newProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var newProjectCmd = &cobra.Command{
	Use:    "new [project_name] [description]",
	Short:  "Create a new project",
	Long:   "Create a new project",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("[project_name] and [description] are required")
		}
		req := newProjectRequest{
			Name:        args[0],
			Description: args[1],
		}
		body, err := postJSON("/api/projects", req)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		var project hashstack.Project
		if err := json.Unmarshal(body, &project); err != nil {
			writeStdErrAndExit("error decoding JSON from server")
		}
		fmt.Println("[+] New project created")
		fmt.Printf("Name: %s\n", project.Name)
		fmt.Printf("Description: %s\n", project.Description)
	},
}

func init() {
	projectCmd.AddCommand(newProjectCmd)
	projectCmd.AddCommand(listProjectCmd)
	projectCmd.AddCommand(getProjectCmd)
	RootCmd.AddCommand(projectCmd)
}
