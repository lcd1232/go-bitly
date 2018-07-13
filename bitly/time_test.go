package bitly

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestJSONDate_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		desc     string
		rawJSON  []byte
		wantTime JSONDate
		wantErr  string
	}{
		{
			desc:     "valid time",
			rawJSON:  []byte(`"2012-12-18T18:14:53+0000"`),
			wantTime: JSONDate(time.Date(2012, 12, 18, 18, 14, 53, 0, time.UTC)),
			wantErr:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var jsonDate = new(JSONDate)
			err := json.Unmarshal(tc.rawJSON, jsonDate)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("want error %v got %v", tc.wantErr, err)
			}
			if !reflect.DeepEqual(tc.wantTime, jsonDate) {
				t.Fatalf("want group %#v got %#v", tc.wantTime, jsonDate)
			}
		})
	}
}
