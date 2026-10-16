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
