package usecase

import (
	"errors"
	"github.com/kozmod/load-operator/apis/load/v1alpha1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type LoadUseCase struct {
	attacker *vegeta.Attacker
	metrics  *vegeta.Metrics
}

func NewLoadUseCase() *LoadUseCase {
	return &LoadUseCase{}
}

func (uc *LoadUseCase) Load(ls v1alpha1.HttpLoadService) error {
	if uc.attacker != nil {
		uc.attacker.Stop()
	}
	uc.attacker = vegeta.NewAttacker()
	var metrics vegeta.Metrics
	uc.metrics = &metrics
	target := vegeta.NewStaticTargeter(vegeta.Target{
		Method: ls.Spec.Target.Method,
		URL:    ls.Spec.Target.URL,
	})
	rate := vegeta.Rate{
		Freq: ls.Spec.RateFreq,
		Per:  ls.Spec.RatePer.Duration,
	}
	duration := ls.Spec.Duration.Duration
	name := ls.Spec.Name
	go func(metrics *vegeta.Metrics) {
		for res := range uc.attacker.Attack(target, rate, duration, name) {
			metrics.Add(res)
		}
	}(uc.metrics)
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
