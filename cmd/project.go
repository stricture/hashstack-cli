package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var (
	glDisplayMulti = false
)

func getUser(user *hashstack.User) {
	path := fmt.Sprintf("/api/users/%d", user.ID)
	if err := getJSON(path, user); err != nil {
		writeStdErrAndExit(err.Error())
	}
}

func getListCount(projectID int64) (int, error) {
	path := fmt.Sprintf("/api/projects/%d/lists", projectID)
	return getTotal(path)
}

func getJobs(projectID int64) ([]hashstack.Job, error) {
	path := fmt.Sprintf("/api/projects/%d/jobs", projectID)
	var jobs []hashstack.Job
	err := getRangeJSON(path, &jobs)
	return jobs, err
}

func getProject(arg string) hashstack.Project {
	var p hashstack.Project
	i, err := strconv.Atoi(arg)
	if err != nil {
		p.Name = arg
		getProjectByName(&p)
	} else {
		p.ID = int64(i)
		getProjectByID(&p)
	}
	return p
}

func getProjectByName(p *hashstack.Project) {
	path := fmt.Sprintf("/api/projects?name=%s", p.Name)
	if err := getJSON(path, &p); err != nil {
		writeStdErrAndExit(err.Error())
	}
}

func getProjectByID(p *hashstack.Project) {
	path := fmt.Sprintf("/api/projects/%d", p.ID)
	if err := getJSON(path, p); err != nil {
		writeStdErrAndExit(err.Error())
	}
}

func displayProject(p hashstack.Project) {
	if p.Name == "" && p.ID != 0 {
		getProjectByID(&p)
	}
	if p.ID == 0 {
		getProjectByName(&p)
	}
	jobs, err := getJobs(p.ID)
	if err != nil {
		writeStdErrAndExit(err.Error())
	}
	var (
		isActive    int
		isPaused    int
		isCompleted int
	)
	for _, j := range jobs {
		if j.IsExhausted {
			isCompleted++
		} else if j.IsActive {
			isActive++
		} else {
			isPaused++
		}
	}
	jobstat := fmt.Sprintf("%d/%d Active %d/%d Complete", isActive, isActive+isPaused, isCompleted, len(jobs))
	listCount, err := getListCount(p.ID)
	if err != nil {
		writeStdErrAndExit(err.Error())
	}
	var owner string

	if glDisplayMulti {
		user := hashstack.User{
			ID: p.OwnerUserID,
		}
		getUser(&user)
		owner = user.Username
	} else {
		owner = p.Owner.Username
	}
	fmt.Printf("ID...............: %d\n", p.ID)
	fmt.Printf("Name.............: %s\n", p.Name)
	fmt.Printf("Jobs.............: %s\n", jobstat)
	fmt.Printf("Lists............: %d\n", listCount)
	fmt.Printf("Owner............: %s\n", owner)
	if glDisplayMulti && len(p.Contributors) > 0 {
		var names []string
		for _, u := range p.Contributors {
			names = append(names, u.Username)
		}
		fmt.Printf("Contributors.....: %s\n", strings.Join(names, ", "))
	}
	fmt.Printf("Last Updated.....: %s\n", humanize.Time(time.Unix(p.UpdatedAt, 0)))
	fmt.Println()
}

func displayProjects() {
	var projects []hashstack.Project
	if err := getRangeJSON("/api/projects", &projects); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
	for _, p := range projects {
		displayProject(p)
	}
}

var projectCmd = &cobra.Command{
	Use:   "projects [name|id]",
	Short: "Display a list of all projects (-h or --help for subcommands)",
	Long: `
Displays a list of your projects. If name or id is provided, details will be displayed for that specific project.
Additional subcommands are available.

Projects are used to organize one or more lists. An exmaple of a project may be AcmeForensicIvestigation2016.
Once a project is created, you will upload your lists using the project's name. You can share projects with many users.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			var p hashstack.Project
			i, err := strconv.Atoi(args[0])
			if err != nil {
				p.Name = args[0]
			} else {
				p.ID = int64(i)
			}
			displayProject(p)
			return
		}
		glDisplayMulti = true
		displayProjects()
	},
}

func deleteProject(arg string) {
	var p hashstack.Project
	i, err := strconv.Atoi(arg)
	if err != nil {
		p.Name = arg
		getProjectByName(&p)
	} else {
		p.ID = int64(i)
	}
	path := fmt.Sprintf("/api/projects/%d", p.ID)
	if err := deleteHTTP(path); err != nil {
		writeStdErrAndExit(err.Error())
	}
	fmt.Println("The project was deleted successfully.")
}

var delProjectCmd = &cobra.Command{
	Use:   "delete <name|id>",
	Short: "Delete a project by name or id",
	Long: `
Delete a project by name or id. Deleting a project will delete any associated lists.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("name or id is required.")
		}
		if ok := promptDelete("this project"); !ok {
			writeStdErrAndExit("Not deleting project.")
		}
		deleteProject(args[0])
	},
}

type projectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var addProjectCmd = &cobra.Command{
	Use:   "add <project_name> <description>",
	Short: "Add a new project",
	Long: `
Add a new project using the provided name and description. The name must be unique across all projects.
The name should not include special characters or spaces. The description is required.
	`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_name and description are required.")
		}
		req := projectRequest{
			Name:        args[0],
			Description: args[1],
		}
		body, err := postJSON("/api/projects", req)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		var project hashstack.Project
		if err := json.Unmarshal(body, &project); err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit(new(jsonServerError).Error())
		}
		displayProject(project)
	},
}

func init() {
	projectCmd.AddCommand(addProjectCmd)
	projectCmd.AddCommand(delProjectCmd)
	RootCmd.AddCommand(projectCmd)
}
