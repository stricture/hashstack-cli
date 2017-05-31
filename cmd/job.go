package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/signal"
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

func getEvents(projectID, jobID int64) []hashstack.AgentEvent {
	var events []hashstack.AgentEvent
	path := fmt.Sprintf("/api/projects/%d/jobs/%d/events", projectID, jobID)
	if err := getJSON(path, &events); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return events
}

func getTasks(projectID, jobID int64) []hashstack.Task {
	var tasks []hashstack.Task
	path := fmt.Sprintf("/api/projects/%d/jobs/%d/tasks", projectID, jobID)
	for {
		if err := getJSON(path, &tasks); err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	return tasks
}

var (
	agentEventTrackTime = time.Now().Unix()
	jobListCrackedCount int64
)

func displayJob(w io.Writer, job hashstack.Job) {
	status := "Running"
	if job.IsExhausted {
		status = "Finished"
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
	tasks := getTasks(job.ProjectID, job.ID)
	events := getEvents(job.ProjectID, job.ID)
	for _, e := range events {
		if e.CreatedAt >= agentEventTrackTime {
			agentEventTrackTime = time.Now().Unix()
			fmt.Fprintf(w, "There was an error returned from an agent: %s!\n\n", e.Buffer)
		}
	}

	var (
		bigTotalSpdCnt        = big.NewInt(0)
		bigTotalSpdMs         = big.NewInt(0)
		bigspeed              = big.NewInt(0)
		bigeta                = big.NewInt(0)
		bigkeyspace           = big.NewInt(0)
		bigkeyspacecomplete   = big.NewInt(0)
		bigkeyspaceinprogress = big.NewInt(0)
		activeDevices         = 0
	)

	debug(fmt.Sprintf("task length: %d", len(tasks)))

	for _, task := range tasks {
		debug(fmt.Sprintf("micro length %d", len(task.Micros)))
		var (
			tsktotal      = big.NewInt(0)
			tskcomplete   = big.NewInt(0)
			tskinprogress = big.NewInt(0)
			tskmod        = big.NewInt(0)
		)
		tsktotal.SetString(task.Keyspace, 10)
		tskcomplete.SetString(task.KeyspaceCompleted, 10)
		tskinprogress.SetString(task.KeyspaceInProgress, 10)
		tskmod.SetString(task.Modifier, 10)
		tsktotal.Mul(tsktotal, tskmod)
		tskcomplete.Mul(tskcomplete, tskmod)
		tskinprogress.Mul(tskinprogress, tskmod)

		bigkeyspace.Add(bigkeyspace, tsktotal)
		bigkeyspacecomplete.Add(bigkeyspacecomplete, tskcomplete)
		bigkeyspaceinprogress.Add(bigkeyspaceinprogress, tskinprogress)

		for x, micro := range task.Micros {
			if time.Now().Add(-2*time.Minute).Unix() > micro.Status.UpdatedAt {
				debug(fmt.Sprintf("micro.id %d is stale", micro.ID))
				continue
			}
			activeDevices++
			xbig := big.NewInt(micro.Status.SpeedCnt)
			bigTotalSpdCnt.Add(bigTotalSpdCnt, xbig)
			debug(fmt.Sprintf("big_speed set to %s", bigTotalSpdCnt.String()))
			xspdbig := big.NewInt(int64(micro.Status.SpeedMS))
			bigTotalSpdMs.Add(bigTotalSpdMs, xspdbig)
			if x == len(task.Micros)-1 {
				bigTotalSpdMs.Div(bigTotalSpdMs, big.NewInt(int64(len(task.Micros))))
			}
		}
	}

	if bigTotalSpdMs.Uint64() != 0 {
		bigspeed = bigTotalSpdMs.Div(bigTotalSpdCnt, bigTotalSpdMs)
		bigspeed.Mul(bigspeed, big.NewInt(1000))
	}

	if bigkeyspace.Uint64() != 0 && bigspeed.Uint64() != 0 {
		debug("speed /s: " + bigspeed.String())
		debug("keyspace: " + bigkeyspace.String())
		bigeta.Div(bigkeyspace, bigspeed)
		debug("eta seconds: " + bigeta.String())
		since := time.Now().Unix() - job.FirstTaskTime
		if bigeta.Int64() > since {
			bigeta.Sub(bigeta, big.NewInt(since))
		}
	}

	var (
		timeETA      = "Undetermined"
		timeStarted  = "Has not started"
		timeFinished = "Unknown"
		timeCreated  string
	)
	if bigeta.Int64() > 0 {
		etaUnix := time.Now().Add(time.Duration(bigeta.Int64()) * time.Second)
		timeETA = fmt.Sprintf("%s (%s)", etaUnix.Format(time.UnixDate), humanize.Time(etaUnix))
	}

	if job.FirstTaskTime != 0 {
		firstTaskUnix := time.Unix(job.FirstTaskTime, 0)
		timeStarted = fmt.Sprintf("%s (%s)", firstTaskUnix.Format(time.UnixDate), humanize.Time(firstTaskUnix))
	}

	if job.LastTaskTime != 0 {
		lastTaskUnix := time.Unix(job.LastTaskTime, 0)
		timeFinished = fmt.Sprintf("%s (%s)", lastTaskUnix.Format(time.UnixDate), humanize.Time(lastTaskUnix))
	}

	createdAtUnix := time.Unix(job.CreatedAt, 0)
	timeCreated = fmt.Sprintf("%s (%s)", createdAtUnix.Format(time.UnixDate), humanize.Time(createdAtUnix))

	strspeed := formatHashRate(bigspeed.Uint64())
	liststat := fmt.Sprintf("%d/%d (%0.2f%%) hashes", list.RecoveredCount, list.DigestCount, percentOf(int(list.RecoveredCount), int(list.DigestCount)))

	fmt.Fprintf(w, "Job.ID..............: %d\n", job.ID)
	fmt.Fprintf(w, "Job.Priority........: %d\n", job.Priority)
	fmt.Fprintf(w, "Job.Name............: %s\n", job.Name)
	fmt.Fprintf(w, "Job.Status..........: %s\n", status)
	fmt.Fprintf(w, "Job.Cracked.........: %s\n", liststat)
	fmt.Fprintf(w, "Job.Progress........: %s/%s (%0.2f%%)\n", bigkeyspacecomplete.String(), bigkeyspace.String(), bigPercentOf(bigkeyspacecomplete, bigkeyspace))
	fmt.Fprintf(w, "Job.Errors..........: %d errors\n", len(events))
	fmt.Fprintf(w, "Hash.Mode...........: %d (%s)\n", mode.HashMode, mode.Algorithm)
	fmt.Fprintf(w, "Hash.Target.........: %s\n", list.Name)
	fmt.Fprintf(w, "Time.Created........: %s\n", timeCreated)
	fmt.Fprintf(w, "Time.Started........: %s\n", timeStarted)
	if status != "Running" {
		fmt.Fprintf(w, "Time.Finished.......: %s\n", timeFinished)
		return
	}
	fmt.Fprintf(w, "Time.Estimated......: %s\n", timeETA)
	fmt.Fprintf(w, "Device.Max..........: %d\n", job.MaxDedicatedDevices)
	fmt.Fprintf(w, "Device.Active.......: %d\n", activeDevices)
	fmt.Fprintf(w, "Device.Speed........: %s\n", strspeed)

	if jobListCrackedCount == 0 && list.RecoveredCount != 0 {
		jobListCrackedCount = list.RecoveredCount
	}
	if list.RecoveredCount > jobListCrackedCount {
		fmt.Fprintf(w, "Use 'hashstack lists cracked %d %d' to view the %d passwords that were cracked\n\n", job.ProjectID, job.ListID, (list.RecoveredCount - jobListCrackedCount))
		jobListCrackedCount = list.RecoveredCount
	}
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
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		fmt.Println("Interrupt caught. Job will continue to run on the server.")
		os.Exit(0)
	}()
	c := time.Tick(5 * time.Second)
	for range c {
		job = getJob(job.ProjectID, job.ID)
		displayJob(os.Stdout, job)
		fmt.Fprintf(os.Stdout, "\nCtrl-C to exit. Job will continue to run.\n\n")
		if job.IsExhausted {
			fmt.Printf("The job is finished. View stats using 'hashstack jobs %d %d'\n\n", job.ProjectID, job.ID)
			break
		}
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
		return jobs[i].CreatedAt < jobs[j].CreatedAt
	})

	for _, j := range jobs {
		displayJob(os.Stdout, j)
		fmt.Println()
	}
}

