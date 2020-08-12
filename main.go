package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func argsHandle(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("need url")
	}
	return nil
}

func handleCmd(cmd *cobra.Command, args []string) error {
	url := args[0]
	cli := NewCLI()

	autoStr := cmd.Flag("auth").Value.String()

	auth := strings.Split(autoStr, ":")
	if len(auth) != 2 {
		return errors.New("parse auth information failed")
	}

	if err := cli.Ping(url, auth); err != nil {
		return err
	}
	return cli.Run()
}

var rootCmd = &cobra.Command{
	Use:  "elasticsql-cli es-url",
	Long: `A simple CLI tools for elasticsearch sql.`,
	Args: argsHandle,
	RunE: handleCmd,
}

func init() {
	rootCmd.PersistentFlags().String("auth", "", "Auth info, e.g username:pwd")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
