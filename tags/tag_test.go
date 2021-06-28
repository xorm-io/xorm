// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTag(t *testing.T) {
	var cases = []struct {
		tag  string
		tags []tag
	}{
		{"not null default '2000-01-01 00:00:00' TIMESTAMP", []tag{
			{
				name: "not",
			},
			{
				name: "null",
			},
			{
				name: "default",
			},
			{
				name: "'2000-01-01 00:00:00'",
			},
			{
				name: "TIMESTAMP",
			},
		},
		},
		{"TEXT", []tag{
			{
				name: "TEXT",
			},
		},
		},
		{"default('2000-01-01 00:00:00')", []tag{
			{
				name: "default",
				params: []string{
					"'2000-01-01 00:00:00'",
				},
			},
		},
		},
		{"json  binary", []tag{
			{
				name: "json",
			},
			{
				name: "binary",
			},
		},
		},
		{"numeric(10, 2)", []tag{
			{
				name:   "numeric",
				params: []string{"10", "2"},
			},
		},
		},
		{"numeric(10, 2) notnull", []tag{
			{
				name:   "numeric",
				params: []string{"10", "2"},
			},
			{
				name: "notnull",
			},
		},
		},
	}

	for _, kase := range cases {
		t.Run(kase.tag, func(t *testing.T) {
			tags, err := splitTag(kase.tag)
			assert.NoError(t, err)
			assert.EqualValues(t, len(tags), len(kase.tags))
			for i := 0; i < len(tags); i++ {
				assert.Equal(t, tags[i], kase.tags[i])
			}
		})
	}
}
