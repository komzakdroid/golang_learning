package main

import "fmt"

func main() {

	age := 26

	if age < 18 {
		fmt.Println("Kichik")
	} else if age >= 18 && age < 65 {
		fmt.Println("Katta, ammo pensiyaga emas")
	} else {
		fmt.Println("Pensioner")
	}

	//For Loop misol: Classic
	for i := 0; i < 5; i++ {
		fmt.Println("Iteratsiya:", i)
	}

	//For While-like
	counter := 0
	for counter < 3 {
		fmt.Println("While-like:", counter)
		counter++
	}

	//Switch misol
	day := "Dushanba"

	switch day {
	case "Dushanba", "Seshanba":
		fmt.Println("Ish haftasi boshlandi")

	case "Shanba":
		fmt.Println("Dam olish")
	default:
		fmt.Println("Oddiy kun")
	}

	//Range misol (array uchun, keyingi darslarda batafsil)
	numbers := []int{10, 20, 30}
	for index, value := range numbers {
		fmt.Println("Index:", index, "Qiymat:", value)
	}

	fmt.Println("Vazifalar:")

	if age < 13 {
		fmt.Println("Bola")
	} else if age >= 13 && age <= 18 {
		fmt.Println("O'smir")
	} else {
		fmt.Println("Katta")
	}

	for j := 1; j <= 10; j++ {
		if j%2 == 0 {
			fmt.Println("Juft Index", j)
		} else {
			fmt.Println("Index", j)
		}
	}

	selectedDay := "Juma"
	switch selectedDay {
	case "Shanba", "Yakshanba":
		fmt.Println("Dam olish")
	default:
		fmt.Println("Ish")
	}

}
