# cny

Maintenance scope: preserve the published API; focus on bug fixes, security and Go compatibility, with no planned API expansion.

Convert RMB amounts (元) to **Chinese uppercase financial numerals** (中文大写), for invoices, contracts, and expense reports.

This library formats amount text only. It does not add the 人民币 prefix, validate document templates or certify legal compliance.

- Module: `github.com/MouXiaoJun/cny`
- Go 1.21+, zero dependencies (standard library only)
- License: [MIT](LICENSE)

```go
s, _ := cny.ConvertCents(1234567890) // exact integer cents
fmt.Println(s) // 壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元玖角
```

## Install

```bash
go get github.com/MouXiaoJun/cny
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/MouXiaoJun/cny"
)

func main() {
	fmt.Println(cny.Convert(1001))    // 壹仟零壹元整
	fmt.Println(cny.Convert(105))     // 壹佰零伍元整
	fmt.Println(cny.Convert(100))     // 壹佰元整
	fmt.Println(cny.Convert(0.05))    // 伍分
	fmt.Println(cny.Convert(1.05))    // 壹元零伍分
	fmt.Println(cny.Convert(-100))    // 负壹佰元整
	fmt.Println(cny.ConvertCents(1234567890)) // 壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元玖角
}
```

## API

| Function | Description |
| --- | --- |
| `Convert(amount float64) (string, error)` | Convert a yuan amount. Rounds to cents (decimal half-up); see [Precision](#precision). |
| `ConvertCents(cents int64) (string, error)` | Convert an exact integer cent amount (1 元 = 100 分). No floating point involved; supports the full `int64` range and never errors. |
| `ErrNotFinite` | `amount` is `NaN` or `±Inf`. |
| `ErrOutOfRange` | Integer part exceeds 20 digits (≈ \|amount\| ≥ 10²⁰). |

Errors support `errors.Is`.

## Rules

Digits: 零 壹 贰 叁 肆 伍 陆 柒 捌 玖; units: 拾 佰 仟 万 亿 万亿 亿亿; currency: 元 角 分 整; sign: 负.

- **Whole part** — read in 4-digit groups from high to low: `10000 → 壹万元整`, `12345678 → 壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元整`.
- **Zeros** — one 零 between non-zero digits, never at the end: `100 → 壹佰元整`, `105 → 壹佰零伍元整`, `1001 → 壹仟零壹元整`, `10000001 → 壹仟万零壹元整`, `100010000 → 壹亿零壹万元整`.
- **Fraction** — 角/分 as present; `整` only when there is no fraction: `1.50 → 壹元伍角`, `1.05 → 壹元零伍分`, `0.05 → 伍分` (no 零元 prefix), `0.55 → 伍角伍分`.
- **Zero amount** — `0`, `0.00`, `-0.0` all give `零元整`.
- **Negative** — prefixed with 负: `-100 → 负壹佰元整`.

### Sourced amount examples

The [PBOC appendix](https://www.pbc.gov.cn/tiaofasi/144941/144957/3601571/2018081309061735002.pdf) (PDF pages 69–70) gives amount-writing examples. A readable cross-reference is the [Beijing transport authority's tender appendix](https://jtw.beijing.gov.cn/xxgk/ztbxx/202005/P020200528617926701667.pdf) (printed page 40). These three examples are regression fixtures, excluding the 人民币 prefix that the caller's template must handle:

| `ConvertCents` input | Amount text |
| --- | --- |
| `140950` | 壹仟肆佰零玖元伍角 |
| `600714` | 陆仟零柒元壹角肆分 |
| `1640902` | 壹万陆仟肆佰零玖元零贰分 |

Checked on 2026-08-30. These fixtures are not a compliance audit. Negative amounts, omission of 零元 below one yuan, and large-number units 万亿/亿亿 are library conventions, not claims that every document accepts them.

## Precision

Inputs with more than 2 decimals are **rounded to cents using the decimal representation** (third decimal ≥ 5 rounds up; applied to the absolute value, sign preserved):

| Input | Output |
| --- | --- |
| `1.005` | 壹元零壹分 |
| `2.675` | 贰元陆角捌分 |
| `0.999` | 壹元整 |
| `-1.005` | 负壹元零壹分 |

Rounding operates on the received float's shortest round-trip decimal form, *not* on `amount*100`. [strconv.FormatFloat](https://pkg.go.dev/strconv#FormatFloat) preserves that float, not the original decimal input or precision lost before the call. There is no cent-exact guarantee for all amounts below `2⁵³/100` yuan: around `2⁴⁶` yuan, adjacent float64 values are already 0.015625 yuan apart. For example, converting `int64(7036874417766401)` cents through `float64(cents)/100` produces an amount ending in 贰分, while `ConvertCents(cents)` ends in 壹分. Use `ConvertCents` directly for exact integer cents.

## Range

The formatting limit is 20 integer digits, including rounding carry; beyond that `ErrOutOfRange` is returned. This is not a guarantee that float64 can represent every amount in that range. `ConvertCents` covers the full `int64` cent range.

## Naming

The core API is named `Convert` — plain "convert an amount to its uppercase form", matching `strconv` conventions. The exact cent path is `ConvertCents`.

## Development

```bash
gofmt -l . && go vet ./... && go test ./... && go test -race ./...
go test -run xxx -fuzz '^FuzzConvert$' -fuzztime 10s .
```

Tests cover whole/fraction/zero/negative/boundary/rounding cases, `Convert`/`ConvertCents` equivalence, and fuzzing (`FuzzConvert`, `FuzzConvertCents`).

## License

[MIT](LICENSE) © 2026 cny contributors
