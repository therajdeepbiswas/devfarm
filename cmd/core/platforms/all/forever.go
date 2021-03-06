package all

import (
	"github.com/dena/devfarm/cmd/core/platforms"
	"sync"
)

func (ps Platforms) ForeverIOS(iosPlan platforms.IOSPlan) error {
	p, err := ps.GetPlatform(iosPlan.Platform)
	if err != nil {
		return err
	}

	runIOS := p.IOSForever()
	return runIOS(iosPlan)
}

func (ps Platforms) ForeverAndroid(iosPlan platforms.AndroidPlan) error {
	p, err := ps.GetPlatform(iosPlan.Platform)
	if err != nil {
		return err
	}

	runAndroid := p.AndroidForever()
	return runAndroid(iosPlan)
}

func (ps Platforms) Forever(plans []platforms.EitherPlan) (ResultTable, error) {
	builder := NewResultTableBuilder()
	var wg sync.WaitGroup

	for platformID, plansForPlatform := range groupByPlatform(plans) {
		p, err := ps.GetPlatform(platformID)
		if err != nil {
			builder.AddErrors(platformID, err)
			continue
		}

		wg.Add(1)
		go func(p platforms.Platform, plansForPlatform []platforms.EitherPlan) {
			forever := p.Forever()

			results, _ := forever(plansForPlatform)
			builder.AddErrors(p.ID(), results...)

			wg.Done()
		}(p, plansForPlatform)
	}

	wg.Wait()
	table := builder.Build()
	return table, table.Err()
}
