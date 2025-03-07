package datex

import (
	"fmt"
	"testing"
)

func TestGuess(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		wantErr  bool
		expected string // 预期格式化后的时间
	}{
		{"标准日期", "2024-01-01", false, "2024-01-01 00:00:00"},
		{"标准日期时间", "2024-01-01 15:04:05", false, "2024-01-01 15:04:05"},
		{"ISO8601带时区", "2024-01-01T15:04:05Z", false, "2024-01-01 15:04:05"},
		{"ISO8601无时区", "2024-01-01T15:04:05", false, "2024-01-01 15:04:05"},
		{"斜杠分隔日期时间", "2024/01/02 15:04:05", false, "2024-01-02 15:04:05"},
		{"无秒格式", "2024-01-02 15:04", false, "2024-01-02 15:04:00"},
		{"斜杠分隔无秒", "2024/01/02 15:04", false, "2024-01-02 15:04:00"},
		{"斜杠分隔仅日期", "2024/01/02", false, "2024-01-02 00:00:00"},
		{"美式日期", "01/02/2024", false, "2024-01-02 00:00:00"},
		// {"欧式日期", "02/01/2024", false, "2024-01-02 00:00:00"},
		{"紧凑格式", "20240102150405", false, "2024-01-02 15:04:05"},
		{"紧凑日期", "20240102", false, "2024-01-02 00:00:00"},
		{"RFC3339", "2024-01-02T15:04:05+08:00", false, "2024-01-02 15:04:05"},
		{"错误格式", "2024-99-99", true, ""},
		{"空字符串", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Guess(tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Guess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				gotFormatted := got.Format("2006-01-02 15:04:05")
				if gotFormatted != tt.expected {
					t.Errorf("Guess() = %v, want %v", gotFormatted, tt.expected)
				}
			}
		})
	}
}

func TestMove(t *testing.T) {
	tests := []struct {
		name       string
		dateStr    string
		movements  []any
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "加一周",
			dateStr:    "2024-01-01",
			movements:  []any{1 * Week},
			wantOutput: "2024-01-08",
			wantErr:    false,
		},
		{
			name:       "加一周后取本周第一天",
			dateStr:    "2024-01-01",
			movements:  []any{1 * Week, FirstDayOfWeek},
			wantOutput: "2024-01-07",
			wantErr:    false,
		},
		{
			name:       "加一周后取中国周第一天",
			dateStr:    "2025-03-07",
			movements:  []any{1 * Week, FirstDayOfCNWeek},
			wantOutput: "2025-03-10",
			wantErr:    false,
		},
		{
			name:       "加一年",
			dateStr:    "2024-01-01",
			movements:  []any{1 * Year},
			wantOutput: "2025-01-01",
			wantErr:    false,
		},
		{
			name:       "减一年",
			dateStr:    "2024-01-01",
			movements:  []any{-1 * Year},
			wantOutput: "2023-01-01",
			wantErr:    false,
		},
		{
			name:       "加一个月",
			dateStr:    "2024-01-31",
			movements:  []any{1 * Month},
			wantOutput: "2024-02-29", // 2024是闰年，1月31日加一个月是2月29日
			wantErr:    false,
		},
		{
			name:       "减一个月",
			dateStr:    "2024-03-31",
			movements:  []any{-1 * Month},
			wantOutput: "2024-02-29", // 2024是闰年，3月31日减一个月是2月29日
			wantErr:    false,
		},
		{
			name:       "加十天",
			dateStr:    "2024-02-20",
			movements:  []any{10 * Day},
			wantOutput: "2024-03-01",
			wantErr:    false,
		},
		{
			name:       "加24小时",
			dateStr:    "2024-01-01 00:00:00",
			movements:  []any{24 * Hour},
			wantOutput: "2024-01-02 00:00:00",
			wantErr:    false,
		},
		{
			name:       "加90分钟",
			dateStr:    "2024-01-01 00:00:00",
			movements:  []any{90 * Minute},
			wantOutput: "2024-01-01 01:30:00",
			wantErr:    false,
		},
		{
			name:       "加3600秒",
			dateStr:    "2024-01-01 00:00:00",
			movements:  []any{3600 * Second},
			wantOutput: "2024-01-01 01:00:00",
			wantErr:    false,
		},
		// 不支持这种调整
		{
			name:       "混合调整-单个表达式",
			dateStr:    "2024-01-01 00:00:00",
			movements:  []any{-1 * Year, 2 * Month, 10 * Day, 12*Hour + 30*Minute + 45*Second},
			wantOutput: "2023-03-11 12:30:45",
			wantErr:    false,
		},
		{
			name:       "混合调整-多个表达式",
			dateStr:    "2024-01-01 00:00:00",
			movements:  []any{-1 * Year, 2 * Month, 10 * Day, 12 * Hour, 30 * Minute, 45 * Second},
			wantOutput: "2023-03-11 12:30:45",
			wantErr:    false,
		},
		{
			name:       "边界测试-月末",
			dateStr:    "2024-01-31",
			movements:  []any{1 * Month, 1 * Month},
			wantOutput: "2024-03-31",
			wantErr:    false,
		},
		{
			name:       "边界测试-2月29日",
			dateStr:    "2024-02-29",
			movements:  []any{1 * Year},
			wantOutput: "2025-02-28", // 2025年不是闰年
			wantErr:    false,
		},
		{
			name:       "不同格式日期-斜杠",
			dateStr:    "2024/01/01",
			movements:  []any{1 * Month},
			wantOutput: "2024-02-01",
			wantErr:    false,
		},
		{
			name:       "不同格式日期-美式",
			dateStr:    "01/02/2024",
			movements:  []any{1 * Month},
			wantOutput: "2024-02-02",
			wantErr:    false,
		},
		{
			name:       "错误格式日期",
			dateStr:    "invalid-date",
			movements:  []any{1 * Day},
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:       "保留输出格式-日期",
			dateStr:    "2024-01-01",
			movements:  []any{1 * Day},
			wantOutput: "2024-01-02",
			wantErr:    false,
		},
		{
			name:       "保留输出格式-日期时间",
			dateStr:    "2024-01-01 12:00:00",
			movements:  []any{1 * Day},
			wantOutput: "2024-01-02 12:00:00",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got = tt.dateStr
			var err error
			got, err = Move(got, tt.movements...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantOutput {
				t.Errorf("Move() = %v, want %v", got, tt.wantOutput)
			}
		})
	}

}

