package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("C:\\Users\\vukco\\OneDrive\\Desktop\\test.txt")
	if err != nil {
		fmt.Errorf("error opening file")
	}
	defer f.Close()

	h := sha256.New()

	_, err = io.Copy(h, f)
	if err != nil {
		fmt.Errorf("error copying file")
	}

	fmt.Printf("type BEFORE sum %T\n", h)
	fmt.Printf("Value %v\n", h)

	xb := h.Sum(nil)
	fmt.Printf("type AFTER sum %T", xb)
	fmt.Printf("%v\n", xb)

	xb = h.Sum(nil)
	fmt.Printf("type AFTER second sum %T", xb)
	fmt.Printf("%v\n", xb)

	xb = h.Sum([]byte("cao"))
	fmt.Printf("type AFTER third sum %T", xb)
	fmt.Printf("%v\n", xb)

	xb = h.Sum(xb)
	fmt.Printf("type AFTER fourth sum %T", xb)
	fmt.Printf("%v\n", xb)
}
