package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := in

	for key := range stages {
		stage := stages[key]
		stageResult := stage(result)
		stageChan := make(Bi)

		go func() {
			defer close(stageChan)
			for {
				select {
				case <-done:
					return
				case v, ok := <-stageResult:
					if !ok {
						return
					}
					stageChan <- v
				}
			}

		}()
		result = stageChan
	}
	return result
}
