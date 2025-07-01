package mbserver

import "sync"

// id in function definition is slave id when use tcp or rtu
type Slaver interface {
	IsSlaveIdValid(id uint8) bool
	DiscreteInputs(id uint8) ([]byte, error)
	Coils(id uint8) ([]byte, error)
	HoldingRegisters(id uint8) ([]uint16, error)
	InputRegisters(id uint8) ([]uint16, error)
	SaveDiscreteInputs(id uint8, b []byte) error
	SaveCoils(id uint8, b []byte) error
	SaveHoldingRegisters(id uint8, b []uint16) error
	SaveInputRegisters(id uint8, b []uint16) error
}

var _ Slaver = new(memorySlaveUint8)

type memorySlaveUint8 struct {
	slaveNum         uint8
	slaveLock        []sync.RWMutex
	discreteInputs   [][]byte
	coils            [][]byte
	holdingRegisters [][]uint16
	inputRegisters   [][]uint16
}

// will create slaveNum slaves, slave id is [1, slaveNum], slaveNumMax is 255, slaveNumMin is 1, if slaveNum > 255, it will use 255, if slaveNum < 1, it will use 1, one uint8, 2^8-1, 255
func NewMemorySlaveUint8(slaveNum uint8) (slaver Slaver) {

	if slaveNum < 1 {
		slaveNum = 1
	} else if slaveNum > 255 {
		slaveNum = 255
	}
	var s = &memorySlaveUint8{
		slaveNum:         slaveNum,
		slaveLock:        make([]sync.RWMutex, slaveNum),
		discreteInputs:   make([][]byte, slaveNum),
		coils:            make([][]byte, slaveNum),
		holdingRegisters: make([][]uint16, slaveNum),
		inputRegisters:   make([][]uint16, slaveNum),
	}
	for i := range s.discreteInputs {
		s.discreteInputs[i] = make([]byte, 65536)
	}
	for i := range s.coils {
		s.coils[i] = make([]byte, 65536)
	}
	for i := range s.holdingRegisters {
		s.holdingRegisters[i] = make([]uint16, 65536)
	}
	for i := range s.inputRegisters {
		s.inputRegisters[i] = make([]uint16, 65536)
	}
	slaver = s
	return
}

func (s *memorySlaveUint8) IsSlaveIdValid(id uint8) bool { return id > 0 && id <= s.slaveNum }

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *memorySlaveUint8) DiscreteInputs(id uint8) (bs []byte, err error) {

	id = s.getRealId(id)
	s.slaveLock[id].RLock()
	bs = CopyBytes(s.discreteInputs[id])
	s.slaveLock[id].RUnlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *memorySlaveUint8) Coils(id uint8) (bs []byte, err error) {

	id = s.getRealId(id)
	s.slaveLock[id].RLock()
	bs = CopyBytes(s.coils[id])
	s.slaveLock[id].RUnlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *memorySlaveUint8) HoldingRegisters(id uint8) (bs []uint16, err error) {

	id = s.getRealId(id)
	s.slaveLock[id].RLock()
	bs = CopyUint16(s.holdingRegisters[id])
	s.slaveLock[id].RUnlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *memorySlaveUint8) InputRegisters(id uint8) (bs []uint16, err error) {

	id = s.getRealId(id)
	s.slaveLock[id].RLock()
	bs = CopyUint16(s.inputRegisters[id])
	s.slaveLock[id].RUnlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *memorySlaveUint8) SaveDiscreteInputs(id uint8, b []byte) (err error) {

	id = s.getRealId(id)
	s.slaveLock[id].Lock()
	s.discreteInputs[id] = b
	s.slaveLock[id].Unlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *memorySlaveUint8) SaveCoils(id uint8, b []byte) (err error) {

	id = s.getRealId(id)
	s.slaveLock[id].Lock()
	s.coils[id] = b
	s.slaveLock[id].Unlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *memorySlaveUint8) SaveHoldingRegisters(id uint8, b []uint16) (err error) {

	id = s.getRealId(id)
	s.slaveLock[id].Lock()
	s.holdingRegisters[id] = b
	s.slaveLock[id].Unlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *memorySlaveUint8) SaveInputRegisters(id uint8, b []uint16) (err error) {

	id = s.getRealId(id)
	s.slaveLock[id].Lock()
	s.inputRegisters[id] = b
	s.slaveLock[id].Unlock()
	return
}

func (s *memorySlaveUint8) getRealId(id uint8) (realId uint8) {

	switch {
	case id > s.slaveNum:
		realId = s.slaveNum - 1
	case id < 1:
		realId = 0
	default:
		realId = id - 1
	}
	return
}
