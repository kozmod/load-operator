package load

import (
	v1 "github.com/kozmod/load-operator/apis/cache/v1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Loader struct {
	metrics *vegeta.Metrics
}

func NewLoader() *Loader {
	return &Loader{}
}

func (uc *Loader) Load(ls v1.LoadService) error {
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	uc.metrics = &metrics
	target := vegeta.NewStaticTargeter(vegeta.Target{
		Method: ls.Spec.Loader.Target.Method,
		URL:    ls.Spec.Loader.Target.URL,
	})
	rate := vegeta.Rate{
		Freq: ls.Spec.Loader.RateFreq,
		Per:  ls.Spec.Loader.RatePer.Duration,
	}
	duration := ls.Spec.Loader.Duration.Duration
	name := ls.Spec.Loader.Name
	for res := range attacker.Attack(target, rate, duration, name) {
		metrics.Add(res)
	}
	return nil
}

func (uc *Loader) Metrics() vegeta.Metrics {
	uc.metrics.Close()
	return *uc.metrics
}
