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

	maxErrors := m
	if m <= 0 {
		maxErrors = 0
	}

	taskCh := make(chan Task)
	errCh := make(chan error, 1)
	done := make(chan struct{}) // bidirectional!
	var wg sync.WaitGroup
	var errHandlerWg sync.WaitGroup

	// Запускаем воркеры
	startWorkers(n, taskCh, errCh, done, &wg)

	// Запускаем обработчик ошибок
	errCount := startErrorHandler(maxErrors, errCh, done, &errHandlerWg)

	// Отправляем задачи
	if err := feedTasks(tasks, taskCh, done); err != nil {
		close(taskCh)
		wg.Wait()
		errHandlerWg.Wait()
		return err
	}

	// Все задачи отправлены — закрываем канал
	close(taskCh)
	wg.Wait()

	// Дожидаемся завершения обработчика ошибок
	close(done)
	errHandlerWg.Wait()

	// Проверяем результат
	if maxErrors == 0 && errCount > 0 {
		return ErrErrorsLimitExceeded
	}
	if maxErrors > 0 && errCount >= maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// startWorkers запускает n воркеров
func startWorkers(n int, taskCh <-chan Task, errCh chan<- error, done <-chan struct{}, wg *sync.WaitGroup) {
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
}

// startErrorHandler запускает горутину, считающую ошибки
func startErrorHandler(maxErrors int, errCh <-chan error, done chan struct{}, wg *sync.WaitGroup) int {
	var errCount int
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCount := 0
		for {
			select {
			case <-errCh:
				errCount++
				if errCount >= maxErrors {
					close(done)
					return
				}
			case <-done:
				return
			}
		}
	}()
	return errCount // возвращаем для проверки, но не используется напрямую
}

// feedTasks отправляет все задачи в канал
func feedTasks(tasks []Task, taskCh chan<- Task, done <-chan struct{}) error {
	for _, task := range tasks {
		select {
		case taskCh <- task:
		case <-done:
			return ErrErrorsLimitExceeded
		}
	}
	return nil
}
