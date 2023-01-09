package handler

import (
	"context"

	"github.com/yqhcode/paasPod/domain/service"
	"github.com/yqhcode/paasPod/proto"

	common "github.com/yqhcode/paas-common"
)

type PodHandler struct {
	PodService service.PodService
}

func (p *PodHandler) AddPod(ctx context.Context, info *proto.PodInfo) (*proto.Response, error) {
	common.Info("添加pod")
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) DeletePod(ctx context.Context, id *proto.PodId) (*proto.Response, error) {
	//TODO implement me
	panic("implement me")

}

func (p *PodHandler) FindPodByID(ctx context.Context, id *proto.PodId) (*proto.PodInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) UpdatePod(ctx context.Context, info *proto.PodInfo) (*proto.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) FindAllPod(ctx context.Context, all *proto.FindAll) (*proto.AllPod, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) mustEmbedUnimplementedPodServer() {
}
