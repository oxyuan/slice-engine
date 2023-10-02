package main

import (
	"flag"
	"fmt"
	"slice-engine/boot"
)

func main() {
	val := flag.String("text", "text", "Text that needs to be recognized and segmented")

	flag.Parse()

	fmt.Printf("%s", *val)
	fmt.Printf("\n--------------\n")

	group := boot.Run(*val)

	for _, v := range *group {
		fmt.Printf("%s\n", v)
	}

}
