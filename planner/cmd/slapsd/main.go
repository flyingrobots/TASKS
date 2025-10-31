package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	execapp "github.com/james/tasks-planner/internal/app/exec"
)

func main() {
	coordPath := flag.String("coord", "./coordinator.json", "Path to coordinator.json artifact")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	svc := execapp.NewDefaultService()
	if err := svc.Run(ctx, *coordPath); err != nil {
		if errors.Is(err, execapp.ErrLoopNotImplemented) {
			fmt.Fprintln(os.Stderr, "slapsd: execution loop not yet implemented (stub)")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "slapsd: %v\n", err)
		os.Exit(1)
	}
}
