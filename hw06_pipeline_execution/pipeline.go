package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	if len(stages) == 0 {
		return in
	}

	// current будет хранить выходной канал предыдущей стадии.
	current := in

	for _, stage := range stages {
		// Создаем промежуточный канал для передачи данных в текущую стадию.
		tempIn := make(chan interface{})

		// Запускаем горутину-передатчик, которая перенаправляет данные из current в tempIn.
		go supportFunc(current, tempIn, done)

		// Применяем текущую стадию к промежуточному каналу.
		stageOutput := stage(tempIn)

		// Создаем канал, который станет выходом для текущего этапа соединения.
		nextCurrent := make(chan interface{})

		// Запускаем горутину-приемник, которая перенаправляет данные из stageOutput в nextCurrent.
		go supportFunc(stageOutput, nextCurrent, done)

		// Обновляем current, чтобы он указывал на канал, из которого будет читать следующая стадия.
		current = nextCurrent
	}

	return current
}

// Вспомогательная функция, которая читает из входного канала `in`
// и отправляет данные в выходной канал `out`, пока не получит сигнал `done`.
func supportFunc(in In, out Bi, done In) {
	defer close(out)
	for {
		select {
		case item, ok := <-in:
			if !ok {
				return
			}
			select {
			case out <- item:
			case <-done:
				return
			}
		case <-done:
			return
		}
	}
}
