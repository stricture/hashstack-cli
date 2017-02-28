package cmd

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

var (
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

type updateRequest struct {
	Priority            int  `json:"priority"`
	MaxDedicatedDevices int  `json:"max_dedicated_devices"`
	IsActive            bool `json:"is_active"`
}

var pauseJobCmd = &cobra.Command{
	Use:    "pause <project_id> <job_id>",
	Short:  "Pauses a job by setting \"is_active\" to false",
	Long:   "Pauses a job by setting \"is_active\" to false",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_id and job_id are required")
		}
		var job hashstack.Job
		path := fmt.Sprintf("/api/projects/%s/jobs/%s", args[0], args[1])
		if err := getJSON(path, &job); err != nil {
			writeStdErrAndExit(err.Error())
		}
		update := updateRequest{
			Priority:            job.Priority,
			MaxDedicatedDevices: job.MaxDedicatedDevices,
			IsActive:            false,
		}
		if _, err := patchJSON(path, &update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("[+] Job paused")
	},
}

var delJobCmd = &cobra.Command{
	Use:    "del <project_id> <job_id>",
	Short:  "Deletes a job by id",
	Long:   "Deletes a job by id",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_id and job_id are required")
		}

		path := fmt.Sprintf("/api/projects/%s/jobs/%s", args[0], args[1])
		var job hashstack.Job
		if err := getJSON(path, &job); err != nil {
			writeStdErrAndExit(err.Error())
		}
		var attack hashstack.Attack
		if err := getJSON(fmt.Sprintf("/api/attacks/%d", job.AttackID), &attack); err != nil {
			writeStdErrAndExit(err.Error())
		}
		if attack.Title == fmt.Sprintf("hashstack-cli-%d-%d-%s", job.ProjectID, job.ListID, job.Name) {
			deleteHTTP(fmt.Sprintf("/api/attacks/%d", job.AttackID))
		}
		if err := deleteHTTP(path); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("[+] Job deleted")
	},
}

var startJobCmd = &cobra.Command{
	Use:    "start <project_id> <job_id>",
	Short:  "Starts a job by setting \"is_active\" to true",
	Long:   "Starts a job by setting \"is_active\" to true",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_id and job_id are required")
		}
		var job hashstack.Job
		path := fmt.Sprintf("/api/projects/%s/jobs/%s", args[0], args[1])
		if err := getJSON(path, &job); err != nil {
			writeStdErrAndExit(err.Error())
		}
		update := updateRequest{
			Priority:            job.Priority,
			MaxDedicatedDevices: job.MaxDedicatedDevices,
			IsActive:            true,
		}
		if _, err := patchJSON(path, &update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("[+] Job started")
	},
}

type jobRequest struct {
	Name                string `json:"name"`
	ListID              int64  `json:"list_id"`
	AttackID            int64  `json:"attack_id"`
	Priority            int    `json:"priority"`
	MaxDedicatedDevices int    `json:"max_dedicated_devices"`
	OpenCLVectorWidth   int    `json:"opencl_vector_width"`
}

type attackStep struct {
	IDX                   int    `json:"idx"`
	AttackMode            int    `json:"attack_mode"`
	WordlistID            int64  `json:"wordlist_id"`
	WordlistCombinationID int64  `json:"wordlist_combination_id"`
	RuleID                int64  `json:"rule_id"`
	RuleBufLeft           string `json:"rule_buf_left"`
	RuleBufRight          string `json:"rule_buf_right"`
	Mask                  string `json:"mask"`
	IsHexCharset          bool   `json:"is_hex_charset"`
	MarkovThreshold       int    `json:"markov_threshold"`
	MarkovHCStatFileID    int64  `json:"markov_hc_stat_file_id"`
	CustomCharset1        string `json:"custom_charset1"`
	CustomCharset2        string `json:"custom_charset2"`
	CustomCharset3        string `json:"custom_charset3"`
	CustomCharset4        string `json:"custom_charset4"`
}

type attackRequest struct {
	Title string       `json:"title"`
	Steps []attackStep `json:"steps"`
}

