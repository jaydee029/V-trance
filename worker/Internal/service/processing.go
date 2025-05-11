package service

import "fmt"

func main() {
	Avgbitrate := map[int]int{
		240:  318,
		360:  506,
		480:  960,
		720:  1726,
		1080: 3399,
	}

	fmt.Println(Avgbitrate)
}
