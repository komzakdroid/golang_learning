package main

import (
	"fmt"
)

// Global variables
var global = "Global scope"

// Simple function
func greet(name string) {
	fmt.Println("Salom", name)
}

// Multiple returns
func addAndSubtract(a, b int) (int, int) {
	sum := a + b // Block scope ichida emas
	diff := a - b
	return sum, diff
}

// Variadic
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Vazifalar
func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func ageCategory(age int) string {
	if age < 13 {
		return "Bola"
	} else if 13 >= age && age <= 18 {
		return "O`smir"
	}
	return "Katta"
}

func average(nums ...int) float64 {
	total := 0

	for _, n := range nums {
		total += n
	}
	return float64(total) / float64(len(nums))
}

func createCounter() (func(), func() int) {
	counter := 0 // Outer scope variable

	increment := func() {
		counter++ // Capture counter
		fmt.Println("Counter incremented:", counter)
	}

	get := func() int {
		return counter
	}

	return increment, get

}

// Anonymous and Closure
func main() {
	greet("Komiljon")

	s, d := addAndSubtract(10, 5)

	fmt.Println("Sum:", s, "Diff:", d)

	total := sum(1, 2, 3, 4)
	fmt.Println("Total:", total)

	//Anonymous function
	func() {
		fmt.Println("Anonymous function")
	}()

	//Closure
	counter := 0
	increment := func() {
		counter++ // Outer scope in capture
		fmt.Println("Counter", counter)
	}
	increment()
	increment()

	// Defer
	defer fmt.Println("Defer: Funksiya tugadi")
	fmt.Println(global) //Global access

	//Scope misol
	if true {
		local := "Block Scope"
		fmt.Println(local)
	}
	// fmt.Println(local) // Error: local bu yerda yo'q

	//Vazifalar
	max := max(10, 12)
	fmt.Println("Katta son:", max)

	ageCategory := ageCategory(26)
	fmt.Println("Yosh kategoriyasi:", ageCategory)

	average := average(10, 11, 12, 13, 14, 15)
	fmt.Println("O'rtacha qiymat:", average)

	// Test qilish
	inc, get := createCounter()

	inc()                                // Counter incremented: 1
	inc()                                // Counter incremented: 2
	fmt.Println("Current value:", get()) // Current value: 2
	inc()                                // Counter incremented: 3
	fmt.Println("Current value:", get()) // Current value: 3
}
