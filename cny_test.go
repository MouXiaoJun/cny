package cny

import (
	"errors"
	"math"
	"strings"
	"testing"
)

// TestConvert 覆盖:整数(0/个/十百千/万/亿)、小数(角/分/整)、
// 零的处理(1001/105/100 等)、负数、四舍五入、大数。
func TestConvert(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		want   string
	}{
		// ---- 整数部分:0、个位、十百千 ----
		{"zero", 0, "零元整"},
		{"zero point zero zero", 0.00, "零元整"},
		{"negative zero", math.Copysign(0, -1), "零元整"},
		{"one", 1, "壹元整"},
		{"five", 5, "伍元整"},
		{"ten", 10, "壹拾元整"},
		{"twelve", 12, "壹拾贰元整"},
		{"twenty", 20, "贰拾元整"},
		{"hundred", 100, "壹佰元整"},
		{"hundred five", 105, "壹佰零伍元整"},
		{"thousand one", 1001, "壹仟零壹元整"},
		{"thousand ten", 1010, "壹仟零壹拾元整"},
		{"hundred ten", 110, "壹佰壹拾元整"},
		{"two thousand", 2000, "贰仟元整"},
		{"nine thousand nine", 9009, "玖仟零玖元整"},

		// ---- 整数部分:万 ----
		{"ten thousand", 10000, "壹万元整"},
		{"ten thousand one", 10001, "壹万零壹元整"},
		{"hundred thousand", 100000, "壹拾万元整"},
		{"million", 1000000, "壹佰万元整"},
		{"ten million", 10000000, "壹仟万元整"},
		{"million one", 1000100, "壹佰万零壹佰元整"},
		{"hundred thousand one hundred", 100100, "壹拾万零壹佰元整"},
		{"ten thousand one", 100001, "壹拾万零壹元整"},

		// ---- 整数部分:亿 ----
		{"hundred million", 100000000, "壹亿元整"},
		{"billion", 1000000000, "壹拾亿元整"},
		{"hundred million one", 100000001, "壹亿零壹元整"},
		{"hundred million ten thousand", 100010000, "壹亿零壹万元整"},
		{"thousand ten thousand", 10010000, "壹仟零壹万元整"},
		{"ten million one", 10000001, "壹仟万零壹元整"},

		// ---- 任务示例 ----
		{"task example", 12345678.90, "壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元玖角"},
		{"task example no cents", 12345678, "壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元整"},

		// ---- 小数:角、分 ----
		{"jiao only", 0.5, "伍角"},
		{"jiao fen", 0.55, "伍角伍分"},
		{"fen only", 0.05, "伍分"},
		{"fen only one", 0.01, "壹分"},
		{"jiao fen one", 0.15, "壹角伍分"},
		{"one yuan jiao", 1.5, "壹元伍角"},
		{"two yuan jiao", 2.2, "贰元贰角"},
		{"one yuan fen", 1.05, "壹元零伍分"},
		{"big jiao fen", 0.99, "玖角玖分"},

		// ---- 负数 ----
		{"negative integer", -100, "负壹佰元整"},
		{"negative with fen", -1.05, "负壹元零伍分"},
		{"negative jiao", -0.5, "负伍角"},
		{"negative fen", -0.05, "负伍分"},
		{"negative ten thousand", -12345.67, "负壹万贰仟叁佰肆拾伍元陆角柒分"},

		// ---- 精度超过 2 位:十进制四舍五入到分 ----
		{"round half up", 1.005, "壹元零壹分"},
		{"round half up big", 2.675, "贰元陆角捌分"},
		{"round down", 1.004, "壹元整"},
		{"round up fen", 0.005, "壹分"},
		{"round down zero", 0.004, "零元整"},
		{"round carry to yuan", 0.999, "壹元整"},
		{"round carry big", 9.995, "壹拾元整"},
		{"round many digits", 1.23456789, "壹元贰角叁分"},
		{"round negative half up", -1.005, "负壹元零壹分"},
		{"round negative down", -1.004, "负壹元整"},
		{"round negative fen", -0.005, "负壹分"},
		{"round 100.005", 100.005, "壹佰元零壹分"},

		// ---- 大数(亿级以上) ----
		{"123456789012", 123456789012, "壹仟贰佰叁拾肆亿伍仟陆佰柒拾捌万玖仟零壹拾贰元整"},
		{"trillion with zero group", 123400000001, "壹仟贰佰叁拾肆亿零壹元整"},
		{"trillion zero group middle", 123400001000, "壹仟贰佰叁拾肆亿零壹仟元整"},
		{"trillion zero group wan", 123400010000, "壹仟贰佰叁拾肆亿零壹万元整"},
		{"wan group zero between", 123400000000, "壹仟贰佰叁拾肆亿元整"},

		// ---- 浮点表示相关的既定行为 ----
		{"0.1 plus 0.2", 0.1 + 0.2, "叁角"},
		{"float jiao fen", 0.29, "贰角玖分"},
		{"float one point two", 1.2, "壹元贰角"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.amount)
			if err != nil {
				t.Fatalf("Convert(%v) unexpected error: %v", tt.amount, err)
			}
			if got != tt.want {
				t.Fatalf("Convert(%v) = %q, want %q", tt.amount, got, tt.want)
			}
		})
	}
}

