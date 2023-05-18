package cron

import (
	"context"
	"github.com/gingerxman/eel/rest_client"
	"github.com/gingerxman/gorm"
	"github.com/pkg/errors"
	"math"
)

type taskInterface interface {
	Run(*TaskContext) error
	GetName() string
	IsEnableTx() bool
}

type Task struct {
	name string
}

func (t *Task) Run(taskContext *TaskContext) error {
	return errors.New("Run not implemented")
}

func (t *Task) GetName() string {
	return t.name
}

func (t *Task) SetName(name string) {
	t.name = name
}

func (t *Task) IsEnableTx() bool {
	return true
}

func NewTask(name string) Task {
	t := Task{name: name}
	return t
}

type pipeInterface interface {
	AddData(data interface{}) error
	GetData() interface{}
	GetCap() int
	GetConsumerCount() int
	RunConsumer(data interface{}, taskCtx *TaskContext)
	EnableParallel() bool
}

type Pipe struct {
	ch    chan interface{}
	chCap int
}

func (p Pipe) GetData() interface{} {
	return <-p.ch
}

func (p Pipe) AddData(data interface{}) error {
	select {
	case p.ch <- data:
	default:
		return errors.New("channel is full")
	}
	return nil
}

func (p Pipe) GetCap() int {
	return p.chCap
}

// GetConsumerCount 消费者数量
// 默认为通道容量十分之一
func (p Pipe) GetConsumerCount() int {
	return int(math.Ceil(float64(p.GetCap()) / 10))
}

func (p Pipe) RunConsumer() error {
	return errors.New("RunConsumer not implemented")
}

// EnableParallel 启用并行，默认启用
func (p Pipe) EnableParallel() bool {
	return true
}

func NewPipe(chCap int) Pipe {
	p := Pipe{}
	p.chCap = chCap
	p.ch = make(chan interface{}, chCap)
	return p
}

type TaskContext struct {
	orm      *gorm.DB
	resource *rest_client.Resource
	ctx      context.Context
}

func (this *TaskContext) Init(ctx context.Context, o *gorm.DB, resource *rest_client.Resource) {
	this.ctx = ctx
	this.orm = o
	this.resource = resource
}

func (this *TaskContext) GetOrm() *gorm.DB {
	if idb := this.ctx.Value("orm"); idb != nil {
		return idb.(*gorm.DB)
	}
	return nil
}

func (this *TaskContext) GetCtx() context.Context {
	return this.ctx
}

func (this *TaskContext) SetCtx(ctx context.Context) {
	this.ctx = ctx
}

func (this *TaskContext) GetResource() *rest_client.Resource {
	return this.resource
}

var managerToken string

func GetManagerResource(ctx context.Context) *rest_client.Resource {
	if managerToken == "" {
		resource := rest_client.NewResource(ctx).LoginAsManager()
		if resource != nil {
			managerToken = resource.CustomJWTToken
		}
	}

	res := rest_client.NewResource(ctx)
	res.CustomJWTToken = managerToken

	return res
}
