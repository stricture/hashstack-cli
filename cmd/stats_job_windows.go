package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	hashstack "github.com/stricture/hashstack-server-core-ng"
)

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
