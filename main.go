package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/bovarysme/memories/attack"
)

var bruteforce bool
var cpuprofile string
var source, dest, ourID, theirID string

func init() {
	flag.BoolVar(&bruteforce, "bruteforce", false, "perform a key-recovery attack on the input file")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write a CPU profile")
	flag.StringVar(&source, "source", "", "path to the input file (e.g. chat-1067048330)")
	flag.StringVar(&dest, "dest", "", "path to the output file")
	flag.StringVar(&ourID, "oid", "", "your MID (e.g. u61726520762e206375746520f09f929c)")
	flag.StringVar(&theirID, "tid", "", "your chat partner's MID")

	flag.Parse()
}

func cmdBruteforce() error {
	if source == "" {
		return errors.New("Error: -source need to be set. See -help for more details.")
	}

	log.Printf("Performing a key-recovery attack on '%s'", source)
	err := attack.Bruteforce(source)

	return err
}

func cmdDecrypt() error {
	if source == "" || ourID == "" || theirID == "" {
		return errors.New("Error: -source, -oid and -tid need to be set. See -help for more details.")
	}

	if dest == "" {
		dest = fmt.Sprintf("%s.sqlite", source)
	}

	err := attack.Decrypt(source, dest, ourID, theirID)

	return err
}

func main() {
	if cpuprofile != "" {
		file, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		err = pprof.StartCPUProfile(file)
		if err != nil {
			log.Fatal(err)
		}

		defer pprof.StopCPUProfile()
	}

	var err error
	if bruteforce {
		err = cmdBruteforce()
	} else {
		err = cmdDecrypt()
	}

	if err != nil {
		log.Fatal(err)
	}
}
