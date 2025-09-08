package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	// Place your code here.
	var result []rune
	runes := []rune(s)

	// Если передана пустая строка.
	if len(runes) == 0 {
		return "", nil
	}

	// Проверка, что первый символ не цифра.
	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	// Идем по строке.
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Проверяем, что текущий символ буква.
		if !unicode.IsDigit(r) {
			result = append(result, r)
			continue
		}

		// Проверяем, что текущий символ цифра.
		if unicode.IsDigit(r) {

			// Предыдущий символ должен быть буквой.
			if i == 0 || unicode.IsDigit(runes[i-1]) {
				return "", ErrInvalidString
			}
			// Преобразуем руну в строку.
			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", ErrInvalidString
			}

			// Удаляем символ, если после него стоит 0.
			if count == 0 {
				if len(result) > 0 {
					result = result[:len(result)-1]
				}
				continue
			}

			// Дублируем предыдущий символ в зависимости от значения
			last := result[len(result)-1]
			for j := 0; j < count-1; j++ {
				result = append(result, last)
			}
		}
	}

	return string(result), nil
}
