package bitly

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"testing"
)

var (
	bitlyLiveTest bool
	bitlyToken    string
	bitlyBaseURL  string
	bitlyClient   *Client
)

func init() {
	bitlyToken = os.Getenv("BITLY_TOKEN")
	bitlyBaseURL = os.Getenv("BITLY_BASE_URL")
	bitlyDebug := strings.ToUpper(os.Getenv("BITLY_DEBUG")) == "TRUE"
	if bitlyBaseURL == "" {
		bitlyBaseURL = defaultBaseURL
	}

	if len(bitlyToken) > 0 {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: bitlyToken})
		tc := oauth2.NewClient(context.Background(), ts)
		bitlyClient = NewClient(tc)
		bitlyClient.BaseURL = bitlyBaseURL
		bitlyClient.UserAgent = fmt.Sprintf("%v +livetest", bitlyClient.UserAgent)
		bitlyClient.Debug = bitlyDebug
	}
}

func TestLive_ListGroups(t *testing.T) {
	if !bitlyLiveTest {
		t.Skip("skipping live test")
	}

	groupsResp, err := bitlyClient.Groups.ListGroups("")
	if err != nil {
		t.Fatalf("Live Groups.ListGroups() returned error: %v", err)
	}
	for _, group := range groupsResp.Groups {
		t.Logf("GUID: %v\n", group.GUID)
		t.Logf("Organization GUID: %v\n", group.OrganizationGUID)
		groupResp, err := bitlyClient.Groups.GetGroup(group.GUID)
		if err != nil {
			t.Fatalf("Live Groups.GetGroup(%v) returned error: %v", group.GUID, err)
		}
		t.Logf("Group %v is active: %v\n", groupResp.GUID, groupResp.IsActive)
	}
}
