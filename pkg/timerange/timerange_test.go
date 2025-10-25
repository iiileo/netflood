package timerange

import (
	"testing"
	"time"
)

func TestParseTimeRanges(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "单个时间段",
			input:   "12:00-13:00",
			wantErr: false,
		},
		{
			name:    "多个时间段",
			input:   "12:00-13:00,14:00-15:00",
			wantErr: false,
		},
		{
			name:    "三个时间段",
			input:   "09:00-10:00,14:00-16:00,20:00-22:00",
			wantErr: false,
		},
		{
			name:    "跨天时间段",
			input:   "23:00-01:00",
			wantErr: false,
		},
		{
			name:    "无效格式 - 缺少结束时间",
			input:   "12:00-",
			wantErr: true,
		},
		{
			name:    "无效格式 - 错误的时间格式",
			input:   "12-13",
			wantErr: true,
		},
		{
			name:    "无效格式 - 小时超出范围",
			input:   "25:00-26:00",
			wantErr: true,
		},
		{
			name:    "无效格式 - 分钟超出范围",
			input:   "12:70-13:00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseTimeRanges(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTimeRanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTimeRangeManager(t *testing.T) {
	tests := []struct {
		name       string
		timeStr    string
		wantEnable bool
		wantErr    bool
	}{
		{
			name:       "空字符串 - 禁用时间控制",
			timeStr:    "",
			wantEnable: false,
			wantErr:    false,
		},
		{
			name:       "有效时间段",
			timeStr:    "12:00-13:00",
			wantEnable: true,
			wantErr:    false,
		},
		{
			name:       "多个有效时间段",
			timeStr:    "12:00-13:00,14:00-15:00",
			wantEnable: true,
			wantErr:    false,
		},
		{
			name:       "无效时间段",
			timeStr:    "invalid",
			wantEnable: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trm, err := NewTimeRangeManager(tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTimeRangeManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && trm.IsEnabled() != tt.wantEnable {
				t.Errorf("NewTimeRangeManager() enabled = %v, want %v", trm.IsEnabled(), tt.wantEnable)
			}
		})
	}
}

func TestTimeRangeManager_IsInRange(t *testing.T) {
	// 测试：当前时间在范围内
	now := time.Now()
	beforeNow := now.Add(-30 * time.Minute)
	afterNow := now.Add(30 * time.Minute)

	timeStr := beforeNow.Format("15:04") + "-" + afterNow.Format("15:04")

	trm, err := NewTimeRangeManager(timeStr)
	if err != nil {
		t.Fatalf("NewTimeRangeManager() error = %v", err)
	}

	if !trm.IsInRange() {
		t.Errorf("IsInRange() = false, want true (current time should be in range)")
	}
}

func TestTimeRangeManager_IsInRange_Disabled(t *testing.T) {
	trm, err := NewTimeRangeManager("")
	if err != nil {
		t.Fatalf("NewTimeRangeManager() error = %v", err)
	}

	if !trm.IsInRange() {
		t.Errorf("IsInRange() = false, want true (should always be true when disabled)")
	}
}

func TestTimeRangeManager_String(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		wantStr string
	}{
		{
			name:    "禁用",
			timeStr: "",
			wantStr: "全天候运行",
		},
		{
			name:    "单个时间段",
			timeStr: "12:00-13:00",
			wantStr: "12:00-13:00",
		},
		{
			name:    "多个时间段",
			timeStr: "12:00-13:00,14:00-15:00",
			wantStr: "12:00-13:00, 14:00-15:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trm, err := NewTimeRangeManager(tt.timeStr)
			if err != nil {
				t.Fatalf("NewTimeRangeManager() error = %v", err)
			}
			if got := trm.String(); got != tt.wantStr {
				t.Errorf("String() = %v, want %v", got, tt.wantStr)
			}
		})
	}
}
