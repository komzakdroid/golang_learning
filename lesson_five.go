package main

import "fmt"

func main() {
	//Array
	var numbers [3]int = [3]int{1, 2, 3}
	fmt.Println("Array:", numbers)
	numbers[1] = 20 //change element
	fmt.Println("Array modified:", numbers)

	//Slice
	slice := []int{10, 20, 30}
	fmt.Println("Slice:", slice)
	slice = append(slice, 40) //Add element
	fmt.Println("Slice appended:", slice)

	//Slice with make
	slice2 := make([]int, 2, 5) //len = 2 , cap =5
	slice2[0] = 100
	slice2[1] = 200
	fmt.Println("Slice2:", slice2, "Len:", len(slice2), "Cap:", cap(slice2))

	//Map
	scores := map[string]int{"Ali": 90, "Vali": 85}
	fmt.Println("Map:", scores)
	scores["Sardor"] = 95  //Add
	delete(scores, "Vali") //Delete
	fmt.Println("Map modified:", scores)

	// Check if key exists
	score, ok := scores["Ali"]
	if ok {
		fmt.Println("Ali's score:", score)
	} else {
		fmt.Println("Ali not found")
	}
}
