package controller

import (
	"fmt"
	"time"
)

func main() {
	var inputDate string
	fmt.Print("Masukkan tanggal (DD/MM/YYYY): ")
	fmt.Scanln(&inputDate)

	date, err := time.Parse("02/01/2006", inputDate)
	if err != nil {
		fmt.Println("Format tanggal salah!")
		return
	}

	daysOfWeek := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	day := int(date.Weekday())

	fmt.Println("Hari", daysOfWeek[day])
}
