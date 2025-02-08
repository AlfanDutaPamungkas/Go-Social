package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int       `json:"limit" validate:"gte=1,lte=20"`
	Offset int       `json:"offset" validate:"gte=0"`
	Sort   string    `json:"sort" validate:"oneof=asc desc"`
	Tags   []string  `json:"tags" validate:"max=5"`
	Search string    `json:"search" validate:"max=100"`
	Since  time.Time `json:"since"`
	Until  time.Time `json:"until"`
}

func (p *PaginatedFeedQuery) Parse(r *http.Request) (*PaginatedFeedQuery, error) {
	q := r.URL.Query()

	limit := q.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}

		p.Limit = l
	}

	offset := q.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return nil, err
		}

		p.Offset = o
	}

	sort := q.Get("sort")
	if sort != "" {
		p.Sort = sort
	}

	tags := q.Get("tags")
	if tags != "" {
		p.Tags = strings.Split(tags, ",")
	}

	search := q.Get("search")
	if search != "" {
		p.Search = search
	}

	since := q.Get("since")
	if since != "" {
		since, err := parseTime(since)
		if err != nil {
			return nil, err
		}

		p.Since = since
	}

	until := q.Get("until")
	if until != "" {
		until, err := parseTime(until)
		if err != nil {
			return nil, err
		}

		p.Until = until
	}

	return p, nil
}

func parseTime(s string) (time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
