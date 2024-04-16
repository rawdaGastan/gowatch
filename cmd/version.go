// Package cmd for farmerbot commands
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "get gowatch latest version and commit",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if len(cmd.Flags().Args()) != 0 {
			return fmt.Errorf("'version' and %v cannot be used together, please use one command at a time", cmd.Flags().Args())
		}

		fmt.Println(Version)
		return nil
	},
}
