package pipeline

type Stage func(<-chan bool, <-chan int) <-chan int

type Pipeline struct {
	stages []Stage
	exit   <-chan bool
}

func NewPipeline(exit <-chan bool, stages ...Stage) *Pipeline {
	return &Pipeline{
		stages: stages,
		exit:   exit,
	}
}

func (p *Pipeline) runStage(stage Stage, source <-chan int) <-chan int {
	return stage(p.exit, source)
}

func (p *Pipeline) Run(source <-chan int) <-chan int {
	var c <-chan int = source
	for i := range p.stages {
		c = p.runStage(p.stages[i], c)
	}
	return c
}
