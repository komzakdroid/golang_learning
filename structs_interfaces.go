package main

import "fmt"

// Struct
type User struct {
	Name  string
	Age   int
	Email string
}

// Embedded Struct
type Admin struct {
	User //Embedded User struct
	Role string
}

// Interface
type Printable interface {
	PrintInfo() string
}

// Method for User
func (u User) PrintInfo() string {
	return fmt.Sprintf("Name: %s, Age: %d, Email: %s", u.Name, u.Age, u.Email)
}

// Method for Admin
func (a Admin) PrintInfo() string {
	return fmt.Sprintf("Name: %s, Age: %d, Email: %s, Role: %s", a.Name, a.Age, a.Email, a.Role)
}

// Interface ishlatish
func printEntity(p Printable) {
	fmt.Println(p.PrintInfo())
}

func main() {
	//Struct yaratish
	user := User{Name: "Ali", Age: 25, Email: "ali@example.com"}
	admin := Admin{User: User{Name: "Vali", Age: 30, Email: "vali@example.com"}, Role: "Admin"}

	//Direct method chaqirish
	fmt.Println(user.PrintInfo())
	fmt.Println(admin.PrintInfo())

	//Interface orqali
	printEntity(user)
	printEntity(admin)

	//Slice of structs
	users := []User{
		{Name: "Sardor", Age: 22, Email: "sardor@example.com"},
		{Name: "Nodir", Age: 28, Email: "nodir@example.com"},
	}

	//Loop bilan print
	for _, u := range users {
		fmt.Println(u.PrintInfo())
	}

}
