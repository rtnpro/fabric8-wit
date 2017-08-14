package resource

// this file contains some paging related utility functions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/jsonapi"
	"github.com/manyminds/api2go"
)

const (
	pageSizeDefault = 20
	pageSizeMax     = 100
)

func computePagingLimits(offsetParam *string, limitParam *int) (offset int, limit int) {
	if offsetParam == nil {
		offset = 0
	} else {
		offsetValue, err := strconv.Atoi(*offsetParam)
		if err != nil {
			offset = 0
		} else {
			offset = offsetValue
		}
	}
	if offset < 0 {
		offset = 0
	}

	if limitParam == nil {
		limit = pageSizeDefault
	} else {
		limit = *limitParam
	}

	if limit <= 0 {
		limit = pageSizeDefault
	} else if limit > pageSizeMax {
		limit = pageSizeMax
	}
	return offset, limit
}

func getPagingLinks(links *jsonapi.Links, path string, resultLen, offset, limit, count int, additionalQuery ...string) (string, string, string, string) {
	var first, prev, next, last string
	format := func(additional []string) string {
		if len(additional) > 0 {
			return "&" + strings.Join(additional, "&")
		}
		return ""
	}

	// prev link
	if offset > 0 && count > 0 {
		var prevStart int
		// we do have a prev link
		if offset <= count {
			prevStart = offset - limit
		} else {
			// the first range that intersects the end of the useful range
			prevStart = offset - (((offset-count)/limit)+1)*limit
		}
		realLimit := limit
		if prevStart < 0 {
			// need to cut the range to start at 0
			realLimit = limit + prevStart
			prevStart = 0
		}
		prev = fmt.Sprintf("%s?page[offset]=%d&page[limit]=%d%s", path, prevStart, realLimit, format(additionalQuery))
	}

	// next link
	nextStart := offset + resultLen
	if nextStart < count {
		// we have a next link
		next = fmt.Sprintf("%s?page[offset]=%d&page[limit]=%d%s", path, nextStart, limit, format(additionalQuery))
	}

	// first link
	var firstEnd int
	if offset > 0 {
		firstEnd = offset % limit // this is where the second page starts
	} else {
		// offset == 0, first == current
		firstEnd = limit
	}
	first = fmt.Sprintf("%s?page[offset]=%d&page[limit]=%d%s", path, 0, firstEnd, format(additionalQuery))

	// last link
	var lastStart int
	if offset < count {
		// advance some pages until touching the end of the range
		lastStart = offset + (((count - offset - 1) / limit) * limit)
	} else {
		// retreat at least one page until covering the range
		lastStart = offset - ((((offset - count) / limit) + 1) * limit)
	}
	realLimit := limit
	if lastStart < 0 {
		// need to cut the range to start at 0
		realLimit = limit + lastStart
		lastStart = 0
	}
	last = fmt.Sprintf("%s?page[offset]=%d&page[limit]=%d%s", path, lastStart, realLimit, format(additionalQuery))
	return first, prev, next, last
}

// ParsePaging parses the paging parameters of a request and returns them in a normalized version (start, limit).
func ParsePaging(r api2go.Request) (int, int, error) {
	var number, size, offset, limit string

	numberQuery, ok := r.QueryParams["page[number]"]
	if ok {
		number = numberQuery[0]
	}
	sizeQuery, ok := r.QueryParams["page[size]"]
	if ok {
		size = sizeQuery[0]
	}
	offsetQuery, ok := r.QueryParams["page[offset]"]
	if ok {
		offset = offsetQuery[0]
	}
	limitQuery, ok := r.QueryParams["page[limit]"]
	if ok {
		limit = limitQuery[0]
	}

	var resultStart int
	var resultLimit int
	if size != "" {
		// PAGE NUMBER AND SIZE MODE
		sizeI, err := strconv.Atoi(size)
		if err != nil {
			return -1, -1, err
		}
		resultStart, resultLimit = computePagingLimits(&number, &sizeI)
	} else {
		// PAGE OFFSET AND LIMIT MODE
		limitI, err := strconv.Atoi(limit)
		if err != nil {
			return -1, -1, err
		}
		offsetI, err := strconv.Atoi(offset)
		if err != nil {
			return -1, -1, err
		}
		resultStart = offsetI
		resultLimit = limitI
	}

	return resultStart, resultLimit, nil
}
