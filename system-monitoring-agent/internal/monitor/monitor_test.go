package monitor

import (
	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the gopsutil functions for testing
type MockCPU struct {
	mock.Mock
}

func (m *MockCPU) Percent(interval float64, percpu bool) ([]float64, error) {
	args := m.Called(interval, percpu)
	return args.Get(0).([]float64), args.Error(1)
}

type MockMem struct {
	mock.Mock
}

func (m *MockMem) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	args := m.Called()
	return args.Get(0).(*mem.VirtualMemoryStat), args.Error(1)
}

// Test for GetCPUUsage
func TestGetCPUUsage(t *testing.T) {
	mockCPU := new(MockCPU)
	cpu.Percent = mockCPU.Percent // Replace the original function with the mock

	mockCPU.On("Percent", 0, false).Return([]float64{50.0}, nil)

	m := NewMonitor()
	usage, err := m.GetCPUUsage()
	assert.NoError(t, err)
	assert.Equal(t, 50.0, usage)

	mockCPU.AssertExpectations(t)
}

// Test for GetMemoryUsage
func TestGetMemoryUsage(t *testing.T) {
	mockMem := new(MockMem)
	mem.VirtualMemory = mockMem.VirtualMemory // Replace the original function with the mock

	mockMem.On("VirtualMemory").Return(&mem.VirtualMemoryStat{Used: 2048}, nil)

	m := NewMonitor()
	usage, err := m.GetMemoryUsage()
	assert.NoError(t, err)
	assert.Equal(t, uint64(2048), usage)

	mockMem.AssertExpectations(t)
}

// Test for LogSystemMetrics
func TestLogSystemMetrics(t *testing.T) {
	// This test would require capturing log output, which can be done using a buffer
	// For simplicity, we will skip the implementation here
}
