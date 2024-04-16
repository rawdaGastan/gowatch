// Package cmd for parsing command line arguments
package cmd

import (
	"os"
	"strings"

	"github.com/rawdaGastan/gowatch/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// goWatchCmd represents the root base command when called without any subcommands.
var goWatchCmd = &cobra.Command{
	Use:   "gowatch",
	Short: "Run gowatch",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		interval, err := cmd.Flags().GetFloat64("interval")
		if err != nil {
			return err
		}

		diff, err := cmd.Flags().GetBool("differences")
		if err != nil {
			return err
		}

		noTitle, err := cmd.Flags().GetBool("no-title")
		if err != nil {
			return err
		}

		exec, err := cmd.Flags().GetBool("exec")
		if err != nil {
			return err
		}

		errExit, err := cmd.Flags().GetBool("errexit")
		if err != nil {
			return err
		}

		chgExit, err := cmd.Flags().GetBool("chgexit")
		if err != nil {
			return err
		}

		beep, err := cmd.Flags().GetBool("beep")
		if err != nil {
			return err
		}

		var execCmd string
		var execArgs []string
		if len(args) > 0 {
			execCmd = args[0]
		}

		if len(args) > 1 {
			execArgs = args[1:]
		}

		updateCmd, err := cmd.Flags().GetString("update")
		if err != nil {
			return err
		}

		var uCmd string
		var uArgs []string
		updateCmdArgs := strings.Split(updateCmd, " ")
		if len(updateCmdArgs) > 0 {
			uCmd = updateCmdArgs[0]
		}

		if len(updateCmdArgs) > 1 {
			uArgs = updateCmdArgs[1:]
		}

		// TODO: validations
		app := app.App{
			Interval: interval,
			Exec:     exec,
			Diff:     diff,
			NoTitle:  noTitle,
			ChgExit:  chgExit,
			ErrExit:  errExit,
			Beep:     beep,

			Cmd:        execCmd,
			Args:       execArgs,
			UpdateCmd:  uCmd,
			UpdateArgs: uArgs,
		}

		app.Run(cmd.Context())

		return nil
	},
}

func Execute() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	goWatchCmd.Root().CompletionOptions.DisableDefaultCmd = true

	goWatchCmd.AddCommand(versionCmd)

	err := goWatchCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func init() {
	goWatchCmd.Flags().StringP("update", "u", "", "update command to be executed when output updates")

	goWatchCmd.Flags().Float64P("interval", "n", 2, "seconds to wait between updates")
	goWatchCmd.Flags().BoolP("differences", "d", false, "highlight changes between updates")
	goWatchCmd.Flags().BoolP("exec", "x", false, "pass command to exec instead of \"sh -c\"")
	goWatchCmd.Flags().BoolP("no-title", "t", false, "turn off header")
	goWatchCmd.Flags().BoolP("chgexit", "g", false, "exit when output from command changes")
	goWatchCmd.Flags().BoolP("errexit", "e", false, "exit if command has a non-zero exit")

	goWatchCmd.Flags().BoolP("beep", "b", false, "beep if command has a non-zero exit")

	// goWatchCmd.Flags().BoolP("precise", "p", false, "attempt run command in precise intervals")
	// goWatchCmd.Flags().BoolP("color", "c", false, "interpret ANSI color and style sequences")
	// goWatchCmd.Flags().BoolP("no-linewrap", "w", false, "Turns off line wrapping and truncates long lines instead.")
}
