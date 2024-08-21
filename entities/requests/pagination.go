package requests

import "strings"

type MetaPaginationRequest struct {
	Page      int    `json:"page" query:"page"`
	Limit     int    `json:"limit" query:"limit"`
	Offset    int    `json:"offset" query:"offset"`
	Order     string `json:"order" query:"order"`
	SortBy    string `json:"sort_by" query:"sort_by"`
	Count     int64  `json:"count,omitempty"`
	TotalPage int    `json:"total_page"`
	Search    string `json:"search" query:"search"`
}

func (p *MetaPaginationRequest) ParsePagination() MetaPaginationRequest {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 10
	}

	offset := (p.Page - 1) * p.Limit
	if offset < 0 {
		offset = 0
	}

	if p.SortBy == "" {
		p.SortBy = "created_at"
	}

	if p.Order == "" {
		p.Order = "DESC"
	} else {
		if strings.ToLower(p.Order) != "asc" || strings.ToLower(p.Order) != "desc" {
			p.Order = "DESC"
		}
	}

	p.Offset = offset
	return *p
}
