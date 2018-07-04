package bitly

import (
	"fmt"
	"os"
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
	_, bitlyDebug := os.LookupEnv("BITLY_DEBUG")
	if bitlyBaseURL == "" {
		bitlyBaseURL = defaultBaseURL
	}

	if len(bitlyToken) > 0 {
		bitlyLiveTest = true
		bitlyClient = NewClient(NewOauthTokenCredentials(bitlyToken))
		bitlyClient.BaseURL = bitlyBaseURL
		bitlyClient.UserAgent = fmt.Sprintf("%v +livetest", bitlyClient.UserAgent)
		bitlyClient.Debug = bitlyDebug
	}
}

func TestLive_ListGroups(t *testing.T) {
	if !bitlyLiveTest {
		t.Skip("skipping live test")
	}

	groupsResponse, err := bitlyClient.Groups.ListGroups("")
	if err != nil {
		t.Fatalf("Live Groups.ListGroups() returned error: %v", err)
	}
	for _, group := range groupsResponse.Groups {
		t.Logf("GUID: %v\n", group.GUID)
		t.Logf("Organization GUID: %v\n", group.OrganizationGUID)
	}
}
