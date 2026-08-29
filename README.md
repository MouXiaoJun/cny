# cny

Convert RMB amounts (元) to **Chinese uppercase financial numerals** (中文大写), for invoices, contracts, and expense reports.

- Module: `github.com/MouXiaoJun/cny`
- Go 1.21+, zero dependencies (standard library only)
- License: [Mulan PSL v2](LICENSE)

```go
s, _ := cny.Convert(12345678.90)
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

## Precision

Inputs with more than 2 decimals are **rounded to cents using the decimal representation** (third decimal ≥ 5 rounds up; applied to the absolute value, sign preserved):

| Input | Output |
| --- | --- |
| `1.005` | 壹元零壹分 |
| `2.675` | 贰元陆角捌分 |
| `0.999` | 壹元整 |
| `-1.005` | 负壹元零壹分 |

Rounding operates on the decimal form of the float (`strconv.FormatFloat`, shortest round-trip), *not* on `amount*100`, so `1.005 → 壹元零壹分` as expected in accounting. Note `float64` only carries ~15–16 significant decimal digits: cent-exactness is only guaranteed for `|amount| < 2⁵³/100 ≈ 9×10¹³` 元. For exact money, use `ConvertCents`.

## Range

Integer part up to 20 digits (incl. rounding carry), i.e. up to `99999999999999999999.99`; beyond that `ErrOutOfRange` is returned. `ConvertCents` covers the full `int64` range.

## Naming

The core API is named `Convert` — plain "convert an amount to its uppercase form", matching `strconv` conventions. The exact cent path is `ConvertCents`.

## Development

```bash
gofmt -l . && go vet ./... && go test ./... && go test -race ./...
go test -run xxx -fuzz '^FuzzConvert$' -fuzztime 10s .
```

Tests cover whole/fraction/zero/negative/boundary/rounding cases, `Convert`/`ConvertCents` equivalence, and fuzzing (`FuzzConvert`, `FuzzConvertCents`).

## License

[Mulan PSL v2](LICENSE) © 2026 cny contributors
