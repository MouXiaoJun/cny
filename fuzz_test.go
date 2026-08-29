package cny

import (
	"errors"
	"math"
	"strings"
	"testing"
)

// FuzzConvert 任意 float64 金额不 panic、不产生非法输出。
// 合法输入(有限且整数部分 ≤ 20 位)必须返回非空字符串;
// 非法输入(NaN/±Inf、超出范围)必须返回已声明的错误哨兵。
func FuzzConvert(f *testing.F) {
	seeds := []float64{
		0, 1, 100, 105, 1001, 10000, 100000, 10000000, 100000000,
		12345678.9, 123456789012, -100, -1.05, 0.05, 0.5, 0.55,
		1.05, 1.005, 2.675, -0.005, 0.999, math.Copysign(0, -1),
		math.NaN(), math.Inf(1), math.Inf(-1), 1e21, math.MaxFloat64,
	}
	for _, v := range seeds {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, v float64) {
		s, err := Convert(v)
		if err != nil {
			if !errors.Is(err, ErrNotFinite) && !errors.Is(err, ErrOutOfRange) {
				t.Fatalf("Convert(%v): unexpected error %v", v, err)
			}
			return
		}
		if s == "" {
			t.Fatalf("Convert(%v): empty result", v)
		}
		for _, r := range s {
			if !strings.ContainsRune("零壹贰叁肆伍陆柒捌玖拾佰仟万亿角分元整负", r) {
				t.Fatalf("Convert(%v) = %q: illegal rune %q", v, s, r)
			}
		}
		if strings.Contains(s, "零零") {
			t.Fatalf("Convert(%v) = %q: consecutive 零", v, s)
		}
	})
}

// FuzzConvertCents 任意 int64 分值不 panic,并验证结构性不变量。
func FuzzConvertCents(f *testing.F) {
	seeds := []int64{0, 1, 5, 50, 55, 100, 105, 1234567890, -1, -100,
		math.MaxInt64, math.MinInt64, 1000000000000}
	for _, v := range seeds {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, cents int64) {
		s, err := ConvertCents(cents)
		if err != nil {
			t.Fatalf("ConvertCents(%d): unexpected error %v", cents, err)
		}
		if s == "" {
			t.Fatalf("ConvertCents(%d): empty result", cents)
		}
		if cents < 0 && !strings.HasPrefix(s, "负") {
			t.Fatalf("ConvertCents(%d) = %q: negative amount missing 负", cents, s)
		}
		if cents >= 0 && strings.HasPrefix(s, "负") {
			t.Fatalf("ConvertCents(%d) = %q: positive amount has 负", cents, s)
		}
		// 结尾结构:分位非零 → 以"分"结尾;分位为零 → 以"角"或"整"结尾。
		switch cents % 100 {
		case 0:
			if !strings.HasSuffix(s, "整") {
				t.Fatalf("ConvertCents(%d) = %q: want suffix 整", cents, s)
			}
		default:
			if cents%10 != 0 {
				if !strings.HasSuffix(s, "分") {
					t.Fatalf("ConvertCents(%d) = %q: want suffix 分", cents, s)
				}
			} else if !strings.HasSuffix(s, "角") {
				t.Fatalf("ConvertCents(%d) = %q: want suffix 角", cents, s)
			}
		}
		if strings.Contains(s, "零零") {
			t.Fatalf("ConvertCents(%d) = %q: consecutive 零", cents, s)
		}
	})
}
