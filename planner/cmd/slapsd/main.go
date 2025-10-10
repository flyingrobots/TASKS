package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	execapp "github.com/james/tasks-planner/internal/app/exec"
)

func main() {
	coordPath := flag.String("coord", "./coordinator.json", "Path to coordinator.json artifact")
	flag.Parse()

	svc := execapp.NewDefaultService()
	if err := svc.Run(context.Background(), *coordPath); err != nil {
		if err == execapp.ErrLoopNotImplemented {
			fmt.Fprintln(os.Stderr, "slapsd: execution loop not yet implemented (stub)")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "slapsd: %v\n", err)
		os.Exit(1)
	}
}
