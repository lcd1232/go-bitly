package bitly

import "fmt"

type GroupsService struct {
	client *Client
}

type groupResponse struct {
	References       map[string]string `json:"references"`
	Name             string            `json:"name"`
	BSDS             []string          `json:"bsds"`
	IsActive         bool              `json:"is_active"`
	Created          string            `json:"created"`
	Modified         string            `json:"modified"`
	OrganizationGUID string            `json:"organization_guid"`
	Role             string            `json:"role"`
	GUID             string            `json:"guid"`
}

func groupPath(GroupGUID string) string {
	if GroupGUID == "" {
		return "/groups"
	} else {
		return fmt.Sprintf("/groups/%s", GroupGUID)
	}
}

type groupsResponse struct {
	Groups []groupResponse `json:"groups"`
}

func (s *GroupsService) ListGroups(OrganizationGUID string) (*groupsResponse, error) {
	path := versioned(groupPath(""))
	groupsResponse := groupsResponse{}

	//path, err :=

	s.client.get(path, groupsResponse)
}
