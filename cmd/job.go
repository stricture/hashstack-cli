package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"sort"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
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

func getTasks(projectID, jobID int64) []hashstack.Task {
	var tasks []hashstack.Task
	path := fmt.Sprintf("/api/projects/%d/jobs/%d/tasks", projectID, jobID)
	if err := getJSON(path, &tasks); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return tasks
}

func displayJob(w io.Writer, job hashstack.Job) {
	status := "Running"
	if job.IsExhausted {
		status = "Complete"
	}
	if !job.IsActive && !job.IsExhausted {
		status = "Paused"
	}
	list := hashstack.List{
		ProjectID: job.ProjectID,
		ID:        job.ListID,
	}
	getListByID(&list)
	mode := getMode(list.HashMode)
	liststat := fmt.Sprintf("%d/%d (%0.2f%%)", list.RecoveredCount, list.DigestCount, percentOf(int(list.RecoveredCount), int(list.DigestCount)))
	firstTime := "has not started"
	if job.FirstTaskTime != 0 {
		firstTime = humanize.Time(time.Unix(job.FirstTaskTime, 0))
	}
	fmt.Fprintf(w, "Name............: %s\n", job.Name)
	fmt.Fprintf(w, "ID..............: %d\n", job.ID)
	fmt.Fprintf(w, "Status..........: %s\n", status)
	fmt.Fprintf(w, "Hash.Type.......: %d (%s)\n", mode.HashMode, mode.Algorithm)
	fmt.Fprintf(w, "Hash.Target.....: %s\n", list.Name)
	fmt.Fprintf(w, "Max Devices.....: %d\n", job.MaxDedicatedDevices)
	fmt.Fprintf(w, "Priority........: %d\n", job.Priority)
	fmt.Fprintf(w, "Time.Created....: %s\n", humanize.Time(time.Unix(job.CreatedAt, 0)))
	fmt.Fprintf(w, "Time.Started....: %s\n", firstTime)
	fmt.Fprintf(w, "Recovered.......: %s\n", liststat)

	tasks := getTasks(job.ProjectID, job.ID)
	var (
		bigTotalSpdCnt = big.NewInt(0)
		bigTotalSpdMs  = big.NewInt(0)
		bigspeed       = big.NewInt(0)
		bigeta         = big.NewInt(0)
		bigprogress    = big.NewInt(0)
	)

	for _, task := range tasks {
		for x, micro := range task.Micros {
			xbig := big.NewInt(micro.Status.SpeedCnt)
			bigTotalSpdCnt.Add(bigTotalSpdCnt, xbig)
			xspdbig := big.NewInt(int64(micro.Status.SpeedMS))
			bigTotalSpdMs.Add(bigTotalSpdMs, xspdbig)
			if x == len(task.Micros)-1 {
				y := big.NewInt(0)
				y.SetString(micro.Status.ProgressTotal, 10)
				bigprogress.Add(bigprogress, y)
			}
		}
	}
	if bigTotalSpdMs.Uint64() != 0 {
		bigspeed = bigTotalSpdMs.Div(bigTotalSpdCnt, bigTotalSpdMs)
	}
	if bigprogress.Uint64() != 0 {
		bigeta.Div(bigprogress, bigspeed)
	}
	bigspeed.Mul(bigspeed, big.NewInt(1000))

	fmt.Fprintf(w, "Speed...........: %s/s\n", humanize.BigComma(bigspeed))
	fmt.Fprintf(w, "ETA.............: %s\n", prettyUptime(bigeta.Int64()))
	fmt.Fprintln(w)
}

func getJob(projectID, jobID int64) hashstack.Job {
	var job hashstack.Job
	path := fmt.Sprintf("/api/projects/%d/jobs/%d", projectID, jobID)
	if err := getJSON(path, &job); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return job
}

func statsJob(job hashstack.Job) {
	displayJob(os.Stdout, job)
	c := time.Tick(5 * time.Second)
	for range c {
		job = getJob(job.ProjectID, job.ID)
		displayJob(os.Stdout, job)
	}
}

func displayJobs(p hashstack.Project) {
	path := fmt.Sprintf("/api/projects/%d/jobs", p.ID)
	var jobs []hashstack.Job
	if err := getRangeJSON(path, &jobs); err != nil {
		writeStdErrAndExit(err.Error())
	}
	if len(jobs) < 1 {
		writeStdErrAndExit("There are no jobs for this project.")
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].CreatedAt > jobs[i].CreatedAt
	})
	for _, j := range jobs {
		displayJob(os.Stdout, j)
		fmt.Println()
	}
}

