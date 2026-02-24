# مقایسه: پروژه تو (ResoDNS) با puredns-master-F (puredns اصلی)

---

## چرا «بیلد/اجرا نمی‌شود» در مقابل «کامل است و بیلد و اجرا دارد»؟

### نقش پوشه `internal`

در این پروژه کد به دو لایه تقسیم شده است:

1. **`main.go` (ورود برنامه)**  
   فقط دو کار می‌کند: یک **context** می‌سازد و دستور **`cmd.Execute(ctx)`** را صدا می‌زند. یعنی تمام منطق CLI (دستورات، پرچم‌ها، اجرای resolve/bruteforce) در جای دیگری است.

2. **پوشه `internal/`**  
   همان «جای دیگر» است و سه بخش اصلی دارد:
   - **`internal/app/cmd`** — تعریف دستورات با Cobra: `resolve`, `bruteforce` (و در puredns اصلی `sponsors`)، پرچم‌ها (`--resolvers`, `-w`, و غیره)، و اتصال به usecaseها.
   - **`internal/app/ctx`** — نگه‌داشتن تنظیمات برنامه (نام، نسخه، آپشن‌ها) که در همهٔ دستورات استفاده می‌شود.
   - **`internal/usecase/resolve`** — منطق واقعی: خواندن دامنه/وردلیست، فراخوانی massdns، فیلتر wildcard، ذخیره نتیجه و غیره.

بدون این فایل‌ها، `main.go` به پکیج‌هایی مثل `resodns/internal/app/cmd` و `resodns/internal/app/ctx` ارجاع می‌دهد که در پروژهٔ تو **وجود ندارند**. کامپایلر Go آن پکیج‌ها را پیدا نمی‌کند و خطای «package not found» می‌دهد، بنابراین **بیلد اصلاً انجام نمی‌شود** و برنامه **اجرا نمی‌شود**.

### در پروژهٔ تو (ResoDNS)

- پوشه **`internal/`** در ریپو نیست (یا حذف شده).
- در نتیجه `main.go` فقط یک «اسکلت» است که به کدی وابسته است که وجود ندارد.
- **نتیجه:** `go build` خطا می‌دهد؛ باینری ساخته نمی‌شود؛ برنامه اجرا نمی‌شود.

### در puredns-master-F

- پوشه **`internal/`** کامل است: `app/cmd`, `app/ctx`, `usecase/resolve`, و غیره.
- `main.go` همان نقش را دارد ولی پکیج‌های `internal` موجودند.
- **نتیجه:** `go build` موفق است؛ باینری `puredns` ساخته می‌شود؛ برنامه اجرا می‌شود و دستورات `resolve` و `bruteforce` کار می‌کنند.

پس وقتی گفته می‌شود «بیلد/اجرا نمی‌شود — پوشه internal وجود ندارد» یعنی دقیقاً همین وابستگی و نبود کد CLI؛ و «کامل است و بیلد و اجرا دارد» یعنی همان پروژه با `internal` کامل است و از نظر بیلد و اجرا سالم است.

---

## ۱. شناسه و نام

| مورد | پروژه تو (puredns-master) | puredns-master-F |
|------|---------------------------|-------------------|
| **نام برنامه** | ResoDNS | puredns |
| **ماژول Go** | `resodns` | `github.com/d3mondev/puredns/v2` |
| **باینری** | `resodns` | `puredns` |

---

## ۲. ساختار کد و بیلد

| مورد | پروژه تو | puredns-master-F |
|------|----------|-------------------|
| **پوشه internal** | ❌ **وجود ندارد** — کد CLI (cmd, ctx) نیست | ✅ دارد: `internal/app`, `internal/usecase`, `internal/pkg` |
| **بیلد** | ❌ شکست می‌خورد (package resodns/internal/app/cmd not found) | ✅ بدون مشکل بیلد می‌شود |
| **دستورات** | همان ایده (resolve, bruteforce) ولی بدون کد اجرایی | resolve, bruteforce, sponsors با Cobra پیاده شده |

پروژه تو بدون کپی کردن (یا بازنویسی) بخش `internal` از F یا جای دیگر **اجرا نمی‌شود**.

