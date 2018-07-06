package bitly

import (
	"fmt"
	"net/url"
	"strings"
)

type GroupsService interface {
	ListGroups(string) (*groupsResponse, error)
	GetGroup(string) (*groupResponse, error)
	GetGroupPreferences(string) (*groupPrefResponse, error)
}

type GroupsClient struct {
	client *Client
}

type groupPrefResponse struct {
	GroupGUID        string `json:"group_guid"`
	DomainPreference string `json:"domain_preference"`
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

type linkResponse struct {
	References map[string]string `json:"references"`
	Archived   bool              `json:"archived"`
	Tags       []string          `json:"tags"`
	CreatedAt  string            `json:"created_at"`
	CreatedBy  string            `json:"created_by"`
	Title      string            `json:"title"`
	LongURL    string            `json:"long_url"`
	ClientID   string            `json:"client_id"`
}

type getBitlinksByGroupResponse struct {
	Pagination Paginate
	Links      []linkResponse
}

// GetBitlinksByGroupQueryParams used by sending query parameters to
type GetBitlinksByGroupQueryParams struct {
	Size            int      `json:"size"`
	Page            int      `json:"page"`
	Keyword         string   `json:"keyword"`
	Query           string   `json:"query"`
	CreatedBefore   int      `json:"created_before"`
	CreatedAfter    int      `json:"created_after"`
	ModifiedAfter   string   `json:"modified_after"`
	Archived        string   `json:"archived"`
	DeepLinks       string   `json:"deeplinks"`
	DomainDeepLinks string   `json:"domain_deeplinks"`
	CampaignGUID    string   `json:"campaign_guid"`
	ChannelGUID     string   `json:"channel_guid"`
	CustomBitlink   string   `json:"custom_bitlink"`
	Tags            []string `json:"tags"`
	EncodingLogin   []string `json:"encoding_login"`
}

func groupPath(GroupGUID string) string {
	return strings.TrimRight(fmt.Sprintf("/groups/%s", GroupGUID), "/")
}

type groupsResponse struct {
	Groups []groupResponse `json:"groups"`
}

func (gc *GroupsClient) ListGroups(OrganizationGUID string) (*groupsResponse, error) {
	path := versioned(groupPath(""))
	groupsResp := &groupsResponse{}

	_, err := gc.client.get(path, groupsResp)
	if err != nil {
		return nil, err
	}

	return groupsResp, nil
}

// GetGroup returns Group info
//
// see - http://dev.bitly.com/v4/#operation/getGroup
func (gc *GroupsClient) GetGroup(GroupGUID string) (*groupResponse, error) {
	path := versioned(groupPath(GroupGUID))
	groupResp := &groupResponse{}

	_, err := gc.client.get(path, groupResp)
	if err != nil {
		return nil, err
	}
	return groupResp, nil
}

// GetGroupPreferences returns Group preferences
//
// see - http://dev.bitly.com/v4/#operation/getGroupPreferences
func (gc *GroupsClient) GetGroupPreferences(GroupGUID string) (*groupPrefResponse, error) {
	path := versioned(groupPath(GroupGUID) + "/preferences")
	groupPrefResp := &groupPrefResponse{}

	_, err := gc.client.get(path, groupPrefResp)
	if err != nil {
		return nil, err
	}
	return groupPrefResp, nil
}

// GetBitlinksByGroup retrieves a paginated collection of Bitlinks for a Group
//
// see - http://dev.bitly.com/v4/#operation/getBitlinksByGroup
func (gc *GroupsClient) GetBitlinksByGroup(GroupGUID string, queryParams *GetBitlinksByGroupQueryParams) error {
	if queryParams != nil {
		q, err := encoder.Encode(queryParams)
		if err != nil {
			return err
		}

	}
	return nil
}
