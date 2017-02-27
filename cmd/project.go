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
	Use:   "projects",
	Short: "Interact with hashstack projects",
	Long: `
Projects are used to organize one or more lists. An exmaple of a project
may be AcmeForensicIvestigation2016. Once a project is created, you will
upload your lists using the project's id.

You can share projects with many users.
	`,
	PreRun: ensureAuth,
	Run:    topUsage,
}

var listProjectCmd = &cobra.Command{
	Use:    "list",
	Short:  "Displays a list of all your projects",
	Long:   "Displays a list of all your projects",
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
	Use:    "show <id>",
	Short:  "Show information about a project by id",
	Long:   "Show information about a project by id",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("id is required")
		}
		var project hashstack.Project
		if err := getJSON(fmt.Sprintf("/api/projects/%s", args[0]), &project); err != nil {
			writeStdErrAndExit(err.Error())
		}
		tbl := uitable.New()
		tbl.AddRow("NAME", "DESCRIPTION", "UPDATED", "OWNER")
		tbl.AddRow(project.Name, project.Description, time.Unix(project.UpdatedAt, 0).String(), project.Owner.Username)
		tbl.AddRow("")
		fmt.Println(tbl)

		fmt.Println("Contributors:")
		tbl2 := uitable.New()
		tbl2.AddRow("ID", "USERNAME")
		for _, u := range project.Contributors {
			tbl2.AddRow(u.ID, u.Username)
		}
		tbl2.AddRow("")
		fmt.Println(tbl2)

		var lists []hashstack.List
		if err := getRangeJSON(fmt.Sprintf("/api/projects/%s/lists", args[0]), &lists); err != nil {
			writeStdErrAndExit(err.Error())
		}
		tbl3 := uitable.New()
		fmt.Println("Lists:")
		tbl3.AddRow("ID", "MODE", "NAME", "RECOVERED")
		for _, l := range lists {
			percent := l.RecoveredCount / l.DigestCount
			tbl3.AddRow(l.ID, l.HashMode, l.Name, fmt.Sprintf("%d/%d (%d%%)", l.RecoveredCount, l.DigestCount, percent))
		}
		fmt.Println(tbl3)
	},
}

var delProjectCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "Delete a project by id",
	Long: `
Delete a project by id. Deleting a project will delete any associated lists.
	`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("id is required")
		}
		if err := deleteHTTP(fmt.Sprintf("/api/projects/%s", args[0])); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Printf("[-] project has been deleted")
	},
}

type newProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var newProjectCmd = &cobra.Command{
	Use:   "new <project_name> <description>",
	Short: "Create a new project",
	Long: `
Create a new project using the provided name and description.
The name must be unique across all projects. The name should
not include special characters or spaces. The description is
required.
	`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_name and description are required")
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
		tbl := uitable.New()
		tbl.AddRow("ID:", project.ID)
		tbl.AddRow("Name:", project.Name)
		tbl.AddRow("Description:", project.Description)
		fmt.Println(tbl)
	},
}

func init() {
	projectCmd.AddCommand(newProjectCmd)
	projectCmd.AddCommand(listProjectCmd)
	projectCmd.AddCommand(getProjectCmd)
	projectCmd.AddCommand(delProjectCmd)
	RootCmd.AddCommand(projectCmd)
}
