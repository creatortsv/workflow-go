/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/flowchart"
	"github.com/spf13/cobra"
)

type (
	config struct {
		Workflow map[string]json.RawMessage `json:"workflow" validate:"required"`
	}
	workflow struct {
		Manager     goData                     `json:"manager"`
		Transitions map[string]json.RawMessage `json:"transitions" validate:"required"`
	}
	transition struct {
		Guards []goData `json:"guards,omitempty"`
		From   []string `json:"from" validate:"required"`
		To     string   `json:"to" validate:"required"`
	}
	goData struct {
		Pkg  string `json:"pkg,omitempty"`
		Type string `json:"type,omitempty"`
		Func string `json:"func,omitempty"`
		Expr string `json:"expr,omitempty"`
	}
)

var configFileName string = "workflow-go.json"

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: graph,
}

func init() {
	rootCmd.AddCommand(graphCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// graphCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// graphCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func graph(cmd *cobra.Command, args []string) {
	config := initConfig()

	for name, wRaw := range config.Workflow {
		f := createFile(name)

		defer f.Close()

		fl := flowchart.NewFlowchart(
			io.Discard,
			flowchart.WithOrientalLeftToRight(),
			flowchart.WithTitle(fmt.Sprintf("%s workflow", name)),
		)

		for name, tRaw := range parseJson(wRaw, &workflow{}).Transitions {
			t := parseJson(tRaw, &transition{})

			for _, state := range t.From {
				fl.LinkWithArrowHeadAndText(state, t.To, name)
			}
		}

		if err := markdown.NewMarkdown(f).
			H2f("%s workflow", name).
			CodeBlocks(markdown.SyntaxHighlightMermaid, fl.String()).
			Build(); err != nil {
			panic(fmt.Errorf("writing markdown file %s.md: %w", name, err))
		}
	}
}

func initConfig() *config {
	d, err := os.ReadFile(configFileName)
	if err != nil {
		panic(fmt.Errorf("reading file %s: %w", configFileName, err))
	}

	return parseJson(d, &config{})
}

func createFile(wn string) *os.File {
	mdName := fmt.Sprintf("%s.md", wn)
	f, err := os.Create(mdName)
	if err != nil {
		panic(fmt.Errorf("creating file %s: %w", mdName, err))
	}

	return f
}

func parseJson[T any](d []byte, v T) T {
	if err := json.Unmarshal(d, v); err != nil {
		panic(fmt.Errorf("parsing file %s: %w", configFileName, err))
	}

	return v
}
