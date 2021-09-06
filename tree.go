package main

type i3TreeNode struct {
	ID               int64             `json:"id"`
	Type             string            `json:"type"`
	Nodes            []*i3TreeNode     `json:"nodes"`
	Name             string            `json:"name"`
	WindowProperties *i3TreeNodeWindow `json:"window_properties"`
	Num              int               `json:"num"`
}

type i3TreeNodeWindow struct {
	Class    string `json:"class"`
	Instance string `json:"instance"`
	Title    string `json:"title"`
}

func (i *i3TreeNode) retrieveWorkspaces() map[string]*i3Workspace {
	ret := make(map[string]*i3Workspace)
	if i.Nodes == nil {
		return ret
	}
	for _, n := range i.Nodes {
		n.retrieveWorkspaceInformation(ret)
	}
	return ret
}

func (i *i3TreeNode) retrieveWorkspaceInformation(v map[string]*i3Workspace) {
	if i.Type != "workspace" {
		if i.Nodes != nil {
			for _, n := range i.Nodes {
				n.retrieveWorkspaceInformation(v)
			}
		}
		return
	}
	v[i.Name] = &i3Workspace{
		Name: i.Name,
		Num:  i.Num,
	}
}

// Retrieve each workspace tree window information
func (i *i3TreeNode) retrieveWorkspacesInformation() map[string][]*i3TreeNodeWindow {
	ret := make(map[string][]*i3TreeNodeWindow)

	if i.Nodes == nil {
		return ret
	}
	for _, node := range i.Nodes {
		node.retrieveNodeWorkspaceInformation(i, ret)
	}
	return ret
}

func (i *i3TreeNode) retrieveNodeWorkspaceInformation(parent *i3TreeNode, v map[string][]*i3TreeNodeWindow) {
	if parent.Type == "workspace" && i.Type == "con" && i.WindowProperties != nil {
		if _, ok := v[parent.Name]; ok {
			v[parent.Name] = append(v[parent.Name], i.WindowProperties)
		} else {
			v[parent.Name] = []*i3TreeNodeWindow{i.WindowProperties}
		}
	}

	if i.Nodes == nil {
		return
	}
	for _, node := range i.Nodes {
		node.retrieveNodeWorkspaceInformation(i, v)
	}
}
