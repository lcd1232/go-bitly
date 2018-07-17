package bitly

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_UserUnmarshal(t *testing.T) {
	testCases := []struct {
		desc     string
		rawJSON  []byte
		wantUser *User
		wantErr  string
	}{
		{
			desc:    "valid json",
			rawJSON: []byte(`{"created":"2012-12-18T18:14:53+0000","modified":"2018-01-20T17:37:52+0000","login":"test","is_active":true,"is_2fa_enabled":false,"name":"test","emails":[{"email":"test@example.com","is_primary":true,"is_verified":true}],"is_sso_user":false}`),
			wantUser: &User{
				Created:      JSONDate(time.Date(2012, 12, 18, 18, 14, 53, 0, time.UTC)),
				Modified:     JSONDate(time.Date(2018, 1, 20, 17, 37, 52, 0, time.UTC)),
				Login:        "test",
				IsActive:     true,
				Is2FAEnabled: false,
				Name:         "test",
				Emails: []Email{
					{
						Email:      "test@example.com",
						IsPrimary:  true,
						IsVerified: true,
					},
				},
				IsSSOUser: false,
			},
			wantErr: "",
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
			if !time.Time(tc.wantUser.Created).Equal(time.Time(u.Created)) {
				t.Fatalf("want created time %#v got %#v", tc.wantUser.Created, u.Created)
			}
			if !time.Time(tc.wantUser.Modified).Equal(time.Time(u.Modified)) {
				t.Fatalf("want modified time %#v got %#v", tc.wantUser.Modified, u.Modified)
			}
			// We do this because reflect.DeepEqual cannot compare time.Time
			emptyTime := time.Time{}
			tc.wantUser.Created = JSONDate(emptyTime)
			u.Created = JSONDate(emptyTime)
			tc.wantUser.Modified = JSONDate(emptyTime)
			u.Modified = JSONDate(emptyTime)
			if !reflect.DeepEqual(tc.wantUser, u) {
				t.Fatalf("want group %#v got %#v", tc.wantUser, u)
			}
		})
	}
}

func TestUserClient_Get(t *testing.T) {
	testCases := []struct {
		desc         string
		responseCode int
		responseBody string
		wantUser     *User
		wantErr      string
	}{
		{
			desc:         "ok response",
			responseCode: http.StatusOK,
			responseBody: `{"created":"2012-12-18T18:14:53+0000","modified":"2018-01-20T17:37:52+0000","login":"test","is_active":true,"is_2fa_enabled":false,"name":"test","emails":[{"email":"test@example.com","is_primary":true,"is_verified":true}],"is_sso_user":false}`,
			wantUser: &User{
				Created:      JSONDate(time.Date(2012, 12, 18, 18, 14, 53, 0, time.UTC)),
				Modified:     JSONDate(time.Date(2018, 1, 20, 17, 37, 52, 0, time.UTC)),
				Login:        "test",
				IsActive:     true,
				Is2FAEnabled: false,
				Name:         "test",
				Emails: []Email{
					{
						Email:      "test@example.com",
						IsPrimary:  true,
						IsVerified: true,
					},
				},
				IsSSOUser: false,
			},
			wantErr: "",
		},
		{
			desc:         "invalid token",
			responseCode: http.StatusForbidden,
			responseBody: `{"message":"FORBIDDEN"}`,
			wantErr:      "403 FORBIDDEN",
			wantUser:     nil,
		},
		{
			desc:         "server error",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"message":"some error"}`,
			wantErr:      "500 some error",
			wantUser:     nil,
		},
		{
			desc:         "temporary unavailable",
			responseCode: http.StatusServiceUnavailable,
			responseBody: `{"message":"unavailable"}`,
			wantErr:      "503 unavailable",
			wantUser:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Fatalf("invalid request method: %q", r.Method)
				}
				if r.URL.Path != "/v4/user" {
					t.Fatalf("invalid request path: %q", r.URL.Path)
				}
				w.WriteHeader(tc.responseCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer s.Close()

			c := NewClient(http.DefaultClient)
			c.BaseURL = s.URL

			u, err := c.User.Get(context.Background())

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("want error %v got %v", tc.wantErr, err)
			}

			// Skip next checks if we wantUser is nil and user nil too
			if tc.wantUser == nil && u == nil {
				return
			}
			if tc.wantUser == nil && u != nil || tc.wantUser != nil && u == nil {
				t.Fatalf("want user %#v got %#v", tc.wantUser, u)
			}
			if !time.Time(tc.wantUser.Created).Equal(time.Time(u.Created)) {
				t.Fatalf("want created time %#v got %#v", tc.wantUser.Created, u.Created)
			}
			if !time.Time(tc.wantUser.Modified).Equal(time.Time(u.Modified)) {
				t.Fatalf("want modified time %#v got %#v", tc.wantUser.Modified, u.Modified)
			}
			// We do this because reflect.DeepEqual cannot compare time.Time
			emptyTime := time.Time{}
			tc.wantUser.Created = JSONDate(emptyTime)
			u.Created = JSONDate(emptyTime)
			tc.wantUser.Modified = JSONDate(emptyTime)
			u.Modified = JSONDate(emptyTime)
			if !reflect.DeepEqual(tc.wantUser, u) {
				t.Fatalf("want user %#v got %#v", tc.wantUser, u)
			}
		})
	}
}
