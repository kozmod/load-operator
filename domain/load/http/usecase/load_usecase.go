package usecase

import (
	"errors"
	"sync"

	"github.com/kozmod/load-operator/apis/load/v1alpha1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type LoadUseCase struct {
	mutex    sync.Mutex
	attacker *vegeta.Attacker
	metrics  *vegeta.Metrics
}

func NewLoadUseCase() *LoadUseCase {
	return &LoadUseCase{}
}

func (uc *LoadUseCase) Load(ls v1alpha1.HttpLoadService) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	if uc.attacker != nil {
		uc.attacker.Stop()
	}
	attacker := vegeta.NewAttacker()
	uc.attacker = attacker
	var metrics vegeta.Metrics
	uc.metrics = &metrics
	targeter := toTargetEncoder(ls.Spec.Target)
	rate := vegeta.Rate{
		Freq: ls.Spec.RateFreq,
		Per:  ls.Spec.RatePer.Duration,
	}
	duration := ls.Spec.Duration.Duration
	name := ls.Spec.Name
	go func(metrics *vegeta.Metrics) {
		for res := range attacker.Attack(targeter, rate, duration, name) {
			metrics.Add(res)
		}
	}(&metrics)
	return nil
}

func (uc *LoadUseCase) Metrics() (vegeta.Metrics, error) {
	if uc.attacker == nil {
		return vegeta.Metrics{}, errors.New("attacker not exists")
	}
	if uc.metrics == nil {
		return vegeta.Metrics{}, errors.New("metrics not exists")
	}
	uc.metrics.Close()
	return *uc.metrics, nil
}

func toTargetEncoder(target v1alpha1.Target) vegeta.Targeter {
	return vegeta.NewStaticTargeter(vegeta.Target{
		Method: target.Method,
		URL:    target.URL,
		Header: target.Header,
	})
}
