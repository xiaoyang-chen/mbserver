package mbserver

import (
	"encoding/json"
	"testing"
)

func isEqual(a interface{}, b interface{}) bool {
	expect, _ := json.Marshal(a)
	got, _ := json.Marshal(b)
	return string(expect) == string(got)
}

// Function 1
func TestReadCoils(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))
	// Set the coil values
	coils, err := s.Coils(1)
	if err != nil {
		t.Errorf("read slave coils fail, err: %s\n", err.Error())
		t.FailNow()
	}
	coils[10] = 1
	coils[11] = 1
	coils[17] = 1
	coils[18] = 1
	if err = s.SaveCoils(1, coils); err != nil {
		t.Errorf("write slave coils fail, err: %s\n", err.Error())
		t.FailNow()
	}

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 6
	frame.Device = 255
	frame.Function = 1
	SetDataWithRegisterAndNumber(&frame, 10, 9)

	var req Request
	req.frame = &frame
	response := s.handle(&req)

	exception := GetException(response)
	if exception != Success {
		t.Errorf("expected Success, got %v", exception.String())
		t.FailNow()
	}
	// 2 bytes, 0b1000011, 0b00000001
	expect := []byte{2, 131, 1}
	got := response.GetData()
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v", expect, got)
	}
}

// Function 2
func TestReadDiscreteInputs(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))
	// Set the discrete input values
	discreteInputs, err := s.DiscreteInputs(1)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
		t.FailNow()
	}
	discreteInputs[0] = 1
	discreteInputs[7] = 1
	discreteInputs[8] = 1
	discreteInputs[9] = 1
	if err = s.SaveDiscreteInputs(1, discreteInputs); err != nil {
		t.Errorf("expected nil, got %v", err)
		t.FailNow()
	}

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 6
	frame.Device = 255
	frame.Function = 2
	SetDataWithRegisterAndNumber(&frame, 0, 10)

	var req Request
	req.frame = &frame
	response := s.handle(&req)

	exception := GetException(response)
	if exception != Success {
		t.Errorf("expected Success, got %v", exception.String())
		t.FailNow()
	}
	expect := []byte{2, 129, 3}
	got := response.GetData()
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v", expect, got)
	}
}

// Function 3
func TestReadHoldingRegisters(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))
	holdingRegisters, err := s.HoldingRegisters(1)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
		t.FailNow()
	}
	holdingRegisters[100] = 1
	holdingRegisters[101] = 2
	holdingRegisters[102] = 65535
	if err = s.SaveHoldingRegisters(1, holdingRegisters); err != nil {
		t.Errorf("expected nil, got %v\n", err)
		t.FailNow()
	}

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 6
	frame.Device = 255
	frame.Function = 3
	SetDataWithRegisterAndNumber(&frame, 100, 3)

	var req Request
	req.frame = &frame
	response := s.handle(&req)
	exception := GetException(response)
	if exception != Success {
		t.Errorf("expected Success, got %v", exception.String())
		t.FailNow()
	}
	expect := []byte{6, 0, 1, 0, 2, 255, 255}
	got := response.GetData()
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v", expect, got)
	}
}

// Function 4
func TestReadInputRegisters(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))
	inputRegisters, err := s.InputRegisters(1)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
		t.FailNow()
	}
	inputRegisters[200] = 1
	inputRegisters[201] = 2
	inputRegisters[202] = 65535
	if err = s.SaveInputRegisters(1, inputRegisters); err != nil {
		t.Errorf("expected nil, got %v\n", err)
		t.FailNow()
	}

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 6
	frame.Device = 255
	frame.Function = 4
	SetDataWithRegisterAndNumber(&frame, 200, 3)

	var req Request
	req.frame = &frame
	response := s.handle(&req)
	exception := GetException(response)
	if exception != Success {
		t.Errorf("expected Success, got %v", exception.String())
		t.FailNow()
	}
	expect := []byte{6, 0, 1, 0, 2, 255, 255}
	got := response.GetData()
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v", expect, got)
	}
}

// Function 5
func TestWriteSingleCoil(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 12
	frame.Device = 255
	frame.Function = 5
	SetDataWithRegisterAndNumber(&frame, 65535, 1024)

	var req Request
	req.frame = &frame
	response := s.handle(&req)
	exception := GetException(response)
	if exception != Success {
		t.Errorf("expected Success, got %v", exception.String())
		t.FailNow()
	}
	expect := 1
	bsGot, err := s.Coils(1)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
		t.FailNow()
	}
	got := bsGot[65535]
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v\n", expect, got)
	}
}

