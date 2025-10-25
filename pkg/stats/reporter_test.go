package stats

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewReporter(t *testing.T) {
	reporter, err := NewReporter("http://example.com/stats")
	if err != nil {
		t.Fatalf("NewReporter() error = %v", err)
	}

	if reporter.apiURL != "http://example.com/stats" {
		t.Errorf("NewReporter() apiURL = %v, want %v", reporter.apiURL, "http://example.com/stats")
	}

	if reporter.hostname == "" {
		t.Errorf("NewReporter() hostname is empty")
	}
}

func TestReporter_Report(t *testing.T) {
	// 创建测试服务器
	var receivedData StatsData
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// 检查Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// 解析请求体
		if err := json.NewDecoder(r.Body).Decode(&receivedData); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建Reporter
	reporter, err := NewReporter(server.URL)
	if err != nil {
		t.Fatalf("NewReporter() error = %v", err)
	}

	// 测试上报
	err = reporter.Report(15.5, 1024.0, "12:00-13:00")
	if err != nil {
		t.Errorf("Report() error = %v", err)
	}

	// 验证接收到的数据
	if receivedData.Speed != 15.5 {
		t.Errorf("Expected speed 15.5, got %f", receivedData.Speed)
	}

	if receivedData.Total != 1024.0 {
		t.Errorf("Expected total 1024.0, got %f", receivedData.Total)
	}

	if receivedData.Time != "12:00-13:00" {
		t.Errorf("Expected time '12:00-13:00', got %s", receivedData.Time)
	}

	if receivedData.Name != reporter.hostname {
		t.Errorf("Expected name '%s', got %s", reporter.hostname, receivedData.Name)
	}
}

func TestReporter_Report_ServerError(t *testing.T) {
	// 创建返回错误的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	reporter, err := NewReporter(server.URL)
	if err != nil {
		t.Fatalf("NewReporter() error = %v", err)
	}

	// 测试上报（应该返回错误）
	err = reporter.Report(15.5, 1024.0, "12:00-13:00")
	if err == nil {
		t.Error("Expected error for server error response, got nil")
	}
}

func TestReporter_GetHostname(t *testing.T) {
	reporter, err := NewReporter("http://example.com/stats")
	if err != nil {
		t.Fatalf("NewReporter() error = %v", err)
	}

	hostname := reporter.GetHostname()
	if hostname == "" {
		t.Error("GetHostname() returned empty string")
	}
}
