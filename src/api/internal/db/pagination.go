package db

import (
	"fmt"
	"strings"
)

// PaginationType represents different pagination methods
type PaginationType int

const (
	CursorPagination PaginationType = iota
	OffsetPagination
	PagePagination
)

// PaginationParams holds common pagination parameters
type PaginationParams struct {
	Type       PaginationType
	Limit      int
	Cursor     string   // For cursor-based pagination
	Offset     int      // For offset-based pagination
	Page       int      // For page-based pagination
	SortFields []string // Fields to sort by
	SortOrder  string   // ASC or DESC
	KeyID      string   // Key ID to use for pagination
}

// BuildPaginationQuery builds the pagination part of SQL query based on pagination type
func BuildPaginationQuery(params PaginationParams) (string, []interface{}, error) {
	var query string
	var args []interface{}

	// Add ORDER BY clause if sort fields are specified
	if len(params.SortFields) > 0 {
		query += fmt.Sprintf(" ORDER BY %s %s",
			strings.Join(params.SortFields, ", "),
			strings.ToUpper(params.SortOrder))
	}

	switch params.Type {
	case CursorPagination:
		if params.Cursor != "" {
			query += fmt.Sprintf(" WHERE %s > ? ", params.KeyID)
			args = append(args, params.Cursor)
		}
		if params.Limit > 0 {
			query += " LIMIT ? "
			args = append(args, params.Limit)
		}

	case OffsetPagination:
		if params.Limit > 0 {
			query += " LIMIT ? "
			args = append(args, params.Limit)
		}
		if params.Offset > 0 {
			query += " OFFSET ? "
			args = append(args, params.Offset)
		}

	case PagePagination:
		if params.Limit > 0 {
			query += " LIMIT ? "
			args = append(args, params.Limit)

			offset := (params.Page - 1) * params.Limit
			if offset > 0 {
				query += " OFFSET ? "
				args = append(args, offset)
			}
		}
	}

	return query, args, nil
}

// PaginationResponse holds the paginated results
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	NextCursor string      `json:"next_cursor,omitempty"`
	PrevCursor string      `json:"prev_cursor,omitempty"`
	Total      int64       `json:"total"`
	Page       int         `json:"page,omitempty"`
	TotalPages int         `json:"total_pages,omitempty"`
	HasMore    bool        `json:"has_more"`
}

// NewPaginationParams creates pagination parameters with defaults
func NewPaginationParams(pType PaginationType) PaginationParams {
	return PaginationParams{
		Type:       pType,
		Limit:      10,             // Default limit
		Page:       1,              // Default page
		SortOrder:  "ASC",          // Default sort order
		SortFields: []string{"id"}, // Default sort field
	}
}
