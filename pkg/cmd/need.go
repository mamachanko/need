package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/mamachanko/need/pkg/need"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

func NewNeedCmd() (cmd *cobra.Command) {
	o := &NeedOptions{}
	cmd = &cobra.Command{
		Use:   "need",
		Short: "Very simple need management",
		Run:   func(cmd *cobra.Command, args []string) { o.Run() },
	}
	cmd.Flags().StringArrayVarP(&o.Files, "file", "f", nil, "File (local path or -) (can be specified multiple times)")
	return
}

type NeedOptions struct {
	Files []string
}

// Run is the entrypoint. It addresses all needs in all given config files.
func (o NeedOptions) Run() {
	succeeded := true

	for _, file := range o.Files {
		var err error
		var fileContent []byte

		switch {
		case file == "-":
			// Issue: This will only read from stdin once.
			// If --file/-f - is given more than once only the first will actually be read.
			// The subsequent reads will result in an empty fileContent.
			fileContent, err = ioutil.ReadAll(os.Stdin)
		default:
			fileContent, err = ioutil.ReadFile(file)
		}
		if err != nil {
			log.Fatalf("error: %v", err)
			os.Exit(1)
		}

		needsCfg := need.NeedsConfig{}
		err = yaml.Unmarshal(fileContent, &needsCfg)
		if err != nil {
			log.Fatalf("error: %v", err)
			os.Exit(1)
		}

		fmt.Printf("\n%s\n", lipgloss.NewStyle().Italic(true).Render(needsCfg.Metadata.Name))

		for _, need := range needsCfg.Spec.Needs {
			if need.Address() != nil {
				succeeded = false
			}
		}

	}

	if succeeded {
		fmt.Println("\nSucceeded")
	} else {
		fmt.Println("\nFailed")
		os.Exit(1)
	}
}
