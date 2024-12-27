package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	//1.1
	ip := [4]byte{127, 0, 0, 1}
	fmt.Println(formatIP(ip))
	//1.2
	evens, err := listEven(1, 10)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Чётные числа:", evens)
	}
	//2
	s := "hello, world!"
	count := countChars(s)

	for char, freq := range count {
		fmt.Printf("%c: %d\n", char, freq)
	}
	//3
	triangle := Triangle{
		A: Point{0, 0},
		B: Point{3, 0},
		C: Point{3, 4},
	}

	circle := Circle{
		Center: Point{0, 0},
		Radius: 5,
	}

	printArea(triangle)
	printArea(circle)

	//4
	values := []float64{1, 2, 3, 4, 5}

	squaredValues := Map(values, square)

	fmt.Println("Срез после применения функции Map:", squaredValues)

}

// 1 Задание

func formatIP(ip [4]byte) string {
	var parts []string
	for _, b := range ip {
		parts = append(parts, strconv.Itoa(int(b)))
	}
	return strings.Join(parts, ".")
}

func listEven(start, end int) ([]int, error) {
	if start > end {
		return nil, errors.New("левая граница больше правой")
	}

	var evens []int
	for i := start; i <= end; i++ {
		if i%2 == 0 {
			evens = append(evens, i)
		}
	}
	return evens, nil
}

// 2 Задание

func countChars(s string) map[rune]int {
	charCount := make(map[rune]int)
	for _, char := range s {
		charCount[char]++
	}
	return charCount
}

// 3 Задание

type Point struct {
	X float64
	Y float64
}

type Segment struct {
	Start Point
	End   Point
}

func (s Segment) Length() float64 {
	return math.Sqrt(math.Pow(s.End.X-s.Start.X, 2) + math.Pow(s.End.Y-s.Start.Y, 2))
}

type Triangle struct {
	A Point
	B Point
	C Point
}

type Circle struct {
	Center Point
	Radius float64
}

func (t Triangle) Area() float64 {
	ab := Segment{Start: t.A, End: t.B}
	bc := Segment{Start: t.B, End: t.C}
	ca := Segment{Start: t.C, End: t.A}

	a := ab.Length()
	b := bc.Length()
	c := ca.Length()

	s := (a + b + c) / 2
	return math.Sqrt(s * (s - a) * (s - b) * (s - c))
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

type Shape interface {
	Area() float64
}

func printArea(s Shape) {
	result := s.Area()
	fmt.Printf("Площадь фигуры: %.2f\n", result)
}

func Map(input []float64, fn func(float64) float64) []float64 {
	result := make([]float64, len(input))
	for i, v := range input {
		result[i] = fn(v)
	}
	return result
}

func square(x float64) float64 {
	return x * x
}
