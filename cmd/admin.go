package cmd

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Executes admin subcommands (-h or --help for more info)",
	Long: `
Executes a subcommand as an administrator (-h or --help for a list of subcommands.
    `,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

type impersonateRequest struct {
	Username string `json:"username"`
}

func getImpersonationToken(username string) string {
	body, err := postJSON("/api/admin/impersonate", impersonateRequest{
		Username: username,
	})
	if err != nil {
		writeStdErrAndExit(err.Error())
	}
	var response tokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		debug(fmt.Sprintf("Error: %s", err.Error()))
		writeStdErrAndExit(new(invalidResponseError).Error())
	}
	return response.Token
}

var adminImpersonateCmd = &cobra.Command{
	Use:   "impersonate <username>",
	Short: "Create an authentication token for another user",
	Long: `
Create an authentication token for another user by username. This
can be useful when troubleshooting as another user.
    `,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("username is required")
			return
		}
		username := args[0]
		token := getImpersonationToken(username)
		fmt.Printf("Authentication token for %s: %s\n", username, token)
	},
}

func displayAdminProject(project hashstack.Project) {
	fmt.Printf("ID........: %d\n", project.ID)
	fmt.Printf("Name......: %s\n", project.Name)
	user := hashstack.User{
		ID: project.OwnerUserID,
	}
	getUser(&user)
	fmt.Printf("Owner.....: %s\n", user.Username)
	fmt.Println()
}

func displayAdminProjects(projects []hashstack.Project) {
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].OwnerUserID < projects[j].OwnerUserID
	})
	for _, p := range projects {
		displayAdminProject(p)
	}
}

func getAdminProjects() []hashstack.Project {
	var projects []hashstack.Project
	if err := getJSON("/api/admin/projects", &projects); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return projects
}

var adminProjectsCmd = &cobra.Command{
	Use:    "projects",
	Short:  "Displays a list of all projects",
	Long:   "Displays a list of all projects",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		projects := getAdminProjects()
		displayAdminProjects(projects)
	},
}

func init() {
	adminCmd.AddCommand(adminImpersonateCmd)
	adminCmd.AddCommand(adminProjectsCmd)
	RootCmd.AddCommand(adminCmd)
}