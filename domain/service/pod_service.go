package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/yqhcode/paasPod/domain/model"
	"github.com/yqhcode/paasPod/domain/repository"
	"github.com/yqhcode/paasPod/proto"

	common "github.com/yqhcode/paas-common"

	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodServicer interface {
	AddPod(*model.Pod) (int64, error)
	DeletePod(int64) error
	UpdatePod(*model.Pod) error
	FindPodByID(int64) (*model.Pod, error)
	FindAllPod() ([]model.Pod, error)
	CreateToK8s(info *proto.PodInfo) error
	DeleteFromK8s(*model.Pod) error
	UpdateToK8s(*proto.PodInfo) error
}

type PodService struct {
	PodRepositoryer repository.PodRepositoryer
	K8sClientSet    *kubernetes.Clientset
	deployment      *v1.Deployment
}

func (p *PodService) AddPod(pod *model.Pod) (int64, error) {
	return p.PodRepositoryer.CreatePod(pod)
}

func (p *PodService) DeletePod(i int64) error {
	return p.PodRepositoryer.DeletePodByID(i)
}

func (p *PodService) UpdatePod(pod *model.Pod) error {
	return p.PodRepositoryer.UpdatePod(pod)
}

func (p *PodService) FindPodByID(i int64) (*model.Pod, error) {
	return p.PodRepositoryer.FindPodByID(i)
}

func (p *PodService) FindAllPod() ([]model.Pod, error) {
	return p.PodRepositoryer.FindAll()
}

func (p *PodService) CreateToK8s(info *proto.PodInfo) error {
	p.SetDeployment(info)
	if _, err := p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Get(context.TODO(), info.PodName, v12.GetOptions{}); err != nil {
		if _, err := p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Create(context.TODO(), p.deployment, v12.CreateOptions{}); err != nil {
			common.Error(err)
			return err
		}
		common.Info("创建成功")
		return nil
	} else {
		common.Error("Pod " + info.PodName + "已经存在")
		return errors.New("Pod " + info.PodName + " 已经存在")
	}
}

func (p *PodService) DeleteFromK8s(pod *model.Pod) error {
	if err := p.K8sClientSet.AppsV1().Deployments(pod.PodNamespace).Delete(context.TODO(), pod.PodName, v12.DeleteOptions{}); err != nil {
		common.Error(err)
		return err
	} else {
		if err := p.PodRepositoryer.DeletePodByID(pod.ID); err != nil {
			common.Error(err)
			return err
		}
		common.Info("删除Pod ID：" + strconv.FormatInt(pod.ID, 10) + " 成功！")
	}
	return nil
}

func (p *PodService) UpdateToK8s(info *proto.PodInfo) error {
	p.SetDeployment(info)
	if _, err := p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Get(context.TODO(), info.PodName, v12.GetOptions{}); err != nil {
		common.Error(err)
		return errors.New("Pod " + info.PodName + " 不存在请先创建")
	} else {
		if _, err := p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Update(context.TODO(), p.deployment, v12.UpdateOptions{}); err != nil {
			common.Error(err)
			return err
		}
		common.Info(info.PodName + " 更新成功")
		return nil
	}
}

func NewPodService(repositoryer repository.PodRepositoryer, clientset *kubernetes.Clientset) PodServicer {
	return &PodService{
		PodRepositoryer: repositoryer,
		K8sClientSet:    clientset,
		deployment:      &v1.Deployment{},
	}
}
func (p *PodService) SetDeployment(podInfo *proto.PodInfo) {
	deployment := &v1.Deployment{}
	deployment.TypeMeta = v12.TypeMeta{
		Kind:       "deployment",
		APIVersion: "v1",
	}
	deployment.ObjectMeta = v12.ObjectMeta{
		Name:      podInfo.PodName,
		Namespace: podInfo.PodNamespace,
		Labels: map[string]string{
			"app-name": podInfo.PodName,
			"author":   "Caplost",
		},
	}
	deployment.Name = podInfo.PodName
	deployment.Spec = v1.DeploymentSpec{
		//副本个数
		Replicas: &podInfo.PodReplicas,
		Selector: &v12.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": podInfo.PodName,
			},
			MatchExpressions: nil,
		},
		Template: v13.PodTemplateSpec{
			ObjectMeta: v12.ObjectMeta{
				Labels: map[string]string{
					"app-name": podInfo.PodName,
				},
			},
			Spec: v13.PodSpec{
				Containers: []v13.Container{
					{
						Name:            podInfo.PodName,
						Image:           podInfo.PodImage,
						Ports:           p.getContainerPort(podInfo),
						Env:             p.getEnv(podInfo),
						Resources:       p.getResources(podInfo),
						ImagePullPolicy: p.getImagePullPolicy(podInfo),
					},
				},
			},
		},
		Strategy:                v1.DeploymentStrategy{},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    nil,
		Paused:                  false,
		ProgressDeadlineSeconds: nil,
	}
	p.deployment = deployment
}
func (p *PodService) getContainerPort(podInfo *proto.PodInfo) (containerPort []v13.ContainerPort) {
	for _, v := range podInfo.PodPort {
		containerPort = append(containerPort, v13.ContainerPort{
			Name:          "port-" + strconv.FormatInt(int64(v.ContainerPort), 10),
			ContainerPort: v.ContainerPort,
			Protocol:      p.getProtocol(v.Protocol),
		})
	}
	return
}
func (p *PodService) getProtocol(protocol string) v13.Protocol {
	switch protocol {
	case "TCP":
		return "TCP"
	case "UDP":
		return "UDP"
	case "SCTP":
		return "SCTP"
	default:
		return "TCP"
	}
}
func (p *PodService) getEnv(podInfo *proto.PodInfo) (envVar []v13.EnvVar) {
	for _, v := range podInfo.PodEnv {
		envVar = append(envVar, v13.EnvVar{
			Name:      v.EnvKey,
			Value:     v.EnvValue,
			ValueFrom: nil,
		})
	}
	return
}
func (p *PodService) getResources(podInfo *proto.PodInfo) (source v13.ResourceRequirements) {
	//最大能够使用多少资源
	source.Limits = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(podInfo.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(podInfo.PodMemoryMax), 'f', 6, 64)),
	}
	//满足最少使用的资源量
	//@TODO 自己实现动态设置
	source.Requests = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(podInfo.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(podInfo.PodMemoryMax), 'f', 6, 64)),
	}
	return
}

func (p *PodService) getImagePullPolicy(podInfo *proto.PodInfo) v13.PullPolicy {
	switch podInfo.PodPullPolicy {
	case "Always":
		return "Always"
	case "Never":
		return "Never"
	case "IfNotPresent":
		return "IfNotPresent"
	default:
		return "Always"
	}
}
