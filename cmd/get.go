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
	"errors"
	"fmt"

	"github.com/aubinlrx/lib/yamlpath"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().IntP("line", "l", 0, "Line number to display path for")
	getCmd.Flags().IntP("column", "c", 0, "Column number to display path for")

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

func filenameFromArgs(args []string) (string, error) {
	var filename string

	if len(args) > 0 {
		filename = strings.TrimSpace(args[0])
	}

	if filename == "" {
		return "", errors.New("Please provide a valid filename")
	}

	return filename, nil
}

func yamlGet(cmd *cobra.Command, args []string) error {
	filename, err := filenameFromArgs(args)
	if err != nil {
		return err
	}

	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var node yaml.Node
	err = yaml.Unmarshal(dat, &node)
	if err != nil {
		return err
	}

	line, _ := cmd.Flags().GetInt("line")
	column, _ := cmd.Flags().GetInt("column")

	// paths := [][]string{}
	// path := make([]string, 0)
	// list(&node, path, &paths, nil, nil, nil)

	// res := make([]string, len(paths))
	// for i, v := range paths {
	// 	res[i] = strings.Join(v[:], ".")
	// }

	path, _ := pathAtPoint(line, column, &node)

	fmt.Println(path)

	return nil
}

func pathAtPoint(line int, col int, node *yaml.Node) (path string, match bool) {
	switch node.Kind {
	case yaml.DocumentNode:
		// Root Node
		for _, childNode := range node.Content {
			a, m := pathAtPoint(line, col, childNode)

			if m == true {
				return a, true
			}
		}
	case yaml.MappingNode:
		// Map Node
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]

			if nodeMatch(line, col, keyNode) {
				return keyNode.Value, true
			}

			valNode := node.Content[i+1]
			a, m := pathAtPoint(line, col, valNode)
			if m == true {
				if a == "" {
					return keyNode.Value, true
				}
				return keyNode.Value + "." + a, true
			}
		}
	case yaml.SequenceNode:
		// Array Node
		fmt.Println("Array not handled")
	case yaml.ScalarNode:
		if nodeMatch(line, col, node) {
			return "", true
		}
	}

	return "", false
}

func nodeMatch(line int, col int, node *yaml.Node) bool {
	return node.Line == line && node.Column <= col && (node.Column+len(node.Value) > col)
}

func list(node *yaml.Node, previousPath []string, paths *[][]string, nodeList []*yaml.Node, parentNode *yaml.Node, parentNodeList []*yaml.Node) {
	path := make([]string, len(previousPath))
	copy(path, previousPath)

	if node.Kind == yaml.MappingNode {
		if previousNode := previousNode(node, nodeList); previousNode != nil {
			path = append(path, previousNode.Value)
			*paths = append(*paths, path)
		}
	} else if parentNode != nil && parentNode.Kind == yaml.SequenceNode {
		if parentPreviousNode := previousNode(parentNode, parentNodeList); parentPreviousNode != nil {
			path = append(path, parentPreviousNode.Value)
			*paths = append(*paths, path)
		}
	} else if nextNode := nextNode(node, nodeList); isNodeKey(node, nextNode) {
		path = append(path, node.Value)
		*paths = append(*paths, path)
	}

	if len(node.Content) == 0 {
		return
	}

	for _, childNode := range node.Content {
		list(childNode, path, paths, node.Content, node, nodeList)
	}
}

func previousNode(node *yaml.Node, nodeList []*yaml.Node) *yaml.Node {
	index := nodeIndex(node, nodeList)

	if index >= 1 {
		return nodeList[index-1]
	}

	return nil
}

func nextNode(node *yaml.Node, nodeList []*yaml.Node) *yaml.Node {
	index := nodeIndex(node, nodeList)

	if index >= 0 && index+1 < len(nodeList) {
		return nodeList[index+1]
	}

	return nil
}

func nodeIndex(node *yaml.Node, nodeList []*yaml.Node) int {
	index := -1

	for nodeIndex, currentNode := range nodeList {
		if currentNode == node {
			index = nodeIndex
			break
		}
	}

	return index
}

func parseChildren(node *yaml.Node) {
	// fmt.Println(node.Value, "line:", node.Line, "column:", node.Column, "type:", nodeType(node))

	if len(node.Content) != 0 {
		var types []string

		for childIndex, childNode := range node.Content {
			parseChildren(childNode)
			var siblingNode *yaml.Node

			if childIndex+1 < len(node.Content) {
				siblingNode = node.Content[childIndex+1]
			}

			if siblingNode != nil && isNodeKey(childNode, siblingNode) {
				fmt.Println("isKey", childNode.Value)
			}

			types = append(types, childNode.Value, nodeType(childNode))
		}

		fmt.Println(types)
	}
}

func isNodeKey(node *yaml.Node, siblingNode *yaml.Node) bool {
	if siblingNode == nil {
		return false
	} else if siblingNode.Kind == yaml.ScalarNode && siblingNode.Line == node.Line {
		return true
	} else {
		return false
	}
}

func nodeType(node *yaml.Node) string {
	var nodeType string

	switch node.Kind {
	case yaml.ScalarNode:
		nodeType = "scalar"
	case yaml.MappingNode:
		nodeType = "mapping"
	default:
		nodeType = "..."
	}

	return nodeType
}
