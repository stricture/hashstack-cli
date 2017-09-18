package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"os"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Executes admin subcommands (-h or --help for more info).",
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
	Short: "Create an authentication token for another user.",
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
		fmt.Println("\nPlace this token in your configuration file.")
	},
}

var adminDelAgentCmd = &cobra.Command{
	Use:   "agent-delete <agent_id>",
	Short: "Deletes an agent by id.",
	Long: `
Delete an agent by id. This does not prevent the agent from
communicating with the server. This is only useful for when an
agent is brought offline and will never return.
    `,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("agent_id is required")
			return
		}
		id := args[0]
		if err := deleteHTTP(fmt.Sprintf("/api/admin/agents/%s", id)); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("Agent has been deleted.")
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
	Short:  "Displays a list of all projects.",
	Long:   "Displays a list of all projects.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		projects := getAdminProjects()
		displayAdminProjects(projects)
	},
}

func displayAdminJob(job hashstack.Job) {
	w := os.Stdout
	firstTime := "has not started"
	if job.FirstTaskTime != 0 {
		firstTime = humanize.Time(time.Unix(job.FirstTaskTime, 0))
	}
	fmt.Fprintf(w, "Name............: %s\n", job.Name)
	fmt.Fprintf(w, "ID..............: %d\n", job.ID)
	fmt.Fprintf(w, "Project.ID......: %d\n", job.ProjectID)
	fmt.Fprintf(w, "Max.Devices.....: %d\n", job.MaxDedicatedDevices)
	fmt.Fprintf(w, "Priority........: %d\n", job.Priority)
	fmt.Fprintf(w, "Time.Created....: %s\n", humanize.Time(time.Unix(job.CreatedAt, 0)))
	fmt.Fprintf(w, "Time.Started....: %s\n", firstTime)
	fmt.Fprintln(w)
}

func displayAdminJobs(jobs []hashstack.Job) {
	for _, j := range jobs {
		displayAdminJob(j)
	}
}

func getAdminJobs() []hashstack.Job {
	var jobs []hashstack.Job
	if err := getJSON("/api/admin/jobs", &jobs); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].CreatedAt > jobs[i].CreatedAt
	})
	return jobs
}

var adminJobsCmd = &cobra.Command{
	Use:    "jobs",
	Short:  "Displays a list of all active jobs.",
	Long:   "Displays a list of all active jobs.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		jobs := getAdminJobs()
		displayAdminJobs(jobs)
	},
}

func getAdminJob(projectID, jobID int64) hashstack.Job {
	var job hashstack.Job
	path := fmt.Sprintf("/api/admin/projects/%d/jobs/%d", projectID, jobID)
	if err := getJSON(path, &job); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return job
}

var adminPauseJobCmd = &cobra.Command{
	Use:    "job-pause <project_id> <job_id>",
	Short:  "Pauses a job by project_id and job_id.",
	Long:   "Pauses a job by project_id and job_id.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_id and job_id are required.")
		}
		projectID, err := strconv.Atoi(args[0])
		if err != nil {
			writeStdErrAndExit("project_id is invalid")
		}
		jobID, err := strconv.Atoi(args[1])
		if err != nil {
			writeStdErrAndExit("job_id is invalid")
		}
		job := getAdminJob(int64(projectID), int64(jobID))
		update := updateRequest{
			Priority:            job.Priority,
			MaxDedicatedDevices: job.MaxDedicatedDevices,
			IsActive:            false,
		}
		path := fmt.Sprintf("/api/admin/projects/%d/jobs/%d", job.ProjectID, job.ID)
		if _, err := patchJSON(path, &update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("The job has been paused.")
	},
}

var adminDelJobCmd = &cobra.Command{
	Use:    "job-delete <project_id> <job_id>",
	Short:  "Deletes a job by project_id and job_id.",
	Long:   "Deletes a job by project_id and job_id.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_id and job_id are required.")
		}
		projectID, err := strconv.Atoi(args[0])
		if err != nil {
			writeStdErrAndExit("project_id is invalid")
		}
		jobID, err := strconv.Atoi(args[1])
		if err != nil {
			writeStdErrAndExit("job_id is invalid")
		}
		job := getAdminJob(int64(projectID), int64(jobID))
		var attack hashstack.Attack
		if err := getJSON(fmt.Sprintf("/api/attacks/%d", job.AttackID), &attack); err != nil {
			writeStdErrAndExit(err.Error())
		}
		if ok := promptDelete("this job"); !ok {
			writeStdErrAndExit("Not deleting job.")
		}
		path := fmt.Sprintf("/api/admin/projects/%d/jobs/%d", job.ProjectID, job.ID)
		if err := deleteHTTP(path); err != nil {
			writeStdErrAndExit(err.Error())
		}
		if attack.Title == fmt.Sprintf("hashstack-cli-%d-%d-%s", job.ProjectID, job.ListID, job.Name) {
			deleteHTTP(fmt.Sprintf("/api/attacks/%d", job.AttackID))
		}
		fmt.Println("The job was successfully deleted.")
	},
}

