package main

import (
	"Masterwow3/docker-netrestore/pkg/netrestore"
	"fmt"
	"os/exec"

	"github.com/docker/cli/cli-plugins/metadata"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

const mainCommandName = "netrestore"

var versionBuild = "0"
var commitSha = "0000000"

func main() {
	version := "1." + versionBuild + ".0"

	plugin.Run(func(dockerCLI command.Cli) *cobra.Command {
		version := &cobra.Command{
			Use:   "version",
			Short: "Show the " + mainCommandName + " version information",
			RunE: func(_ *cobra.Command, _ []string) error {
				fmt.Fprintln(dockerCLI.Out(), "Version: "+version)
				fmt.Fprintln(dockerCLI.Out(), "CommitSHA: "+commitSha)
				return nil
			},
		}

		var registerAutorun bool
		cmd := &cobra.Command{
			Use:   mainCommandName,
			Short: "Restores the HNS network to prevent Docker networks from disappearing",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				return plugin.PersistentPreRunE(cmd, args)
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				if registerAutorun {
					return registerWindowsTask()
				}

				return netrestore.FixNetwork()
			},
		}

		flags := cmd.Flags()
		flags.BoolVar(&registerAutorun, "register-autorun", false, "Registers the "+mainCommandName+" command with the Windows Task Scheduler to automatically restore networks after each system reboot.")

		cmd.AddCommand(version)
		return cmd
	},
		metadata.Metadata{
			SchemaVersion: "0.1.0",
			Vendor:        "Masterwow3",
			Version:       version,
		})
}

func registerWindowsTask() error {
	taskName := "DockerNetrestoreTask"
	taskCmd := `cmd /c "docker netrestore"`
	schedule := "ONSTART"
	user := "SYSTEM"

	cmd := exec.Command("schtasks",
		"/Create",
		"/TN", taskName,
		"/TR", taskCmd,
		"/SC", schedule,
		"/RU", user,
		"/RL", "HIGHEST",
		"/F", // überschreibt bestehenden Task ohne Nachfrage
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error creating task: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("Task created successfully.")
	return nil
}
