package utils

import (
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	Page      int
	Limit     int
	Search    string
	SortBy    string
	SortOrder string
	Filters   map[string]string
}

func ParseQuery(v url.Values) Query {
	q := Query{Page: 1, Limit: 10, Search: strings.TrimSpace(v.Get("search")), SortBy: "created_at", SortOrder: "desc", Filters: make(map[string]string)}
	if n, err := strconv.Atoi(v.Get("page")); err == nil {
		q.Page = n
	}
	if n, err := strconv.Atoi(v.Get("limit")); err == nil {
		q.Limit = n
	}
	if v.Get("sortBy") != "" {
		q.SortBy = v.Get("sortBy")
	}
	if s := strings.ToLower(v.Get("sortOrder")); s == "asc" || s == "desc" {
		q.SortOrder = s
	}

	for key, value := range v {
		switch key {
		case "search", "page", "limit", "sortBy", "sortOrder":
			continue
		default:
			if len(value) > 0 {
				val := strings.TrimSpace(value[0])
				if val != "" {
					q.Filters[key] = val
				}
			}
		}
	}

	return q
}

func (q Query) Offset() int {
	return (q.Page - 1) * q.Limit
}