// Function 6
func TestWriteHoldingRegister(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 12
	frame.Device = 255
	frame.Function = 6
	SetDataWithRegisterAndNumber(&frame, 5, 6)

	var req Request
	req.frame = &frame
	response := s.handle(&req)
	exception := GetException(response)
	if exception != Success {
		t.Errorf("expected Success, got %v", exception.String())
		t.FailNow()
	}
	expect := 6
	bsGot, err := s.HoldingRegisters(1)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
		t.FailNow()
	}
	got := bsGot[5]
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v\n", expect, got)
	}
}

// Function 15
func TestWriteMultipleCoils(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 12
	frame.Device = 255
	frame.Function = 15
	SetDataWithRegisterAndNumberAndBytes(&frame, 1, 2, []byte{3})

	var req Request
	req.frame = &frame
	response := s.handle(&req)
	exception := GetException(response)
	if exception != Success {
		t.Errorf("expected Success, got %v", exception.String())
		t.FailNow()
	}
	expect := []byte{1, 1}
	bsGot, err := s.Coils(1)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
		t.FailNow()
	}
	got := bsGot[1:3]
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v\n", expect, got)
	}
}

// Function 16
func TestWriteHoldingRegisters(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 12
	frame.Device = 255
	frame.Function = 16
	SetDataWithRegisterAndNumberAndValues(&frame, 1, 2, []uint16{3, 4})

	var req Request
	req.frame = &frame
	response := s.handle(&req)
	exception := GetException(response)
	if exception != Success {
		t.Errorf("expected Success, got %v", exception.String())
		t.FailNow()
	}
	expect := []uint16{3, 4}
	bsGot, err := s.HoldingRegisters(1)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
		t.FailNow()
	}
	got := bsGot[1:3]
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v\n", expect, got)
	}
}

func TestBytesToUint16(t *testing.T) {
	bytes := []byte{1, 2, 3, 4}
	got := BytesToUint16(bytes)
	expect := []uint16{258, 772}
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v\n", expect, got)
	}
}

func TestUint16ToBytes(t *testing.T) {
	values := []uint16{1, 2, 3}
	got := Uint16ToBytes(values)
	expect := []byte{0, 1, 0, 2, 0, 3}
	if !isEqual(expect, got) {
		t.Errorf("expected %v, got %v\n", expect, got)
	}
}

func TestOutOfBounds(t *testing.T) {
	s := NewServer(NewMemorySlaveUint8(1))

	var frame TCPFrame
	frame.TransactionIdentifier = 1
	frame.ProtocolIdentifier = 0
	frame.Length = 6
	frame.Device = 255

	var req Request
	req.frame = &frame

	// bits
	SetDataWithRegisterAndNumber(&frame, 65535, 2)

	frame.Function = 1
	response := s.handle(&req)
	exception := GetException(response)
	if exception != IllegalDataAddress {
		t.Errorf("expected IllegalDataAddress, got %v", exception.String())
	}

	frame.Function = 2
	response = s.handle(&req)
	exception = GetException(response)
	if exception != IllegalDataAddress {
		t.Errorf("expected IllegalDataAddress, got %v", exception.String())
	}

	SetDataWithRegisterAndNumberAndBytes(&frame, 65535, 2, []byte{3})
	frame.Function = 15
	response = s.handle(&req)
	exception = GetException(response)
	if exception != IllegalDataAddress {
		t.Errorf("expected IllegalDataAddress, got %v", exception.String())
	}

	// registers
	SetDataWithRegisterAndNumber(&frame, 65535, 2)

	frame.Function = 3
	response = s.handle(&req)
	exception = GetException(response)
	if exception != IllegalDataAddress {
		t.Errorf("expected IllegalDataAddress, got %v", exception.String())
	}

	frame.Function = 4
	response = s.handle(&req)
	exception = GetException(response)
	if exception != IllegalDataAddress {
		t.Errorf("expected IllegalDataAddress, got %v", exception.String())
	}

	SetDataWithRegisterAndNumberAndValues(&frame, 65535, 2, []uint16{0, 0})
	frame.Function = 16
	response = s.handle(&req)
	exception = GetException(response)
	if exception != IllegalDataAddress {
		t.Errorf("expected IllegalDataAddress, got %v", exception.String())
	}
}
