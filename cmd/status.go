package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func displayStats() {
	var stats hashstack.ClusterStats
	if err := getJSON("/api/stats", &stats); err != nil {
		writeStdErrAndExit(err.Error())
	}

	var (
		gpuload string
		cpuload string
	)

	if stats.GPULoadMax > 0 {
		gpuload = fmt.Sprintf("%d/%d (%0.2f%%)", stats.GPULoad, stats.GPULoadMax, percentOf(int(stats.GPULoad), int(stats.GPULoadMax)))
	} else {
		gpuload = fmt.Sprintf("%d/%d (0%%)", stats.GPULoad, stats.GPULoadMax)
	}
	if stats.CPULoadMax > 0 {
		cpuload = fmt.Sprintf("%d/%d (%0.2f%%)", stats.CPULoad, stats.CPULoadMax, percentOf(int(stats.CPULoad), int(stats.CPULoadMax)))
	} else {
		cpuload = fmt.Sprintf("%d/%d (0%%)", stats.CPULoad, stats.CPULoadMax)
	}
	fmt.Printf("Jobs...........: %d Active, %d Paused\n", stats.ActiveJobCount, stats.PausedJobCount)
	fmt.Printf("Node.Count.....: %d\n", stats.AgentCount)
	fmt.Printf("GPU.Count......: %d\n", stats.GPUCount)
	fmt.Printf("CPU.Count......: %d\n", stats.CPUCount)
	fmt.Printf("Load.GPU.......: %s\n", gpuload)
	fmt.Printf("Load.CPU.......: %s\n", cpuload)
	fmt.Printf("Temp.GPU.......: %dC - %dC\n", stats.LowestGPUTemp, stats.HighestGPUTemp)
	fmt.Printf("Temp.CPU.......: %dC - %dC\n", stats.LowestCPUTemp, stats.HighestCPUTemp)
	var agents []hashstack.Agent
	if err := getRangeJSON("/api/agents", &agents); err != nil {
		writeStdErrAndExit(err.Error())
	}
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].CreatedAt > agents[j].CreatedAt
	})
	for j, agent := range agents {
		for i, d := range agent.Devices {
			fmt.Printf("Agent.#%d.Dev.#%d..........: %s, %d Mhz, %d%% load, %dC, %d%% Fan\n",
				j+1,
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

var statusCmd = &cobra.Command{
	Use:    "status",
	Short:  "Displays information about the Hashstack cluster",
	Long:   "Displays information about the Hashstack cluster",
	PreRun: ensureAuth,
	Run: func(cmd *cobra.Command, args []string) {
		displayStats()
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