var newJobCmd = &cobra.Command{
	Use:    "new <project_id> <list_id> <name> <dictionary|mask>",
	Short:  "Create a new job for the provided project and list",
	Long:   "Create a new job for the provided project and list",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 4 {
			writeStdErrAndExit("missing required arugment")
		}
		var (
			projectIDStr = args[0]
			listIDStr    = args[1]
			name         = args[2]
		)
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			writeStdErrAndExit("invalid project_id")
		}
		listID, err := strconv.Atoi(listIDStr)
		if err != nil {
			writeStdErrAndExit("invalid list_id")
		}
		step := attackStep{
			IDX:        0,
			AttackMode: flAttackMode,
		}
		switch flAttackMode {
		case 0:
			var (
				wordlistFile hashstack.File
				ruleFile     hashstack.File
			)
			if err := getJSON(fmt.Sprintf("/api/wordlists?filename=%s", args[3]), &wordlistFile); err != nil {
				debug(err.Error())
				writeStdErrAndExit("provided dictionary does not exist on the server")
			}
			step.WordlistID = wordlistFile.ID
			if flRulesFile != "" {
				if err := getJSON(fmt.Sprintf("/api/rules?filename=%s", flRulesFile), &ruleFile); err != nil {
					debug(err.Error())
					writeStdErrAndExit("provided rule file does not exist on the server")
				}
				step.RuleID = ruleFile.ID
			}
		case 1:
			if len(args) < 4 {
				writeStdErrAndExit("two dictionary files are required for a combination attack")
			}
			if flRuleLeft != "" {
				step.RuleBufLeft = flRuleLeft
			}
			if flRuleRight != "" {
				step.RuleBufRight = flRuleRight
			}
			var (
				wordlistFile    hashstack.File
				combinationFile hashstack.File
			)
			if err := getJSON(fmt.Sprintf("/api/wordlists?filename=%s", args[3]), &wordlistFile); err != nil {
				debug(err.Error())
				writeStdErrAndExit("provided dictionary does not exist on the server")
			}
			if err := getJSON(fmt.Sprintf("/api/wordlists?filename=%s", args[4]), &combinationFile); err != nil {
				debug(err.Error())
				writeStdErrAndExit("provided combination dictionary does not exist on the server")
			}
			step.WordlistID = wordlistFile.ID
			step.WordlistCombinationID = combinationFile.ID

		case 3:
			step.Mask = args[3]
			step.CustomCharset1 = flCustomCharset1
			step.CustomCharset2 = flCustomCharset2
			step.CustomCharset3 = flCustomCharset3
			step.CustomCharset4 = flCustomCharset4
			step.IsHexCharset = flIsHexCharset
		case 6:
			if len(args) < 4 {
				writeStdErrAndExit("a dictionary file and mask are required for this attack mode")
			}
			var wordlistFile hashstack.File
			if err := getJSON(fmt.Sprintf("/api/wordlists?filename=%s", args[3]), &wordlistFile); err != nil {
				debug(err.Error())
				writeStdErrAndExit("provided dictionary does not exist on the server")
			}
			step.WordlistID = wordlistFile.ID
			step.Mask = args[4]
		case 7:
			if len(args) < 4 {
				writeStdErrAndExit("a mask and dictionary file are required for this attack mode")
			}
			var wordlistFile hashstack.File
			if err := getJSON(fmt.Sprintf("/api/wordlists?filename=%s", args[4]), &wordlistFile); err != nil {
				debug(err.Error())
				writeStdErrAndExit("provided dictionary does not exist on the server")
			}
			step.WordlistID = wordlistFile.ID
			step.Mask = args[3]
		default:
			writeStdErrAndExit("invalid attack-mode")
		}
		attack := attackRequest{
			Title: fmt.Sprintf("hashstack-cli-%d-%d-%s", projectID, listID, name),
			Steps: []attackStep{step},
		}
		data, err := postJSON("/api/attacks", &attack)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		debug("uploaded temporary attack plan")
		var plan hashstack.Attack
		if err := json.Unmarshal(data, &plan); err != nil {
			debug(err.Error())
			writeStdErrAndExit("there was an error decoding a response from the server")
		}

		job := jobRequest{
			Name:                name,
			ListID:              int64(listID),
			AttackID:            plan.ID,
			Priority:            flPriority,
			MaxDedicatedDevices: flMaxDedicatedDevices,
			OpenCLVectorWidth:   flOpenCLVectorWidth,
		}
		data, err = postJSON(fmt.Sprintf("/api/projects/%s/jobs", projectIDStr), &job)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println(string(data))
	},
}

func init() {
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
	jobCmd.AddCommand(pauseJobCmd)
	jobCmd.AddCommand(startJobCmd)
	jobCmd.AddCommand(delJobCmd)
	RootCmd.AddCommand(jobCmd)
}
