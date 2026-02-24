# ResoDNS — توضیح مختصر برای Release

متن زیر را می‌توانی در GitHub Release (یا هر جایی که توضیح کوتاه می‌خواهی) استفاده کنی.

---

## English (برای GitHub Release)

**ResoDNS** is a CLI tool for fast mass DNS resolution and subdomain discovery. It resolves large lists of domains and bruteforces subdomains using wordlists, then filters out wildcard results using trusted resolvers (e.g. 8.8.8.8).

**Features:** resolve from file/stdin, subdomain bruteforce (single or multiple domains), wildcard detection, rate limiting, template expansion (`[1-100]`, `[file:wordlist.txt]`), JSON/output file options, scope filtering.

**Built on:** [massdns](https://github.com/blechschmidt/massdns).  
**Platform:** Linux (or Cygwin on Windows).  
**Requires:** Go to build; massdns binary at runtime.

---

## متن کوتاه‌تر (یک پاراگراف)

ResoDNS performs mass DNS resolution and subdomain bruteforce with wildcard filtering and trusted-resolver validation. Resolve domains from a file or bruteforce subdomains with a wordlist; results are validated against trusted DNS to drop wildcards. Supports rate limiting, templates, JSON/output files, and scope filtering. Built on massdns; runs on Linux (or Cygwin on Windows).

---

## فارسی (خلاصه)

ابزار خط فرمان برای حل انبوه DNS و کشف ساب‌دامین با وردلیست. خروجی با رزولورهای معتبر فیلتر می‌شود تا پاسخ‌های wildcard حذف شوند. قابلیت‌ها: resolve از فایل، bruteforce ساب‌دامین، تشخیص wildcard، محدودیت نرخ، قالب عددی/فایل، خروجی JSON و فایل، فیلتر scope. مبتنی بر massdns؛ پلتفرم: لینوکس (یا Cygwin در ویندوز).
