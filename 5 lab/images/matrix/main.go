package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
)

// Пример ядра свёртки для размытия (матрица 3x3)
var kernel = [3][3]float64{
	{0.0625, 0.125, 0.0625},
	{0.125, 0.25, 0.125},
	{0.0625, 0.125, 0.0625},
}

// Применяет фильтр к одной строке пикселей
func applyConvolutionRow(src draw.Image, dst draw.Image, y int, wg *sync.WaitGroup) {
	defer wg.Done()

	bounds := src.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// Обрабатываем каждый пиксель строки
	for x := 0; x < width; x++ {
		var rSum, gSum, bSum float64

		for ky := -1; ky <= 1; ky++ {
			for kx := -1; kx <= 1; kx++ {
				neighborX := x + kx
				neighborY := y + ky

				// Проверяем, что соседний пиксель находится в пределах изображения
				if neighborX >= 0 && neighborX < width && neighborY >= 0 && neighborY < height {
					neighborColor := src.At(neighborX, neighborY).(color.RGBA)
					weight := kernel[ky+1][kx+1]

					rSum += float64(neighborColor.R) * weight
					gSum += float64(neighborColor.G) * weight
					bSum += float64(neighborColor.B) * weight
				}
			}
		}

		// Приводим значения к uint8
		dst.Set(x, y, color.RGBA{
			R: uint8(rSum),
			G: uint8(gSum),
			B: uint8(bSum),
			A: src.At(x, y).(color.RGBA).A,
		})
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
	srcImg := image.NewRGBA(img.Bounds())
	draw.Draw(srcImg, img.Bounds(), img, image.Point{}, draw.Src)

	// Создаем пустое изображение для результата
	dstImg := image.NewRGBA(img.Bounds())

	// 4. Параллельная обработка строк изображения
	var wg sync.WaitGroup
	for y := 0; y < img.Bounds().Max.Y; y++ {
		wg.Add(1)
		go applyConvolutionRow(srcImg, dstImg, y, &wg)
	}

	wg.Wait()

	// 5. Сохраняем обработанное изображение
	outputFile, err := os.Create("output_convolution.png")
	if err != nil {
		fmt.Println("Ошибка при создании выходного файла:", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, dstImg)
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return
	}

	fmt.Println("Изображение успешно обработано и сохранено в 'output_convolution.png'")
}
