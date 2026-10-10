package main

import (
	"UpYnab/internal/config"
	"UpYnab/internal/doctor"
	"context"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s doctor\n", os.Args[0])
		os.Exit(2)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load configuration: %v\n", err)
		os.Exit(2)
	}

	command := os.Args[1]
	switch command {
	case "doctor":
		err := doctor.Run(context.Background(), cfg, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Doctor failed: %v\n", err)
			os.Exit(2)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		fmt.Fprintf(os.Stderr, "Usage: %s doctor\n", os.Args[0])
		os.Exit(2)
	}
}
