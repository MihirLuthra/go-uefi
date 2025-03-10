package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/foxboron/go-uefi/authenticode"
	"github.com/foxboron/go-uefi/efi/util"
)

func main() {
  signingTimeStr := flag.String("signing-time", "", "Signing time in ISO 8601 format (e.g., 2023-03-09T12:34:56Z)")
  inputFile := flag.String("input-file", "", "Input File")
  outputFile := flag.String("output-file", "", "Output File")
  signatureToAttach := flag.String("attach-signature", "", "Signature to attach given in a file")
  cert := flag.String("cert", "", "Certificate")
	flag.Parse()

  if *inputFile == "" {
    fmt.Println("Missing --input-file")
  }

  if *signingTimeStr == "" {
    fmt.Println("Missing --signing-time")
  }

  signingTime, err := time.Parse(time.RFC3339, *signingTimeStr)
  if err != nil {
    fmt.Println("Error parsing signing time:", err)
    return
  }

  b, err := os.ReadFile(*inputFile)
  if err != nil {
    log.Fatal(err)
  }

  binary, err := authenticode.Parse(bytes.NewReader(b))
  if err != nil {
    log.Fatal(err)
  }

  if *signatureToAttach == "" {
    tbs, err := binary.ToBeSignedData(signingTime)
    if err != nil {
      log.Fatal(err)
    }

    encodedString := base64.StdEncoding.EncodeToString(tbs.Tbs)
    fmt.Println(encodedString)
  } else {
    if *outputFile == "" {
      log.Fatal("Missing --outputFile")
    }

    if *cert == "" {
      log.Fatal("Missing --cert")
    }

    sig, err := os.ReadFile(*signatureToAttach)
    if err != nil {
      log.Fatal(err)
    }

    println(len(sig))

    cert, err := util.ReadCertFromFile(*cert)
    if err != nil {
      log.Fatal(err)
    }

    signedData, err := binary.MakeSignedData(cert, signingTime, sig)
    if err != nil {
      log.Fatal(err)
    }

    fmt.Printf("Signed Data len: %d\n", len(signedData))

    if err := binary.AppendSignature(signedData); err != nil {
      log.Fatal(err)
    }

    if err = os.WriteFile(*outputFile, binary.Bytes(), 0644); err != nil {
      log.Fatal(err)
    }
  }
}
