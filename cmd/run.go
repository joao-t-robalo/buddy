package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dreadster3/buddy/models"
	"github.com/spf13/cobra"
)

var listCommands bool

func init() {
	runCmd.Flags().BoolVarP(&listCommands, "list", "l", false, "List all available commands")
	buddyCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [flags] [command]",
	Short: "Run a predefined command",
	Long:  `Run a predefined command`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if listCommands || len(args) == 0 {
			buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
			if err != nil {
				return err
			}

			for commandName := range buddyConfig.Scripts {
				fmt.Println(commandName, "->", buddyConfig.Scripts[commandName])
			}

			return nil
		}

		commandName := args[0]
		commandArgs := args[1:]

		buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
		if err != nil {
			return err
		}

		commandTemplate, ok := buddyConfig.Scripts[commandName]
		if !ok {
			return fmt.Errorf("Command %s not found", commandName)
		}

		command := replaceArgs(commandTemplate, commandArgs)

		execCommand := exec.Command("sh", "-c", command)
		execCommand.Stdout = os.Stdout
		execCommand.Stderr = os.Stderr

		err = execCommand.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func replaceArgs(commandTemplate string, args []string) string {
	for i, arg := range args {
		placeholder := fmt.Sprintf("${%d}", i+1)
		commandTemplate = strings.ReplaceAll(commandTemplate, placeholder, arg)
	}
	return commandTemplate
}