// Package template provides dynamic domain template expansion for bruteforce-style
// subdomain generation. Templates support numeric ranges and wordlist placeholders.
//
// Syntax:
//   - [start-end]       numeric range, e.g. [1-9] or [1-1000]
//   - [start-end:step]  numeric range with step, e.g. [0-100:10]
//   - [file:path]       one value per line from file, e.g. [file:words.txt]
//
// All placeholders are expanded in order; multiple placeholders produce a Cartesian
// product. Example: ex[1-2]change-[file:a.txt].example.com with a.txt containing
// "x" and "y" yields ex1change-x.example.com, ex1change-y.example.com,
// ex2change-x.example.com, ex2change-y.example.com.
package template
