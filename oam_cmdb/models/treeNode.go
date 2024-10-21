package models

import (
	fn "OAM/util"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type TreeGridNode[T HostTreeGridNode | FunTreeGridNode | MenuNode] struct {
	Id       int
	Text     string
	State    string
	checked  bool
	Children []*T `json:",omitempty"`
}

type HostTreeGridNode struct {
	TreeGridNode[HostTreeGridNode]
	OriginalId int
	Uri        string
	Port       int `json:",omitempty"`
	Desc       string
	Type       uint8
	CreateTime time.Time
}

type FunTreeGridNode struct {
	TreeGridNode[FunTreeGridNode]
	FunCode  string
	FunType  int
	FunOrder uint16
	FunUrl   string
	MenuCls  string
	ParentId int
}

type MenuNode struct {
	Text     string
	State    string
	Url      string
	IconCls  string
	Children []*MenuNode `json:",omitempty"`
}

const (
	state_open  = "open"
	state_close = "closed"
)

func (parent *HostTreeGridNode) AddNode(child *HostTreeGridNode) {
	var baseId = parent.Id
	if baseId == 0 {
		baseId = 1
	}
	if len(parent.Children) == 0 {
		child.Id = baseId * 10 //生成一个唯一ID,不用数据库中ID,因为hostId,AppId可能重复,导致treegrid行无法选中
		children := []*HostTreeGridNode{child}
		parent.Children = children
	} else {
		child.Id = baseId*10 + len(parent.Children) + 1
		parent.Children = append(parent.Children, child)
	}
}

func (parent HostTreeGridNode) GetChild(originalId int) *HostTreeGridNode {
	for _, node := range parent.Children {
		if node.OriginalId == originalId {
			return node
		}
	}
	return nil
}

func (parent *HostTreeGridNode) AddHostNode(child Host) *HostTreeGridNode {
	uri := fn.JoinStr("/", child.PublicIp, child.InternalIp)
	if uri[0] == '/' {
		uri = uri[1:]
	}
	childNode := HostTreeGridNode{
		TreeGridNode: TreeGridNode[HostTreeGridNode]{Text: child.HostName, State: state_open},
		OriginalId:   child.HostId,
		Uri:          uri,
		Port:         child.SshPort,
		CreateTime:   child.CreateTime,
		Type:         0,
	}
	/* if child.ServiceSoftwares != "" {
		childNode.Desc = "已安装软件:" + child.ServiceSoftwares
	} */
	parent.AddNode(&childNode)
	return &childNode
}

func (parent *HostTreeGridNode) AddAppNode(child AppInfo) *HostTreeGridNode {
	childNode := HostTreeGridNode{
		TreeGridNode: TreeGridNode[HostTreeGridNode]{Text: child.AppName},
		OriginalId:   child.AppId,
		Uri:          child.AppDir,
		Port:         child.AppPort,
		CreateTime:   child.CreateTime,
		Type:         1,
	}
	if child.AppUrl != "" {
		childNode.Desc = child.AppUrl
	}
	parent.AddNode(&childNode)
	return &childNode
}

func BuildFunTree(funs []*Fun, parentId int, roleFunIds []int) []*FunTreeGridNode {
	var nodes []*FunTreeGridNode
	for _, f := range funs {
		if f.ParentId == parentId {
			node := f.ToTreeNode()
			if len(roleFunIds) > 0 && slices.Contains(roleFunIds, f.FunId) {
				node.checked = true
			}
			//node.MenuCls = f.MenuClass
			childs := BuildFunTree(funs, f.FunId, roleFunIds)
			if childs != nil {
				node.Children = childs
			}

			if f.FunLevel < 2 || node.Children == nil {
				node.State = state_open
			} else {
				node.State = state_close
			}
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (f *Fun) ToTreeNode() *FunTreeGridNode {
	node := new(FunTreeGridNode)
	node.Id = f.FunId
	node.Text = f.FunName
	node.FunCode = f.FunCode
	node.FunType = f.FunType
	node.FunUrl = f.FunUrl
	node.FunOrder = f.FunOrder
	node.ParentId = f.ParentId
	return node
}

func BuildMenuTree(funs []*Fun, parentId int) []*MenuNode {
	var nodes []*MenuNode
	for _, f := range funs {
		if f.ParentId == parentId {
			node := f.ToMenuNode()
			childs := BuildMenuTree(funs, f.FunId)
			if childs != nil {
				node.Children = childs
			}
			nodes = append(nodes, node)
		}
	}

	return nodes

}

func (f *Fun) ToMenuNode() *MenuNode {
	node := new(MenuNode)
	node.Text = f.FunName
	node.IconCls = f.MenuClass
	if len(f.FunUrl) > 0 {
		splitIdx := strings.Index(f.FunUrl, ",")
		if splitIdx > 1 {
			node.Url = f.FunUrl[0:splitIdx]
		} else {
			node.Url = f.FunUrl
		}
	}
	if f.FunLevel <= 1 && f.FunOrder == 1 {
		node.State = state_open
	} else {
		node.State = state_close
	}

	return node
}
