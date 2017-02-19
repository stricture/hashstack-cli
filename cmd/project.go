package cmd

import "github.com/spf13/cobra"

var projectCmd = &cobra.Command{
	Use:    "project",
	Short:  "Returns a list of your projects",
	Long:   "Returns a list of your projects",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var newProjectCmd = &cobra.Command{
	Use:    "new [project_name] [description]",
	Short:  "Create a new project",
	Long:   "Create a new project",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	projectCmd.AddCommand(newProjectCmd)
	RootCmd.AddCommand(projectCmd)
}
