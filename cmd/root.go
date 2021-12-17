/*
Copyright © 2021 Max Brauer <mamachanko>

*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Options struct {
	Files    []string
	FailFast bool
}

// NeedsConfig holds the fields parsed from the Needs configuration file (needs.yaml).
type NeedsConfig struct {
	// APIVersion is the version of the configuration.
	APIVersion string `yaml:"apiVersion" yamltags:"required"`

	// Kind is always `Needs`. Defaults to `Needs`.
	Kind string `yaml:"kind" yamltags:"required"`

	// Metadata holds additional information about the config.
	Metadata Metadata `yaml:"metadata,omitempty"`

	// Dependencies describes a list of other required configs for the current config.
	Spec NeedsSpec `yaml:"spec" yamltags:"required"`
}

// Metadata holds an optional name of the Needs.
type Metadata struct {
	// Name is an identifier of the Needs.
	Name string `yaml:"name,omitempty"`
}

type NeedsSpec struct {
	Needs []Need
}

type Need struct {
	Name       string
	FulfillCmd string `yaml:"fulfillCmd"`
	AssessCmd  string `yaml:"assessCmd"`
	Help       string
}

var options Options
var bold = lipgloss.NewStyle().Bold(true)
var italic = lipgloss.NewStyle().Italic(true)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "need",
	Short: "Very simple need fulfillment and assessment",
	Run: func(cmd *cobra.Command, args []string) {
		succeeded := true

		for _, file := range options.Files {
			var err error
			fileContent, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatalf("error: %v", err)
				os.Exit(1)
			}

			needsCfg := NeedsConfig{}
			err = yaml.Unmarshal(fileContent, &needsCfg)
			if err != nil {
				log.Fatalf("error: %v", err)
				os.Exit(1)
			}

			fmt.Printf("\n%s\n", italic.Render(needsCfg.Metadata.Name))

			for _, need := range needsCfg.Spec.Needs {
				fmt.Printf("\n%s\n", bold.Render(need.Name))
				fmt.Println("  Assessing ...")
				assessBeforeCmd := exec.Command("bash", "-euo", "pipefail", "-c", need.AssessCmd)
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				assessBeforeCmd.Stdout = &stdout
				assessBeforeCmd.Stderr = &stderr
				err = assessBeforeCmd.Run()
				if err == nil {
					fmt.Println("  Fulfilled")
				} else {
					if options.FailFast {
						if need.Help != "" {
							fmt.Print("  Help:")
							var renderedHelp string
							renderedHelp, err = glamour.Render(need.Help, "dark")
							fmt.Print(renderedHelp)
						} else {
							var renderedMissingHelp string
							renderedMissingHelp, err = glamour.Render("_Sorry, no help._", "dark")
							fmt.Print(renderedMissingHelp)
						}
						succeeded = false
						os.Exit(1)
					}

					fmt.Println("  Unfulfilled")
					fmt.Printf("    stdout> %s\n", stdout.String())
					fmt.Printf("    stderr> %s\n", stderr.String())

					fmt.Println("  Fulfilling ...")
					fulfillCmd := exec.Command("bash", "-euo", "pipefail", "-c", need.FulfillCmd)
					var fulfillStdout bytes.Buffer
					var fulfillStderr bytes.Buffer
					fulfillCmd.Stdout = &fulfillStdout
					fulfillCmd.Stderr = &fulfillStderr
					err = fulfillCmd.Run()
					if err == nil {
						fmt.Printf("  Done\n")
					} else {
						fmt.Println("  Failed")
						fmt.Printf("    stdout> %s\n", fulfillStdout.String())
						fmt.Printf("    stderr> %s\n", fulfillStderr.String())
						succeeded = false
						break
					}
					fmt.Println("  Assessing again ...")
					assessAfterCmd := exec.Command("bash", "-euo", "pipefail", "-c", need.AssessCmd)
					var assessAfterStdout bytes.Buffer
					var assessAfterStderr bytes.Buffer
					assessAfterCmd.Stdout = &assessAfterStdout
					assessAfterCmd.Stderr = &assessAfterStderr
					err = assessAfterCmd.Run()
					if err == nil {
						fmt.Printf("  Fulfilled\n")
					} else {
						fmt.Println("  Still unfulfilled")
						fmt.Printf("    stdout> %s\n", assessAfterStdout.String())
						fmt.Printf("    stderr> %s\n", assessAfterStderr.String())
						if need.Help != "" {
							fmt.Print("  Help:")
							var out string
							out, err = glamour.Render(need.Help, "dark")
							fmt.Print(out)
						} else {
							fmt.Printf("  Help: \n    <Sorry, no help.>\n")
						}

						succeeded = false
					}
				}
			}

		}
		if succeeded {
			fmt.Println("\nSucceeded")
		} else {
			fmt.Println("\nFailed")
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringArrayVarP(&options.Files, "file", "f", nil, "File (local path or -) (can be specified multiple times)")
	rootCmd.Flags().BoolVarP(&options.FailFast, "fail-fast", "x", false, "Fail upon encountering an error")
}
