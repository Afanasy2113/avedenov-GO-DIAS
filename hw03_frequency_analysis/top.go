package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func Top10(s string) []string {
	// Place your code here.

	// Разбиваем текст на отдельные слова.
	words := strings.Fields(s)

	// Считаем кол-во слов.
	wCount := make(map[string]int)

	// Считаем, сколько раз встретилось каждое слово.
	for _, word := range words {

		// Нижний регистр.
		lowerCase := strings.ToLower(word)

		var trimWord string

		// Корректно обрабатываем дефис: "-" - не слово, "----" - слово
		if onlyDefis(lowerCase) {
			trimWord = lowerCase
		} else {

			// Убираем знаки препинания по краям.
			start := 0
			for start < len(lowerCase) && unicode.IsPunct(rune(lowerCase[start])) {
				start++
			}

			end := len(lowerCase)
			for end > start && unicode.IsPunct(rune(lowerCase[end-1])) {
				end--
			}

			trimWord = lowerCase[start:end]
		}

		// Проверка слова на валидность.
		if trimWord != "" && trimWord != "-" {
			wCount[trimWord]++
		}
	}

	// Слайс "слово-частота".
	type wSlice struct {
		word  string
		count int
	}

	var wSlices []wSlice
	for word, count := range wCount {
		wSlices = append(wSlices, wSlice{word, count})
	}

	// Сортируем слайс.
	sort.Slice(wSlices, func(i, j int) bool {
		if wSlices[i].count == wSlices[j].count {
			return wSlices[i].word < wSlices[j].word
		}
		return wSlices[i].count > wSlices[j].count
	})

	//Выбираем топ-10 слов.
	result := make([]string, 0, 10)
	for i := 0; i < len(wSlices) && i < 10; i++ {
		result = append(result, wSlices[i].word)
	}

	//Вывод в консоль для отладки.
	fmt.Println("======Дебаг==========")
	for i := 0; i < len(wSlices) && i < 10; i++ {
		fmt.Printf("%d. '%s' (%d)\n", i+1, wSlices[i].word, wSlices[i].count)
	}
	fmt.Println("=====================")
	return result
}

// Корректная обработка дефисов.
func onlyDefis(s string) bool {
	for _, r := range s {
		if r != '-' {
			return false
		}
	}
	return s != ""
}
