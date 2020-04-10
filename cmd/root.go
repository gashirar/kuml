/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/gashirar/kuml/pkg/plantuml"
	"github.com/spf13/cobra"
	"os"

	"github.com/gashirar/kuml/pkg/resource"
)

// rootCmd represents the base command when called without any subcommands

var rootCmd = &cobra.Command{
	Use:   "kuml [FILE | DIRECTORY]",
	Short: "Kuml is a Manifest visualization tool.",
	Long:  `Kuml is a misualization tool that outputs PlantUML from Kubernetes YAML Manifest.`,

	Run: func(cmd *cobra.Command, args []string) {
		yamlByteSlice := resource.ReadYaml(false, args...)
		apiResourceList := resource.NewAPIResourceList(yamlByteSlice)
		pUml := plantuml.NewPlantUML(apiResourceList)
		pUml.Render()

	},

	Args: cobra.ExactArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("show-link-label", "s", false, "Display the label of Link between Elements.")
}

func initConfig() {

}
