package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}
	if n <= 0 {
		return ErrErrorsLimitExceeded
	}

	// Интерпретация: m <= 0 → "нельзя ни одной ошибки"
	maxErrors := m
	if m <= 0 {
		maxErrors = 0
	}

	taskCh := make(chan Task)
	errCh := make(chan error, 1) // небольшой буфер для избежания блокировки
	done := make(chan struct{})
	var wg sync.WaitGroup
	var errHandlerWg sync.WaitGroup

	// Запускаем n воркеров
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-taskCh:
					if !ok {
						return
					}
					if err := task(); err != nil {
						select {
						case errCh <- err:
						case <-done:
							return
						}
					}
				case <-done:
					return
				}
			}
		}()
	}

	// Горутина для подсчёта ошибок
	errCount := 0
	errHandlerWg.Add(1)
	go func() {
		defer errHandlerWg.Done()
		for {
			select {
			case <-errCh:
				errCount++
				if errCount >= maxErrors { // если достигли лимита — останавливаемся
					close(done)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Отправляем задачи
	for _, task := range tasks {
		select {
		case taskCh <- task:
		case <-done:
			// Остановились из-за ошибок — закрываем канал задач
			close(taskCh)
			wg.Wait()
			errHandlerWg.Wait()
			return ErrErrorsLimitExceeded
		}
	}
	// Все задачи отправлены — закрываем канал
	close(taskCh)
	wg.Wait()

	// Дожидаемся завершения обработчика ошибок
	close(done)
	errHandlerWg.Wait()

	// Проверяем, не превышен ли лимит
	if errCount >= maxErrors && maxErrors > 0 {
		return ErrErrorsLimitExceeded
	}
	if maxErrors == 0 && errCount > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
