package datex

import (
	"fmt"
	"time"
)

// 定义时间单位常量，使用特定的数值便于计算
type TimeUnit int64

const (
	Second time.Duration = time.Second
	Minute time.Duration = time.Minute
	Hour   time.Duration = time.Hour
	Day    TimeUnit      = 1000000
	Month  TimeUnit      = 100000000
	Year   TimeUnit      = 10000000000
)

// Guess 函数尝试按优先级从上到下解析字符串时间
func Guess(dateStr string) (time.Time, error) {
	// 按优先级从高到低尝试解析
	formats := []string{
		"2006-01-02",           // 仅日期
		"2006-01-02 15:04:05",  // 标准日期时间格式
		"2006-01-02T15:04:05Z", // ISO8601
		"2006-01-02T15:04:05",  // ISO8601 无时区
		"2006/01/02 15:04:05",  // 斜杠分隔日期时间
		"2006-01-02 15:04",     // 无秒
		"2006/01/02 15:04",     // 斜杠分隔无秒
		"2006/01/02",           // 斜杠分隔仅日期
		"01/02/2006",           // 美式日期
		//"02/01/2006",         // 欧式日期, 和上面一样，猜不出哪个是月，故不支持
		"20060102150405", // 紧凑格式
		"20060102",       // 紧凑日期
		time.RFC3339,     // RFC3339
		time.RFC3339Nano, // RFC3339带纳秒
		time.RFC1123,     // RFC1123
		time.RFC1123Z,    // RFC1123带时区
		time.RFC850,      // RFC850
		time.RFC822,      // RFC822
		time.RFC822Z,     // RFC822带时区
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// 如果都无法解析，返回空
	return time.Time{}, fmt.Errorf("日期格式错误")
}

// 获取月份的最后一天
func lastDayOfMonth(year int, month time.Month) int {
	// 获取下个月的第一天，然后减去一天
	firstDayOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDayOfNextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}

// 处理月份边界问题
func adjustMonthBoundary(t time.Time, years, months, days int) time.Time {
	// 记录原始日期的信息
	originalYear := t.Year()
	originalMonth := t.Month()
	originalDay := t.Day()

	// 检查原始日期是否是月末
	lastDayOfOriginalMonth := lastDayOfMonth(originalYear, originalMonth)
	isLastDay := originalDay == lastDayOfOriginalMonth

	// 计算目标年月
	targetYear := originalYear + years
	targetMonth := originalMonth + time.Month(months)

	// 调整月份溢出（例如13月变成下一年的1月）
	for targetMonth > 12 {
		targetYear++
		targetMonth -= 12
	}
	for targetMonth < 1 {
		targetYear--
		targetMonth += 12
	}

	// 确定目标日
	var targetDay int
	if isLastDay {
		// 如果原始日期是月末，则目标日期也应该是月末
		targetDay = lastDayOfMonth(targetYear, targetMonth)
	} else {
		// 如果不是月末，则尝试保持原始日期的日
		lastDayOfTargetMonth := lastDayOfMonth(targetYear, targetMonth)
		if originalDay > lastDayOfTargetMonth {
			// 如果原始日期的日大于目标月份的最后一天，则使用目标月份的最后一天
			targetDay = lastDayOfTargetMonth
		} else {
			targetDay = originalDay
		}
	}

	// 创建新的日期
	newTime := time.Date(targetYear, targetMonth, targetDay,
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

	// 再调整天数
	if days != 0 {
		newTime = newTime.AddDate(0, 0, days)
	}

	return newTime
}

func Move(date string, movements ...any) (string, error) {
	// 使用Guess函数尝试解析日期
	t, err := Guess(date)
	if err != nil {
		return "", err
	}

	// 记录原始日期是否包含时间部分
	hasTime := len(date) > 10

	// 处理所有的时间调整
	for _, m := range movements {
		// 提取各个时间单位的调整值
		switch m := any(m).(type) {
		case TimeUnit:
			years := int((m / Year) % 100)
			months := int((m / Month) % 100)
			days := int((m / Day) % 100)
			// 使用自定义函数处理年月日的调整，特别是月份边界问题
			if years != 0 || months != 0 || days != 0 {
				t = adjustMonthBoundary(t, years, months, days)
			}
		case time.Duration:
			t = t.Add(m)
		}
	}

	// 根据输入日期格式决定输出格式
	if hasTime {
		return t.Format("2006-01-02 15:04:05"), nil
	}
	return t.Format("2006-01-02"), nil
}
