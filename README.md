# ResoDNS

Mass DNS resolution and subdomain bruteforce with wildcard filtering and trusted-resolver validation. Built on [massdns](https://github.com/blechschmidt/massdns).

This project started as a fork of [puredns](https://github.com/d3mondev/puredns) but has since diverged significantly in architecture and implementation.

**Platform:** Linux (or WSL / Cygwin / Docker on Windows).

---

## Main features (ویژگی‌های اصلی)

- **Resolve from file** — Resolve a list of domains from a `.txt` file (one domain per line).
- **Resolve via pipeline** — Pass a single domain or a list of domains through stdin:  
  `echo example.com | resodns resolve -` or `cat domains.txt | resodns resolve -`
- **Bruteforce with wordlist** — Subdomain bruteforce with a wordlist; wordlist can be from file or stdin (`cat wordlist.txt | resodns bruteforce - example.com`).
- **Bruteforce with template** — Numeric ranges and steps: `[1-100]`, `[1-100:10]`, and `[file:wordlist.txt]`; combine with wordlist. *(Template support in CLI is planned; the engine exists in `pkg/template`.)*
- **Wildcard control** — Number of wildcard tests (`-n, --wildcard-tests`), threads for wildcard filtering (`-t, --threads`), and batch size (`--wildcard-batch`).
- **Resolvers** — Public resolvers file (`-r, --resolvers`) and trusted resolvers file (`--resolvers-trusted`); control who resolves and who validates.

---

## Other features (سایر قابلیت‌ها)

| Feature | Description |
|--------|-------------|
| **Rate limiting** | `-l, --rate-limit` for public resolvers; `--rate-limit-trusted` for trusted (QPS). |
| **Output files** | `-w, --write` (found domains), `--write-wildcards` (wildcard roots), `--write-massdns` (raw massdns output). |
| **Skip steps** | `--skip-sanitize`, `--skip-wildcard-filter`, `--skip-validation` to speed up or debug. |
| **Multiple domains (bruteforce)** | `-d, --domains <file>` — bruteforce against a list of domains from a file. |
| **massdns path** | `-b, --bin` — path to massdns binary when not in PATH. |

*Note: Options like `--scope`, `--scope-file`, `--unique`, `--write-json`, and full CLI integration of `-T, --template` / `--template-max` are documented for completeness but may not yet be implemented in the current build; the template engine exists in code.*

---

## Windows and massdns

**massdns does not provide an official Windows build.** ResoDNS itself builds and runs on Windows (`go build -o resodns.exe .`), but it needs the massdns binary to perform DNS resolution. On native Windows you have two practical options:

1. **WSL (Windows Subsystem for Linux)** — Install a Linux distro in WSL, build or use massdns there, and run resodns from the same environment (e.g. from a Windows drive mount like `/mnt/c/...`). This is the recommended way.
2. **Cygwin** — Build massdns from source inside Cygwin and point resodns to it with `--bin`.
3. **Docker** — Run ResoDNS and massdns inside a Linux container; use volume mounts for wordlists and output.

So: **ResoDNS on Windows is usable only together with a Linux-like environment (WSL/Cygwin/Docker) where massdns can run.**

---

## What you need for full execution (نیازها برای اجرای کامل)

| Requirement | Description |
|-------------|-------------|
| **1. Go** | To build ResoDNS: `go build -o resodns .` (or `resodns.exe` on Windows). |
| **2. massdns** | A working **massdns** binary. ResoDNS calls it to do the actual DNS lookups. No official Windows build; on **WSL/Linux** either use the included `massdns/bin/massdns` (if it runs on your glibc) or **build from source:** `git clone https://github.com/blechschmidt/massdns && cd massdns && make` then point ResoDNS to it with `-b /path/to/massdns`. |
| **3. Public resolvers file** | A `.txt` file with one DNS resolver IP per line (e.g. `massdns/lists/resolvers.txt`). Pass with `-r resolvers.txt`. |
| **4. (Optional) Trusted resolvers** | For wildcard validation. Default 8.8.8.8, 8.8.4.4; or use `--resolvers-trusted trusted.txt`. |

**Summary:** Install Go → build ResoDNS → have a working massdns binary + a resolvers list. Then you can run `resodns resolve` and `resodns bruteforce` fully.

**If the included `massdns/bin/massdns` fails in WSL (e.g. GLIBC version), build massdns from source in WSL:**  
`bash scripts/build-massdns-wsl.sh` (requires `gcc`; install with `sudo apt install build-essential` if needed).

---

## Installation

**Dependency:** [massdns](https://github.com/blechschmidt/massdns) binary. This repo includes a Linux binary at `massdns/bin/massdns`; it may require a recent glibc (e.g. Ubuntu 24). Otherwise build massdns from source on your system (see table above). The public resolvers list is in `massdns/lists/resolvers.txt` (or use any text file with one resolver per line).

**Build resodns:** From the project root with Go installed:

- **Linux / WSL:** `make build` or `go build -o resodns .`
- **Windows (native):** `go build -o resodns.exe .` — the binary runs on Windows but still requires massdns from WSL/Cygwin/Docker (see above).

**Using massdns:** If you run from the project root and the binary is at `massdns/bin/massdns`, you usually don't need `--bin`. On Windows, use WSL/Cygwin and pass the path to the massdns binary with `-b` or `--bin` if needed. For a newer massdns, replace the binary or pass the new path with **`--bin`**.

---

## Commands

| Command | Description |
|--------|-------------|
| `resodns resolve <file>` | Resolve domains from a file (one per line). Stdin is supported. |
| `resodns bruteforce <wordlist> <domain>` | Bruteforce subdomains with a wordlist against one domain. |
| `resodns bruteforce <wordlist> -d <domains-file>` | Same with multiple domains from a file. |

More details: `resodns --help`, `resodns resolve --help`, `resodns bruteforce --help`.

---

## Command examples

### resolve

Resolve a list of domains from a file or stdin.

```bash
# From file
resodns resolve domains.txt

# From stdin (e.g. pipe)
cat domains.txt | resodns resolve -

# With custom resolvers and write results to file
resodns resolve domains.txt -r massdns/lists/resolvers.txt -w results.txt
```

### bruteforce

Bruteforce subdomains for one or more domains using a wordlist. The wordlist can be from a file or **stdin** (use `-` as the first argument when piping).

```bash
# One domain
resodns bruteforce wordlist.txt example.com

# Wordlist from stdin (pipe)
cat wordlist.txt | resodns bruteforce - example.com

# Multiple domains from file
resodns bruteforce wordlist.txt -d domains.txt

# Save found subdomains to file
resodns bruteforce wordlist.txt example.com -w found.txt

# With template: numeric range [1-100] and step [0-50:10]
resodns bruteforce wordlist.txt example.com -T "sub[1-5].example.com" --template-max 1000
```

---

## Options (switches)

All options can be used with `resolve` and `bruteforce` unless noted. Short and long forms are listed. **Each option below includes a full example.**

---

### `-r, --resolvers <file>`

**Description:** File containing public DNS resolvers (one IP per line). Used by massdns for bulk resolution.

**Example:**

```bash
resodns resolve domains.txt -r massdns/lists/resolvers.txt
resodns bruteforce wordlist.txt example.com --resolvers ./my-resolvers.txt
```

---

### `--resolvers-trusted <file>`

**Description:** File containing trusted DNS resolvers (e.g. 8.8.8.8, 8.8.4.4). Used for wildcard validation. Default: 8.8.8.8, 8.8.4.4 if not set.

**Example:**

```bash
resodns bruteforce wordlist.txt example.com --resolvers-trusted trusted.txt
resodns resolve domains.txt -r resolvers.txt --resolvers-trusted ~/.config/resodns/resolvers-trusted.txt
```

---

### `-b, --bin <path>`

**Description:** Path to the massdns binary. Use when massdns is not in PATH or you use a custom build (e.g. on Windows/Cygwin).

**Example:**

```bash
resodns resolve domains.txt -b /usr/local/bin/massdns
resodns bruteforce wordlist.txt example.com --bin ./massdns/bin/massdns
```

---

### `-w, --write <file>`

**Description:** Write found/valid domains to this file (one per line).

**Example:**

```bash
resodns resolve domains.txt -w resolved.txt
resodns bruteforce wordlist.txt example.com -w subdomains.txt
```

---

### `--write-wildcards <file>`

**Description:** Write detected wildcard (root) domains to this file.

**Example:**

```bash
resodns bruteforce wordlist.txt example.com --write-wildcards wildcards.txt
resodns resolve domains.txt -w out.txt --write-wildcards wildcard-roots.txt
```

---

### `--write-massdns <file>`

**Description:** Save raw massdns output to this file (for debugging or later processing).

**Example:**

```bash
resodns resolve domains.txt --write-massdns massdns-raw.txt
resodns bruteforce wordlist.txt example.com --write-massdns massdns.log
```

---

### `--write-json <file>`

**Description:** Write found domains as JSON to this file.

**Example:**

```bash
resodns resolve domains.txt --write-json results.json
resodns bruteforce wordlist.txt example.com -w out.txt --write-json out.json
```

---

### `--scope <domain-list>` / `--scope-file <file>`

**Description:** Only output domains that fall within the given scope. `--scope` accepts a comma-separated list of domains; `--scope-file` accepts a file with one domain per line.

**Example:**

```bash
resodns resolve domains.txt --scope example.com,test.example.com
resodns resolve domains.txt --scope-file allowed-domains.txt
resodns bruteforce wordlist.txt example.com -w out.txt --scope-file scope.txt
```

---

### `--unique`

**Description:** Deduplicate output so each domain is printed (or written) only once.

**Example:**

```bash
resodns resolve domains.txt --unique -w resolved.txt
resodns bruteforce wordlist.txt example.com --unique
```

---

### `-l, --rate-limit <qps>`

**Description:** Queries per second (QPS) limit for public resolvers (massdns). Helps avoid rate limits and abuse. `0` = no limit.

**Example:**

```bash
resodns resolve domains.txt -r resolvers.txt -l 100
resodns bruteforce wordlist.txt example.com --rate-limit 500
```

---

### `--rate-limit-trusted <qps>`

**Description:** QPS limit for trusted resolvers (used in wildcard detection).

**Example:**

```bash
resodns bruteforce wordlist.txt example.com --rate-limit-trusted 50
resodns resolve domains.txt --rate-limit-trusted 20
```

---

### `-t, --threads <n>`

**Description:** Number of threads (workers) for wildcard filtering. Affects speed of bruteforce when wildcard detection is used.

**Example:**

```bash
resodns bruteforce wordlist.txt example.com -t 20
resodns bruteforce wordlist.txt -d domains.txt --threads 10
```

---

### `-n, --wildcard-tests <n>`

**Description:** Number of random subdomain tests used for wildcard / load-balancing detection. More tests = more accurate but slower.

**Example:**

```bash
resodns bruteforce wordlist.txt example.com -n 5
resodns bruteforce wordlist.txt example.com --wildcard-tests 10
```

---

### `--wildcard-batch <n>`

**Description:** Number of subdomains per batch when running wildcard checks.

**Example:**

```bash
resodns bruteforce wordlist.txt example.com --wildcard-batch 100
resodns bruteforce wordlist.txt -d domains.txt --wildcard-batch 50
```

---

### `-T, --template <template-string>`

**Description:** Use a dynamic domain template instead of a static wordlist/domain. Placeholders: `[start-end]`, `[start-end:step]`, `[file:path]`. Multiple placeholders are expanded in order (Cartesian product). Used with **bruteforce**.

**Template syntax:**

| Placeholder | Meaning | Example |
|-------------|---------|--------|
| `[start-end]` | Numeric range | `[1-9]` → 1,2,…,9 |
| `[start-end:step]` | Numeric range with step | `[0-20:5]` → 0,5,10,15,20 |
| `[file:path]` | One value per line from file | `[file:words.txt]` |

**Examples:**

```bash
# Numeric subdomains 1 to 100
resodns bruteforce - -T "sub[1-100].example.com" example.com

# Step: 0, 10, 20, ..., 100
resodns bruteforce - -T "api[0-100:10].example.com" example.com

# Values from file (wordlist as part of template)
resodns bruteforce - -T "[file:wordlist.txt].example.com" example.com

# Cartesian: [1-2] × [file:a.txt] → sub1-x.example.com, sub1-y.example.com, sub2-x.example.com, sub2-y.example.com
resodns bruteforce - -T "sub[1-2]-[file:a.txt].example.com" example.com
```

---

### `--template-max <n>`

**Description:** Maximum number of domains to generate from template expansion. Use to cap size when template can produce very many combinations. `0` or negative = no limit.

**Example:**

```bash
resodns bruteforce - -T "[file:huge.txt].example.com" example.com --template-max 10000
resodns bruteforce - -T "sub[1-1000].example.com" example.com --template-max 5000
```

---

## Default resolver files

ResoDNS looks for resolver files in this order:

- `resolvers.txt` in the current directory
- `~/.config/resodns/resolvers.txt` (public resolvers)
- `~/.config/resodns/resolvers-trusted.txt` (trusted resolvers)

If not found, you must pass `-r` (and optionally `--resolvers-trusted`).

---

## Combined example

Full example using several options together:

```bash
resodns bruteforce massdns/lists/Subdomain.txt example.com \
  -r massdns/lists/resolvers.txt \
  --resolvers-trusted trusted.txt \
  -w found.txt \
  --write-wildcards wildcards.txt \
  -l 200 \
  -t 15
```

---

## License and disclaimer

Using this tool to attack targets without authorization is illegal. You are responsible for complying with applicable laws. This repository is licensed under the MIT License.
