/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"task-lineage-diagram/dot"
	"task-lineage-diagram/reader"

	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"
)

var configFile string
var rootDir string
var outputFile string
var reachabilityOutputFile string
var noReach bool
var format string
var size string
var group bool
var color bool
var layout string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tld",
	Short: "A simple golang module for rendering task lineage diagrams",
	Long:  `A simple golang module for rendering task lineage diagrams`,
	Run: func(cmd *cobra.Command, args []string) {
		layout := getLayout(layout)
		format := getFormat(format)
		fmt.Printf("Reading task yaml files from \"%s\" ...\n", rootDir)
		config, err := reader.ReadConfig(configFile)
		tasks, err := reader.ReadTasks(rootDir)
		if err != nil {
			log.Fatal(err)
		}
		dot.Render(tasks, outputFile, config, format, layout, group, color, size, reachabilityOutputFile, noReach)
		fmt.Println("Done.")
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.task-lineage-diagram.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&configFile, "config", "k", "config.yaml", "Path for yaml config file.")
	rootCmd.Flags().StringVarP(&rootDir, "input", "i", ".", "Root directory for task yaml files.")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "graph", "Output file path.")
	rootCmd.Flags().StringVarP(&reachabilityOutputFile, "reach", "r", "reachability.json", "Output file path for the reachability analysis report.")
	rootCmd.Flags().BoolVarP(&noReach, "no-reach", "n", false, "Turn this on to skip the reachability analysis.")
	rootCmd.Flags().StringVarP(&format, "format", "f", "svg", "Output file format, one of [svg, dot, png, jpg].")
	rootCmd.Flags().BoolVarP(&group, "group", "g", false, "Group Layer. If turned on, nodes under the same layer are grouped together, which means they are placed next to each other if possible.")
	rootCmd.Flags().StringVarP(&layout, "layout", "l", "dot", "Graph Layout. Currently support [circo, dot, fdp, neato, osage, patchwork, sfdp, twopi].")
	rootCmd.Flags().BoolVarP(&color, "color", "c", false, "Color mode. If turned on, the output is colored.")
	rootCmd.Flags().StringVarP(&size, "size", "s", "fhd", "Graph size. Currently only support [fhd, a3].")
}

func getLayout(input string) graphviz.Layout {
	switch {
	case strings.ToLower(input) == "circo":
		return graphviz.CIRCO
	case strings.ToLower(input) == "dot":
		return graphviz.DOT
	case strings.ToLower(input) == "fdp":
		return graphviz.FDP
	case strings.ToLower(input) == "neato":
		return graphviz.NEATO
	case strings.ToLower(input) == "osage":
		return graphviz.OSAGE
	case strings.ToLower(input) == "patchwork":
		return graphviz.PATCHWORK
	case strings.ToLower(input) == "sfdp":
		return graphviz.SFDP
	case strings.ToLower(input) == "twopi":
		return graphviz.TWOPI
	default:
		return graphviz.DOT
	}
}

func getFormat(input string) graphviz.Format {
	switch {
	case strings.ToLower(input) == "svg":
		return graphviz.SVG
	case strings.ToLower(input) == "png":
		return graphviz.PNG
	case strings.ToLower(input) == "jpg":
		return graphviz.JPG
	case strings.ToLower(input) == "dot":
		return graphviz.XDOT
	default:
		return graphviz.XDOT
	}
}
