package datex

import (
	"fmt"
	"strings"
	"time"
)

// 定义时间单位常量，使用特定的数值便于计算
type moveType string
type dayUnit int64
type weekUnit int64
type monthUnit int64
type yearUnit int64

const (
	Second time.Duration = time.Second
	Minute time.Duration = time.Minute
	Hour   time.Duration = time.Hour
	Day    dayUnit       = 1
	Week   weekUnit      = 1
	Month  monthUnit     = 1
	Year   yearUnit      = 1

	FirstDayOfMonth  moveType = "FirstDayOfMonth"
	LastDayOfMonth   moveType = "LastDayOfMonth"
	FirstDayOfYear   moveType = "FirstDayOfYear"
	LastDayOfYear    moveType = "LastDayOfYear"
	FirstDayOfWeek   moveType = "FirstDayOfWeek"
	LastDayOfWeek    moveType = "LastDayOfWeek"
	FirstDayOfCNWeek moveType = "FirstDayOfCNWeek"
	LastDayOfCNWeek  moveType = "LastDayOfCNWeek"

	EQ  = "=="
	NE  = "!="
	GT  = ">"
	GE  = ">="
	LT  = "<"
	LE  = "<="
	MEQ = "~="
)

func Guess(dateStr string) (time.Time, error) {
	t, _, err := guess(dateStr)
	return t, err
}

// Guess 函数尝试按优先级从上到下解析字符串时间
func guess(dateStr string) (time.Time, string, error) {
	// 去除可能的空白字符
	dateStr = strings.TrimSpace(dateStr)

	if dateStr == "" {
		return time.Time{}, "", fmt.Errorf("日期字符串为空")
	}

	// 根据长度尝试解析
	var formats []string
	switch len(dateStr) {
	case 10:
		formats = []string{
			"2006-01-02", // 标准日期
			"2006/01/02", // 斜杠分隔日期
			"01/02/2006", // 美式日期
		}

	case 16:
		formats = []string{
			"2006-01-02 15:04",
			"2006/01/02 15:04",
		}

	case 19:
		formats = []string{
			"2006-01-02 15:04:05",
			"2006/01/02 15:04:05",
			"2006-01-02T15:04:05",
		}

	case 20:
		formats = []string{
			time.RFC3339, // ISO8601带时区
		}

	case 24, 25:
		formats = []string{
			time.RFC3339, // 带时区的ISO8601
			time.RFC1123Z,
		}

	case 14:
		formats = []string{
			"20060102150405", // 紧凑格式
		}

	case 8:
		formats = []string{
			"20060102", // 紧凑日期
		}

	case 28, 29, 30:
		formats = []string{
			time.RFC3339Nano,
		}

	case 22, 23:
		formats = []string{
			time.RFC850,
			time.RFC1123,
		}

	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, format, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("日期格式错误: %s", dateStr)
}

func Move[T string | *string | time.Time | *time.Time](date T, movements ...any) T {
	if len(movements) == 0 {
		return date
	}
	// 使用Guess函数尝试解析日期
	var flag bool = true
	var d string
	var isString bool
	var isPointer bool
	switch s := any(date).(type) {
	case string:
		d = s
		isString = true
	case *string:
		if s == nil {
			return date
		}
		d = *s
		isString = true
		isPointer = true
	case *time.Time:
		if s == nil {
			return date
		}
		isPointer = true
	}

	var t time.Time
	var err error
	var format string

	if isString {
		t, format, err = guess(d)
		if err != nil {
			var zero T
			return zero
		}
	} else {
		t = any(date).(time.Time)
	}

	// 处理所有的时间调整
	var years, months, days int
	var duration int64
	var mType moveType
	for _, m := range movements {
		// 提取各个时间单位的调整值
		switch m := any(m).(type) {
		case yearUnit:
			years += int(m)
		case monthUnit:
			months += int(m)
		case dayUnit:
			days += int(m)
		case weekUnit:
			days += int(m) * 7
		case time.Duration:
			duration += int64(m)
		case moveType:
			mType = m
		case bool:
			flag = flag && m
		}
	}
	if !flag {
		return date
	}
	// 使用自定义函数处理年月日的调整，特别是月份边界问题
	if years != 0 || months != 0 || days != 0 {
		t = adjustMonthBoundary(t, years, months, days)
	}
	// 如果需要按照mType调整
	if mType != "" {
		switch mType {
		case FirstDayOfMonth:
			t = time.Date(t.Year(), t.Month(), 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		case LastDayOfMonth:
			t = t.AddDate(0, 1, -t.Day())
		case FirstDayOfYear:
			t = time.Date(t.Year(), 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		case LastDayOfYear:
			t = time.Date(t.Year(), 12, 31, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		case FirstDayOfWeek:
			t = t.AddDate(0, 0, -int(t.Weekday()))
		case FirstDayOfCNWeek:
			// 对于中国周（从周一开始），需要特殊处理
			// 如果是周日(0)，需要回退6天；否则回退到周一
			offset := int(t.Weekday())
			if offset == 0 {
				offset = 6
			} else {
				offset -= 1
			}
			t = t.AddDate(0, 0, -offset)
		case LastDayOfWeek:
			t = t.AddDate(0, 0, 6-int(t.Weekday()))
		case LastDayOfCNWeek:
			// 对于中国周（从周一开始），需要特殊处理
			// 如果是周日(0)，需要回退6天；否则回退到周一
			offset := int(t.Weekday())
			if offset == 0 {
				offset = 6
			} else {
				offset -= 1
			}
			t = t.AddDate(0, 0, 6-int(t.Weekday()))
		}
	}
	if duration != 0 {
		t = t.Add(time.Duration(duration))
	}

	if isString {
		// 根据输入日期格式决定输出格式
		res := t.Format(format)
		if isPointer {
			*(any(date).(*string)) = res
			return date
		}
		return any(res).(T)
	} else {
		if isPointer {
			*(any(date).(*time.Time)) = t
			return date
		}
		return any(t).(T)
	}
}

func When[L string | time.Time, R string | time.Time](left L, operator string, right R) bool {
	cache := compareHolder.Get()

	str, isStr := any(left).(string)
	var leftTime time.Time
	// 获取左值
	if isStr {
		leftTime = cache.GetOrSetFunc(str, func() time.Time {
			t, err := Guess(str)
			if err != nil {
				return time.Time{}
			}
			return t
		})
	} else {
		leftTime = any(left).(time.Time)
	}
	// 获取右值
	str, isStr = any(right).(string)
	var rightTime time.Time
	if isStr {
		rightTime = cache.GetOrSetFunc(str, func() time.Time {
			t, err := Guess(str)
			if err != nil {
				return time.Time{}
			}
			return t
		})
	} else {
		rightTime = any(right).(time.Time)
	}
	switch operator {
	case GT:
		return leftTime.After(rightTime)
	case GE:
		return leftTime.After(rightTime) || leftTime.Equal(rightTime)
	case LT:
		return leftTime.Before(rightTime)
	case LE:
		return leftTime.Before(rightTime) || leftTime.Equal(rightTime)
	case EQ:
		return leftTime.Equal(rightTime)
	case NE:
		return !leftTime.Equal(rightTime)
	case MEQ:
		return leftTime.Year() == rightTime.Year() && leftTime.Month() == rightTime.Month() && leftTime.Day() == rightTime.Day()
	default:
		return false
	}
}
