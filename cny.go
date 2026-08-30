// Package cny 将人民币金额(元)转换为中文大写数字,适用于发票、合同、报销单等财务文书。
// 本包仅生成金额文本,不添加"人民币"前缀,不校验票据模板或法律合规性。
//
// 基本用法:
//
//	s, err := cny.Convert(12345678.90)      // 壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元玖角
//	s, err := cny.ConvertCents(1234567890)  // 同上,以分为单位,精确无浮点误差
//
// 规则要点(完整规则见 README_zh.md):
//
//   - 大写数字:零壹贰叁肆伍陆柒捌玖拾佰仟万亿角分整负
//   - 金额按十进制四舍五入到分;无角分时以"整"结尾
//   - 数字中的零按"只读一个零"处理,如 1001 → 壹仟零壹元整
//   - 负数以"负"开头;0 → 零元整
//   - NaN/±Inf 返回 ErrNotFinite;整数部分超过 20 位返回 ErrOutOfRange
package cny

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// 大写数字与位置单位。
var (
	// digitUpper 大写数字,下标 0..9。
	digitUpper = [...]string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
	// posUnit 每个四位组内 个/拾/佰/仟 的单位。
	posUnit = [...]string{"", "拾", "佰", "仟"}
	// groupUnit 四位一组的分组单位:个(空)、万、亿、万亿、亿亿。
	groupUnit = [...]string{"", "万", "亿", "万亿", "亿亿"}
)

// 错误哨兵。
var (
	// ErrNotFinite 表示金额不是有限数值(NaN 或 ±Inf),无法转换。
	ErrNotFinite = errors.New("cny: 金额必须是有限数值(不支持 NaN 或 ±Inf)")
	// ErrOutOfRange 表示金额绝对值超出支持范围(整数部分超过 20 位,即 |金额| ≥ 10^20)。
	ErrOutOfRange = errors.New("cny: 金额超出支持范围(整数部分最多 20 位)")
)

// maxIntDigits 是支持的整数部分最大位数(含四舍五入进位后),对应 9999 亿亿。
const maxIntDigits = 20

// Convert 将金额(元,float64)转换为人民币大写字符串。
//
// 精度:按十进制四舍五入到分——只看小数点后第三位,≥5 进一、<5 舍去;
// 舍入作用于绝对值,符号保持不变(如 -1.005 → 负壹元零壹分)。
// 与 "float64 乘以 100 后取整" 不同,本实现基于收到的浮点值的最短可往返十进制表示,
// 因此 1.005 → 壹元零壹分、2.675 → 贰元陆角捌分。
//
// 返回错误:
//   - ErrNotFinite:amount 为 NaN 或 ±Inf。
//   - ErrOutOfRange:四舍五入后整数部分超过 20 位(≈ |amount| ≥ 10^20)。
//
// 注意:最短可往返表示不能恢复传入前已丢失的金额信息。例如 2^46 元附近,
// 相邻 float64 的间隔已超过 1 分;不能以 2^53/100 元作为分精度保证。
// 需要精确到分的金额请直接使用 ConvertCents,不要先经过 float64 换算。
func Convert(amount float64) (string, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return "", ErrNotFinite
	}
	// 取收到的浮点值的最短可往返十进制表示(如 1.005 → "1.005"),不恢复原始输入精度。
	s := strconv.FormatFloat(amount, 'f', -1, 64)
	negative := strings.HasPrefix(s, "-")
	if negative {
		s = s[1:]
	}
	intPart, fracPart, _ := strings.Cut(s, ".")
	cents, carry := roundFracToCents(fracPart)
	if carry {
		intPart = incrementDecimal(intPart)
	}
	return convertParts(intPart, cents, negative)
}

// ConvertCents 将整数分值(1 元 = 100 分)转换为人民币大写字符串。
//
// 与 Convert 不同,分值为整数,不存在浮点误差,适合以分为单位存储的系统
// (数据库 BIGINT、第三方支付接口等)。支持 int64 全范围。
func ConvertCents(cents int64) (string, error) {
	s := strconv.FormatInt(cents, 10)
	negative := strings.HasPrefix(s, "-")
	if negative {
		s = s[1:]
	}
	var intDigits string
	var c int
	switch {
	case len(s) > 2:
		intDigits = s[:len(s)-2]
		c = int(s[len(s)-2]-'0')*10 + int(s[len(s)-1]-'0')
	case len(s) == 2:
		intDigits = "0"
		c = int(s[0]-'0')*10 + int(s[1]-'0')
	default:
		intDigits = "0"
		c = int(s[0] - '0')
	}
	return convertParts(intDigits, c, negative)
}

