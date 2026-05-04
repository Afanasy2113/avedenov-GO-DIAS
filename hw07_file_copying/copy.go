package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	// Открываем исходный файл.
	src, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	// Получаем информацию о файле.
	srcInfo, err := src.Stat()
	if err != nil {
		return fmt.Errorf("ошибка получения информации о файле: %w", err)
	}

	// Проверяем, что это обычный файл.
	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("это необычный файл")
	}

	fileSize := srcInfo.Size()

	// Проверка размера offset.
	if offset > fileSize {
		return fmt.Errorf("offset (%d) больше, чем размер файла (%d): %w", offset, fileSize, ErrOffsetExceedsFileSize)
	}

	if _, err := src.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("ошибка поиска offset: %w", err)
	}

	var bytesToCopy int64
	if limit == 0 {
		bytesToCopy = fileSize - offset
	} else {
		bytesToCopy = limit
		if bytesToCopy > fileSize-offset {
			bytesToCopy = fileSize - offset
		}
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {

		}
	}(dst)

	// Прогресс бар.
	progressThreshold := bytesToCopy / 100
	if progressThreshold == 0 {
		progressThreshold = 1
	}

	buf := make([]byte, 64*1024)
	var totalCopy int64
	var nextProgress int64 = progressThreshold

	// Начальное состояние.
	fmt.Fprintf(os.Stderr, "\rКопирование: 0%%")

	for totalCopy < bytesToCopy {
		toRead := int64(len(buf))
		if remaining := bytesToCopy - totalCopy; toRead > remaining {
			toRead = remaining
		}

		// Копируем.
		n, err := src.Read(buf[:toRead])
		if err != nil && err != io.EOF {
			return fmt.Errorf("ошибка чтения: %w", err)
		}

		// Пишем прочитанные данные.
		if n > 0 {
			if _, err := dst.Write(buf[:n]); err != nil {
				return fmt.Errorf("ошибка записи: %w", err)
			}
			totalCopy += int64(n)

			// Обновляем прогресс-бар.
			if bytesToCopy > 0 && totalCopy >= nextProgress {
				percent := int64(100 * totalCopy / bytesToCopy)
				if percent > 100 {
					percent = 100
				}
				fmt.Fprintf(os.Stderr, "\rКопирование: %d%%", percent)
				nextProgress += progressThreshold
			}
		}

		if err == io.EOF {
			break
		}
	}
	return nil
}
