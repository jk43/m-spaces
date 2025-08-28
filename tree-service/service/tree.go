package service

import (
	"fmt"
	"regexp"
	"strings"
	"tree/models"

	"github.com/moly-space/molylibs/utils"
)

/*

Example of a tree structure

A
| - AA
|    | - AAA
|    | - AAB
|    |    | - AABA
|    |    | - AABB
| - AB
|    | - ABA
|		 |    | - ABAA
|		 |		| - ABAB
|		 |    |    | - ABABA
|		 |    |    | - ABABB
|    | - ABB
| - AC

*/

type RetrieveType = string

const (
	Single    RetrieveType = "single"
	All       RetrieveType = "all"
	Ancestors RetrieveType = "ancestors"
	Children  RetrieveType = "children"
	Sibling   RetrieveType = "sibling"
	Parent    RetrieveType = "parent"
)

type TreeRequest struct {
	Id          uint               `json:"id"`
	Slug        string             `json:"slug"`
	OrgID       string             `json:"organization_id"`
	Label       string             `json:"label"`
	View        string             `json:"view"`
	Description string             `json:"description"`
	ParentID    uint               `json:"parent_id"`
	RootID      uint               `json:"group_id"`
	TreeOptions models.TreeOptions `json:"options"`
	UserID      string             `json:"user_id"`
	db          models.DatabaseRepo
	//src         *models.Tree
}

type BatchTreeRequest struct {
	Label       string             `json:"label"`
	Description string             `json:"description"`
	View        string             `json:"view"`
	Children    []BatchTreeRequest `json:"children"`
}

func (t *TreeRequest) SetDB(db models.DatabaseRepo) {
	t.db = db
}

func (t *TreeRequest) ValidateOrg(slug string, orgID string) error {
	_, err := t.db.GetTreeWithSlug(orgID, slug)
	if err != nil {
		return err
	}
	return nil
}

func (t *TreeRequest) Insert() (uint, string, error) {
	if t.Label == "" {
		return 0, "", fmt.Errorf("label is required")
	}
	//utils.TermDebugging(`t.db`, t.db)
	slug := getSlug(t.Label)
	tail := ""
	i := 0
	for {
		_, err := t.db.GetTreeAttributeWithSlug(slug + tail)
		if err == nil {
			if i != 0 {
				tail = fmt.Sprintf("-%d", i)
			}
		} else {
			break
		}
		i++
	}
	slug = slug + tail
	tree := &models.Tree{}
	// add child
	if t.Slug != "" {
		parent, err := t.db.GetTreeWithSlug(t.OrgID, t.Slug)
		if err != nil {
			return 0, "", err
		}
		tree.ParentID = &parent.ID
		tree.RootID = parent.RootID
	}
	tree, err := t.db.InsertTree(tree)
	if err != nil {
		return 0, "", err
	}
	_, err = t.db.InsertTreeAttribute(&models.TreeAttribute{
		Label:          t.Label,
		Slug:           slug,
		View:           t.View,
		Description:    t.Description,
		TreeID:         tree.ID,
		OrganizationID: t.OrgID,
		CreatedBy:      t.UserID,
	})
	if err != nil {
		t.db.DeleteTreeWithID(tree.ID)
		return 0, "", err
	}
	if t.Slug == "" {
		tree.RootID = tree.ID
		err = t.db.UpdateTreeWithID(tree.ID, tree)
	} else {
		tree.RootID = t.RootID
		tree.ParentID = &t.ParentID
	}
	if err != nil {
		return 0, "", err
	}
	return tree.ID, slug, nil
}

