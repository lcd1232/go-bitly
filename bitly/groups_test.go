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

			c := NewClient(NewOauthTokenCredentials("bitly-token"))
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
