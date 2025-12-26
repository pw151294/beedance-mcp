package timeutils

import (
	"beedance-mcp/api/tools"
	"time"
)

type Step string

const SecondStep Step = "SECOND"
const MinuteStep Step = "MINUTE"
const HourStep Step = "HOUR"
const DayStep Step = "DAY"

// 根据开始时间距离当前时间的间隔，分配合适的步长
// 1.将start转换成time.Time类型，如果转换失败，start默认设置成time.Now().Add(-24 * time.Hour)，返回"DAY"
// 2.计算start和time.Now()之间的间隔：
// 1）小于1分钟，返回SecondStep
// 2）小于1小时，返回MinuteStep
// 3）小于一天，返回HourStep
// 4）小于七天，返回DayStep
// 5）大于七天，start设置成time.Now().Add(-7*time.Day)，返回"DAY"
func parseStep(start *string) Step {
	var startTime time.Time
	var err error

	// 1. 尝试在本地时区解析时间
	startTime, err = time.ParseInLocation(time.DateTime, *start, time.Local)
	if err != nil {
		// 如果解析失败，设置默认 start 时间，并用该时间继续后续计算
		newStart := time.Now().Add(-24 * time.Hour)
		*start = newStart.Format(time.DateTime)
		startTime = newStart
	}

	// 2. 计算时间间隔
	now := time.Now()
	duration := now.Sub(startTime)

	// 根据间隔返回步长
	if duration < time.Minute {
		return SecondStep
	}
	if duration < time.Hour {
		return MinuteStep
	}
	if duration < 24*time.Hour {
		return HourStep
	}
	if duration < 7*24*time.Hour {
		return DayStep
	}

	// 大于等于七天
	newStart := time.Now().Add(-7 * 24 * time.Hour)
	*start = newStart.Format(time.DateTime)
	return DayStep
}

func BuildDuration(start string) (tools.Duration, error) {
	step := parseStep(&start)
	end := time.Now()

	// 根据Step对start还有end进行格式化：
	startTime, err := time.ParseInLocation(time.DateTime, start, time.Local)
	if err != nil {
		return tools.Duration{}, err
	}

	var layout string
	switch step {
	case SecondStep:
		layout = "2006-01-02 150405"
	case MinuteStep:
		layout = "2006-01-02 1504"
	case HourStep:
		layout = "2006-01-02 15"
	case DayStep:
		layout = "2006-01-02"
	}

	duration := tools.Duration{}
	duration.Start = startTime.Format(layout)
	duration.End = end.Format(layout)
	duration.Step = string(step)
	return duration, nil
}
