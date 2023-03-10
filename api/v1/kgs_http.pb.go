// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.3
// - protoc             v3.19.2
// source: v1/kgs.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationKGSGetKeys = "/v1.KGS/GetKeys"

type KGSHTTPServer interface {
	GetKeys(context.Context, *GetKeysRequest) (*GetKeysReply, error)
}

func RegisterKGSHTTPServer(s *http.Server, srv KGSHTTPServer) {
	r := s.Route("/")
	r.POST("/api/v1/keys", _KGS_GetKeys0_HTTP_Handler(srv))
}

func _KGS_GetKeys0_HTTP_Handler(srv KGSHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetKeysRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationKGSGetKeys)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetKeys(ctx, req.(*GetKeysRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetKeysReply)
		return ctx.Result(200, reply)
	}
}

type KGSHTTPClient interface {
	GetKeys(ctx context.Context, req *GetKeysRequest, opts ...http.CallOption) (rsp *GetKeysReply, err error)
}

type KGSHTTPClientImpl struct {
	cc *http.Client
}

func NewKGSHTTPClient(client *http.Client) KGSHTTPClient {
	return &KGSHTTPClientImpl{client}
}

func (c *KGSHTTPClientImpl) GetKeys(ctx context.Context, in *GetKeysRequest, opts ...http.CallOption) (*GetKeysReply, error) {
	var out GetKeysReply
	pattern := "/api/v1/keys"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationKGSGetKeys))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
