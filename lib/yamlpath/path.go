package yamlpath

import (
	"gopkg.in/yaml.v3"
)

type Path struct {
	In []byte
}

func (yamlPath *Path) List(sep string) (paths []*NodePath, err error) {
	node, err := yamlPath.getNode()
	if err != nil {
		return nil, err
	}

	list(&node, &paths, nil)

	return paths, nil
}

func list(node *yaml.Node, paths *[]*NodePath, previousNodePath *NodePath) {
	l := 0
	if previousNodePath != nil {
		l = len(previousNodePath.Paths)
	}

	path := make([]string, l)
	if previousNodePath != nil {
		copy(path, previousNodePath.Paths)
	}

	switch node.Kind {
	case yaml.DocumentNode:
		for _, childNode := range node.Content {
			list(childNode, paths, nil)
		}

	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]

			keyNodePath := NodePath{
				Paths:  append(path, keyNode.Value),
				Column: keyNode.Column,
				Line:   keyNode.Line,
				Node:   keyNode,
			}

			*paths = append(*paths, &keyNodePath)

			valNode := node.Content[i+1]

			valNodePath := NodePath{
				Paths:  append(path, keyNode.Value),
				Column: valNode.Column,
				Line:   valNode.Line,
				Node:   valNode,
			}

			list(valNode, paths, &valNodePath)
		}
	}
}

func (yamlPath *Path) PathAtPoint(line int, col int, sep string) (path string, err error) {
	node, err := yamlPath.getNode()
	if err != nil {
		return "", err
	}

	path, _ = pathAtPoint(&node, line, col, sep)

	return path, nil
}

func pathAtPoint(node *yaml.Node, line int, col int, separator string) (path string, match bool) {
	switch node.Kind {
	case yaml.DocumentNode:
		// Root Node
		for _, childNode := range node.Content {
			p, m := pathAtPoint(childNode, line, col, separator)

			if m == true {
				return p, true
			}
		}
	case yaml.MappingNode:
		// Map Node
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]

			if nodeMatch(keyNode, line, col) {
				return keyNode.Value, true
			}

			valNode := node.Content[i+1]
			p, m := pathAtPoint(valNode, line, col, separator)

			if m == true {
				if p == "" {
					return keyNode.Value, true
				}

				return keyNode.Value + separator + p, true
			}
		}
	case yaml.ScalarNode:
		if nodeMatch(node, line, col) {
			return "", true
		}
	}

	return "", false
}

func nodeMatch(node *yaml.Node, line int, col int) bool {
	return node.Line == line && node.Column <= col && (node.Column+len(node.Value) > col)
}

func (yp *Path) getNode() (yaml.Node, error) {
	var node yaml.Node

	err := yaml.Unmarshal(yp.In, &node)
	if err != nil {
		return node, err
	}

	return node, nil
}
