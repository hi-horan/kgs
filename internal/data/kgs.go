package data

import (
	"context"

	"kgs/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type kgsRepo struct {
	data *Data
	log  *log.Helper
}

func NewKGSRepo(data *Data, logger log.Logger) biz.KGSRepo {
	return &kgsRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}


func (r *kgsRepo) GetKeys(ctx context.Context, count int64) ([]string, error) {
    return r.data.Getkeys(count)
}
