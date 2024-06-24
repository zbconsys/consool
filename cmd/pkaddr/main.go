package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/zbconsys/consool/pkg/tools"
	"log"
	"strings"
)

func main() {
	keysList := flag.String("keys", "", "comma delimited list of public keys")
	flag.Parse()

	if *keysList == "" {
		log.Fatalln("no public keys specified")
	}

	buff := bytes.NewBufferString("[ADDRESS] \t [PUBLIC_KEY]\n")

	keys := strings.Split(*keysList, ",")
	for _, key := range keys {
		eoaAddr, err := tools.GetAddressFromPublicKey(strings.TrimSpace(key))
		if err != nil {
			log.Fatalln("invalid public key:", err)
		}

		buff.WriteString(fmt.Sprintf("%s \t %s\n", eoaAddr, strings.TrimSpace(key)))
	}

	fmt.Println(buff.String())
}
