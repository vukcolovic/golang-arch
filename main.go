package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
)


func main() {
	msg := "Vuk Colovic"
	password := "ilovedogs"
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Panic("couldn't bcrypt password")
	}
	bs = bs[:16]

	rslt, err := enDecode(bs, msg)

	fmt.Println("Message: ", msg,  " Encoded:  ", string(rslt))

	rslt2, err := enDecode(bs, string(rslt))

	fmt.Println("Message: ", msg,  " Decoded with same method:  ", string(rslt2))

	wtr := &bytes.Buffer{}
	encWriter, err := encryptWriter(wtr, bs)

	_, err = io.WriteString(encWriter, msg)
	if err != nil {
		log.Panic("writte error %w", err)
	}

	fmt.Println("Same but other way: ", wtr.String())

}

func enDecode(key []byte, input string) ([]byte, error) {
	// cipher is algoritham, symetrical encription
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't newCipher %w", err)
	}

	//inicialization vector
	iv := make([]byte, aes.BlockSize)

	//create cipher
	s := cipher.NewCTR(b, iv)

	buff := &bytes.Buffer{}
	sw := cipher.StreamWriter{
		S: s,
		W: buff,
	}

	_, err = sw.Write([]byte(input))
	if err != nil {
		return nil, fmt.Errorf("couldn't sw.Write to streamwritter %w", err)
	}

	return buff.Bytes(), nil
}

func encryptWriter(w io.Writer, key []byte) (io.Writer, error) {
	// cipher is algoritham, symetrical encription
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't newCipher %w", err)
	}

	//inicialization vector
	iv := make([]byte, aes.BlockSize)

	//create cipher
	s := cipher.NewCTR(b, iv)

	return cipher.StreamWriter{
		S: s,
		W: w,
	}, nil
}