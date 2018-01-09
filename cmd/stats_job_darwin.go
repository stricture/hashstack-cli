package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	hashstack "github.com/stricture/hashstack-server-core-ng"
)

func statsJob(job hashstack.Job) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)
	go func() {
		switch <-ch {
		case os.Interrupt:
			fmt.Println("Interrupt caught. Job will continue to run on the server.")
			os.Exit(0)
		case syscall.SIGQUIT:
			fmt.Println("Quit caught. Job will be removed from the server.")
			deleteJob(job)
			os.Exit(0)
		}

	}()
	c := time.Tick(5 * time.Second)
	for range c {
		job = getJob(job.ProjectID, job.ID)
		displayJob(os.Stdout, job)
		fmt.Fprintf(os.Stdout, "\nCtrl-C to exit, the job will continue to run. Ctrl-\\ to abort, the job will be removed.\n\n")
		if job.IsExhausted {
			fmt.Printf("The job is finished. View stats using 'hashstack jobs %d %d'.\n\n", job.ProjectID, job.ID)
			fmt.Println("Lists may continue to be updated with recovered plains after the job has finished.")
			break
		}
	}
}