var jobCmd = &cobra.Command{
	Use:   "jobs [project_name|project_id] [job_id]",
	Short: "Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).",
	Long: `
Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands). If no project is
provided, then all jobs for all projects will be displayed.
`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			projects := getProjects()
			for _, project := range projects {
				fmt.Printf("Project.ID......: %d\n", project.ID)
				fmt.Printf("Project.Name....: %s\n", project.Name)
				displayJobs(project)
			}
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
	Short:  "Pauses a job by project_name|project_id and job_id.",
	Long:   "Pauses a job by project_name|project_id and job_id.",
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
		fmt.Println("The job has been paused.")
	},
}

var errorJobCmd = &cobra.Command{
	Use:   "errors <project_name|project_id> <job_id>.",
	Short: "Displays errors for a job by project_name|project_id and job_id.",
	Long: `
Displays errors for a job by project_name|project_id and job_id.
	`,
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
		events := getEvents(project.ID, job.ID)
		agentMap := make(map[int64]hashstack.Agent)
		for _, e := range events {
			agent, ok := agentMap[e.AgentID]
			if !ok {
				agent = getAgent(e.AgentID)
				agentMap[e.AgentID] = agent
			}
			fmt.Printf("Agent.ID..............: %d\n", agent.ID)
			fmt.Printf("Agent.Host............: %s\n", agent.Hostname)
			fmt.Printf("Agent.IP.Address......: %s\n", agent.IPAddress)
			fmt.Printf("Error.Time............: %s\n", humanize.Time(time.Unix(e.UpdatedAt, 0)))
			fmt.Printf("Error.Message.........: %s\n", e.Buffer)
			fmt.Printf("\n")
		}
	},
}

