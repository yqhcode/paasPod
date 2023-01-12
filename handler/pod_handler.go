package handler

import (
	"context"
	"github.com/yqhcode/paasPod/domain/model"
	"github.com/yqhcode/paasPod/domain/service"
	"github.com/yqhcode/paasPod/proto"
	"strconv"

	common "github.com/yqhcode/paas-common"
)

type PodHandler struct {
	PodService service.PodService
}

func (p *PodHandler) FindPodByID(ctx context.Context, id *proto.PodId) (rsp *proto.PodInfo, err error) {
	podModel, err := p.PodService.FindPodByID(id.Id)
	if err != nil {
		common.Error(err)
		return nil, err
	}
	err = common.SwapTo(podModel, rsp)
	if err != nil {
		common.Error(err)
		return nil, err
	}
	return rsp, nil
}

func (p *PodHandler) AddPod(ctx context.Context, info *proto.PodInfo) (rsp *proto.Response, err error) {
	common.Info("添加pod")
	podModel := &model.Pod{}
	if err := common.SwapTo(info, podModel); err != nil {
		common.Error(err)
		return nil, err
	}

	if err := p.PodService.CreateToK8s(info); err != nil {
		common.Error(err)
		return nil, err
	} else {
		podID, err := p.PodService.AddPod(podModel)
		if err != nil {
			common.Error(err)
			return nil, err
		}
		common.Info("Pod 添加成功数据库ID号为：" + strconv.FormatInt(podID, 10))
		rsp.Msg = "Pod 添加成功数据库ID号为：" + strconv.FormatInt(podID, 10)
		return rsp, nil
	}
}

func (p *PodHandler) DeletePod(ctx context.Context, id *proto.PodId) (rsp *proto.Response, err error) {
	podModel, err := p.PodService.FindPodByID(id.Id)
	if err != nil {
		common.Error(err)
		return nil, err
	}
	if err := p.PodService.DeleteFromK8s(podModel); err != nil {
		common.Error(err)
		return nil, err
	}
	return rsp, nil

}

func (p *PodHandler) UpdatePod(ctx context.Context, info *proto.PodInfo) (rsp *proto.Response, err error) {
	//先更新k8s中的pod信息
	err = p.PodService.UpdateToK8s(info)
	if err != nil {
		common.Error(err)
		return nil, err
	}
	//查询数据库中的pod
	podModel, err := p.PodService.FindPodByID(info.Id)
	if err != nil {
		common.Error(err)
		return nil, err
	}
	err = common.SwapTo(info, podModel)
	if err != nil {
		common.Error(err)
		return nil, err
	}
	p.PodService.UpdatePod(podModel)
	return rsp, nil
}

func (p *PodHandler) FindAllPod(ctx context.Context, all *proto.FindAll) (allPods *proto.AllPod, err error) {
	allPod, err := p.PodService.FindAllPod()
	if err != nil {
		common.Error(err)
		return nil, err
	}
	//整理格式
	for _, v := range allPod {
		podInfo := &proto.PodInfo{}
		err := common.SwapTo(v, podInfo)
		if err != nil {
			common.Error(err)
			return nil, err
		}
		allPods.PodInfo = append(allPods.PodInfo, podInfo)
	}
	return allPods, nil
}

func (p *PodHandler) mustEmbedUnimplementedPodServer() {
}
