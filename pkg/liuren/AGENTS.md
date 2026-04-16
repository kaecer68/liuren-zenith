# pkg/liuren — 大六壬核心引擎

**Lines**: ~2200 (5 files)  
**Complexity**: HIGH — Core divination algorithms  
**Entry**: `engine.go:Calculate()`

---

## WHERE TO LOOK

| Component | File | Lines | Purpose |
|-----------|------|-------|---------|
| Main engine | `engine.go` | ~1160 | `Engine.Calculate()` — full divination flow |
| Nine Methods | `jiuzongmen.go` | ~550 | `JiuZongMenCalculator` — trigram algorithms |
| Types | `types.go` | ~221 | Stems, Branches, Generals, Pan structures |
| Calculator | `calculator.go` | ~250 | Helper calculations (天地盤/四課/遁干) |
| Tests | `calculator_test.go` | TBD | Unit tests |

---

## CODE MAP

| Symbol | Type | Role |
|--------|------|------|
| `Engine.Calculate()` | Method | Main entry — 14-step divination flow |
| `JiuZongMenCalculator.CalculateSanChuan()` | Method | Nine Methods selector |
| `DivinationResult` | Struct | Full output (Four Courses, Three Transmissions, KeTi) |
| `FourKe` / `SanChuan` | Struct | Four Courses / Three Transmissions data |
| `HeavenlyGeneral` | Type (enum) | 12 Celestial Generals (貴人，螣蛇，朱雀...) |
| `Stem` / `Branch` | Type (enum) | 10 Stems, 12 Branches (0-based indexing) |
| `CalculateDunGan()` | Function | 遁干計算（旬首甲順佈十天干） |

---

## DIVINATION FLOW (Engine.Calculate)

```
1. Get calendar data (lunar-zenith) → stems/branches/solar terms
2. Parse pillars → Stem/Branch structs
3. Determine day/night → 貴人 selection
4. Calculate Month General → from solar term index (at mid-terms 中氣)
5. Calculate Heaven/Earth Pan → 月將加時
6. Find 貴人 position → place 12 Generals
7. Calculate Four Courses → 日干支之寄宮
8. Calculate Three Transmissions → Nine Methods (九宗門)
9. Calculate 空亡 → XunKong void branches (60-Jia cycle)
10. Calculate 神煞 → 驛馬，桃花，天馬
11. Assemble result → DivinationResult struct
12. Determine KeTi → 伏吟/返吟/九宗門 classification
13. Judge XiongJi → 吉凶 assessment
14. Calculate YongShen → 用神 by question type
```

---

## CORRECTED ALGORITHMS

### 月將計算 — `CalculateMonthGeneral()`
月將以**中氣**換將，非節氣。
- 中氣索引 1,3,5,...,23 分別對應月將 亥→子
- 節氣（偶數索引）沿用前一個中氣的月將
- 立春(0) 沿用大寒(23) 的月將 子

### 天地盤 — `GetDiPan()` / `CalculateTianPan()`
**地盤索引約定**（順時針，從亥開始）：
```
  巳(6) 午(7) 未(8) 申(9)
  辰(5)       酉(10)
  卯(4)       戌(11)
  寅(3) 丑(2) 子(1) 亥(0)
```
- `diPan[0] = 亥`, `diPan[1] = 子`, ..., `diPan[11] = 戌`

**天盤**：先找到占時在地盤中的索引位置，將月將置於該位置，其餘順時針按地支數值順序（0=子, 1=丑, ..., 11=亥）填充。

### 貴人順逆 — `CalculateTianJiang()`
順逆由**貴人所在地盤位置**決定，非晝夜：
- 陽位（亥子丑寅卯辰）→ 順行
- 陰位（巳午未申酉戌）→ 逆行

### 空亡 — `calculateXunKong()`
根據完整日干支計算旬首：
- 旬首地支 = (日支 - 日干 + 12) % 12
- 空亡 = (旬首 + 10) % 12, (旬首 + 11) % 12

### 遁干 — `CalculateDunGan()`
1. 計算日干支旬首地支
2. 找到旬首地支在天盤中的位置
3. 該位置放「甲」，順時針依次放乙、丙、丁...（循環十天干）

