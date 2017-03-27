package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func getAgent(id int64) hashstack.Agent {
	var agent hashstack.Agent
	path := fmt.Sprintf("/api/agents/%d", id)
	if err := getJSON(path, &agent); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return agent
}

func displayAgent(a hashstack.Agent) {
	agent := getAgent(a.ID)
	memstat := fmt.Sprintf("%s/%s (%2.f%%)",
		humanize.Bytes(uint64(agent.MemoryUsed)),
		humanize.Bytes(uint64(agent.MemoryTotal)),
		percentOf(int(agent.MemoryUsed), int(agent.MemoryTotal)))

	fmt.Printf("ID..............: %d\n", agent.ID)
	fmt.Printf("Host............: %s\n", agent.Hostname)
	fmt.Printf("IP.Address......: %s\n", agent.IPAddress)
	fmt.Printf("Uptime..........: %s\n", prettyUptime(agent.Uptime))
	fmt.Printf("Last.Seen.......: %s\n", humanize.Time(time.Unix(agent.CheckinAt, 0)))
	fmt.Printf("Memory..........: %s\n", memstat)
	for i, d := range agent.Devices {
		fmt.Printf("Dev.#%d..........: %s, %d Mhz, %d%% load, %dC, %d%% Fan\n",
			i+1,
			d.Name,
			d.CurrentClockFrequency,
			d.Load,
			d.Temperature,
			d.FanSpeed)
	}
	fmt.Println()
}

func displayAgents() {
	var agents []hashstack.Agent
	if err := getRangeJSON("/api/agents", &agents); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].CreatedAt > agents[j].CreatedAt
	})
	for _, a := range agents {
		displayAgent(a)
	}
}

var agentCmd = &cobra.Command{
	Use:   "agents [id]",
	Short: "Display a list of agents connected to the cluster.",
	Long: `
Display a list of agents connected to the cluster. If an id is provided, only information
on that agent will be displayed.`,
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			displayAgents()
		case 1:
			i, err := strconv.Atoi(args[0])
			if err != nil {
				writeStdErrAndExit(fmt.Sprintf("%s is not a valid agent id.", args[0]))
			}
			displayAgent(hashstack.Agent{
				ID: int64(i),
			})
		default:
			cmd.Usage()
		}
	},
}

func init() {
	RootCmd.AddCommand(agentCmd)
}
