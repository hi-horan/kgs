package service

import (
	"context"

	v1 "kgs/api/v1"
	"kgs/internal/biz"
)

// GreeterService is a greeter service.
type KGSService struct {
	v1.UnimplementedKGSServer

	uc *biz.KGSUsecase
}

func NewKGSService(uc *biz.KGSUsecase) *KGSService {
	return &KGSService{uc: uc}
}

func (s *KGSService) GetKeys(ctx context.Context, in *v1.GetKeysRequest) (*v1.GetKeysReply, error) {
	keys, err := s.uc.GetKeys(ctx, in.Count)
	if err != nil {
		return nil, err
	}
	return &v1.GetKeysReply{Keys: keys}, nil
}
