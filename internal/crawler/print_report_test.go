package crawler

import (
	"reflect"
	"testing"
)

func TestSortPages(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]*PageData
		expected []Page
	}{
		{
			name: "order count descending",
			input: map[string]*PageData{
				"url1": {LinkCount: 5},
				"url2": {LinkCount: 1},
				"url3": {LinkCount: 3},
				"url4": {LinkCount: 10},
				"url5": {LinkCount: 7},
			},
			expected: []Page{
				{URL: "url4", Count: 10},
				{URL: "url5", Count: 7},
				{URL: "url1", Count: 5},
				{URL: "url3", Count: 3},
				{URL: "url2", Count: 1},
			},
		},
		{
			name: "alphabetize",
			input: map[string]*PageData{
				"d": {LinkCount: 1},
				"a": {LinkCount: 1},
				"e": {LinkCount: 1},
				"b": {LinkCount: 1},
				"c": {LinkCount: 1},
			},
			expected: []Page{
				{URL: "a", Count: 1},
				{URL: "b", Count: 1},
				{URL: "c", Count: 1},
				{URL: "d", Count: 1},
				{URL: "e", Count: 1},
			},
		},
		{
			name: "order count then alphabetize",
			input: map[string]*PageData{
				"d": {LinkCount: 2},
				"a": {LinkCount: 1},
				"e": {LinkCount: 3},
				"b": {LinkCount: 1},
				"c": {LinkCount: 2},
			},
			expected: []Page{
				{URL: "e", Count: 3},
				{URL: "c", Count: 2},
				{URL: "d", Count: 2},
				{URL: "a", Count: 1},
				{URL: "b", Count: 1},
			},
		},
		{
			name:     "empty map",
			input:    map[string]*PageData{},
			expected: []Page{},
		},
		{
			name:     "nil map",
			input:    nil,
			expected: []Page{},
		},
		{
			name: "one key",
			input: map[string]*PageData{
				"url1": {LinkCount: 1},
			},
			expected: []Page{
				{URL: "url1", Count: 1},
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := sortPages(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