// TestConvertCents 以整数分(1 元 = 100 分)为单位,精确无浮点误差。
func TestConvertCents(t *testing.T) {
	tests := []struct {
		name  string
		cents int64
		want  string
	}{
		{"zero", 0, "零元整"},
		{"one fen", 1, "壹分"},
		{"five fen", 5, "伍分"},
		{"jiao", 50, "伍角"},
		{"jiao fen", 55, "伍角伍分"},
		{"one yuan", 100, "壹元整"},
		{"one yuan five", 105, "壹元零伍分"},
		{"one hundred", 10000, "壹佰元整"},
		{"task example", 1234567890, "壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元玖角"},
		{"negative", -100, "负壹元整"},
		{"negative fen", -1, "负壹分"},
		{"big", 123456789012, "壹拾贰亿叁仟肆佰伍拾陆万柒仟捌佰玖拾元壹角贰分"},
		{"max int64", math.MaxInt64, "玖亿亿贰仟贰佰叁拾叁万亿柒仟贰佰零叁亿陆仟捌佰伍拾肆万柒仟柒佰伍拾捌元零柒分"},
		{"min int64", math.MinInt64, "负玖亿亿贰仟贰佰叁拾叁万亿柒仟贰佰零叁亿陆仟捌佰伍拾肆万柒仟柒佰伍拾捌元零捌分"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertCents(tt.cents)
			if err != nil {
				t.Fatalf("ConvertCents(%d) unexpected error: %v", tt.cents, err)
			}
			if got != tt.want {
				t.Fatalf("ConvertCents(%d) = %q, want %q", tt.cents, got, tt.want)
			}
		})
	}
}

// TestConvertErrors 错误路径:NaN/±Inf 与超出支持范围。
func TestConvertErrors(t *testing.T) {
	for _, amount := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := Convert(amount); !errors.Is(err, ErrNotFinite) {
			t.Fatalf("Convert(%v) error = %v, want ErrNotFinite", amount, err)
		}
	}

	for _, amount := range []float64{1e20, -1e20, 1e21, -1e21, 1e30, math.MaxFloat64} {
		if _, err := Convert(amount); !errors.Is(err, ErrOutOfRange) {
			t.Fatalf("Convert(%v) error = %v, want ErrOutOfRange", amount, err)
		}
	}

	// ConvertCents 永远不报错(int64 全范围都支持)。
	for _, cents := range []int64{math.MaxInt64, math.MinInt64, 0} {
		if _, err := ConvertCents(cents); err != nil {
			t.Fatalf("ConvertCents(%d) unexpected error: %v", cents, err)
		}
	}
}

// TestRoundFracToCents 直接单测四舍五入辅助函数。
func TestRoundFracToCents(t *testing.T) {
	tests := []struct {
		frac      string
		wantCents int
		wantCarry bool
	}{
		{"", 0, false},
		{"0", 0, false},
		{"9", 90, false},
		{"5", 50, false},
		{"05", 5, false},
		{"99", 99, false},
		{"10", 10, false},
		{"994", 99, false},
		{"995", 0, true},
		{"999", 0, true},
		{"005", 1, false},
		{"004", 0, false},
		{"9999999999", 0, true},
		{"30000000000000004", 30, false},
	}
	for _, tt := range tests {
		cents, carry := roundFracToCents(tt.frac)
		if cents != tt.wantCents || carry != tt.wantCarry {
			t.Fatalf("roundFracToCents(%q) = (%d, %v), want (%d, %v)",
				tt.frac, cents, carry, tt.wantCents, tt.wantCarry)
		}
	}
}

