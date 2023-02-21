package plan

import (
	"sync"
	"testing"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

func TestStepSecuence(t *testing.T) {
	desc := "Testing plan steps secuence"

	t.Run(desc, func(t *testing.T) {
		//t.Log(desc)
		var wg sync.WaitGroup

		stepPlan1 := make(chan struct{})
		step1 := NewStep(&image.Image{}, "plan1", stepPlan1)

		stepPlan2 := make(chan struct{})
		step2 := NewStep(&image.Image{}, "plan2", stepPlan2)

		stepPlan3 := make(chan struct{})
		step3 := NewStep(&image.Image{}, "plan3", stepPlan3)

		stepPlan4 := make(chan struct{})
		step4 := NewStep(&image.Image{}, "plan4", stepPlan4)

		step1.Subscribe(stepPlan2)
		step1.Subscribe(stepPlan3)
		step3.Subscribe(stepPlan4)

		stepFunc := func(step *Step) {
			wg.Add(1)
			step.Wait()
			step.Notify()
			wg.Done()
		}

		go stepFunc(step1)
		go stepFunc(step2)
		go stepFunc(step3)
		go stepFunc(step4)

		step1.sync <- struct{}{}
		wg.Wait()
	})
}
