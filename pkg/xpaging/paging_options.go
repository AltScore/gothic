package xpaging

import (
	"github.com/AltScore/gothic/pkg/xvalidator"
)

type PagingOptions struct {
	Offset int64 `query:"offset" json:"offset" validate:"min=0"`
	Limit  int64 `query:"limit"  json:"limit" validate:"min=0,max=100"`
}

func (po PagingOptions) Validate() error {
	return xvalidator.Struct(po)
}

func (po PagingOptions) Normalized() PagingOptions {
	offset := po.Offset
	if offset < 0 {
		offset = 0
	}

	limit := po.Limit
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return PagingOptions{
		Offset: offset,
		Limit:  limit,
	}
}
