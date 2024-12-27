package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"time"
)

// Функция для применения фильтра "оттенки серого"
func filter(img draw.Image) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Получаем значение текущего пикселя
			originalColor := img.At(x, y).(color.RGBA)

			// Рассчитываем среднее значение цветовых каналов (для оттенков серого)
			gray := uint8((uint32(originalColor.R) + uint32(originalColor.G) + uint32(originalColor.B)) / 3)

			// Устанавливаем новый цвет для пикселя
			img.Set(x, y, color.RGBA{R: gray, G: gray, B: gray, A: originalColor.A})
		}
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

	// 4. Замеряем время выполнения фильтра
	start := time.Now()
	filter(drawImg)
	duration := time.Since(start)
	fmt.Printf("Преобразование изображения заняло: %v\n", duration)

	outputFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Ошибка при создании выходного файла:", err)
		return
	}
	defer outputFile.Close()

	// 6. Сохраняем обработанное изображение
	err = png.Encode(outputFile, drawImg)
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return
	}

	fmt.Println("Изображение успешно обработано и сохранено в 'output.png'")
}
