package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
	"time"
)

// Функция для обработки одной строки пикселей в параллельном режиме
func filterParallel(img draw.Image, y int, wg *sync.WaitGroup) {
	defer wg.Done()

	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		// Получаем значение текущего пикселя
		originalColor := img.At(x, y).(color.RGBA)

		// Рассчитываем среднее значение цветовых каналов (для оттенков серого)
		gray := uint8((uint32(originalColor.R) + uint32(originalColor.G) + uint32(originalColor.B)) / 3)

		// Устанавливаем новый цвет для пикселя
		img.Set(x, y, color.RGBA{R: gray, G: gray, B: gray, A: originalColor.A})
	}
}

func main() {
	// 1. Открываем входной файл
	inputFile, err := os.Open("input.png")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer inputFile.Close()

	// 2. Декодируем изображение
	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return
	}

	// 3. Преобразуем изображение в редактируемый тип
	// Создаем новое изображение, совместимое с редактируемым форматом
	drawImg := image.NewRGBA(img.Bounds())
	draw.Draw(drawImg, img.Bounds(), img, image.Point{}, draw.Src)

	// 4. Замеряем время выполнения фильтрации в параллельном режиме
	start := time.Now()

	var wg sync.WaitGroup
	bounds := drawImg.Bounds()

	// Создаем горутину для каждой строки
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go filterParallel(drawImg, y, &wg)
	}

	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Параллельное преобразование изображения заняло: %v\n", duration)

	// 5. Создаем выходной файл
	outputFile, err := os.Create("output_parallel.png")
	if err != nil {
		fmt.Println("Ошибка при создании выходного файла:", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, drawImg)
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return
	}

	fmt.Println("Изображение успешно обработано и сохранено в 'output_parallel.png'")
}