var updateJobCmd = &cobra.Command{
	Use:   "update <project_name|project_id> <job_id>.",
	Short: "Updates a job by project_name|project_id and job_id.",
	Long: `
Updates a job by project_name|project_id and job_id. Can be used to update
priority and/or max-devices.
	`,
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
			Priority:            flPriority,
			MaxDedicatedDevices: flMaxDedicatedDevices,
			IsActive:            job.IsActive,
		}
		path := fmt.Sprintf("/api/projects/%d/jobs/%d", project.ID, job.ID)
		if _, err := patchJSON(path, &update); err != nil {
			writeStdErrAndExit(err.Error())
		}
		fmt.Println("The job has been updated.")
	},
}

var delJobCmd = &cobra.Command{
	Use:    "delete <project_name|project_id> <job_id>",
	Short:  "Deletes a job by project_name|project_id and job_id.",
	Long:   "Deletes a job by project_name|project_id and job_id.",
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
		if ok := promptDelete("this job"); !ok {
			writeStdErrAndExit("Not deleting job.")
		}
		path := fmt.Sprintf("/api/projects/%d/jobs/%d", project.ID, job.ID)
		if err := deleteHTTP(path); err != nil {
			writeStdErrAndExit(err.Error())
		}
		if attack.Title == fmt.Sprintf("hashstack-cli-%d-%d-%s", job.ProjectID, job.ListID, job.Name) {
			deleteHTTP(fmt.Sprintf("/api/attacks/%d", job.AttackID))
		}
		fmt.Println("The job was successfully deleted.")
	},
}

var startJobCmd = &cobra.Command{
	Use:    "start <project_name|project_id> <job_id>",
	Short:  "Starts a job by project_name|project_id and job_id.",
	Long:   "Starts a job by project_name|project_id and job_id.",
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
		fmt.Println("The job has been started.")
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
	Use:   "add <project_name|project_id> <list_name|list_id> <name> <wordlist|mask>",
	Short: "Add a job for the provided project and list.",
	Long: `Add a job for the provided project and list.


Attack Modes:
0 | Straight
1 | Combination
3 | Brute-force
6 | Hybrid Wordlist + Mask
7 | Hybrid Mask + Wordlist
`,
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
			deleteHTTP(fmt.Sprintf("/api/attacks/%d", plan.ID))
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
	addJobCmd.PersistentFlags().IntVarP(&flAttackMode, "attack-mode", "a", 0, "Attack mode, see references above")
	addJobCmd.PersistentFlags().BoolVar(&flIsHexCharset, "hex-charset", false, "Assume charset is given in hex")
	addJobCmd.PersistentFlags().StringVar(&flMarkovHcstat, "markov-hcstat", "", "Specify hcstat file to use")
	addJobCmd.PersistentFlags().IntVarP(&flMarkovThreshold, "markov-threshold", "t", 0, "Threshold X when to stop accepting new markov-chains")
	addJobCmd.PersistentFlags().IntVar(&flOpenCLVectorWidth, "opencl-vector-width", 0, "Manual override OpenCL vector-width to X")
	addJobCmd.PersistentFlags().IntVar(&flPriority, "priority", 1, "The priority for this job 1-100")
	addJobCmd.PersistentFlags().IntVar(&flMaxDedicatedDevices, "max-devices", 0, "Maximum devices across the entire cluster to use, 0 is unlimited")
	addJobCmd.PersistentFlags().StringVarP(&flRuleLeft, "rule-left", "j", "", "Single rule applied to each word from left wordlist")
	addJobCmd.PersistentFlags().StringVarP(&flRuleRight, "rule-right", "k", "", "Single rule applied to each word from left wordlist")
	addJobCmd.PersistentFlags().StringVarP(&flRulesFile, "rules-file", "r", "", "Rule file to be applied to each word from wordlists")
	addJobCmd.PersistentFlags().StringVarP(&flCustomCharset1, "custom-charset1", "1", "", "User-defined charset ?1")
	addJobCmd.PersistentFlags().StringVarP(&flCustomCharset2, "custom-charset2", "2", "", "User-defined charset ?2")
	addJobCmd.PersistentFlags().StringVarP(&flCustomCharset3, "custom-charset3", "3", "", "User-defined charset ?3")
	addJobCmd.PersistentFlags().StringVarP(&flCustomCharset4, "custom-charset4", "4", "", "User-defined charset ?4")
	updateJobCmd.PersistentFlags().IntVar(&flPriority, "priority", 1, "The priority for this job 1-100")
	updateJobCmd.PersistentFlags().IntVar(&flMaxDedicatedDevices, "max-devices", 0, "Maximum devices across the entire cluster to use, 0 is unlimited")
	jobCmd.AddCommand(addJobCmd)
	jobCmd.AddCommand(pauseJobCmd)
	jobCmd.AddCommand(startJobCmd)
	jobCmd.AddCommand(updateJobCmd)
	jobCmd.AddCommand(delJobCmd)
	jobCmd.AddCommand(errorJobCmd)
	RootCmd.AddCommand(jobCmd)
}