var jobCmd = &cobra.Command{
	Use:    "jobs <project_name|project_id> [job_id]",
	Short:  "Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands)",
	Long:   "Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands)",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		switch len(args) {
		case 1:
			project := getProject(args[0])
			displayJobs(project)
		case 2:
			i, err := strconv.Atoi(args[1])
			if err != nil {
				writeStdErrAndExit("The provided job_id is not valid.")
			}
			project := getProject(args[0])
			statsJob(getJob(project.ID, int64(i)))
		default:
			cmd.Usage()
		}
	},
}

type updateRequest struct {
	Priority            int  `json:"priority"`
	MaxDedicatedDevices int  `json:"max_dedicated_devices"`
	IsActive            bool `json:"is_active"`
}

var pauseJobCmd = &cobra.Command{
	Use:    "pause <project_name|project_id> <job_id>",
	Short:  "Pauses a job by project_name|project_id and job_id",
	Long:   "Pauses a job by project_name|project_id and job_id",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_name|project_id and job_id are required.")
		}
		project := getProject(args[0])
		i, err := strconv.Atoi(args[1])
		if err != nil {
			writeStdErrAndExit("job_id is invalid")
		}
		job := getJob(project.ID, int64(i))
		update := updateRequest{
			Priority:            job.Priority,
			MaxDedicatedDevices: job.MaxDedicatedDevices,
			IsActive:            false,
		}
		path := fmt.Sprintf("/api/projects/%d/jobs/%d", project.ID, job.ID)
		if _, err := patchJSON(path, &update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("THe job has been paused.")
	},
}

var delJobCmd = &cobra.Command{
	Use:    "delete <project_name|project_id> <job_id>",
	Short:  "Deletes a job by project_name|project_id and job_id",
	Long:   "Deletes a job by project_name|project_id and job_id",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_name|project_id and job_id are required.")
		}
		project := getProject(args[0])
		i, err := strconv.Atoi(args[1])
		if err != nil {
			writeStdErrAndExit("job_id is invalid")
		}
		job := getJob(project.ID, int64(i))
		var attack hashstack.Attack
		if err := getJSON(fmt.Sprintf("/api/attacks/%d", job.AttackID), &attack); err != nil {
			writeStdErrAndExit(err.Error())
		}
		if attack.Title == fmt.Sprintf("hashstack-cli-%d-%d-%s", job.ProjectID, job.ListID, job.Name) {
			deleteHTTP(fmt.Sprintf("/api/attacks/%d", job.AttackID))
		}
		path := fmt.Sprintf("/api/projects/%d/jobs/%d", project.ID, job.ID)
		if err := deleteHTTP(path); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("The job was successfully deleted.")
	},
}

