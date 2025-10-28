// Copyright 2018 The dipnet-core Authors
// This file is part of the dipnet-core library.
//
// The dipnet-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The dipnet-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the dipnet-core library. If not, see <http://www.gnu.org/licenses/>.

package accounts

import (
	"testing"
)

func TestURLParsing(t *testing.T) {
	t.Parallel()
	url, err := parseURL("https://dipnet.org")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "dipnet.org" {
		t.Errorf("expected: %v, got: %v", "dipnet.org", url.Path)
	}

	for _, u := range []string{"dipnet.org", ""} {
		if _, err = parseURL(u); err == nil {
			t.Errorf("input %v, expected err, got: nil", u)
		}
	}
}

func TestURLString(t *testing.T) {
	t.Parallel()
	url := URL{Scheme: "https", Path: "dipnet.org"}
	if url.String() != "https://dipnet.org" {
		t.Errorf("expected: %v, got: %v", "https://dipnet.org", url.String())
	}

	url = URL{Scheme: "", Path: "dipnet.org"}
	if url.String() != "dipnet.org" {
		t.Errorf("expected: %v, got: %v", "dipnet.org", url.String())
	}
}

func TestURLMarshalJSON(t *testing.T) {
	t.Parallel()
	url := URL{Scheme: "https", Path: "dipnet.org"}
	json, err := url.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(json) != "\"https://dipnet.org\"" {
		t.Errorf("expected: %v, got: %v", "\"https://dipnet.org\"", string(json))
	}
}

func TestURLUnmarshalJSON(t *testing.T) {
	t.Parallel()
	url := &URL{}
	err := url.UnmarshalJSON([]byte("\"https://dipnet.org\""))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "dipnet.org" {
		t.Errorf("expected: %v, got: %v", "https", url.Path)
	}
}

func TestURLComparison(t *testing.T) {
	t.Parallel()
	tests := []struct {
		urlA   URL
		urlB   URL
		expect int
	}{
		{URL{"https", "dipnet.org"}, URL{"https", "dipnet.org"}, 0},
		{URL{"http", "dipnet.org"}, URL{"https", "dipnet.org"}, -1},
		{URL{"https", "dipnet.org/a"}, URL{"https", "dipnet.org"}, 1},
		{URL{"https", "abc.org"}, URL{"https", "dipnet.org"}, -1},
	}

	for i, tt := range tests {
		result := tt.urlA.Cmp(tt.urlB)
		if result != tt.expect {
			t.Errorf("test %d: cmp mismatch: expected: %d, got: %d", i, tt.expect, result)
		}
	}
}
