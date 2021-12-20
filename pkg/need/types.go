package need

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"os/exec"
)

var bold = lipgloss.NewStyle().Bold(true)
var italic = lipgloss.NewStyle().Italic(true)

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

// Address first assesses, then fulfills if necessary and assesses again.
func (n *Need) Address() (err error) {
	fmt.Printf("\n%s\n", bold.Render(n.Name))
	if n.Assess() == nil {
		return
	}
	err = n.Fulfill()
	if err != nil {
		n.RenderHelp()
		return
	}
	err = n.Assess()
	if err != nil {
		n.RenderHelp()
	}
	return
}

// Assess runs the AssessCmd and returns an error if it failed.
func (n *Need) Assess() (err error) {
	fmt.Println("  Assessing ...")
	assessCmd := exec.Command("bash", "-euo", "pipefail", "-c", n.AssessCmd)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	assessCmd.Stdout = &stdout
	assessCmd.Stderr = &stderr
	err = assessCmd.Run()
	if err == nil {
		fmt.Println("  Fulfilled")
	} else {
		fmt.Println("  Unfulfilled")
		fmt.Printf("    stdout> %s\n", stdout.String())
		fmt.Printf("    stderr> %s\n", stderr.String())
	}
	return
}

// Fulfill runs the FulfillCmd and returns an error if it failed.
func (n *Need) Fulfill() (err error) {
	fmt.Println("  Fulfilling ...")
	fulfillCmd := exec.Command("bash", "-euo", "pipefail", "-c", n.FulfillCmd)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	fulfillCmd.Stdout = &stdout
	fulfillCmd.Stderr = &stderr
	err = fulfillCmd.Run()
	if err == nil {
		fmt.Printf("  Done\n")
	} else {
		fmt.Println("  Failed")
		fmt.Printf("    stdout> %s\n", stdout.String())
		fmt.Printf("    stderr> %s\n", stderr.String())
	}
	return
}

// RenderHelp returns the terminal-friendly rendered Help markdown.
func (n *Need) RenderHelp() (err error) {
	var renderedHelp string
	helpMarkdown := n.Help
	if helpMarkdown == "" {
		helpMarkdown = "_Sorry, no help._"
	}
	renderedHelp, err = glamour.Render(helpMarkdown, "dark")
	if err != nil {
		fmt.Print("Help: _<Failed to render help. Is it valid Markdown?>_")
		return
	}
	fmt.Print("  Help:")
	fmt.Print(renderedHelp)
	return
}
