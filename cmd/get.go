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
	"os"
	"strings"

	"github.com/aubinlrx/yaml-path/lib/filename"
	"github.com/aubinlrx/yaml-path/lib/yamlpath"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().IntP("line", "l", 0, "Line number to display path for")
	getCmd.Flags().IntP("column", "c", 0, "Column number to display path for")
	getCmd.Flags().StringP("separator", "s", ".", "Path separator (default: .)")
	getCmd.Flags().BoolP("remove-first-key", "r", false, "Remove the first key from the path")

	getCmd.MarkFlagRequired("line")
	getCmd.MarkFlagRequired("column")
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [FILE]",
	Short: "Get the yaml path for specific line, position",
	Long:  ``,
	RunE:  yamlGet,
}

func yamlGet(cmd *cobra.Command, args []string) error {
	filename, err := filename.GetFromArgs(args)
	if err != nil {
		return err
	}

	dat, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	line, _ := cmd.Flags().GetInt("line")
	col, _ := cmd.Flags().GetInt("column")
	sep, _ := cmd.Flags().GetString("separator")

	yp := yamlpath.Path{In: dat}
	path, err := yp.PathAtPoint(line, col, sep)
	if err != nil {
		return err
	}

	if cmd.Flag("remove-first-key").Changed {
		// path key.key_2.key_3.key_4
		// remove the first key when flag is set
		path = removeFirstKey(path, sep)
	}

	fmt.Println(path)

	return nil
}

func removeFirstKey(path, sep string) string {
	parts := strings.SplitN(path, sep, 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
