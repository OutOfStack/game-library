package model

// SortOrder - type for query sort order
type SortOrder string

// SortOrder values
const (
	AscendingSortOrder  SortOrder = "ASC"
	DescendingSortOrder SortOrder = "DESC"
)

// OrderBy type of ordering
type OrderBy struct {
	Field string
	Order SortOrder
}
