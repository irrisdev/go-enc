# go-enc

A simple CLI tool for encrypting and decrypting files using AES-GCM encryption.

## Installation

```bash
go install github.com/irrisdev/go-enc@latest
```

## Usage

### Encrypt a file

```bash
go-enc encrypt -f myfile.txt -p "your-strong-passphrase"
```

This creates `myfile.txt.genc` (encrypted file).

### Decrypt a file

```bash
go-enc decrypt -f myfile.txt.genc -p "your-strong-passphrase"
```

This restores the original `myfile.txt`.

### Options

- `--delete-origin` - Remove original file after encryption
- `-o, --outpath` - Specify custom output path for decryption

```bash
# Encrypt and delete original
go-enc encrypt -f secret.pdf -p "passphrase" --delete-origin

# Decrypt to specific location
go-enc decrypt -f secret.pdf.genc -p "passphrase" -o /path/to/output.pdf
```

## Requirements

- Passphrase must be at least 10 characters
- Uses AES-GCM with Argon2id key derivation

## License

[GPL-3.0](LICENSE)

