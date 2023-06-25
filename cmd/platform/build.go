package main

import (
	"log"

	"github.com/spf13/cobra"
	"platform.prodigy9.co/builder"
	"platform.prodigy9.co/config"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds current directory",
	Run:   runBuild,
}

func runBuild(cmd *cobra.Command, args []string) {
	cfg, err := config.Configure(".")
	if err != nil {
		log.Fatalln(err)
	}

	jobs, err := builder.JobsFromArgs(cfg, args)
	if err != nil {
		log.Fatalln(err)
	}

	if err := builder.Build(cfg, jobs...); err != nil {
		log.Fatalln(err)
	}
}
