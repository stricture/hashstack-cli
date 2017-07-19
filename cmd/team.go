package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var (
	flTeamName          string
	flTeamDescription   string
	flTeamOwnerUsername string
	flTeamOwnerUserID   int
)

func displayTeam(isMulti bool, team hashstack.Team) {
	user := hashstack.User{
		ID: team.OwnerUserID,
	}
	getUser(&user)
	fmt.Printf("ID...............: %d\n", team.ID)
	fmt.Printf("Name.............: %s\n", team.Name)
	fmt.Printf("Owner............: %s\n", user.Username)
	if isMulti {
		return
	}
	var names []string
	for _, u := range team.Contributors {
		names = append(names, u.Username)
	}
	fmt.Printf("Description......: %s\n", team.Description)
	fmt.Printf("Members..........: %s\n", strings.Join(names, ", "))
}

func displayTeams(teams []hashstack.Team) {
	for _, t := range teams {
		displayTeam(true, t)
		fmt.Println()
	}
}

func getTeams() []hashstack.Team {
	var teams []hashstack.Team
	if err := getRangeJSON("/api/teams", &teams); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	return teams
}

func getTeam(team hashstack.Team) hashstack.Team {
	var path string
	if team.Name == "" && team.ID != 0 {
		path = fmt.Sprintf("/api/teams/%d", team.ID)
	} else {
		path = fmt.Sprintf("/api/teams?name=%s", team.Name)
	}
	if err := getJSON(path, &team); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return team
}

var teamCmd = &cobra.Command{
	Use:   "teams [name|id]",
	Short: "Display a list of all teams (-h or --help for subcommands).",
	Long: `
Displays a list of your teams. If a name or id is provided, details will be displayed for that specific team.
Additional subcommands are available.addHCStatCmd

Teams are used to provide access to projects by adding and remove indiviual users from a team.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			var team hashstack.Team
			i, err := strconv.Atoi(args[0])
			if err != nil {
				team.Name = args[0]
			} else {
				team.ID = int64(i)
			}
			displayTeam(false, getTeam(team))
			return
		}
		displayTeams(getTeams())
	},
}

type teamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var addTeamCmd = &cobra.Command{
	Use:   "add <name> <description>",
	Short: "Add a new team with the provided name and description.",
	Long: `
Add a new team with the provided name and description. Teams can be used to control
access to projects.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("name and description are required.")
		}

		req := teamRequest{
			Name:        args[0],
			Description: args[1],
		}
		body, err := postJSON("/api/teams", req)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		var team hashstack.Team
		if err := json.Unmarshal(body, &team); err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit(new(jsonServerError).Error())
		}
		displayTeam(false, team)
	},
}

var delTeamCmd = &cobra.Command{
	Use:   "delete <name|id>",
	Short: "Delete a team by name or id.",
	Long: `
Delete a team by name or id.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			writeStdErrAndExit("name or id is reuqired.")
		}
		if ok := promptDelete("this team"); !ok {
			writeStdErrAndExit("Not deleting team.")
		}
		var team hashstack.Team
		teamID, err := strconv.Atoi(args[0])
		if err != nil {
			team.Name = args[0]
		} else {
			team.ID = int64(teamID)
		}
		team = getTeam(team)
		path := fmt.Sprintf("/api/teams/%d", team.ID)
		if err := deleteHTTP(path); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("The team was deleted successfully.")
	},
}

type updateTeamReq struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsActive     bool   `json:"is_active"`
	OwnerUserID  int64  `json:"owner_user_id"`
	Contributors []updateContributor
}

var updateTeamCmd = &cobra.Command{
	Use:   "update <name|id>",
	Short: "Modifies a team with by it's name or id.",
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

var addMemberTeamCmd = &cobra.Command{
	Use:   "add-member <team_name|team_id> <username|user_id>",
	Short: "Adds a user to a team by the team's name or id and the user's username or id.",
	Long: `
Adds a user to a team by the team's name or id and the user's username.
	`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("team id or name and user username or id are required.")
		}
		var team hashstack.Team
		teamID, err := strconv.Atoi(args[0])
		if err != nil {
			team.Name = args[0]
		} else {
			team.ID = int64(teamID)
		}
		team = getTeam(team)
		var user hashstack.User
		userID, err := strconv.Atoi(args[1])
		if err != nil {
			user.Username = args[1]
		} else {
			user.ID = int64(userID)
		}
		getUser(&user)
		team.Contributors = append(team.Contributors, user)
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
		if _, err := patchJSON(fmt.Sprintf("/api/teams/%d", team.ID), update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("User has been added to the team.")
	},
}

var removeMemberTeamCmd = &cobra.Command{
	Use:   "remove-member <team_name|team_id> <username>",
	Short: "Removes a user from a team by the team's name or id and the user's username.",
	Long: `
Removes a user from a team by the team's name or id and the user's username.
	`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("team id or name and user username or id are required.")
		}
		var team hashstack.Team
		teamID, err := strconv.Atoi(args[0])
		if err != nil {
			team.Name = args[0]
		} else {
			team.ID = int64(teamID)
		}
		team = getTeam(team)
		var user hashstack.User
		userID, err := strconv.Atoi(args[1])
		if err != nil {
			user.Username = args[1]
		} else {
			user.ID = int64(userID)
		}
		getUser(&user)
		var contribs []updateContributor
		for _, c := range team.Contributors {
			if c.ID == user.ID {
				continue
			}
			contribs = append(contribs, updateContributor{ID: c.ID})
		}
		update := updateTeamReq{
			Name:         team.Name,
			Description:  team.Description,
			OwnerUserID:  team.OwnerUserID,
			IsActive:     true,
			Contributors: contribs,
		}
		if _, err := patchJSON(fmt.Sprintf("/api/teams/%d", team.ID), update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("User has been removed from the team.")
	},
}

func init() {
	updateTeamCmd.PersistentFlags().StringVar(&flTeamName, "name", "", "Sets team's name")
	updateTeamCmd.PersistentFlags().StringVar(&flTeamDescription, "description", "", "Sets the team's description")
	updateTeamCmd.PersistentFlags().StringVar(&flTeamOwnerUsername, "owner-name", "", "Sets the team's owner to the username provided")
	updateTeamCmd.PersistentFlags().IntVar(&flTeamOwnerUserID, "owner-id", 0, "Sets the team's owner to the id provided")
	teamCmd.AddCommand(addTeamCmd)
	teamCmd.AddCommand(delTeamCmd)
	teamCmd.AddCommand(updateTeamCmd)
	teamCmd.AddCommand(addMemberTeamCmd)
	teamCmd.AddCommand(removeMemberTeamCmd)
	RootCmd.AddCommand(teamCmd)
}