// TestIncrementDecimal 直接单测十进制加一辅助函数。
func TestIncrementDecimal(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"0", "1"},
		{"9", "10"},
		{"99", "100"},
		{"123", "124"},
		{"129", "130"},
		{"199", "200"},
	}
	for _, tt := range tests {
		if got := incrementDecimal(tt.in); got != tt.want {
			t.Fatalf("incrementDecimal(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// TestIntDigitsToUpper 直接单测整数部分转换(万/亿分组与零规则)。
func TestIntDigitsToUpper(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"0", ""},
		{"1", "壹"},
		{"10", "壹拾"},
		{"20", "贰拾"},
		{"100", "壹佰"},
		{"1001", "壹仟零壹"},
		{"10001", "壹万零壹"},
		{"10000001", "壹仟万零壹"},
		{"100010000", "壹亿零壹万"},
		{"123400000001", "壹仟贰佰叁拾肆亿零壹"},
		{"9999", "玖仟玖佰玖拾玖"},
		{"99999999", "玖仟玖佰玖拾玖万玖仟玖佰玖拾玖"},
	}
	for _, tt := range tests {
		if got := intDigitsToUpper(tt.in); got != tt.want {
			t.Fatalf("intDigitsToUpper(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// TestConvertFloatCentsAgree Convert(float64(cents)/100) 与 ConvertCents(cents) 结果一致。
// 在小范围内验证 float64 路径与整数分路径的等价性。
func TestConvertFloatCentsAgree(t *testing.T) {
	for c := int64(-100000); c <= 100000; c++ {
		want, err := ConvertCents(c)
		if err != nil {
			t.Fatalf("ConvertCents(%d) error: %v", c, err)
		}
		got, err := Convert(float64(c) / 100)
		if err != nil {
			t.Fatalf("Convert(%d/100) error: %v", c, err)
		}
		if got != want {
			t.Fatalf("Convert(%d/100) = %q, ConvertCents(%d) = %q", c, got, c, want)
		}
	}
}

// TestConvertInvariants 输出结构不变量:只含合法字符、无连续"零"、"整"只出现在结尾。
func TestConvertInvariants(t *testing.T) {
	allowed := "零壹贰叁肆伍陆柒捌玖拾佰仟万亿角分元整负"
	check := func(s string, amount float64) {
		t.Helper()
		if s == "" {
			t.Fatalf("Convert(%v) returned empty string", amount)
		}
		for _, r := range s {
			if !strings.ContainsRune(allowed, r) {
				t.Fatalf("Convert(%v) = %q contains illegal rune %q", amount, s, r)
			}
		}
		if strings.Contains(s, "零零") {
			t.Fatalf("Convert(%v) = %q contains consecutive 零", amount, s)
		}
		if strings.Count(s, "整") > 1 {
			t.Fatalf("Convert(%v) = %q contains more than one 整", amount, s)
		}
		if strings.Contains(s, "整") && !strings.HasSuffix(s, "整") {
			t.Fatalf("Convert(%v) = %q has 整 not at end", amount, s)
		}
		if strings.Count(s, "负") > 1 {
			t.Fatalf("Convert(%v) = %q contains more than one 负", amount, s)
		}
		if strings.Contains(s, "负") && !strings.HasPrefix(s, "负") {
			t.Fatalf("Convert(%v) = %q has 负 not at start", amount, s)
		}
	}

	// 整数扫描。
	for i := -200000; i <= 200000; i++ {
		s, err := Convert(float64(i))
		if err != nil {
			t.Fatalf("Convert(%d) error: %v", i, err)
		}
		check(s, float64(i))
	}
	// 分数扫描。
	for i := 0; i <= 10000; i++ {
		s, err := Convert(float64(i) / 100)
		if err != nil {
			t.Fatalf("Convert(%d/100) error: %v", i, err)
		}
		check(s, float64(i)/100)
	}
	// 大数抽查。
	for _, v := range []float64{123456789012, 999999999999, 987654321098765, 1e17, 1e18, -1e17} {
		s, err := Convert(v)
		if err != nil {
			t.Fatalf("Convert(%v) error: %v", v, err)
		}
		check(s, v)
	}
}
