package main

type i3TreeNode struct {
	ID               int64             `json:"id"`
	Type             string            `json:"type"`
	Nodes            []*i3TreeNode     `json:"nodes"`
	Name             string            `json:"name"`
	WindowProperties *i3TreeNodeWindow `json:"window_properties"`
}

type i3TreeNodeWindow struct {
	Class    string `json:"class"`
	Instance string `json:"instance"`
	Title    string `json:"title"`
}

// Retrieve each workspace tree window information
func (i *i3TreeNode) retrieveWorkspacesInformation() map[int64][]string {
	ret := make(map[int64][]string)

	if i.Nodes == nil {
		return ret
	}
	for _, node := range i.Nodes {
		node.retrieveNodeWorkspaceInformation(i, ret)
	}
	return ret
}

func (i *i3TreeNode) retrieveNodeWorkspaceInformation(parent *i3TreeNode, v map[int64][]string) {
	if parent.Type == "workspace" && i.Type == "con" && i.WindowProperties != nil {
		if _, ok := v[parent.ID]; ok {
			v[parent.ID] = append(v[parent.ID], i.WindowProperties.Class)
		} else {
			v[parent.ID] = []string{i.WindowProperties.Class}
		}
	}

	if i.Nodes == nil {
		return
	}
	for _, node := range i.Nodes {
		node.retrieveNodeWorkspaceInformation(i, v)
	}
}