var adminStartJobCmd = &cobra.Command{
	Use:    "job-start <project_id> <job_id>",
	Short:  "Starts a job by project_id and job_id.",
	Long:   "Starts a job by project_id and job_id.",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_id and job_id are required.")
		}
		projectID, err := strconv.Atoi(args[0])
		if err != nil {
			writeStdErrAndExit("project_id is invalid")
		}
		jobID, err := strconv.Atoi(args[1])
		if err != nil {
			writeStdErrAndExit("job_id is invalid")
		}
		job := getAdminJob(int64(projectID), int64(jobID))
		update := updateRequest{
			Priority:            job.Priority,
			MaxDedicatedDevices: job.MaxDedicatedDevices,
			IsActive:            true,
		}
		path := fmt.Sprintf("/api/admin/projects/%d/jobs/%d", job.ProjectID, job.ID)
		if _, err := patchJSON(path, &update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("The job has been started.")
	},
}

func getAdminTeams() []hashstack.Team {
	var teams []hashstack.Team
	if err := getJSON("/api/admin/teams", &teams); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	return teams
}

var adminTeamCmd = &cobra.Command{
	Use:   "teams",
	Short: "Display a list of all teams.",
	Long: `
Displays a list of all teams.

Teams are used to provide access to projects by adding and remove indiviual users from a team.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		displayTeams(getAdminTeams())
	},
}
var adminUpdateTeamCmd = &cobra.Command{
	Use:   "update-team <name|id>",
	Short: "Modifies a team by it's name or id.",
	Long: `
Modifies a team by it's name or id. All values are optional and only provided options will be modified.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("name or id is reuqired.")
		}
		var team hashstack.Team
		teamID, err := strconv.Atoi(args[0])
		if err != nil {
			team.Name = args[0]
		} else {
			team.ID = int64(teamID)
		}
		team = getTeam(team)
		var contribs []updateContributor
		for _, c := range team.Contributors {
			contribs = append(contribs, updateContributor{ID: c.ID})
		}
		update := updateTeamReq{
			Name:         team.Name,
			Description:  team.Description,
			IsActive:     true,
			OwnerUserID:  team.OwnerUserID,
			Contributors: contribs,
		}
		if flTeamName != "" {
			update.Name = flTeamName
		}
		if flTeamDescription != "" {
			update.Description = flTeamDescription
		}
		if flTeamOwnerUsername != "" {
			user := hashstack.User{
				Username: flTeamOwnerUsername,
			}
			getUser(&user)
			update.OwnerUserID = user.ID
		}
		if flTeamOwnerUserID != 0 {
			user := hashstack.User{
				ID: int64(flTeamOwnerUserID),
			}
			getUser(&user)
			update.OwnerUserID = user.ID
		}
		data, err := patchJSON(fmt.Sprintf("/api/teams/%d", team.ID), update)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		if err := json.Unmarshal(data, &team); err != nil {
			writeStdErrAndExit(err.Error())
		}
		displayTeam(false, team)
	},
}

func init() {
	adminCmd.AddCommand(adminImpersonateCmd)
	adminCmd.AddCommand(adminProjectsCmd)
	adminCmd.AddCommand(adminJobsCmd)
	adminCmd.AddCommand(adminPauseJobCmd)
	adminCmd.AddCommand(adminStartJobCmd)
	adminCmd.AddCommand(adminDelJobCmd)
	adminCmd.AddCommand(adminDelAgentCmd)
	adminUpdateTeamCmd.PersistentFlags().StringVar(&flTeamName, "name", "", "Sets team's name")
	adminUpdateTeamCmd.PersistentFlags().StringVar(&flTeamDescription, "description", "", "Sets the team's description")
	adminUpdateTeamCmd.PersistentFlags().StringVar(&flTeamOwnerUsername, "owner-name", "", "Sets the team's owner to the username provided")
	adminUpdateTeamCmd.PersistentFlags().IntVar(&flTeamOwnerUserID, "owner-id", 0, "Sets the team's owner to the id provided")
	adminCmd.AddCommand(adminUpdateTeamCmd)
	adminCmd.AddCommand(adminTeamCmd)
	RootCmd.AddCommand(adminCmd)
}
