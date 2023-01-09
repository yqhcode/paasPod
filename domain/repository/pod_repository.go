package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/yqhcode/paasPod/domain/model"
)

//创建需要实现的接口
type PodRepositoryer interface {
	//初始化表
	InitTable() error
	//根据ID查找数据
	FindPodByID(int64) (*model.Pod, error)
	//创建一条 Pod 数据
	CreatePod(*model.Pod) (int64, error)
	//根据ID删除一条 Pod 数据
	DeletePodByID(int64) error
	//修改一条数据
	UpdatePod(*model.Pod) error
	//查找Pod所有数据
	FindAll() ([]model.Pod, error)
}

type PodRepository struct {
	mysqlDB *gorm.DB
}

//初始化3个表
func (p *PodRepository) InitTable() error {
	return p.mysqlDB.CreateTable(&model.Pod{}, &model.PodPort{}, &model.PodEnv{}).Error
}

func (p *PodRepository) FindPodByID(i int64) (*model.Pod, error) {
	pod := &model.Pod{}
	return pod, p.mysqlDB.Preload("PodEnv").Preload("PodPort").First(pod, i).Error
}

func (p *PodRepository) CreatePod(pod *model.Pod) (int64, error) {
	return pod.ID, p.mysqlDB.Create(pod).Error
}

func (p *PodRepository) DeletePodByID(i int64) error {
	tx := p.mysqlDB.Begin()
	defer func() {
		if recover() != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	//彻底删除 POD 信息

	if err := p.mysqlDB.Where("id=?", i).Delete(&model.Pod{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//彻底删除 podenv 信息
	if err := p.mysqlDB.Where("id = ?", i).Delete(&model.PodEnv{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//彻底删除 podport 信息
	if err := p.mysqlDB.Where("id = ?", i).Delete(&model.PodPort{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//更新Pod信息
func (p *PodRepository) UpdatePod(pod *model.Pod) error {
	return p.mysqlDB.Model(pod).Update(pod).Error
}

//获取结果集合
func (p *PodRepository) FindAll() (podAll []model.Pod, err error) {
	return podAll, p.mysqlDB.Find(&podAll).Error
}

//创建 PodRepository
func NewPodRepository(DB *gorm.DB) PodRepositoryer {
	return &PodRepository{mysqlDB: DB}
}
