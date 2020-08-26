package usecase

import (
	"context"
	"sync"

	"github.com/go-logr/logr"
	"github.com/kozmod/load-operator/apis/load/v1alpha1"
	"github.com/kozmod/load-operator/domain/load/http/usecase"
	"github.com/kozmod/load-operator/domain/metrics/usecase/internal/get"
	"github.com/kozmod/load-operator/domain/metrics/usecase/internal/print"
	"github.com/kozmod/load-operator/domain/metrics/usecase/internal/schedule"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var once sync.Once

var suc scheduleUC
var luc loadUC = usecase.NewLoadUseCase()

func Schedule(ctx context.Context, ms v1alpha1.MetricsService) error {
	err := suc.Schedule(ctx, ms)
	if err != nil {
		return err
	}
	return nil
}

func Load(ls v1alpha1.HttpLoadService) error {
	err := luc.Load(ls)
	if err != nil {
		return err
	}
	return nil
}

func InitScheduleUseCase(cl client.Client, l logr.Logger) {
	once.Do(func() {
		config := rest.AnonymousClientConfig(&rest.Config{Host: "127.0.0.1:57927"}) //todo local test
		//config := rest.InClusterConfig() //todo in cluster
		//if err != nil {
		//	fmt.Println("rest error")
		//	panic(err.Error())
		//}
		suc = schedule.NewScheduleUseCase(
			print.New(get.New(config, cl, l), luc, l),
			l)
	})
}
