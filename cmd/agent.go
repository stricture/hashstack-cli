package cmd

import (
	"encoding/binary"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func getAgentByID(id int64) hashstack.Agent {
	var agent hashstack.Agent
	path := fmt.Sprintf("/api/agents/%d", id)
	if err := getJSON(path, &agent); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return agent
}

func getAgent(uuid string) hashstack.Agent {
	var agent hashstack.Agent
	path := fmt.Sprintf("/api/agents/%s", uuid)
	if err := getJSON(path, &agent); err != nil {
		writeStdErrAndExit(err.Error())
	}
	return agent
}

func isOnline(agent hashstack.Agent) bool {
	return time.Now().Add(-5*time.Minute).Unix() < agent.CheckinAt
}

func displayAgent(a hashstack.Agent) {
	agent := getAgent(a.UUID)
	memstat := fmt.Sprintf("%s/%s (%2.f%%)",
		humanize.Bytes(uint64(agent.MemoryUsed)),
		humanize.Bytes(uint64(agent.MemoryTotal)),
		percentOf(int(agent.MemoryUsed), int(agent.MemoryTotal)))

	online := "Offline"
	if isOnline(agent) {
		online = "Online"
	}

	fmt.Printf("ID..............: %s\n", agent.UUID)
	fmt.Printf("Host............: %s\n", agent.Hostname)
	fmt.Printf("IP.Address......: %s\n", agent.IPAddress)
	fmt.Printf("Status..........: %s\n", online)
	fmt.Printf("Last.Seen.......: %s\n", humanize.Time(time.Unix(agent.CheckinAt, 0)))
	if isOnline(agent) {
		fmt.Printf("Uptime..........: %s\n", prettyUptime(agent.Uptime))
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
	}
	fmt.Println()
}

var (
	flAgentSortOrder      string
	flAgentShowOnlineOnly bool
)

func displayAgents() {
	var agents []hashstack.Agent
	if err := getRangeJSON("/api/agents", &agents); err != nil {
		writeStdErrAndExit(err.Error())
	}
	if len(agents) < 1 {
		writeStdErrAndExit("There are no agents in the cluster!")
	}
	switch flAgentSortOrder {
	case "created_at":
		sort.Slice(agents, func(i, j int) bool {
			return agents[i].CreatedAt > agents[j].CreatedAt
		})
	case "hostname":
		sort.Slice(agents, func(i, j int) bool {
			si := agents[i].Hostname
			sj := agents[j].Hostname
			silower := strings.ToLower(si)
			sjlower := strings.ToLower(sj)
			if silower == sjlower {
				return si < sj
			}
			return silower < sjlower
		})
	default:
		sort.Slice(agents, func(i, j int) bool {
			first := net.ParseIP(agents[i].IPAddress).To4()
			second := net.ParseIP(agents[j].IPAddress).To4()
			return binary.BigEndian.Uint32(first) < binary.BigEndian.Uint32(second)
		})
	}
	for _, a := range agents {
		if flAgentShowOnlineOnly && !isOnline(a) {
			continue
		}
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
			displayAgent(hashstack.Agent{
				UUID: args[0],
			})
		default:
			cmd.Usage()
		}
	},
}

func init() {
	agentCmd.PersistentFlags().StringVar(&flAgentSortOrder, "sort", "ip", "Sort order for agents. Values can be 'ip', 'hostname', or 'created_at'")
	agentCmd.PersistentFlags().BoolVar(&flAgentShowOnlineOnly, "show-online", false, "Show only online agents.")
	RootCmd.AddCommand(agentCmd)
}
