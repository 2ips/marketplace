package port

import "fmt"

type PageRequest struct {
	page     uint64
	pageSize uint64
}

type ErrInvalidPage struct {
	Value uint64
}

type ErrInvalidPageSize struct {
	Value uint64
}

type ErrPageOutOfRange struct {
	Value      uint64
	TotalPages uint64
}

func (e *ErrPageOutOfRange) Error() string {
	return fmt.Sprintf("page must be at most %d, got %d", e.TotalPages, e.Value)
}

func (e *ErrInvalidPage) Error() string {
	return fmt.Sprintf("page must be non-zero, got %d", e.Value)
}

func (e *ErrInvalidPageSize) Error() string {
	return fmt.Sprintf("page size must be-non zero, got %d", e.Value)
}

func NewPageRequest(page uint64, pageSize uint64) (PageRequest, error) {
	if page == 0 {
		return PageRequest{}, &ErrInvalidPage{Value: page}
	}

	if pageSize == 0 {
		return PageRequest{}, &ErrInvalidPageSize{Value: pageSize}
	}

	return PageRequest{page: page, pageSize: pageSize}, nil
}

func (p PageRequest) Page() uint64 {
	return p.page
}

func (p PageRequest) PageSize() uint64 {
	return p.pageSize
}

func (p PageRequest) Bounds(total uint64) (uint64, uint64, error) {
	totalPages := (total + p.pageSize - 1) / p.pageSize
	if p.page-1 >= totalPages {
		return 0, 0, &ErrPageOutOfRange{Value: p.page, TotalPages: totalPages}
	}

	offset := (p.page - 1) * p.pageSize
	end := min(offset+p.pageSize, total)

	return offset, end, nil
}
