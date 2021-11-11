package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, s := range stages {
		out = channelWrap(out, done)
		out = s(out)
	}
	return out
}

func convHook(c chan interface{}) <-chan interface{} {
	return c
}

func channelWrap(in In, done In) Out {
	tmp := make(chan interface{})
	go func(in In) {
		defer close(tmp)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				tmp <- v
			case <-done:
				return
			}
		}
	}(in)
	return convHook(tmp)
}
