package main

import (
	"crypto"
	"fmt"
	"log"
	"os"
  "encoding/hex"
	"bytes"

	"github.com/foxboron/go-uefi/authenticode"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("authenticode-digest: [input]")
		os.Exit(1)
	}

  b, err := os.ReadFile(os.Args[1])
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  binary, err := authenticode.Parse(bytes.NewReader(b))
  if err != nil {
    log.Fatal(err)
  }


  buffer := binary.Hash(crypto.SHA256)


  fmt.Println(len(buffer))

  encodedString := hex.EncodeToString(buffer)

  fmt.Println("Base64 Encoded Content:", encodedString)

	// if err = os.WriteFile(args[1], file.Bytes(), 0644); err != nil {
	// 	log.Fatal(err)
	// }
}
