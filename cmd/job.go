package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	flHashMode            int
	flAttackMode          int
	flIsHexCharset        bool
	flMarkovHcstat        string
	flMarkovThreshold     int
	flOpenCLVectorWidth   int
	flPriority            int
	flMaxDedicatedDevices int
	flRuleLeft            string
	flRuleRight           string
	flRulesFile           string
	flCustomCharset1      string
	flCustomCharset2      string
	flCustomCharset3      string
	flCustomCharset4      string
)
var jobCmd = &cobra.Command{
	Use:    "jobs",
	Short:  "Subcommands can be used to interact with hashstack jobs",
	Long:   "Subcommands can be used to interact with hashstack jobs",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("try hashstack-cli lists -h")
	},
}

var newJobCmd = &cobra.Command{
	Use:    "new project_id list_id [dictionary|mask]",
	Short:  "Create a new job for the provided project and list",
	Long:   "Create a new job for the provided project and list",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	newJobCmd.PersistentFlags().IntVarP(&flHashMode, "hash-type", "m", 0, "Hash-type, see hashstack-cli modes for references")
	newJobCmd.PersistentFlags().IntVarP(&flAttackMode, "attack-mode", "a", 0, "Attack-mode, see references below")
	newJobCmd.PersistentFlags().BoolVar(&flIsHexCharset, "hex-charset", false, "Assume charset is given in hex")
	newJobCmd.PersistentFlags().StringVar(&flMarkovHcstat, "markov-hcstat", "", "Specify hcstat file to use")
	newJobCmd.PersistentFlags().IntVarP(&flMarkovThreshold, "markov-threshold", "t", 0, "Threshold X when to stop accepting new markov-chains")
	newJobCmd.PersistentFlags().IntVar(&flOpenCLVectorWidth, "opencl-vector-width", 0, "Manual override OpenCL  vector-width to X")
	newJobCmd.PersistentFlags().IntVar(&flPriority, "priority", 1, "The priority for this job 1-100")
	newJobCmd.PersistentFlags().IntVar(&flMaxDedicatedDevices, "max-devices", 0, "Maximum devices across the entire cluster to use, 0 is unlimited")
	newJobCmd.PersistentFlags().StringVarP(&flRuleLeft, "rule-left", "j", "", "Single rule applied to each word from left wordlist")
	newJobCmd.PersistentFlags().StringVarP(&flRuleRight, "rule-right", "k", "", "Single rule applied to each word from left wordlist")
	newJobCmd.PersistentFlags().StringVarP(&flRulesFile, "rules-file", "r", "", "Rule file to be applied to each word from wordlists")
	newJobCmd.PersistentFlags().StringVarP(&flCustomCharset1, "custom-charset1", "1", "", "User-defined charset ?1")
	newJobCmd.PersistentFlags().StringVarP(&flCustomCharset2, "custom-charset2", "2", "", "User-defined charset ?2")
	newJobCmd.PersistentFlags().StringVarP(&flCustomCharset3, "custom-charset3", "3", "", "User-defined charset ?3")
	newJobCmd.PersistentFlags().StringVarP(&flCustomCharset4, "custom-charset4", "4", "", "User-defined charset ?4")
	jobCmd.AddCommand(newJobCmd)
	RootCmd.AddCommand(jobCmd)
}
