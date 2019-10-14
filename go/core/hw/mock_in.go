package hw

import "github.com/G0WCZ/cwc/config"

type MockIn struct {
	Config      *config.Config
	adapterName string
	keyer       bool
	bitValue    bool
	ditValue    bool
	dahValue    bool
}

func NewMockIn(config *config.Config, adapterName string) MorseIn {
	return &MockIn{
		Config:      config,
		adapterName: adapterName,
		keyer:       KEYER == adapterName,
		bitValue:    false,
		ditValue:    false,
		dahValue:    false,
	}
}

func (m *MockIn) Open() error {
	m.bitValue = false
	m.ditValue = false
	m.dahValue = false
}

func (m *MockIn) ConfigChanged() error {
}

func (m *MockIn) Bit() bool {
	return m.bitValue
}

func (m *MockIn) Dit() bool {
	return m.ditValue
}

func (m *MockIn) Dah() bool {
	return m.dahValue
}

func (m *MockIn) Close() error {
}

func (m *MockIn) UseKeyer() bool {
	return m.keyer
}

func (m *MockIn) SetBit(bit bool) {
	m.bitValue = bit
}

func (m *MockIn) SetDit(bit bool) {
	m.ditValue = bit
}

func (m *MockIn) SetDah(bit bool) {
	m.dahValue = bit
}