var startJobCmd = &cobra.Command{
	Use:    "start <project_name|project_id> <job_id>",
	Short:  "Starts a job by project_name|project_id and job_id",
	Long:   "Starts a job by project_name|project_id and job_id",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			writeStdErrAndExit("project_name|project_id and job_id are required.")
		}
		project := getProject(args[0])
		i, err := strconv.Atoi(args[1])
		if err != nil {
			writeStdErrAndExit("The job_id provided is not valid.")
		}
		job := getJob(project.ID, int64(i))
		update := updateRequest{
			Priority:            job.Priority,
			MaxDedicatedDevices: job.MaxDedicatedDevices,
			IsActive:            true,
		}
		path := fmt.Sprintf("/api/projects/%d/jobs/%d", project.ID, job.ID)
		if _, err := patchJSON(path, &update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("job has been started")
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

var addJobCmd = &cobra.Command{
	Use:    "add <project_name|project_id> <list_name|list_id> <name> <wordlist|mask>",
	Short:  "Add a job for the provided project and list",
	Long:   "Add a job for the provided project and list",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 4 {
			writeStdErrAndExit("Missing required argument.")
		}
		var (
			name = args[2]
		)
		project := getProject(args[0])
		list := getList(project.ID, args[1])

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
				debug(fmt.Sprintf("Error: %s", err.Error()))
				writeStdErrAndExit("The provided wordlist does not exist on the server.")
			}
			step.WordlistID = wordlistFile.ID
			if flRulesFile != "" {
				if err := getJSON(fmt.Sprintf("/api/rules?filename=%s", flRulesFile), &ruleFile); err != nil {
					debug(fmt.Sprintf("Error: %s", err.Error()))
					writeStdErrAndExit("The provided rule file does not exist on the server.")
				}
				step.RuleID = ruleFile.ID
			}
		case 1:
			if len(args) < 4 {
				writeStdErrAndExit("Two wordlist files are required for a combination attack.")
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
				debug(fmt.Sprintf("Error: %s", err.Error()))
				writeStdErrAndExit("The provided wordlist does not exist on the server.")
			}
			if err := getJSON(fmt.Sprintf("/api/wordlists?filename=%s", args[4]), &combinationFile); err != nil {
				debug(fmt.Sprintf("Error: %s", err.Error()))
				writeStdErrAndExit("The provided combination wordlist does not exist on the server.")
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
				writeStdErrAndExit("A wordlist file and mask are required for this attack mode.")
			}
			var wordlistFile hashstack.File
			if err := getJSON(fmt.Sprintf("/api/wordlists?filename=%s", args[3]), &wordlistFile); err != nil {
				debug(fmt.Sprintf("Error: %s", err.Error()))
				writeStdErrAndExit("The provided wordlist does not exist on the server.")
			}
			step.WordlistID = wordlistFile.ID
			step.Mask = args[4]
		case 7:
			if len(args) < 4 {
				writeStdErrAndExit("A mask and wordlist file are required for this attack mode.")
			}
			var wordlistFile hashstack.File
			if err := getJSON(fmt.Sprintf("/api/wordlists?filename=%s", args[4]), &wordlistFile); err != nil {
				debug(fmt.Sprintf("Error: %s", err.Error()))
				writeStdErrAndExit("The provided wordlist does not exist on the server.")
			}
			step.WordlistID = wordlistFile.ID
			step.Mask = args[3]
		default:
			writeStdErrAndExit("The attack-mode provided is not valid.")
		}
		attack := attackRequest{
			Title: fmt.Sprintf("hashstack-cli-%d-%d-%s", project.ID, list.ID, name),
			Steps: []attackStep{step},
		}
		data, err := postJSON("/api/attacks", &attack)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		debug("uploaded temporary attack plan")
		var plan hashstack.Attack
		if err := json.Unmarshal(data, &plan); err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit(new(jsonServerError).Error())
		}

		jobreq := jobRequest{
			Name:                name,
			ListID:              list.ID,
			AttackID:            plan.ID,
			Priority:            flPriority,
			MaxDedicatedDevices: flMaxDedicatedDevices,
			OpenCLVectorWidth:   flOpenCLVectorWidth,
		}
		data, err = postJSON(fmt.Sprintf("/api/projects/%d/jobs", project.ID), &jobreq)
		if err != nil {
			writeStdErrAndExit(err.Error())
		}
		var job hashstack.Job
		if err := json.Unmarshal(data, &job); err != nil {
			debug(fmt.Sprintf("Error: %s", err.Error()))
			writeStdErrAndExit(new(jsonServerError).Error())
		}
		statsJob(job)
	},
}

func init() {
	addJobCmd.PersistentFlags().IntVarP(&flAttackMode, "attack-mode", "a", 0, "Attack-mode, see references below")
	addJobCmd.PersistentFlags().BoolVar(&flIsHexCharset, "hex-charset", false, "Assume charset is given in hex")
	addJobCmd.PersistentFlags().StringVar(&flMarkovHcstat, "markov-hcstat", "", "Specify hcstat file to use")
	addJobCmd.PersistentFlags().IntVarP(&flMarkovThreshold, "markov-threshold", "t", 0, "Threshold X when to stop accepting new markov-chains")
	addJobCmd.PersistentFlags().IntVar(&flOpenCLVectorWidth, "opencl-vector-width", 0, "Manual override OpenCL  vector-width to X")
	addJobCmd.PersistentFlags().IntVar(&flPriority, "priority", 1, "The priority for this job 1-100")
	addJobCmd.PersistentFlags().IntVar(&flMaxDedicatedDevices, "max-devices", 0, "Maximum devices across the entire cluster to use, 0 is unlimited")
	addJobCmd.PersistentFlags().StringVarP(&flRuleLeft, "rule-left", "j", "", "Single rule applied to each word from left wordlist")
	addJobCmd.PersistentFlags().StringVarP(&flRuleRight, "rule-right", "k", "", "Single rule applied to each word from left wordlist")
	addJobCmd.PersistentFlags().StringVarP(&flRulesFile, "rules-file", "r", "", "Rule file to be applied to each word from wordlists")
	addJobCmd.PersistentFlags().StringVarP(&flCustomCharset1, "custom-charset1", "1", "", "User-defined charset ?1")
	addJobCmd.PersistentFlags().StringVarP(&flCustomCharset2, "custom-charset2", "2", "", "User-defined charset ?2")
	addJobCmd.PersistentFlags().StringVarP(&flCustomCharset3, "custom-charset3", "3", "", "User-defined charset ?3")
	addJobCmd.PersistentFlags().StringVarP(&flCustomCharset4, "custom-charset4", "4", "", "User-defined charset ?4")
	jobCmd.AddCommand(addJobCmd)
	jobCmd.AddCommand(pauseJobCmd)
	jobCmd.AddCommand(startJobCmd)
	jobCmd.AddCommand(delJobCmd)
	RootCmd.AddCommand(jobCmd)
}
