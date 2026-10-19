package config

import (
	"reflect"
	"testing"
)

func TestParseStreams(t *testing.T) {
	cases := []struct {
		name       string
		streamsEnv string
		streamFall string
		want       []string
	}{
		{"multi from REDIS_STREAMS", "homerun,releases", "ignored", []string{"homerun", "releases"}},
		{"whitespace trimmed", " homerun , releases ", "", []string{"homerun", "releases"}},
		{"empty entries dropped", "homerun,,releases,", "", []string{"homerun", "releases"}},
		{"legacy REDIS_STREAM single", "", "homerun", []string{"homerun"}},
		{"empty streams env falls back to legacy", "", "messages", []string{"messages"}},
		{"both unset returns hardcoded default", "", "", []string{"messages"}},
		{"streams with only whitespace falls back to legacy", " , , ", "legacy", []string{"legacy"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseStreams(tc.streamsEnv, tc.streamFall)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseStreams(%q, %q) = %v, want %v", tc.streamsEnv, tc.streamFall, got, tc.want)
			}
		})
	}
}

func TestParseBasePath(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"", "", false},
		{"/", "", false},
		{"/core", "/core", false},
		{"core", "/core", false},
		{"/core/", "/core", false},
		{" /homerun/core ", "/homerun/core", false},
		{"/a//b", "", true},
		{"/../core", "", true},
		{"/core\"><script>", "", true},
		{"/core?x=1", "", true},
	}
	for _, tc := range cases {
		got, err := ParseBasePath(tc.in)
		if (err != nil) != tc.wantErr {
			t.Errorf("ParseBasePath(%q) error = %v, wantErr %v", tc.in, err, tc.wantErr)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseBasePath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
