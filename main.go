package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var name [2]string
	n, _ := fmt.Fscanln(in, &name[0], &name[1])
	fmt.Println(n, name)
}
