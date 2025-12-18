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

		// Горутина, которая читает из предыдущего канала (current) и отправляет в tempIn.
		go func(prevIn In, targetIn Bi) {
			defer close(targetIn)
			for {
				select {
				case item, ok := <-prevIn:
					if !ok {
						// Канал предыдущей стадии закрыт.
						return
					}
					// Пытаемся отправить в промежуточный канал.
					select {
					case targetIn <- item:
					case <-done:
						// Сигнал остановки получен, выходим.
						return
					}
				case <-done:
					// Сигнал остановки получен, выходим.
					return
				}
			}
		}(current, tempIn)

		// Применяем текущую стадию к промежуточному каналу.
		stageOutput := stage(tempIn)

		// Создаем канал, который станет выходом для текущего этапа соединения.
		nextCurrent := make(chan interface{})

		// Горутина, которая читает из выхода стадии (stageOutput) и отправляет в nextCurrent.
		go func(stageOut In, finalOut Bi) {
			defer close(finalOut)
			for {
				select {
				case item, ok := <-stageOut:
					if !ok {
						// Канал стадии закрыт.
						return
					}
					// Отправика в финальный канал для следующей стадии.
					select {
					case finalOut <- item:
					case <-done:
						// Сигнал остановки получен, выходим.
						return
					}
				case <-done:
					// Сигнал остановки получен, выходим.
					return
				}
			}
		}(stageOutput, nextCurrent)

		// Обновляем current, чтобы он указывал на канал, из которого будет читать следующая стадия.
		current = nextCurrent
	}

	return current
}
