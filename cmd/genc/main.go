package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/irrisdev/go-enc/genc"
)

const MinPassLen = 10

type Config struct {
	Encrypt      bool
	Decrypt      bool
	DeleteOrigin bool
	OutPath      string
	Passphrase   string
	File         string
}

func main() {

	args := ParseFlags()

	// must choose exactly one mode
	if args.Encrypt == args.Decrypt {
		log.Fatal("must specify exactly one option of -encrypt or -decrypt")
	}

	// file path required
	if args.File == "" {
		log.Fatal("file path required")
	}

	info, err := os.Stat(args.File)
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
	if args.Passphrase == "" {
		log.Fatal("passphrase required")
	}

	if len(args.Passphrase) < MinPassLen {
		log.Fatalf("minimum passphrase length: %d chars", MinPassLen)
	}

	if args.Encrypt {
		if args.DeleteOrigin {
			genc.Encrypt(args.Passphrase, args.File, args.DeleteOrigin)
		} else {
			genc.Encrypt(args.Passphrase, args.File)
		}
	}

	if args.Decrypt {
		if args.OutPath == "" {
			genc.Decrypt(args.Passphrase, args.File)

		} else {
			genc.Decrypt(args.Passphrase, args.File, args.OutPath)
		}
	}

}

func ParseFlags() *Config {
	cfg := &Config{}

	flag.BoolVar(&cfg.Encrypt, "encrypt", false, "use encryption")
	flag.BoolVar(&cfg.Encrypt, "e", false, "shorthand for -encrypt")

	flag.BoolVar(&cfg.Decrypt, "decrypt", false, "use decryption")
	flag.BoolVar(&cfg.Decrypt, "d", false, "shorthand for -decrypt")

	flag.BoolVar(&cfg.DeleteOrigin, "delete-origin", false, "remove origin file after encryption")
	flag.BoolVar(&cfg.DeleteOrigin, "do", false, "shorthand for -delete-origin")

	flag.StringVar(&cfg.OutPath, "outpath", "", "enter outfile path")
	flag.StringVar(&cfg.OutPath, "o", "", "shorthand for -outpath")

	flag.StringVar(&cfg.Passphrase, "passphrase", "", "enter encryption passphrase")
	flag.StringVar(&cfg.Passphrase, "p", "", "shorthand for -passphrase")
	flag.StringVar(&cfg.Passphrase, "pass", "", "shorthand for -passphrase")

	flag.StringVar(&cfg.File, "file", "", "enter path to file")
	flag.StringVar(&cfg.File, "f", "", "shorthand for -file")

	flag.Parse()
	return cfg
}
