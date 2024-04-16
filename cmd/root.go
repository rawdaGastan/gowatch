// Package cmd for parsing command line arguments
package cmd

import (
	"os"

	"github.com/rawdaGastan/gowatch/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// goWatchCmd represents the root base command when called without any subcommands
var goWatchCmd = &cobra.Command{
	Use:   "gowatch",
	Short: "Run gowatch",
	RunE: func(cmd *cobra.Command, args []string) error {
		interval, err := cmd.Flags().GetUint64("interval")
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

		exec, err := cmd.Flags().GetString("exec")
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

		app := app.App{
			Interval: interval,
			Diff:     diff,
			NoTitle:  noTitle,
			ChgExit:  chgExit,
			ErrExit:  errExit,
			Beep:     beep,

			Cmd: exec,
			// Args:       args,
			// UpdateCmd:  "",
			// UpdateArgs: args,
		}

		app.Run(cmd.Context())

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootcmd.Flags().
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
	var defaultInterval uint64 = 2

	goWatchCmd.Flags().Uint64P("interval", "n", defaultInterval, "seconds to wait between updates")
	goWatchCmd.Flags().BoolP("differences", "d", false, "highlight changes between updates")
	goWatchCmd.Flags().BoolP("no-title", "t", false, "turn off header")
	goWatchCmd.Flags().BoolP("chgexit", "g", false, "exit when output from command changes")
	goWatchCmd.Flags().BoolP("errexit", "e", false, "exit if command has a non-zero exit")

	goWatchCmd.Flags().BoolP("beep", "b", false, "beep if command has a non-zero exit")

	goWatchCmd.Flags().BoolP("precise", "p", false, "attempt run command in precise intervals")
	goWatchCmd.Flags().BoolP("color", "c", false, "interpret ANSI color and style sequences")
	goWatchCmd.Flags().BoolP("--no-linewrap", "w", false, "Turns off line wrapping and truncates long lines instead.")
	goWatchCmd.Flags().StringP("exec", "x", "", "pass command to exec instead of \"sh -c\"")
}
