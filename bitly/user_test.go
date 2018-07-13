package bitly

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func Test_UserUnmarshal(t *testing.T) {
	testCases := []struct {
		desc     string
		rawJSON  []byte
		wantUser *User
		wantErr  string
	}{
		{
			desc:     "valid json",
			rawJSON:  []byte(`{"created":"2012-12-18T18:14:53+0000","modified":"2018-01-20T17:37:52+0000","login":"test","is_active":true,"is_2fa_enabled":false,"name":"test","emails":[{"email":"test@example.com","is_primary":true,"is_verified":true}],"is_sso_user":false}`),
			wantUser: &User{},
			wantErr:  "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			u := &User{}
			err := json.Unmarshal(tc.rawJSON, u)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("want error %v got %v", tc.wantErr, err)
			}
			if !reflect.DeepEqual(tc.wantUser, u) {
				t.Fatalf("want group %#v got %#v", tc.wantUser, u)
			}
		})
	}
}
