package bitly

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"strings"
)

type GroupsService interface {
	ListGroups(string) (*groupsResponse, error)
	GetGroup(string) (*groupResponse, error)
	GetGroupPreferences(string) (*groupPrefResponse, error)
	GetBitlinksByGroup(GroupGUID string, queryParams *GetBitlinksByGroupQueryParams) (*getBitlinksByGroupResponse, error)
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

type deepLink struct {
	BitLink     string `json:"bitlink"`
	InstallURL  string `json:"install_url"`
	Created     string `json:"created"`
	AppURIPath  string `json:"app_uri_path"`
	Modified    string `json:"modified"`
	InstallType string `json:"install_type"`
	AppGUID     string `json:"app_guid"`
	GUID        string `json:"guid"`
	OS          string `json:"os"`
}

type linkResponse struct {
	References map[string]string `json:"references"`
	Archived   bool              `json:"archived"`
	Tags       []string          `json:"tags"`
	CreatedAt  string            `json:"created_at"`
	CreatedBy  string            `json:"created_by"`
	Title      string            `json:"title"`
	DeepLinks  []deepLink        `json:"deep_links"`
	LongURL    string            `json:"long_url"`
	ClientID   string            `json:"client_id"`
}

type getBitlinksByGroupObject struct {
	url      string
	Resp     *getBitlinksByGroupResponse
	isLoaded bool
	client   *Client
}

func (o *getBitlinksByGroupObject) Next() bool {
	if nextURL := o.Resp.Pagination.Next; nextURL != "" {
		o.url = nextURL
		o.isLoaded = false
		return true
	}
	return false
}

func (o *getBitlinksByGroupObject) Prev() bool {
	if prevURL := o.Resp.Pagination.Prev; prevURL != "" {
		o.url = prevURL
		o.isLoaded = false
		return true
	}
	return false
}

func (o *getBitlinksByGroupObject) Get() error {
	if o.isLoaded {
		return nil
	}
	_, err := o.client.get(o.url, o.Resp)
	return err
}

type getBitlinksByGroupResponse struct {
	Pagination Paginate
	Links      []linkResponse
}

// GetBitlinksByGroupQueryParams used by sending query parameters to
type GetBitlinksByGroupQueryParams struct {
	Size            int         `url:"size,omitempty"`
	Page            int         `url:"page,omitempty"`
	Keyword         string      `url:"keyword,omitempty"`
	Query           string      `url:"query,omitempty"`
	CreatedBefore   int         `url:"created_before,omitempty"`
	CreatedAfter    int         `url:"created_after,omitempty"`
	ModifiedAfter   string      `url:"modified_after,omitempty"`
	Archived        queryOption `url:"archived,omitempty"`
	DeepLinks       queryOption `url:"deeplinks,omitempty"`
	DomainDeepLinks queryOption `url:"domain_deeplinks,omitempty"`
	CampaignGUID    string      `url:"campaign_guid,omitempty"`
	ChannelGUID     string      `url:"channel_guid,omitempty"`
	CustomBitlink   queryOption `url:"custom_bitlink,omitempty"`
	Tags            []string    `url:"tags,omitempty"`
	EncodingLogin   []string    `url:"encoding_login,omitempty"`
}

type listGroupsParams struct {
	OrganizationGUID string `url:"organization_guid,omitempty"`
}

func groupPath(GroupGUID string) string {
	return strings.TrimRight(fmt.Sprintf("/groups/%s", GroupGUID), "/")
}

type groupsResponse struct {
	Groups []groupResponse `json:"groups"`
}

func (gc *GroupsClient) ListGroups(OrganizationGUID string) (*groupsResponse, error) {
	q, err := query.Values(&listGroupsParams{OrganizationGUID})
	if err != nil {
		return nil, err
	}

	path, err := buildURL(versioned(groupPath("")), q)
	if err != nil {
		return nil, err
	}

	groupsResp := &groupsResponse{}

	_, err = gc.client.get(path, groupsResp)
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
func (gc *GroupsClient) GetBitlinksByGroup(GroupGUID string, queryParams *GetBitlinksByGroupQueryParams) (*getBitlinksByGroupResponse, error) {
	getBitlinksByGroupResp := &getBitlinksByGroupResponse{}
	var path string
	if queryParams != nil {
		q, err := query.Values(queryParams)
		if err != nil {
			return nil, err
		}
		path, err = buildURL(versioned(groupPath(GroupGUID)+"/bitlinks"), q)
		if err != nil {
			return nil, err
		}
	} else {
		path = versioned(groupPath(GroupGUID) + "/bitlinks")
	}
	_, err := gc.client.get(path, getBitlinksByGroupResp)
	if err != nil {
		return nil, err
	}

	return getBitlinksByGroupResp, nil
}
