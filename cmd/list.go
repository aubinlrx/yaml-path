/*
Copyright © 2020 Aubin LORIEUX

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
	"io/ioutil"
	"strconv"

	"github.com/aubinlrx/yaml-path/lib/filename"
	"github.com/aubinlrx/yaml-path/lib/yamlpath"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("separator", "s", ".", "Path separator (default: .)")
	listCmd.Flags().BoolP("line", "l", false, "Add Line to output")
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [FILE]",
	Short: "List yaml path of the files",
	Long:  ``,
	RunE:  yamlList,
}

func yamlList(cmd *cobra.Command, args []string) error {
	filename, err := filename.GetFromArgs(args)
	if err != nil {
		return err
	}

	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	sep, _ := cmd.Flags().GetString("separator")
	line, _ := cmd.Flags().GetBool("line")

	yp := yamlpath.Path{In: dat}
	paths, err := yp.List(sep)
	if err != nil {
		return err
	}

	for _, p := range paths {
		res := p.Path()

		if line {
			res += " #" + strconv.Itoa(p.Line)
		}

		fmt.Println(res)
	}

	return nil
}
