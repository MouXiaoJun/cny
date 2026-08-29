package cny

import "fmt"

func ExampleConvert() {
	s, err := Convert(12345678.90)
	fmt.Println(s, err)
	// Output: 壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元玖角 <nil>
}

func ExampleConvert_zero() {
	s, _ := Convert(0)
	fmt.Println(s)
	// Output: 零元整
}

func ExampleConvert_negative() {
	s, _ := Convert(-100)
	fmt.Println(s)
	// Output: 负壹佰元整
}

func ExampleConvertCents() {
	// 以分为单位,精确无浮点误差。
	s, err := ConvertCents(1234567890)
	fmt.Println(s, err)
	// Output: 壹仟贰佰叁拾肆万伍仟陆佰柒拾捌元玖角 <nil>
}
