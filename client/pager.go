package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"net/url"
	"strconv"
)

type Pager struct {
	pagination resource.Pagination
}

func (p *Pager) HasNextPage() bool {
	return p.pagination.Next.Href != ""
}

func (p *Pager) NextPage(opts *ListOptions) bool {
	if !p.HasNextPage() {
		return false
	}

	qs, err := newQuerystringReader(p.pagination.Next.Href)
	if err != nil {
		return false
	}
	opts.Page = qs.Int(PageField)
	opts.PerPage = qs.Int(PerPageField)
	return true
}

func (p *Pager) HasPreviousPage() bool {
	return p.pagination.Previous.Href != ""
}

func (p *Pager) PreviousPage(opts *ListOptions) bool {
	if !p.HasPreviousPage() {
		return false
	}

	qs, err := newQuerystringReader(p.pagination.Previous.Href)
	if err != nil {
		return false
	}
	opts.Page = qs.Int(PageField)
	opts.PerPage = qs.Int(PerPageField)
	return true
}

type querystringReader struct {
	qs url.Values
}

func newQuerystringReader(pageURL string) (*querystringReader, error) {
	u, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}
	return &querystringReader{
		qs: u.Query(),
	}, nil
}

func (r querystringReader) String(key string) string {
	return r.qs.Get(key)
}

func (r querystringReader) Int(key string) int {
	i, _ := strconv.Atoi(r.qs.Get(key))
	return i
}
