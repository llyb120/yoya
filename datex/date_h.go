package datex

import (
	"time"
)

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
