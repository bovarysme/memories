package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bovarysme/memories/crypto"
)

var source, dest, ourID, theirID string

func init() {
	flag.StringVar(&source, "source", "", "path to the input file (e.g. chat-1067048330)")
	flag.StringVar(&dest, "dest", "", "path to the output file")
	flag.StringVar(&ourID, "oid", "", "your MID (e.g. u529a3d0285ef0aa49e713aeac1d2bafb)")
	flag.StringVar(&theirID, "tid", "", "your chat partner's MID")

	flag.Parse()
}

func main() {
	if source == "" || ourID == "" || theirID == "" {
		log.Fatal("Error: -source, -oid and -tid need to be set. See -help for more details.")
	}

	if dest == "" {
		dest = fmt.Sprintf("%s.sqlite", source)
	}

	log.Printf("Decrypting '%s' to '%s'\n", source, dest)
	err := crypto.Decrypt(source, dest, ourID, theirID)
	if err != nil {
		log.Fatal(err)
	}
}
