package bitly

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGroupsClient_ListGroups(t *testing.T) {
	testCases := []struct {
		desc             string
		responseCode     int
		responseBody     string
		organizationGUID string
		wantErr          string
		wantResult       *groupsResponse
	}{
		{
			desc:             "ok response",
			responseCode:     http.StatusOK,
			responseBody:     `{"groups":[{"created":"2012-12-18T18:14:53+0000","modified":"2016-11-11T21:04:26+0000","bsds":[],"guid":"BcciiJcGgDF","organization_guid":"OssccSr9D4j","name":"test","is_active":true,"role":"org-admin","references":{"organization":"https://api-ssl.bitly.com/v4/organizations/OssccSr9D4j"}}]}`,
			organizationGUID: "",
			wantErr:          "",
			wantResult: &groupsResponse{
				[]groupResponse{
					{
						Created:          "2012-12-18T18:14:53+0000",
						Modified:         "2016-11-11T21:04:26+0000",
						BSDS:             []string{},
						GUID:             "BcciiJcGgDF",
						OrganizationGUID: "OssccSr9D4j",
						Name:             "test",
						IsActive:         true,
						Role:             "org-admin",
						References: map[string]string{
							"organization": "https://api-ssl.bitly.com/v4/organizations/OssccSr9D4j",
						},
					},
				},
			},
		},
		{
			desc:             "ok empty response",
			responseCode:     http.StatusOK,
			responseBody:     `{"groups":[]}`,
			organizationGUID: "test",
			wantErr:          "",
			wantResult: &groupsResponse{
				[]groupResponse{},
			},
		},
		{
			desc:         "invalid token",
			responseCode: http.StatusForbidden,
			responseBody: `{"message":"FORBIDDEN"}`,
			wantErr:      "403 FORBIDDEN",
			wantResult:   nil,
		},
		{
			desc:         "server error",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"message":"some error"}`,
			wantErr:      "500 some error",
			wantResult:   nil,
		},
		{
			desc:         "temporary unavailable",
			responseCode: http.StatusServiceUnavailable,
			responseBody: `{"message":"unavailable"}`,
			wantErr:      "503 unavailable",
			wantResult:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Fatalf("invalid request method: %q", r.Method)
				}
				if r.URL.Path != "/v4/groups" {
					t.Fatalf("invalid request path: %q", r.URL.Path)
				}
				w.WriteHeader(tc.responseCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer s.Close()

			c := NewClient(http.DefaultClient)
			c.BaseURL = s.URL

			got, err := c.Groups.ListGroups(tc.organizationGUID)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("want error %v got %v", tc.wantErr, err)
			}
			if !reflect.DeepEqual(tc.wantResult, got) {
				t.Fatalf("want group %#v got %#v", tc.wantResult, got)
			}
		})
	}
}

func TestGroupsClient_GetBitlinksByGroup(t *testing.T) {
	testCases := []struct {
		desc         string
		responseCode int
		responseBody string
		groupGUID    string
		queryParams  *GetBitlinksByGroupQueryParams
		wantURL      string
		wantErr      string
		wantResult   *getBitlinksByGroupResponse
	}{
		{
			desc:         "ok response",
			responseCode: http.StatusOK,
			responseBody: `{"links":[{"created_at":"2012-12-18T20:16:46+0000","id":"bit.ly/F3zBa5","link":"http://bit.ly/F3zBa5","custom_bitlinks":[],"long_url":"http://example.com/","title":"Example.com Main Page","archived":false,"created_by":"test","client_id":"36b72d37f23e9e247e0aa40083841c92163c5c2f","tags":[],"deeplinks":[],"references":{"group":"https://api-ssl.bitly.com/v4/groups/BcciiJsSgCZ"}},{"created_at":"2012-12-18T18:15:00+0000","id":"on.natgeo.com/WmsHnP","link":"http://on.natgeo.com/WmsHnP","custom_bitlinks":[],"long_url":"http://animals.nationalgeographic.com/animals/fish/pufferfish/","title":"All about Pufferfish","archived":false,"created_by":"test","client_id":"36b72d37f23e9e247e0aa40083841c92163c5c2f","tags":[],"deeplinks":[],"references":{"group":"https://api-ssl.bitly.com/v4/groups/BcciiJsSgCZ"}}],"pagination":{"prev":"","next":"","size":50,"page":1,"total":2}}`,
			groupGUID:    "BcciiJcGgDF",
			queryParams:  nil,
			wantURL:      "/v4/groups/BcciiJcGgDF/bitlinks",
			wantErr:      "",
			wantResult: &getBitlinksByGroupResponse{
				Pagination: Paginate{
					Total: 2,
					Size:  50,
					Prev:  "",
					Page:  1,
					Next:  "",
				},
				Links: []linkResponse{
					{
						References: map[string]string{"group": "https://api-ssl.bitly.com/v4/groups/BcciiJsSgCZ"},
						Archived:   false,
						Tags:       []string{},
						CreatedAt:  "2012-12-18T20:16:46+0000",
						CreatedBy:  "test",
						Title:      "Example.com Main Page",
						LongURL:    "http://example.com/",
						ClientID:   "36b72d37f23e9e247e0aa40083841c92163c5c2f",
					},
					{
						References: map[string]string{"group": "https://api-ssl.bitly.com/v4/groups/BcciiJsSgCZ"},
						Archived:   false,
						Tags:       []string{},
						CreatedAt:  "2012-12-18T18:15:00+0000",
						CreatedBy:  "test",
						Title:      "All about Pufferfish",
						LongURL:    "http://animals.nationalgeographic.com/animals/fish/pufferfish/",
						ClientID:   "36b72d37f23e9e247e0aa40083841c92163c5c2f",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Fatalf("invalid request method: %q", r.Method)
				}
				if r.URL.Path != tc.wantURL {
					t.Fatalf("invalid request path: %q", r.URL.Path)
				}
				w.WriteHeader(tc.responseCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer s.Close()

			c := NewClient(http.DefaultClient)
			c.BaseURL = s.URL

			got, err := c.Groups.GetBitlinksByGroup(tc.groupGUID, tc.queryParams)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("want error %v got %v", tc.wantErr, err)
			}
			if !reflect.DeepEqual(tc.wantResult, got) {
				t.Fatalf("want group %#v got %#v", tc.wantResult, got)
			}
		})
	}
}
