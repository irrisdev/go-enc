package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/irrisdev/go-enc/genc"
)

const MinPassLen = 10

func main() {

	// encrypt flag
	encrypt := flag.Bool("encrypt", false, "use encryption")
	flag.BoolVar(encrypt, "e", false, "shorthand for -encrypt")

	// decrypt flag
	decrypt := flag.Bool("decrypt", false, "use decryption")
	flag.BoolVar(decrypt, "d", false, "shorthand for -decrypt")

	// outfile flag
	outfile := flag.String("outpath", "", "enter outfile path")
	flag.StringVar(outfile, "o", "", "shortand for outfile path -outpath")

	// passphrase flag
	passphrase := flag.String("passphrase", "", "enter encryption passphrase")
	flag.StringVar(passphrase, "pass", "", "shorthand for -passphrase")
	flag.StringVar(passphrase, "p", "", "shorthand for -passphrase")

	// path flag
	path := flag.String("file", "", "enter path to file")
	flag.StringVar(path, "f", "", "shorthand for -file")

	flag.Parse()

	// must choose exactly one mode
	if *encrypt == *decrypt {
		log.Fatal("must specify exactly one option of -encrypt or -decrypt")
	}

	// file path required
	if *path == "" {
		log.Fatal("file path required")
	}

	info, err := os.Stat(*path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("file does not exist")
		}
		log.Fatalf("failed to stat file: %v", err)
	}

	if info.IsDir() {
		log.Fatal("directories are not supported")
	}

	// passphrase required
	if *passphrase == "" {
		log.Fatal("passphrase required")
	}

	if len(*passphrase) < MinPassLen {
		log.Fatalf("minimum passphrase length: %d chars", MinPassLen)
	}

	if *encrypt {
		genc.Encrypt(*passphrase, *path)
	}

	if *decrypt {
		if *outfile == "" {
			genc.Decrypt(*passphrase, *path)

		} else {
			genc.Decrypt(*passphrase, *path, *outfile)
		}
	}

}