// ABABA = ABABA
// ABA = ABA
func (t *TreeRequest) GetNode() (*models.Tree, error) {
	var err error
	node, err := t.db.GetTreeWithSlug(t.OrgID, t.Slug)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (t *TreeRequest) GetAllNodes() (*models.Tree, error) {
	node, err := t.db.GetTreeWithSlug(t.OrgID, t.Slug)
	if err != nil {
		return nil, err
	}
	collectNodes(node, t.db)
	return node, nil
}

func collectNodes(node *models.Tree, db models.DatabaseRepo) {
	nodes, _ := db.GetChildrenWithSlug(node.Attributes.OrganizationID, node.Attributes.Slug)
	for _, n := range nodes {
		node.Children = append(node.Children, n)
		collectNodes(n, db)
	}
}

func (t *TreeRequest) GetAllNodesAsArray() ([][]*models.Tree, error) {
	asts, err := t.GetAncestors()
	if err != nil {
		return nil, err
	}
	for _, ast := range asts {
		utils.TermDebugging(`ast`, ast)
	}
	output := [][]*models.Tree{}
	collectNodesAsArray(&output, asts[len(asts)-1], t.db)
	return output, nil
}

func collectNodesAsArray(o *[][]*models.Tree, node *models.Tree, db models.DatabaseRepo) {
	nodes, _ := db.GetChildrenWithSlug(node.Attributes.OrganizationID, node.Attributes.Slug)
	if len(nodes) == 0 {
		t := TreeRequest{
			OrgID: node.Attributes.OrganizationID,
			Slug:  node.Attributes.Slug,
		}
		t.SetDB(db)
		asts, _ := t.GetAncestors()
		*o = append(*o, asts)
	}
	for _, n := range nodes {
		collectNodesAsArray(o, n, db)
	}
}

func (t *TreeRequest) GetChildren() ([]*models.Tree, error) {
	nodes, err := t.db.GetChildrenWithSlug(t.OrgID, t.Slug)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (t *TreeRequest) GetAncestors() ([]*models.Tree, error) {
	// Get the ancestors of the given node
	// and return the list of ancestors
	node, err := t.db.GetTreeWithSlug(t.OrgID, t.Slug)
	if err != nil {
		return nil, err
	}
	output := []*models.Tree{}
	output = append(output, node)
	for {
		if node.ParentID == nil {
			break
		}
		node, err = t.db.GetTreeWithID(t.OrgID, *node.ParentID)
		if err != nil {
			return nil, err
		}
		output = append(output, node)
	}
	// reverse the order
	for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}
	return output, nil
}

// ABABA = true
// ABA = false
func (t *TreeRequest) IsEnd() (bool, error) {
	// Check if the given node is a leaf node
	// and return true if it is a leaf node
	node, err := t.db.GetChildrenWithSlug(t.OrgID, t.Slug)
	if err != nil {
		return false, err
	}
	if len(node) == 0 {
		return true, nil
	}
	return false, nil
}

func (t *TreeRequest) GetSibling() ([]*TreeRequest, error) {
	// Get the siblings of the given node
	// and return the list of siblings
	return nil, nil
}

func (t TreeRequest) GetItems(f func()) ([]*TreeRequest, error) {

	return nil, nil
}

func (t *TreeRequest) DeleteTree() error {
	tree, err := t.db.GetTreeWithSlug(t.OrgID, t.Slug)
	if err != nil {
		return err
	}
	err = t.db.DeleteTree(tree)
	if err != nil {
		return err
	}
	return nil
}

func (t *TreeRequest) BatchProcess(tree *BatchTreeRequest, slug string) error {
	for _, child := range tree.Children {
		t.Label = child.Label
		t.View = ""
		t.Description = child.Description
		t.Slug = slug
		_, slug, err := t.Insert()
		if err != nil {
			return err
		}
		t.BatchProcess(&child, slug)
	}
	return nil
}

func getSlug(input string) string {
	reg, err := regexp.Compile(`[^\p{L}\p{N}\s]+`)
	if err != nil {
		fmt.Println("Regex error:", err)
		return ""
	}
	cleanStr := reg.ReplaceAllString(input, "")
	cleanStr = strings.ReplaceAll(cleanStr, " ", "-")
	cleanStr = strings.ToLower(cleanStr)

	return cleanStr
}
