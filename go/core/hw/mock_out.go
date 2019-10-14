package hw

import "github.com/G0WCZ/cwc/config"

type MockOut struct {
	Config          *config.Config
	adapterName     string
	currentBitValue bool
}

func (m *MockOut) Open() error {
	m.currentBitValue = false
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
	return
}
