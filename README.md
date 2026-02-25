
# ResoDNS

**Mass DNS Resolution & Subdomain Bruteforce with Wildcard Filtering**

Built on **massdns** — optimized for large-scale enumeration.

Linux / WSL / Docker

---

## Overview

ResoDNS is a high‑performance DNS resolver and subdomain bruteforcer.

### Features

* Fast bulk domain resolution
* Subdomain bruteforce from wordlists
* Wildcard DNS detection & filtering
* Trusted resolver validation
* Template-based dynamic subdomain generation
* Full stdin / stdout piping support

---

## Installation

### Requirements

* Go
* massdns binary
* Resolvers file (one IP per line)

### Build

```bash
git clone https://github.com/l4tr0d3ctism/ResoDNS.git
cd ResoDNS
go build -o resodns .
```

### WSL Note

If massdns fails with GLIBC errors:

```bash
bash scripts/build-massdns-wsl.sh
```

---

## Usage

### Resolve

```bash
# From file
./resodns resolve domains.txt -r resolvers.txt

# From stdin
cat domains.txt | ./resodns resolve - -r resolvers.txt

# Save output
./resodns resolve domains.txt -r resolvers.txt -w found.txt
```

---

### Bruteforce

```bash
# Single domain
./resodns bruteforce wordlist.txt example.com -r resolvers.txt

# Multiple domains
./resodns bruteforce wordlist.txt -d domains.txt -r resolvers.txt

# Save output
./resodns bruteforce wordlist.txt example.com -r resolvers.txt -w found.txt
```

---

## Pipelining Examples

```bash
subfinder -d example.com -silent | ./resodns resolve - -r resolvers.txt

amass enum -passive -d example.com | ./resodns resolve - -r resolvers.txt

cat wordlist.txt | ./resodns bruteforce - example.com -r resolvers.txt
```

---

## Template Syntax

Generate subdomains dynamically.

| Pattern            | Description      |
| ------------------ | ---------------- |
| `[1-5]`            | Numeric range    |
| `[0-20:5]`         | Range with step  |
| `[file:words.txt]` | Values from file |

Examples:

```bash
./resodns bruteforce --template "api[1-5].example.com" -r resolvers.txt

./resodns bruteforce --template "[file:prefixes.txt].example.com" -r resolvers.txt

./resodns bruteforce --template "app[1-3]-[file:env.txt].example.com" -r resolvers.txt
```

Multiple placeholders expand as Cartesian product.

---

## Key Options

### Resolvers

* `-r, --resolvers` → Public resolvers file
* `--resolvers-trusted` → Trusted resolvers (default: 8.8.8.8, 8.8.4.4)
* `-b, --bin` → Path to massdns binary

### Performance

* `-l, --rate-limit`
* `--rate-limit-trusted`
* `-t, --threads`

### Wildcard

* `-n, --wildcard-tests`
* `--wildcard-batch`
* `--skip-wildcard-filter`

### Output

* `-w, --write`
* `--write-wildcards`
* `--write-massdns`

---

## Windows

massdns has no official Windows build.

Recommended: use WSL.

Alternative:

* Cygwin
* Docker

---

## License

MIT License

Use only on authorized targets.

