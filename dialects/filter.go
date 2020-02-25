// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dialects

import (
	"fmt"
	"strings"

	"xorm.io/xorm/schemas"
)

// Filter is an interface to filter SQL
type Filter interface {
	Do(sql string, dialect Dialect, table *schemas.Table) string
}

// QuoteFilter filter SQL replace ` to database's own quote character
type QuoteFilter struct {
}

func (s *QuoteFilter) Do(sql string, dialect Dialect, table *schemas.Table) string {
	quoter := dialect.Quoter()
	if quoter.IsEmpty() {
		return sql
	}

	prefix, suffix := quoter[0][0], quoter[1][0]
	raw := []byte(sql)
	for i, cnt := 0, 0; i < len(raw); i = i + 1 {
		if raw[i] == '`' {
			if cnt%2 == 0 {
				raw[i] = prefix
			} else {
				raw[i] = suffix
			}
			cnt++
		}
	}
	return string(raw)

}

// SeqFilter filter SQL replace ?, ? ... to $1, $2 ...
type SeqFilter struct {
	Prefix string
	Start  int
}

func convertQuestionMark(sql, prefix string, start int) string {
	var buf strings.Builder
	var beginSingleQuote bool
	var index = start
	for _, c := range sql {
		if !beginSingleQuote && c == '?' {
			buf.WriteString(fmt.Sprintf("%s%v", prefix, index))
			index++
		} else {
			if c == '\'' {
				beginSingleQuote = !beginSingleQuote
			}
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

func (s *SeqFilter) Do(sql string, dialect Dialect, table *schemas.Table) string {
	return convertQuestionMark(sql, s.Prefix, s.Start)
}
