package timeutils

import (
	"testing"
	"time"
)

func TestBuildDuration(t *testing.T) {
	// 注意：'now' 在每个 t.Run 中定义，以确保时间准确性

	testCases := []struct {
		name                string
		offset              time.Duration // 使用偏移量代替固定的时间字符串
		isInvalid           bool          // 标记是否为无效输入
		expectedStep        Step
		expectedLayout      string
		expectStartModified bool
		modifiedDuration    time.Duration
	}{
		{
			name:                "Invalid start time",
			isInvalid:           true,
			expectedStep:        MinuteStep,
			expectedLayout:      "2006-01-02 1504",
			expectStartModified: true,
			modifiedDuration:    -30 * time.Minute,
		},
		{
			name:           "Duration less than 1 minute",
			offset:         -30 * time.Second,
			expectedStep:   SecondStep,
			expectedLayout: "2006-01-02 150405",
		},
		{
			name:           "Duration less than 1 hour",
			offset:         -30 * time.Minute,
			expectedStep:   MinuteStep,
			expectedLayout: "2006-01-02 1504",
		},
		{
			name:           "Duration less than 1 day",
			offset:         -12 * time.Hour,
			expectedStep:   HourStep,
			expectedLayout: "2006-01-02 15",
		},
		{
			name:           "Duration less than 7 days",
			offset:         -3 * 24 * time.Hour,
			expectedStep:   DayStep,
			expectedLayout: "2006-01-02",
		},
		{
			name:                "Duration more than 7 days",
			offset:              -10 * 24 * time.Hour,
			expectedStep:        DayStep,
			expectedLayout:      "2006-01-02",
			expectStartModified: true,
			modifiedDuration:    -7 * 24 * time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 在每个测试用例开始时获取当前时间，以保证时间相对准确
			now := time.Now()

			var currentInputStart string
			if tc.isInvalid {
				currentInputStart = "invalid-time-string"
			} else {
				currentInputStart = now.Add(tc.offset).Format(time.DateTime)
			}

			duration, err := BuildDuration(currentInputStart)
			if err != nil {
				t.Fatalf("BuildDuration failed: %v", err)
			}

			if duration.Step != string(tc.expectedStep) {
				t.Errorf("Expected step %s, but got %s", tc.expectedStep, duration.Step)
			}

			var expectedStartTime time.Time
			if tc.expectStartModified {
				expectedStartTime = now.Add(tc.modifiedDuration)
			} else {
				expectedStartTime, _ = time.Parse(time.DateTime, currentInputStart)
			}

			expectedStartStr := expectedStartTime.Format(tc.expectedLayout)
			// 由于 now 的细微差别，对修改后的 start 时间进行近似比较
			if tc.expectStartModified {
				// 比较到秒级精度
				gotStart, _ := time.Parse(tc.expectedLayout, duration.Start)
				wantStart, _ := time.Parse(tc.expectedLayout, expectedStartStr)
				if gotStart.Truncate(time.Second).Unix() != wantStart.Truncate(time.Second).Unix() {
					t.Errorf("Expected start time to be around %s, but got %s", expectedStartStr, duration.Start)
				}
			} else if duration.Start != expectedStartStr {
				t.Errorf("Expected start time %s, but got %s", expectedStartStr, duration.Start)
			}

			expectedEndStr := now.Format(tc.expectedLayout)
			if duration.End != expectedEndStr {
				t.Errorf("Expected end time %s, but got %s", expectedEndStr, duration.End)
			}
		})
	}
}
