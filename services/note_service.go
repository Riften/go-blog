package services

import (
	"encoding/json"
	"errors"
	logging "github.com/ipfs/go-log"
	"go-blog/common"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var log = logging.Logger("services")

// NoteService define the note service related interface.
// NoteService should be thread safe.
type NoteService interface {
	// Fetch fetch a node with relative path.
	// It would return a copy of original node with limited node info if copy is true.
	Fetch(relative string, copy bool) *NoteTreeNode	// Return nil if not found.

	// FetchAll is equal to call Fetch with "relative=/" and "copy=false".
	FetchAll() *NoteTreeNode
	LoadFromDisk() error
	WriteBack() error	// Write the note service to store.
	// Upload
	// WriteBack
}

type NoteTreeNode struct {
	Links  map[string]*NoteTreeNode
	parent *NoteTreeNode
	//Data   *NoteInfo
	IsDir bool
	Name string
	RawPath string
	RenderedPath string
	Abstract string
}

type RefreshOption struct {
	Recursive bool  // Whether to refresh recursively.
	Render bool		// Whether to render html files.
	OverWrite bool	// Whether to over-write rendered files.
	CopyOthers bool	// Whether to copy non-markdown files to cache directory.
}

var DefaultAddOption = &RefreshOption{
	Recursive:  true,
	Render:     true,
	OverWrite:  false,
	CopyOthers: true,
}

func (n *NoteTreeNode) ToJsonString() string {
	res, err := json.MarshalIndent(n, "", "\t")
	if err != nil {
		log.Error("Error when convert NoteInfo to json: ", err)
	}
	return string(res)
}

func (n *NoteTreeNode) ToJsonBytes() []byte {
	res, err := json.MarshalIndent(n, "", "\t")
	if err != nil {
		log.Error("Error when convert NoteInfo to json: ", err)
	}
	return res
}

func (n *NoteTreeNode) ToJsonWrite(writer io.Writer) (int, error) {
	res, err := json.MarshalIndent(n, "", "\t")
	if err != nil {
		log.Error("Error when convert NoteInfo to json: ", err)
		return 0, err
	}
	return writer.Write(res)
}

func (n *NoteTreeNode) LightCopy() *NoteTreeNode {
	return &NoteTreeNode{
		IsDir:        n.IsDir,
		Name:         n.Name,
		RawPath:      n.RawPath,
		RenderedPath: n.RenderedPath,
		Abstract:     n.Abstract,
	}
}

// Add add a new node for node n from a relative path.
// For example:
//		n: root/a
//		current tree:
//			- root
//				- a
//					-b
//		relative: b/c/d
//		Add would first try to reach as deep as much from "n:root/a" which is node "root/a/b".
//		Then it would create node "root/a/b/c" and "root/a/b/c/d".
//		If "d" is a directory and recursive is true, Add would recursively add child nodes.
//		If render is true, Add would render markdown files into cache directory derived from parent node.
func (n *NoteTreeNode) Add(relative string, name string, option *RefreshOption) error {
	entries := strings.Split(relative, "/")
	current := n
	var tmpNode *NoteTreeNode
	var ok bool
	var err error
	if len(entries) > 1 {
		for _, entry := range entries[0 : len(entries)-1] {
			if entry != "" {
				tmpNode, ok = current.Links[entry]
				if !ok {
					tmpNode, err = current.deriveNode(entry, entry, true)
					if err != nil {
						return err
					}
				}
				current = tmpNode
			}
		}
	}

	entry := entries[len(entries)-1]
	if entry != "" {
		tmpNode, ok = current.Links[entry]
		if !ok {
			tmpNode, err = current.deriveNode(entry, name, true)
			if err != nil {
				return err
			}
		}
		current = tmpNode
	}

	if option.Render {
		if !current.IsDir {
			return common.MdRenderFile(current.RawPath, current.RenderedPath)
		} else {
			return os.Mkdir(current.RenderedPath, os.ModePerm)
		}
	}
	return nil
}

// deriveNode generate a new node from given node.
// Link n and new node if link is true.
func (n *NoteTreeNode) deriveNode(relativePath string, name string, link bool) (*NoteTreeNode, error) {
	if strings.Contains(relativePath, "/") {
		return nil, errors.New(relativePath + " has '/'")
	}

	if n.hasNode(name) {
		return nil, errors.New(n.getPath() + "/" + name + " already exists")
	}

	rawPath := filepath.Join(n.RawPath, relativePath)
	fi, err := os.Stat(rawPath)
	if err != nil {
		return nil, err
	}

	var renderPath = filepath.Join(n.RenderedPath, relativePath)
	if !fi.IsDir() {
		if filepath.Ext(renderPath) != ".md" {
			return nil, errors.New(rawPath + " is not a markdown file")
		}
		renderPath = common.ChExt(renderPath, ".html")
	}

	newNode := &NoteTreeNode{
		Links:  make(map[string]*NoteTreeNode),
		parent: nil,

		IsDir:        fi.IsDir(),
		Name:         name,
		RawPath:      rawPath,
		RenderedPath: renderPath,
		Abstract:     "",
	}
	if link {
		newNode.parent = n
		n.Links[name] = newNode
	}
	return newNode, nil
}

func (n *NoteTreeNode) hasNode(name string) bool {
	_, ok := n.Links[name]
	return ok
}

func (n *NoteTreeNode) getPath() string {
	if n.parent == nil {
		return "/" + n.Name
	} else {
		return n.parent.getPath() + "/" + n.Name
	}
}

// Note Service implemented based on file system.
// No database is required.
type fsNoteService struct {
	root *NoteTreeNode

	lock sync.RWMutex
}

func NewFsNoteService(cacheDir string, rootDir string) NoteService {
	return &fsNoteService{
		root: &NoteTreeNode{
			Links:        make(map[string]*NoteTreeNode),
			parent:       nil,
			IsDir:        true,
			Name:         ".",
			RawPath:      rootDir,
			RenderedPath: cacheDir,
			Abstract:     "",
		},
		lock: sync.RWMutex{},
	}
}

func (ns *fsNoteService) Fetch(relative string, copy bool) *NoteTreeNode {
	ns.lock.RLock()
	defer ns.lock.RUnlock()
	entries := strings.Split(relative, "/")
	rawNode := ns.root.walkTo(entries, 0)
	if copy {
		newNode := rawNode.LightCopy()
		if rawNode.Links != nil {
			newNode.Links = make(map[string]*NoteTreeNode)
			for name, tmpNode := range rawNode.Links {
				newNode.Links[name] = tmpNode.LightCopy()
			}
		}
		if rawNode.parent != nil {
			newNode.parent = rawNode.parent.LightCopy()
		}
		return newNode
	} else {
		return rawNode
	}
}

func (ns *fsNoteService) FetchAll() *NoteTreeNode{
	return ns.root
}

// walkTo may return nil while can not find the entry.
func (n *NoteTreeNode) walkTo(entries []string, index int) *NoteTreeNode {
	for index < len(entries) && entries[index] == "" {
		index++
	}
	if index == len(entries) -1 {
		return n
	}else if index < len(entries) {
		nextNode, ok := n.Links[entries[index]]
		if !ok {
			return nil
		} else {
			return nextNode.walkTo(entries, index+1)
		}
	} else {
		return nil
	}
}

func (ns *fsNoteService) LoadFromDisk() error {
	return nil
}

func (ns *fsNoteService) WriteBack() error {
	return nil
}