// roundFracToCents 将小数部分字符串四舍五入到分。
// 返回分值(0..99)与是否向整数部分进位。
func roundFracToCents(frac string) (cents int, carry bool) {
	if frac == "" {
		return 0, false
	}
	var d0, d1 int
	if len(frac) >= 1 {
		d0 = int(frac[0] - '0')
	}
	if len(frac) >= 2 {
		d1 = int(frac[1] - '0')
	}
	c := d0*10 + d1
	if len(frac) >= 3 && frac[2] >= '5' {
		c++
		if c == 100 {
			return 0, true
		}
	}
	return c, false
}

// incrementDecimal 将非负十进制数字符串加一(如 "0"→"1"、"99"→"100")。
func incrementDecimal(s string) string {
	digits := []byte(s)
	i := len(digits) - 1
	for i >= 0 && digits[i] == '9' {
		digits[i] = '0'
		i--
	}
	if i < 0 {
		return "1" + string(digits)
	}
	digits[i]++
	return string(digits)
}

// convertParts 由整数部分数字串(可带前导零)、分值(0..99)与符号组装最终大写字符串。
func convertParts(intDigits string, cents int, negative bool) (string, error) {
	intDigits = strings.TrimLeft(intDigits, "0")
	if intDigits == "" {
		intDigits = "0"
	}
	if len(intDigits) > maxIntDigits {
		return "", ErrOutOfRange
	}
	if intDigits == "0" && cents == 0 {
		// 0、0.00、-0.0 以及四舍五入后归零的金额统一为"零元整"。
		return "零元整", nil
	}
	var sb strings.Builder
	if negative {
		sb.WriteString("负")
	}
	hasYuan := false
	if intPart := intDigitsToUpper(intDigits); intPart != "" {
		sb.WriteString(intPart)
		sb.WriteString("元")
		hasYuan = true
	}
	if cents == 0 {
		sb.WriteString("整")
		return sb.String(), nil
	}
	jiao, fen := cents/10, cents%10
	if jiao > 0 {
		sb.WriteString(digitUpper[jiao])
		sb.WriteString("角")
	}
	if fen > 0 {
		// 仅在整数部分非零(已写"元")且角位为零时补"零",如 1.05 → 壹元零伍分;
		// 0.05 → 伍分,整数为零时不写"零元"也不补"零"。
		if hasYuan && jiao == 0 {
			sb.WriteString("零")
		}
		sb.WriteString(digitUpper[fen])
		sb.WriteString("分")
	}
	return sb.String(), nil
}

// intDigitsToUpper 将非负整数数字串(无符号、可有前导零)转为中文大写;
// 数值为零时返回空串,由调用方处理。
//
// 零的处理:非零数字之间出现连续零时只读一个"零";末尾的零不读;
// 每个四位组(万/亿/万亿/亿亿)若自身非零则补分组单位。
func intDigitsToUpper(s string) string {
	if s == "0" {
		return ""
	}
	n := len(s)
	var sb strings.Builder
	inZero := false  // 是否处于连续零区间
	started := false // 是否已输出过非零数字
	for i := 0; i < n; i++ {
		pos := n - 1 - i // 从个位(0)起的位号
		d := s[i] - '0'
		if d == 0 {
			if started {
				inZero = true
			}
		} else {
			if inZero {
				sb.WriteString("零")
				inZero = false
			}
			sb.WriteString(digitUpper[d])
			sb.WriteString(posUnit[pos%4])
			started = true
		}
		// 每处理完一个四位组(位号 4、8、12、16),若该组非零则补分组单位。
		if pos > 0 && pos%4 == 0 {
			if g := pos / 4; g < len(groupUnit) && groupNonZero(s, n, pos) {
				sb.WriteString(groupUnit[g])
			}
		}
	}
	return sb.String()
}

// groupNonZero 判断数字串 s 中位号 pos..pos+3(不足四位则到最高位为止)是否存在非零数字。
func groupNonZero(s string, n, pos int) bool {
	for p := pos; p < pos+4 && p < n; p++ {
		if s[n-1-p] != '0' {
			return true
		}
	}
	return false
}
