package mbserver

import "encoding/binary"

// ReadCoils function 1, reads coils from internal memory.
func ReadCoils(s *Server, frame Framer) ([]byte, *Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	if endRegister > 65535 {
		return []byte{}, &IllegalDataAddress
	}
	dataSize := numRegs / 8
	if (numRegs % 8) != 0 {
		dataSize++
	}
	data := make([]byte, 1+dataSize)
	data[0] = byte(dataSize)

	coils := s.Coils(frame.Addr())
	for i, value := range coils[register:endRegister] {
		if value != 0 {
			shift := uint(i) % 8
			data[1+i/8] |= byte(1 << shift)
		}
	}
	return data, &Success
}

// ReadDiscreteInputs function 2, reads discrete inputs from internal memory.
func ReadDiscreteInputs(s *Server, frame Framer) ([]byte, *Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	if endRegister > 65535 {
		return []byte{}, &IllegalDataAddress
	}
	dataSize := numRegs / 8
	if (numRegs % 8) != 0 {
		dataSize++
	}
	data := make([]byte, 1+dataSize)
	data[0] = byte(dataSize)

	discreteInputs := s.DiscreteInputs(frame.Addr())
	for i, value := range discreteInputs[register:endRegister] {
		if value != 0 {
			shift := uint(i) % 8
			data[1+i/8] |= byte(1 << shift)
		}
	}
	return data, &Success
}

// ReadHoldingRegisters function 3, reads holding registers from internal memory.
func ReadHoldingRegisters(s *Server, frame Framer) ([]byte, *Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	if endRegister > 65536 {
		return []byte{}, &IllegalDataAddress
	}

	holdingRegisters := s.HoldingRegisters(frame.Addr())
	return append([]byte{byte(numRegs * 2)}, Uint16ToBytes(holdingRegisters[register:endRegister])...), &Success
}

// ReadInputRegisters function 4, reads input registers from internal memory.
func ReadInputRegisters(s *Server, frame Framer) ([]byte, *Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	if endRegister > 65536 {
		return []byte{}, &IllegalDataAddress
	}

	inputRegisters := s.InputRegisters(frame.Addr())
	return append([]byte{byte(numRegs * 2)}, Uint16ToBytes(inputRegisters[register:endRegister])...), &Success
}

// WriteSingleCoil function 5, write a coil to internal memory.
func WriteSingleCoil(s *Server, frame Framer) ([]byte, *Exception) {
	register, value := registerAddressAndValue(frame)
	// TODO Should we use 0 for off and 65,280 (FF00 in hexadecimal) for on?
	if value != 0 {
		value = 1
	}

	coils := s.Coils(frame.Addr())
	coils[register] = byte(value)
	s.SaveCoils(frame.Addr(), coils)

	return frame.GetData()[0:4], &Success
}

// WriteHoldingRegister function 6, write a holding register to internal memory.
func WriteHoldingRegister(s *Server, frame Framer) ([]byte, *Exception) {
	register, value := registerAddressAndValue(frame)

	holdingRegisters := s.HoldingRegisters(frame.Addr())
	holdingRegisters[register] = value
	s.SaveHoldingRegisters(frame.Addr(), holdingRegisters)

	return frame.GetData()[0:4], &Success
}

// WriteMultipleCoils function 15, writes holding registers to internal memory.
func WriteMultipleCoils(s *Server, frame Framer) ([]byte, *Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	valueBytes := frame.GetData()[5:]

	if endRegister > 65536 {
		return []byte{}, &IllegalDataAddress
	}

	// TODO This is not correct, bits and bytes do not always align
	//if len(valueBytes)/2 != numRegs {
	//	return []byte{}, &IllegalDataAddress
	//}

	coils := s.Coils(frame.Addr())
	bitCount := 0
	for i, value := range valueBytes {
		for bitPos := uint(0); bitPos < 8; bitPos++ {
			coils[register+(i*8)+int(bitPos)] = bitAtPosition(value, bitPos)
			bitCount++
			if bitCount >= numRegs {
				break
			}
		}
		if bitCount >= numRegs {
			break
		}
	}

	s.SaveCoils(frame.Addr(), coils)

	return frame.GetData()[0:4], &Success
}

// WriteHoldingRegisters function 16, writes holding registers to internal memory.
func WriteHoldingRegisters(s *Server, frame Framer) ([]byte, *Exception) {
	register, numRegs, _ := registerAddressAndNumber(frame)
	valueBytes := frame.GetData()[5:]
	var exception *Exception
	var data []byte

	if len(valueBytes)/2 != numRegs {
		exception = &IllegalDataAddress
	}

	holdingRegisters := s.HoldingRegisters(frame.Addr())

	// Copy data to memroy
	values := BytesToUint16(valueBytes)
	valuesUpdated := copy(holdingRegisters[register:], values)
	if valuesUpdated == numRegs {
		exception = &Success
		data = frame.GetData()[0:4]
	} else {
		exception = &IllegalDataAddress
	}

	s.SaveHoldingRegisters(frame.Addr(), holdingRegisters)

	return data, exception
}

// BytesToUint16 converts a big endian array of bytes to an array of unit16s
func BytesToUint16(bytes []byte) []uint16 {
	values := make([]uint16, len(bytes)/2)

	for i := range values {
		values[i] = binary.BigEndian.Uint16(bytes[i*2 : (i+1)*2])
	}
	return values
}

// Uint16ToBytes converts an array of uint16s to a big endian array of bytes
func Uint16ToBytes(values []uint16) []byte {
	bytes := make([]byte, len(values)*2)

	for i, value := range values {
		binary.BigEndian.PutUint16(bytes[i*2:(i+1)*2], value)
	}
	return bytes
}

func bitAtPosition(value uint8, pos uint) uint8 {
	return (value >> pos) & 0x01
}

func SlaveOperate(fn func(*Server, Framer) ([]byte, *Exception)) func(*Server, Framer) ([]byte, *Exception) {
	return func(s *Server, f Framer) ([]byte, *Exception) {
		if !s.IsSlaveIdValid(f.Addr()) {
			return f.GetData(), &GatewayPathUnavailable
		}
		return fn(s, f)
	}
}
