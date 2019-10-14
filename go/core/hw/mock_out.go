package hw

import "github.com/G0WCZ/cwc/config"

type MockOut struct {
	Config   *config.Config
	bitValue bool
}

func NewMockOut(config *config.Config) MorseIn {
	return &MockIn{
		Config:   config,
		bitValue: false,
	}
}

func (m *MockOut) Open() error {
	m.bitValue = false
}

func (m *MockOut) SetBit(bool) {
	panic("implement me")
}

func (m *MockOut) SetToneOut(bool) {
	panic("implement me")
}

func (m *MockOut) SetStatusLED(bool) {
	panic("implement me")
}

func (m *MockOut) Close() error {
	panic("implement me")
}

func (m *MockOut) Bit() bool {
	return m.bitValue
}

func (m *MockOut) SetBitValue(bit bool) {
	m.bitValue = bit
}
