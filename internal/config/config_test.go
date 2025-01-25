package config_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hizla/hizla/internal/config"
)

type testConfig string

func (t *testConfig) FromEnviron() error {
	return nil
}

func TestLoad(t *testing.T) {
	testCases := []struct {
		name    string
		whence  int
		data    string
		want    any
		wantErr bool
	}{
		{"environ", config.FromEnviron, "", new(testConfig), false},
		{"json", config.FromJSON, `""`, new(testConfig), false},
		{"toml", config.FromTOML, "", new(testConfig), true},
		{"invalid", -1, "", nil, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := new(testConfig)
			if err := config.Load(got, tc.whence, strings.NewReader(tc.data)); (err != nil) != tc.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Load() c = %#v; want %#v", got, tc.want)
				return
			}
		})
	}
}
