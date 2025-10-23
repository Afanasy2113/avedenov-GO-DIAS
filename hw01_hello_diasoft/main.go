package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

// Функция переворота текста.
func reverseHello(text string) string {
	return reverse.String(text)
}

// Запуск программы.
func main() {
	// Выводим в консоль перевернутый текст.
	fmt.Println(reverseHello("Hello, DIASOFT!"))
}