// TestMoveEdgeCases 测试一些特殊边界情况
func TestMoveEdgeCases(t *testing.T) {
	// 测试闰秒、夏令时等特殊情况
	// 注意：Go的时间库会自动处理夏令时和闰秒

	// 测试非常大的时间调整
	t.Run("大数值调整", func(t *testing.T) {
		got, err := Move("2024-01-01", 50*Year, 60*Month)
		if err != nil {
			t.Errorf("Move() error = %v", err)
			return
		}
		// 50年加60个月 = 55年
		expected := "2079-01-01"
		if got != expected {
			t.Errorf("Move() = %v, want %v", got, expected)
		}
	})

	// 测试负数和零
	t.Run("零调整", func(t *testing.T) {
		got, err := Move("2024-01-01", 0*Year)
		if err != nil {
			t.Errorf("Move() error = %v", err)
			return
		}
		expected := "2024-01-01"
		if got != expected {
			t.Errorf("Move() = %v, want %v", got, expected)
		}
	})

	// 测试多个时区
	t.Run("时区处理", func(t *testing.T) {
		got, err := Move("2024-01-01T12:00:00+08:00", 1*Day)
		if err != nil {
			t.Errorf("Move() error = %v", err)
			return
		}
		// 检查结果应该保留时区信息
		if len(got) < 10 {
			t.Errorf("Move() result too short = %v", got)
		}
	})
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		left     string
		operator string
		right    string
		want     bool
		wantErr  bool
	}{
		{"相等", "2024-01-01", EQ, "2024-01-01", true, false},
		{"不相等", "2024-01-01", NE, "2024-01-02", true, false},
		{"大于", "2024-01-01", GT, "2024-01-01", false, false},
		{"大于等于", "2024-01-01", GE, "2024-01-01", true, false},
		{"小于", "2024-01-01", LT, "2024-01-02", true, false},
		{"小于等于", "2024-01-01", LE, "2024-01-01", true, false},
		{"错误格式", "2024-01-01", EQ, "invalid-date", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compare(tt.left, tt.operator, tt.right)
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}

	if Compare("2024-02-03 20:00", ">=", "2024-02-02") {
		fmt.Println("a jian")
	}
}

// 性能测试
func BenchmarkGuess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Guess("2024-01-01 15:04:05")
	}
}

func BenchmarkMove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Move("2024-01-01", 1*Year, 2*Month, 3*Day)
	}
}