### 地支相剋 — `isKe()`
正統五行相剋表：
- 子水→午火
- 丑土→子水、亥水
- 寅木→丑土、辰土、未土、戌土
- 卯木→辰土、戌土、丑土、未土
- 辰土→子水、亥水
- 巳火→申金、酉金
- 午火→申金、酉金
- 未土→亥水
- 申金→寅木、卯木
- 酉金→卯木、寅木
- 戌土→亥水
- 亥水→巳火、午火

---

## NINE METHODS (九宗門) PRIORITY

`jiuzongmen.go:CalculateSanChuan()` — sequential fallback:

1. **賊克法** (ZeiKe) — Below/above clash (highest priority)
2. **比用法** (BiYong) — Yin/yang match with day stem
3. **涉害法** (SheHai) — Deepest harm count
4. **遙克法** (YaoKe) — Remote clash
5. **昴星法** (MaoXing) — 酉/卯 special case
6. **別責法** (BieZe) — 剛日取干合上神，柔日取支前三合上神
7. **八專法** (BaZhuan) — 八專日（甲寅、乙卯、丁未、己未、庚申、辛酉、戊午、癸丑），陽日取辰上神，陰日取戌上神
8. **伏吟法** (FuYin) — Heaven/Earth identical, 取刑神
9. **返吟法** (FanYin) — Heaven/Earth opposite (fallback)

### 三刑表（伏吟法用）
- 無恩之刑：寅刑巳，巳刑申，申刑寅
- 恃勢之刑：丑刑戌，戌刑未，未刑丑
- 無禮之刑：子刑卯，卯刑子
- 自刑：辰、午、酉、亥

---

## CONVENTIONS

### Indexing
- **Stems**: 0-based (甲=0, 乙=1, ..., 癸=9)
- **Branches**: 0-based (子=0, 丑=1, ..., 亥=11)
- **Generals**: 0-based (貴人=0, 螣蛇=1, ..., 天后=11)
- **Earth Plate**: 0=亥(bottom-right), 1=子, 2=丑, 3=寅, 4=卯, 5=辰, 6=巳, 7=午, 8=未, 9=申, 10=酉, 11=戌

### Naming
- **Exported types**: PascalCase (`DivinationResult`, `HeavenlyGeneral`)
- **Constants**: Mixed (English names with Chinese comments)
- **Error messages**: Traditional Chinese + English terms

### Time Handling
- Input: `time.Time` (local timezone: Asia/Taipei)
- Hour pillar: Calculated from hour (子=23-1, 丑=1-3, ...)
- Day/night: Determined from time (for 貴人 selection)

---

## ANTI-PATTERNS

| Pattern | Why Forbidden |
|---------|---------------|
| Hardcoding stems/branches | Use `Stem`/`Branch` enums, not raw ints |
| Skipping lunar-zenith | Must use `CalendarDataSource` interface |
| Direct array access | Use helper methods (`getStemElement`, `getBranchElement`) |
| Modifying types.go constants | Add new constants, don't change existing indices |
| Using solar terms instead of mid-terms for month general | Month general changes at 中氣, not 節氣 |
| Using day/night for tianjiang direction | Direction depends on noble's earth-plate position |

---

## KEY ALGORITHMS

### 十神 (ShiShen) — `engine.go:getShiShen()`
Combines stem element + yin/yang to determine relationship:
- 比肩/劫財 (same element)
- 食神/傷官 (I generate)
- 正財/偏財 (I overcome)
- 正官/七殺 (overcomes me)
- 正印/偏印 (generates me)

### 貴人 Finding — `FindGuiRen()`
Day stem + day/night → 貴人 position (0-11):
- 晝貴 (day): Yang noble
- 夜貴 (night): Yin noble

### 空亡 (XunKong) — `calculateXunKong()`
Six-Jia void branches — two branches per 10-day cycle, computed from full day pillar.

---

## TESTING

```bash
# Run engine tests
go test -v ./pkg/liuren/...

# Single file
go test -v ./pkg/liuren/calculator_test.go
```

---

## DEPENDENCIES

| Package | Purpose |
|---------|---------|
| `pkg/client` | `CalendarDataSource` interface (lunar-zenith client) |
| `time` | Time handling, timezone (Asia/Taipei) |

---

**Lines**: 2200+ | **Files**: 5 | **Last Updated**: 2026-04-16
