package util

const (
	defaultPage     = 1
	defaultPageSize = 10
)

func ToOffsetLimit(page, size int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}

	if size <= 0 {
		size = defaultPageSize
	}

	limit := size
	offset := (page - 1) * size

	return offset, limit
}
