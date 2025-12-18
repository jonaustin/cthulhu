package main

import "testing"

func TestParseFloorSize(t *testing.T) {
	cases := []struct {
		in      string
		wantW   int
		wantH   int
		wantErr bool
	}{
		{in: "16x16", wantW: 16, wantH: 16},
		{in: " 16X32 ", wantW: 16, wantH: 32},
		{in: "4x16", wantErr: true},
		{in: "16x4", wantErr: true},
		{in: "999x16", wantErr: true},
		{in: "16x999", wantErr: true},
		{in: "abc", wantErr: true},
		{in: "16", wantErr: true},
		{in: "16x", wantErr: true},
	}

	for _, tc := range cases {
		gotW, gotH, err := parseFloorSize(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("parseFloorSize(%q): expected error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("parseFloorSize(%q): unexpected error: %v", tc.in, err)
		}
		if gotW != tc.wantW || gotH != tc.wantH {
			t.Fatalf("parseFloorSize(%q): expected %dx%d, got %dx%d", tc.in, tc.wantW, tc.wantH, gotW, gotH)
		}
	}
}
