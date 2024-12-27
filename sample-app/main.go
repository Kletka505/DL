package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	fmt.Println("#1")
	hello("Vasya")
	fmt.Println("#2")
	fmt.Println(printEven(1, 5))
	fmt.Println("#3")
	fmt.Println(apply(1, 5, "+"))

}
func hello(s string) {
	fmt.Println("Hello,", s)
}
func printEven(a, b int64) error {
	if a < b {
		if a%2 == 0 {
			for i := a; i <= b; i += 2 {
				fmt.Println(i)
			}
		} else {
			for i := a + 1; i <= b; i += 2 {
				fmt.Println(i)
			}
		}
		return nil
	} else {
		return fmt.Errorf("левая граница должна быть меньше правой")
	}
}
func apply(a float64, b float64, op string) (float64, error) {
	if op == "-" {
		s := a - b
		return s, nil
	}
	if op == "+" {
		s := a + b
		return s, nil
	}
	if op == "*" {
		s := a * b
		return s, nil
	}
	if op == "/" {
		s := a / b
		return s, nil
	} else {
		return 0, fmt.Errorf("неверный оператор")
	}
}
