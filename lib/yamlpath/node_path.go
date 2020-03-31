package yamlpath

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type NodePath struct {
	Paths  []string
	Column int
	Line   int
	Node   *yaml.Node
}

func (nodePath *NodePath) FullPath() string {
	return nodePath.Path() + ":" + strconv.Itoa(nodePath.Line) + ":" + strconv.Itoa(nodePath.Column)
}

func (nodePath *NodePath) Path() string {
	return strings.Join(nodePath.Paths[:], ".")
}
