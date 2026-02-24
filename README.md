<p align="center">
  <strong>ResoDNS</strong>
</p>
<p align="center">
  <em>Mass DNS resolution & subdomain bruteforce with wildcard filtering</em>
</p>
<p align="center">
  <a href="https://github.com/blechschmidt/massdns">massdns</a> · 
  <a href="https://github.com/d3mondev/puredns">Fork of puredns</a> · 
  Linux / WSL / Cygwin / Docker
</p>

---

## Table of contents

- [Overview](#overview)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Commands](#commands)
- [Options](#options)
- [Template syntax](#template-syntax)
- [Windows & WSL](#windows--wsl)
- [License](#license)

---

## Overview

ResoDNS resolves large domain lists and bruteforces subdomains using wordlists, then filters wildcards and validates results with trusted resolvers. Built on [massdns](https://github.com/blechschmidt/massdns). This project started as a fork of [puredns](https://github.com/d3mondev/puredns) and has since diverged in architecture and implementation.

| Feature | Description |
|---------|-------------|
| **Resolve** | Domains from file or stdin (pipeline). |
| **Bruteforce** | Subdomains from wordlist (file or stdin) against one or many domains. |
| **Wildcard filtering** | Configurable tests, threads, and batch size. |
| **Resolvers** | Public resolvers file + optional trusted resolvers file. |
| **Output** | Write domains, wildcard roots, and raw massdns output to files. |

---

## Requirements

| # | Requirement | Notes |
|---|-------------|--------|
| 1 | **Go** | To build: `go build -o resodns .` |
| 2 | **massdns** | No official Windows build. On WSL/Linux use included binary or [build from source](https://github.com/blechschmidt/massdns). |
| 3 | **Resolvers file** | One resolver IP per line (e.g. `massdns/lists/resolvers.txt`). Use `-r`. |
| 4 | **Trusted resolvers** | Optional. Default 8.8.8.8, 8.8.4.4; or `--resolvers-trusted <file>`. |

> **WSL:** If `massdns/bin/massdns` fails (e.g. GLIBC), run `bash scripts/build-massdns-wsl.sh` (needs `gcc`; `sudo apt install build-essential` if required).

---

## Installation

**1. Clone and build**

```bash
git clone https://github.com/l4tr0d3ctism/ResoDNS.git
cd ResoDNS
go build -o resodns .
```

- **Linux / WSL:** `make build` or `go build -o resodns .`
- **Windows:** `go build -o resodns.exe .` (massdns still needs WSL/Cygwin/Docker)

**2. Resolvers**

Use `massdns/lists/resolvers.txt` or your own file. Pass with `-r <file>`.

**3. massdns path**

If the binary is not in PATH or you use a custom path: `-b ./massdns/bin/massdns` (or your path).

---

## Quick start

```bash
# Resolve domains from file
./resodns resolve domains.txt -r massdns/lists/resolvers.txt

# Resolve from pipeline
echo "example.com" | ./resodns resolve - -r massdns/lists/resolvers.txt

# Bruteforce one domain
./resodns bruteforce wordlist.txt example.com -r massdns/lists/resolvers.txt -w found.txt

# Bruteforce with wordlist from stdin
cat wordlist.txt | ./resodns bruteforce - example.com -r massdns/lists/resolvers.txt
```

---

## Commands

| Command | Description |
|---------|-------------|
| `resodns resolve <file>` | Resolve domains from file. Use `-` for stdin. |
| `resodns bruteforce <wordlist> <domain>` | Bruteforce subdomains for one domain. |
| `resodns bruteforce <wordlist> -d <domains-file>` | Bruteforce for multiple domains from file. |

Run `resodns --help`, `resodns resolve --help`, or `resodns bruteforce --help` for details.

---

## Options

Common options for both `resolve` and `bruteforce`. Short and long forms available.

| Option | Description |
|--------|-------------|
| `-r, --resolvers <file>` | Public DNS resolvers (one IP per line). |
| `--resolvers-trusted <file>` | Trusted resolvers for wildcard validation (default: 8.8.8.8, 8.8.4.4). |
| `-b, --bin <path>` | Path to massdns binary. |
| `-w, --write <file>` | Write found domains to file. |
| `--write-wildcards <file>` | Write wildcard roots to file. |
| `--write-massdns <file>` | Write raw massdns output. |
| `-l, --rate-limit <qps>` | QPS limit for public resolvers (0 = unlimited). |
| `--rate-limit-trusted <qps>` | QPS limit for trusted resolvers. |
| `-t, --threads <n>` | Threads for wildcard filtering. |
| `-n, --wildcard-tests <n>` | Number of tests for wildcard/load-balancing detection. |
| `--wildcard-batch <n>` | Subdomains per batch for wildcard checks (0 = unlimited). |
| `-d, --domains <file>` | (Bruteforce) File with domains to bruteforce. |
| `--skip-sanitize` | Do not sanitize domains. |
| `--skip-wildcard-filter` | Skip wildcard detection. |
| `--skip-validation` | Skip validation with trusted resolvers. |

**Examples**

```bash
# Resolve and save
./resodns resolve domains.txt -r resolvers.txt -w results.txt

# Bruteforce with rate limit and output files
./resodns bruteforce wordlist.txt example.com -r resolvers.txt -w found.txt \
  --write-wildcards wildcards.txt -l 200 -t 15
```

**Default resolver paths (if `-r` not given):**  
`resolvers.txt` (current dir) → `~/.config/resodns/resolvers.txt` → `~/.config/resodns/resolvers-trusted.txt`

---

## Template syntax

*(Template engine exists in `pkg/template`; full CLI integration may be pending.)*

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `[start-end]` | Numeric range | `[1-9]` → 1,2,…,9 |
| `[start-end:step]` | Range with step | `[0-20:5]` → 0,5,10,15,20 |
| `[file:path]` | One value per line from file | `[file:words.txt]` |

Multiple placeholders are expanded in order (Cartesian product).

---

## Windows & WSL

massdns has **no official Windows build**. ResoDNS builds and runs on Windows, but it needs massdns for DNS lookups.

| Option | Description |
|--------|-------------|
| **WSL** | Recommended. Use Linux in WSL; build or use massdns there. Project path in WSL: `/mnt/c/Users/<user>/.../puredns-master`. |
| **Cygwin** | Build massdns from source in Cygwin; point ResoDNS with `--bin`. |
| **Docker** | Run ResoDNS and massdns in a Linux container; mount wordlists and output. |

On native Windows, use one of the above so massdns is available.

---

## License

Use only on targets you are authorized to test. You are responsible for complying with applicable laws.

**MIT License.**