---

## ۳. قابلیت‌های اضافه در پروژه تو

| قابلیت | پروژه تو | puredns-master-F |
|--------|----------|-------------------|
| **پکیج template** | ✅ دارد — `[1-100]`, `[1-100:10]`, `[file:path]` برای ساخت دامنه از قالب | ❌ ندارد |
| **پکیج fileoperation (appendword)** | ✅ دارد | در F فقط appendlines، copy، cat، readlines، writelines، countlines، fileexists — appendword به این شکل نیست |

در F فقط progressbar و چیزهای داخلی resolve از کلمه «template» استفاده می‌کنند، نه قالب دامنه مثل ResoDNS.

---

## ۴. ورودی از stdin (پایپ)

| دستور | پروژه تو | puredns-master-F |
|--------|----------|-------------------|
| **resolve** | در README: stdin پشتیبانی می‌شود (با `-`) | ✅ پشتیبانی می‌شود؛ در root مثال دارد: `cat domains.txt \| puredns resolve` |
| **bruteforce — وردلیست از stdin** | ✅ **پشتیبانی می‌شود**: `cat wordlist.txt \| resodns bruteforce - example.com` با `app.HasStdin()` | ✅ **پشتیبانی می‌شود** |
| **bruteforce — دامنه از stdin** | خیر | دامنه از آرگومان یا `-d domains.txt` است |

هر دو پروژه وردلیست را از stdin (پایپ) پشتیبانی می‌کنند.

---

## ۵. تست و CI/CD

| مورد | پروژه تو | puredns-master-F |
|------|----------|-------------------|
| **فایل‌های *_test.go** | ❌ حذف شده‌اند | ✅ همه پکیج‌ها تست دارند |
| **GitHub Actions** | ❌ ندارد | ✅ دارد: `build.yml`, `release.yml`, اکشن `setup-massdns` |
| **CHANGELOG / .gitignore / .vscode** | بخشی دارد | ✅ دارد (مثلاً CHANGELOG.md, .gitignore, .vscode/settings.json) |

---

## ۶. مستندات و README

| مورد | پروژه تو | puredns-master-F |
|------|----------|-------------------|
| **README** | ResoDNS، لینک فورک از puredns، جدول سوئیچ‌ها و مثال‌های کامل | puredns، لوگو، badges، Getting Started، FAQ، Sponsorship |
| **توضیح سوئیچ‌ها** | ✅ هر سوئیچ با مثال کامل | خلاصه‌تر |
| **سینتکس template** | ✅ مستند شده | وجود ندارد |

---

## ۷. پکیج‌های مشترک (pkg)

هر دو این پکیج‌ها را دارند (با تفاوت‌های جزئی در فایل‌ها):

- `massdns` — اجرای massdns، LineReader، callback، stdouthandler، resolver  
- `wildcarder` — تشخیص wildcard، gather، detectiontask، dnscache، answercache، clientdns  
- `threadpool`  
- `progressbar`  
- `fileoperation` (در تو appendword هم هست)  
- `filetest`  
- `procreader`  
- `shellexecutor`  

فقط پروژه تو **`pkg/template`** دارد.

---

## جمع‌بندی

| جنبه | پروژه تو | puredns-master-F |
|------|----------|-------------------|
| **اجراشدنی بودن** | ✅ با internal بیلد و اجرا می‌شود | ✅ کامل و قابل اجرا |
| **نام و برند** | ResoDNS، فورک puredns | puredns اصلی |
| **قالب دامنه (template)** | ✅ دارد | ❌ ندارد |
| **ورود وردلیست از stdin در bruteforce** | ✅ پشتیبانی می‌شود | ✅ دارد |
| **تست و CI** | تست حذف شده، CI نیست | تست و GitHub Actions دارد |
| **مستندات سوئیچ‌ها** | قوی‌تر و با مثال | معمولی |

پروژه تو اکنون با پوشه **internal** اضافه‌شده بیلد و اجرا می‌شود و وردلیست از stdin در bruteforce پشتیبانی می‌شود.
