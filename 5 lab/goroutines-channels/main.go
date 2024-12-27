package main

import (
	"fmt"
	"sync"
)

// Функция count читает числа из канала и обрабатывает их
func count(numbers <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Сообщаем WaitGroup о завершении работы горутины

	for num := range numbers {
		// Возведение числа в квадрат
		squared := num * num
		fmt.Printf("Число: %d, Квадрат: %d\n", num, squared)
	}
	fmt.Println("Канал закрыт, работа завершена.")
}

func main() {
	// Создаем канал для передачи чисел
	numbers := make(chan int)

	// Создаем WaitGroup для синхронизации
	var wg sync.WaitGroup

	// Добавляем 1 в WaitGroup, так как запускаем одну горутину
	wg.Add(1)

	// Запускаем функцию count в отдельной горутине
	go count(numbers, &wg)

	// Отправляем несколько чисел в канал
	for i := 1; i <= 5; i++ {
		numbers <- i
	}

	// Закрываем канал
	close(numbers)

	// Ожидаем завершения горутины count
	wg.Wait()
	fmt.Println("Программа завершена.")
}
