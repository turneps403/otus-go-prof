package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func convHook(c chan interface{}) <-chan interface{} {
	return c
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// chan in the middle
	// tmp := make(chan interface{})
	// go func() {
	// 	defer close(tmp)
	// 	for {
	// 		select {
	// 		case v, ok := <-in:
	// 			if !ok {
	// 				fmt.Println("in not OK, return")
	// 				return
	// 			}
	// 			fmt.Printf("send value %v from in to tmp\n", v)
	// 			tmp <- v
	// 		case <-done:
	// 			fmt.Printf("got signal from done\n")
	// 			return
	// 		}
	// 	}
	// }()

	//out := convHook(tmp)
	out := in
	for _, s := range stages {
		tmp := make(chan interface{})
		go func(in In) {
			defer close(tmp)
			for {
				select {
				case v, ok := <-in:
					if !ok {
						// fmt.Println("in not OK, return")
						return
					}
					// fmt.Printf("send value %v from in to tmp\n", v)
					tmp <- v
				case <-done:
					// fmt.Printf("got signal from done\n")
					return
				}
			}
		}(out)
		out = convHook(tmp)
		out = s(out)
	}
	return out
}
