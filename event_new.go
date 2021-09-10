package main

import (
	"fmt"
	"net"
	"strings"
)

var (
	windowClassShortNames = map[string]string{
		"gimp":          "✎ Gimp",
		"clockify":      "🗒  Clockify",
		"google-chrome": "◎ Chrome",
		"st":            "▱ Terminal",
		"discord":       "🗪 Discord",
		"spotify":       "🎵 Spotify",
	}
	windowNameShortNames = map[string]string{
		"vim":  "▤ Vim",
		"gimp": "✎ Gimp",
	}
)

type eventNewHandler struct {
	conn  net.Conn
	chans *eventHandlerChannels
}

func (e *eventNewHandler) handle(msg *i3Message, event *i3MessageEvent) error {
	tree, err := retrievei3Tree(e.chans.getTreeCh, e.conn)
	if err != nil {
		return err
	}

	treeData, err := tree.decodeTreeData()
	if err != nil {
		return err
	}
	contents := treeData.retrieveWorkspacesInformation()
	workspaces := treeData.retrieveWorkspaces()
	inUse := make([]int, len(workspaces))

	// Loop tree data to rename workspaces based on contents
	for workspace, content := range contents {
		nameList := retrieveWorkspaceNamesFromContent(content)
		if len(nameList) <= 0 {
			continue
		}

		// Retrieve workspace data
		workspaceData, ok := workspaces[workspace]
		if !ok {
			continue
		}

		var cmd string
		if workspaceData.Num < 0 {
			n := findWorkspaceNumberNotInUse(inUse, workspaces)
			inUse = append(inUse, n)

			if fmt.Sprintf("%d:%d %s", n, n, strings.Join(nameList, " | ")) == workspace {
				continue
			}
			cmd = fmt.Sprintf("rename workspace \"%s\" to \"%d:%d %s\"", workspace, n, n, strings.Join(nameList, " | "))
		} else {
			if fmt.Sprintf("%d:%d %s", workspaceData.Num, workspaceData.Num, strings.Join(nameList, " | ")) == workspace {
				continue
			}
			cmd = fmt.Sprintf("rename workspace \"%s\" to \"%d:%d %s\"", workspace, workspaceData.Num, workspaceData.Num, strings.Join(nameList, " | "))
		}
		if err := sendi3Command(cmd, e.conn); err != nil {
			return err
		}
	}
	return nil
}

func findWorkspaceNumberNotInUse(inUse []int, workspaces map[string]*i3Workspace) int {
mainLoop:
	for i := 1; i < 9; i++ {
		for _, n := range inUse {
			if n == i {
				continue mainLoop
			}
		}
		for _, w := range workspaces {
			if w.Num == i {
				continue mainLoop
			}
		}
		return i
	}
	return 0
}

// Retrieve the new workspace name
func retrieveWorkspaceNamesFromContent(content []*i3TreeNodeWindow) []string {
	ret := make([]string, 0, len(content))
	for _, w := range content {
		if v, ok := windowNameShortNames[strings.ToLower(w.Title)]; ok {
			ret = append(ret, v)
			continue
		}
		if v, ok := windowClassShortNames[strings.ToLower(w.Class)]; ok {
			ret = append(ret, v)
			continue
		}
	}
	return ret
}
