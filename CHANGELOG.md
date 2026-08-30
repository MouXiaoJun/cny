# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2026-08-30

### Fixed / Changed

- Correct floating-point precision and rounding documentation; add amount-format regression examples.
- Keep ConvertCents as the exact integer-cent API.
- Align LICENSE and current documentation with MIT, as confirmed by the maintainer; preserve existing copyright notices.

## [1.0.0] - 2026-08-29

### Added
- 初始版本发布:人民币金额转中文大写(财务大写),适用于发票、合同、报销单等场景
- 核心 API `Convert(amount float64) (string, error)`:
  - 整数部分支持 0、个/十/百/千、万、亿、万亿、亿亿(整数部分最多 20 位)
  - 小数部分支持角、分;无角分时以"整"结尾;不足壹角如 0.05 → 伍分
  - 数字中的零按"只读一个零"规则处理,如 1001 → 壹仟零壹元整、105 → 壹佰零伍元整
  - 负数以"负"开头,如 -100 → 负壹佰元整
  - 精度超过 2 位时按十进制四舍五入到分(只看绝对值,符号保留)
  - 0 / 0.00 / -0.0 统一为"零元整";NaN/±Inf 返回 ErrNotFinite
  - 整数部分超过 20 位(≈ |金额| ≥ 10^20)返回 ErrOutOfRange
- 精确路径 `ConvertCents(cents int64) (string, error)`:以分为单位,无浮点误差,支持 int64 全范围
- 错误哨兵 `ErrNotFinite`、`ErrOutOfRange`,支持 errors.Is
- 单元测试覆盖整数/小数/零/负数/舍入/大数/边界,`Convert` 与 `ConvertCents` 小范围等价性校验
- 模糊测试 `FuzzConvert`、`FuzzConvertCents`(随机输入不 panic、输出结构不变量)
- 基准测试 `BenchmarkConvert`、`BenchmarkConvertCents`
- 零依赖,仅标准库,go.mod 声明 go 1.21
- CI:test(vet/build/race/coverage/fuzz)、bench(benchstat 对比)、release(tag 触发)
