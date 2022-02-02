package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)


type KGSRepo interface {
	GetKeys(ctx context.Context, count int64) ([]string, error)
}

type KGSUsecase struct {
	repo KGSRepo
	log  *log.Helper
}

func NewKGSUsecase(repo KGSRepo, logger log.Logger) *KGSUsecase {
	return &KGSUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *KGSUsecase) GetKeys(ctx context.Context, count int64) ([]string, error) {
	uc.log.WithContext(ctx).Infof("GetKeys: %v", count)
	return uc.repo.GetKeys(ctx, count)
}
