package cmd

import (
	"fmt"
	"log"
	"os"

	chutils "creativehashtags.com/wordle/utils"
	"github.com/spf13/cobra"
)

var config = chutils.Config{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "guesser",
	Short: "An application to guess words based on a pattern",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("guesser called")
		wrds, err := chutils.ReadWords()
		if err != nil {
			log.Fatalf("Error opening file %v\n", err)
		}
		wrds = chutils.FiveLetterWords(wrds)
		// let's see what folks give us
		rep := chutils.Report(wrds, config)
		fmt.Println(rep)
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
	rootCmd.Flags().StringVarP(&config.Include, "include", "i", "", "The letters to include")
	rootCmd.Flags().StringVarP(&config.Exclude, "exclude", "e", "", "The letters to exclude")
	rootCmd.Flags().StringVarP(&config.Pattern, "pattern", "p", "", "A string of letters that are in the right position")
	rootCmd.Flags().StringVarP(&config.AntiPattern, "antipattern", "a", "", "A string of letters that are in the wrong position")
}
