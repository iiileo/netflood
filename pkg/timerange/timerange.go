package timerange

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeRange 时间段
type TimeRange struct {
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
}

// TimeRangeManager 时间段管理器
type TimeRangeManager struct {
	ranges  []TimeRange
	enabled bool
}

// NewTimeRangeManager 创建时间段管理器
func NewTimeRangeManager(timeStr string) (*TimeRangeManager, error) {
	if timeStr == "" {
		return &TimeRangeManager{enabled: false}, nil
	}

	ranges, err := parseTimeRanges(timeStr)
	if err != nil {
		return nil, err
	}

	return &TimeRangeManager{
		ranges:  ranges,
		enabled: true,
	}, nil
}

// parseTimeRanges 解析时间段字符串
// 格式: "12:00-13:00,14:00-15:00"
func parseTimeRanges(timeStr string) ([]TimeRange, error) {
	parts := strings.Split(timeStr, ",")
	var ranges []TimeRange

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 分割开始和结束时间
		times := strings.Split(part, "-")
		if len(times) != 2 {
			return nil, fmt.Errorf("无效的时间段格式: %s (应为 HH:MM-HH:MM)", part)
		}

		startTime := strings.TrimSpace(times[0])
		endTime := strings.TrimSpace(times[1])

		// 解析开始时间
		startHour, startMinute, err := parseTime(startTime)
		if err != nil {
			return nil, fmt.Errorf("解析开始时间失败 %s: %v", startTime, err)
		}

		// 解析结束时间
		endHour, endMinute, err := parseTime(endTime)
		if err != nil {
			return nil, fmt.Errorf("解析结束时间失败 %s: %v", endTime, err)
		}

		ranges = append(ranges, TimeRange{
			StartHour:   startHour,
			StartMinute: startMinute,
			EndHour:     endHour,
			EndMinute:   endMinute,
		})
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("没有有效的时间段")
	}

	return ranges, nil
}

// parseTime 解析时间字符串 (HH:MM)
func parseTime(timeStr string) (int, int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("无效的时间格式: %s (应为 HH:MM)", timeStr)
	}

	hour, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("无效的小时: %s", parts[0])
	}

	minute, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("无效的分钟: %s", parts[1])
	}

	if hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("小时必须在 0-23 之间: %d", hour)
	}

	if minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("分钟必须在 0-59 之间: %d", minute)
	}

	return hour, minute, nil
}

// IsInRange 检查当前时间是否在允许的时间段内
func (tm *TimeRangeManager) IsInRange() bool {
	if !tm.enabled {
		return true // 未启用时间控制，始终返回true
	}

	now := time.Now()
	currentMinutes := now.Hour()*60 + now.Minute()

	for _, r := range tm.ranges {
		startMinutes := r.StartHour*60 + r.StartMinute
		endMinutes := r.EndHour*60 + r.EndMinute

		// 处理跨天的情况（例如 23:00-01:00）
		if endMinutes < startMinutes {
			// 跨天情况：如果当前时间在开始时间之后或结束时间之前
			if currentMinutes >= startMinutes || currentMinutes < endMinutes {
				return true
			}
		} else {
			// 正常情况：当前时间在开始和结束之间
			if currentMinutes >= startMinutes && currentMinutes < endMinutes {
				return true
			}
		}
	}

	return false
}

// WaitUntilNextRange 等待到下一个时间段开始
// 返回等待的持续时间，如果已经在时间段内则返回0
func (tm *TimeRangeManager) WaitUntilNextRange() time.Duration {
	if !tm.enabled {
		return 0
	}

	if tm.IsInRange() {
		return 0
	}

	now := time.Now()
	currentMinutes := now.Hour()*60 + now.Minute()
	minWait := 24 * 60 // 最大等待时间（分钟）

	for _, r := range tm.ranges {
		startMinutes := r.StartHour*60 + r.StartMinute

		var waitMinutes int
		if startMinutes > currentMinutes {
			// 今天的时间段
			waitMinutes = startMinutes - currentMinutes
		} else {
			// 明天的时间段
			waitMinutes = (24*60 - currentMinutes) + startMinutes
		}

		if waitMinutes < minWait {
			minWait = waitMinutes
		}
	}

	return time.Duration(minWait) * time.Minute
}

// GetNextRangeStart 获取下一个时间段的开始时间
func (tm *TimeRangeManager) GetNextRangeStart() time.Time {
	if !tm.enabled || tm.IsInRange() {
		return time.Now()
	}

	wait := tm.WaitUntilNextRange()
	return time.Now().Add(wait)
}

// IsEnabled 返回是否启用了时间控制
func (tm *TimeRangeManager) IsEnabled() bool {
	return tm.enabled
}

// String 返回时间段的字符串表示
func (tm *TimeRangeManager) String() string {
	if !tm.enabled {
		return "全天候运行"
	}

	var parts []string
	for _, r := range tm.ranges {
		parts = append(parts, fmt.Sprintf("%02d:%02d-%02d:%02d",
			r.StartHour, r.StartMinute, r.EndHour, r.EndMinute))
	}
	return strings.Join(parts, ", ")
}
