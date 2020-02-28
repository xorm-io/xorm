// Copyright 2020 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schemas

import (
	"strings"
)

// Quoter represents two quote characters
type Quoter [2]string

// CommonQuoter represetns a common quoter
var CommonQuoter = Quoter{"`", "`"}

func (q Quoter) IsEmpty() bool {
	return q[0] == "" && q[1] == ""
}

func (q Quoter) Quote(s string) string {
	var buf strings.Builder
	q.QuoteTo(&buf, s)
	return buf.String()
}

func (q Quoter) Replace(sql string, newQuoter Quoter) string {
	if q.IsEmpty() {
		return sql
	}

	if newQuoter.IsEmpty() {
		var buf strings.Builder
		for i := 0; i < len(sql); i = i + 1 {
			if sql[i] != q[0][0] && sql[i] != q[1][0] {
				_ = buf.WriteByte(sql[i])
			}
		}
		return buf.String()
	}

	prefix, suffix := newQuoter[0][0], newQuoter[1][0]
	var buf strings.Builder
	for i, cnt := 0, 0; i < len(sql); i = i + 1 {
		if cnt == 0 && sql[i] == q[0][0] {
			_ = buf.WriteByte(prefix)
			cnt = 1
		} else if cnt == 1 && sql[i] == q[1][0] {
			_ = buf.WriteByte(suffix)
			cnt = 0
		} else {
			_ = buf.WriteByte(sql[i])
		}
	}
	return buf.String()
}

func (q Quoter) ReverseQuote(s string) string {
	reverseQuoter := Quoter{q[1], q[0]}
	return reverseQuoter.Quote(s)
}

// Trim removes quotes from s
func (q Quoter) Trim(s string) string {
	if len(s) < 2 {
		return s
	}

	var buf strings.Builder
	for i := 0; i < len(s); i++ {
		switch {
		case i == 0 && s[i:i+1] == q[0]:
		case i == len(s)-1 && s[i:i+1] == q[1]:
		case s[i:i+1] == q[1] && s[i+1] == '.':
		case s[i:i+1] == q[0] && s[i-1] == '.':
		default:
			buf.WriteByte(s[i])
		}
	}
	return buf.String()
}

func (q Quoter) Join(a []string, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	var b strings.Builder
	b.Grow(n)
	for i, s := range a {
		if i > 0 {
			b.WriteString(sep)
		}
		if q[0] != "" && s != "*" {
			b.WriteString(q[0])
		}
		b.WriteString(strings.TrimSpace(s))
		if q[1] != "" && s != "*" {
			b.WriteString(q[1])
		}
	}
	return b.String()
}

func (q Quoter) JoinWrite(b *strings.Builder, a []string, sep string) error {
	if len(a) == 0 {
		return nil
	}

	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b.Grow(n)
	for i, s := range a {
		if i > 0 {
			if _, err := b.WriteString(sep); err != nil {
				return err
			}
		}
		if q[0] != "" && s != "*" && s[0] != '`' {
			if _, err := b.WriteString(q[0]); err != nil {
				return err
			}
		}
		if _, err := b.WriteString(strings.TrimSpace(s)); err != nil {
			return err
		}
		if q[1] != "" && s != "*" && s[0] != '`' {
			if _, err := b.WriteString(q[1]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (q Quoter) Strings(s []string) []string {
	var res = make([]string, 0, len(s))
	for _, a := range s {
		res = append(res, q.Quote(a))
	}
	return res
}

func (q Quoter) QuoteTo(buf *strings.Builder, value string) {
	if q.IsEmpty() {
		buf.WriteString(value)
		return
	}

	prefix, suffix := q[0][0], q[1][0]
	lastCh := 0 // 0 prefix, 1 char, 2 suffix
	i := 0
	for i < len(value) {
		// start of a token; might be already quoted
		if value[i] == '.' {
			_ = buf.WriteByte('.')
			lastCh = 1
			i++
		} else if value[i] == prefix || value[i] == '`' {
			// Has quotes; skip/normalize `name` to prefix+name+sufix
			var ch byte
			if value[i] == prefix {
				ch = suffix
			} else {
				ch = '`'
			}
			_ = buf.WriteByte(prefix)
			i++
			lastCh = 0
			for ; i < len(value) && value[i] != ch && value[i] != ' '; i++ {
				_ = buf.WriteByte(value[i])
				lastCh = 1
			}
			_ = buf.WriteByte(suffix)
			lastCh = 2
			i++
		} else if value[i] == ' ' {
			if lastCh != 2 {
				_ = buf.WriteByte(suffix)
				lastCh = 2
			}

			// a AS b or a b
			for ; i < len(value); i++ {
				if value[i] != ' ' && value[i-1] == ' ' && (len(value) > i+1 && !strings.EqualFold(value[i:i+2], "AS")) {
					break
				}

				_ = buf.WriteByte(value[i])
				lastCh = 1
			}
		} else {
			// Requires quotes
			_ = buf.WriteByte(prefix)
			for ; i < len(value) && value[i] != '.' && value[i] != ' '; i++ {
				_ = buf.WriteByte(value[i])
				lastCh = 1
			}
			_ = buf.WriteByte(suffix)
			lastCh = 2
		}
	}
}
