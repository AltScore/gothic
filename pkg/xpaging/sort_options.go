package xpaging

import (
	"github.com/AltScore/gothic/pkg/xerrors"
	"regexp"
	"strings"
)

const (
	DirectionAsc  SortDirection = 1
	DirectionDesc SortDirection = -1
)

type SortDirection int

func directionFromString(order string) SortDirection {
	if order == "desc" {
		return DirectionDesc
	}

	return DirectionAsc
}

type SortEntry struct {
	FieldName string
	Direction SortDirection
}

func NewSortEntryFromString(fieldName string, order string) SortEntry {
	return SortEntry{
		FieldName: fieldName,
		Direction: directionFromString(order),
	}
}

func NewSortEntry(fieldName string, direction SortDirection) SortEntry {
	return SortEntry{
		FieldName: fieldName,
		Direction: direction,
	}
}

type SortOptions []SortEntry

func (s SortOptions) IsEmpty() bool {
	return len(s) == 0
}

func NewSortOptionsFromString(sorts string, defaultOption ...string) (SortOptions, error) {

	sortEntries := strings.Split(sorts, ",")

	// if no sort options are provided, use the default
	if len(sorts) == 0 {
		return makeSortOptions(defaultOption)
	}

	return makeSortOptions(sortEntries)
}

func makeSortOptions(sortEntries []string) (SortOptions, error) {
	sortOptions := SortOptions{}
	for _, sortOption := range sortEntries {
		trimmedOption := strings.TrimSpace(sortOption)
		sortEntry, err := extractOption(trimmedOption)
		if err != nil {
			return nil, err
		}
		sortOptions = append(sortOptions, sortEntry)
	}
	return sortOptions, nil
}

func extractOption(option string) (SortEntry, error) {

	var sortResult SortEntry

	var sortOptionRegex = regexp.MustCompile(`^([-+!])?(\w+)(?::(asc|desc))?$`)
	parts := sortOptionRegex.FindStringSubmatch(option)
	if len(parts) == 0 {
		return sortResult, xerrors.NewInvalidArgumentError("QueryParamFormat", "SortBy Param allows these formats: field:asc, field:desc, -field, field")
	}
	if parts[1] == "-" || parts[1] == "!" {
		sortResult = NewSortEntry(parts[2], DirectionDesc)
	} else {
		sortResult = NewSortEntryFromString(parts[2], parts[3])
	}

	return sortResult, nil
}
