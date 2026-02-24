# ResoDNS

Mass DNS resolution and subdomain bruteforce with wildcard filtering and trusted-resolver validation. Built on [massdns](https://github.com/blechschmidt/massdns).

This project started as a fork of [puredns](https://github.com/d3mondev/puredns) but has since diverged significantly in architecture and implementation.

**Platform:** Linux (or Cygwin on Windows).

---

## Installation

**Dependency:** [massdns](https://github.com/blechschmidt/massdns) binary. This repo includes a Linux binary at `massdns/bin/massdns`; use it on Linux. The public resolvers list is in `massdns/lists/resolvers.txt` (or use any text file with one resolver per line).

**Build resodns:** From the project root with Go installed:

- `make build` — outputs `resodns`
- or `go build -o resodns .`

**Using massdns:** If you run from the project root and the binary is at `massdns/bin/massdns`, you usually don't need `--bin`. On Windows (Cygwin) build massdns yourself and pass its path with `-b` or `--bin`. For a newer massdns, replace the binary or pass the new path with **`--bin`**.

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

# Quiet: only print resolved domains
resodns resolve domains.txt -q
```

### bruteforce

Bruteforce subdomains for one or more domains using a wordlist.

```bash
# One domain
resodns bruteforce wordlist.txt example.com

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
resodns bruteforce wordlist.txt example.com --unique -q
```

---

### `-q, --quiet`

**Description:** Only print domains (minimal output, no progress or extra messages).

**Example:**

```bash
resodns resolve domains.txt -q
resodns bruteforce wordlist.txt example.com --quiet -w found.txt
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
  -t 15 \
  -q
```

---

## License and disclaimer

Using this tool to attack targets without authorization is illegal. You are responsible for complying with applicable laws. This repository is licensed under the MIT License.
