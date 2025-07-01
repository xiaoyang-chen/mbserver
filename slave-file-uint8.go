package mbserver

import (
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
)

var _ Slaver = new(fileSlaveUint8)

type fileSlaveUint8 struct {
	slaveNum     uint8
	slaveLock    []sync.RWMutex
	fileStoreDir string
}

// will create slaveNum slaves, slave id is [1, slaveNum], slaveNumMax is 255, slaveNumMin is 1, if slaveNum > 255, it will use 255, if slaveNum < 1, it will use 1, one uint8, 2^8-1, 255; if fileStoreDir is "", will use "./file-slave"
func NewFileSlaveUint8(slaveNum uint8, fileStoreDir string) (slaver Slaver) {

	if slaveNum < 1 {
		slaveNum = 1
	} else if slaveNum > 255 {
		slaveNum = 255
	}
	if fileStoreDir == "" {
		fileStoreDir = "./file-slave"
	}
	slaver = &fileSlaveUint8{
		slaveNum:     slaveNum,
		slaveLock:    make([]sync.RWMutex, slaveNum),
		fileStoreDir: fileStoreDir,
	}
	return
}

func (s *fileSlaveUint8) IsSlaveIdValid(id uint8) bool { return id > 0 && id <= s.slaveNum }

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *fileSlaveUint8) DiscreteInputs(id uint8) (bs []byte, err error) {

	id = s.getRealId(id)
	s.slaveLock[id].RLock()
	bs, err = s.localStorageFileRead(fmt.Sprintf("%s/%d-discreteInputs", s.fileStoreDir, id+1))
	s.slaveLock[id].RUnlock()
	if err == nil && len(bs) < 65536 {
		var newBs = make([]byte, 65536)
		copy(newBs, bs)
		bs = newBs
	}
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *fileSlaveUint8) Coils(id uint8) (bs []byte, err error) {

	id = s.getRealId(id)
	s.slaveLock[id].RLock()
	bs, err = s.localStorageFileRead(fmt.Sprintf("%s/%d-coils", s.fileStoreDir, id+1))
	s.slaveLock[id].RUnlock()
	if err == nil && len(bs) < 65536 {
		var newBs = make([]byte, 65536)
		copy(newBs, bs)
		bs = newBs
	}
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *fileSlaveUint8) HoldingRegisters(id uint8) (bs []uint16, err error) {

	id = s.getRealId(id)
	var bsFileContent []byte
	s.slaveLock[id].RLock()
	bsFileContent, err = s.localStorageFileRead(fmt.Sprintf("%s/%d-holdingRegisters", s.fileStoreDir, id+1))
	s.slaveLock[id].RUnlock()
	if err == nil {
		if bs = BytesToUint16(bsFileContent); len(bs) < 65536 {
			var newBs = make([]uint16, 65536)
			copy(newBs, bs)
			bs = newBs
		}
	}
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *fileSlaveUint8) InputRegisters(id uint8) (bs []uint16, err error) {

	id = s.getRealId(id)
	var bsFileContent []byte
	s.slaveLock[id].RLock()
	bsFileContent, err = s.localStorageFileRead(fmt.Sprintf("%s/%d-inputRegisters", s.fileStoreDir, id+1))
	s.slaveLock[id].RUnlock()
	if err == nil {
		if bs = BytesToUint16(bsFileContent); len(bs) < 65536 {
			var newBs = make([]uint16, 65536)
			copy(newBs, bs)
			bs = newBs
		}
	}
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *fileSlaveUint8) SaveDiscreteInputs(id uint8, b []byte) (err error) {

	id = s.getRealId(id)
	s.slaveLock[id].Lock()
	_, err = s.localStorageWrite(s.fileStoreDir, fmt.Sprintf("%s/%d-discreteInputs", s.fileStoreDir, id+1), b)
	s.slaveLock[id].Unlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *fileSlaveUint8) SaveCoils(id uint8, b []byte) (err error) {

	id = s.getRealId(id)
	s.slaveLock[id].Lock()
	_, err = s.localStorageWrite(s.fileStoreDir, fmt.Sprintf("%s/%d-coils", s.fileStoreDir, id+1), b)
	s.slaveLock[id].Unlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *fileSlaveUint8) SaveHoldingRegisters(id uint8, b []uint16) (err error) {

	id = s.getRealId(id)
	s.slaveLock[id].Lock()
	_, err = s.localStorageWrite(s.fileStoreDir, fmt.Sprintf("%s/%d-holdingRegisters", s.fileStoreDir, id+1), Uint16ToBytes(b))
	s.slaveLock[id].Unlock()
	return
}

// if id not in [1, s.slaveNum], id > s.slaveNum => id will use slave ${slaveNum}, id < 1 => id will use slave 1
func (s *fileSlaveUint8) SaveInputRegisters(id uint8, b []uint16) (err error) {

	id = s.getRealId(id)
	s.slaveLock[id].Lock()
	_, err = s.localStorageWrite(s.fileStoreDir, fmt.Sprintf("%s/%d-inputRegisters", s.fileStoreDir, id+1), Uint16ToBytes(b))
	s.slaveLock[id].Unlock()
	return
}

func (s *fileSlaveUint8) getRealId(id uint8) (realId uint8) {

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

func (s *fileSlaveUint8) localStorageFileRead(filePath string) (bsFileContent []byte, err error) {

	if bsFileContent, err = os.ReadFile(filePath); err == nil {
		if lenBsFileContent := len(bsFileContent); lenBsFileContent > 0 {
			var newBs = make([]byte, hex.DecodedLen(lenBsFileContent))
			if _, err = hex.Decode(newBs, bsFileContent); err == nil {
				bsFileContent = newBs
				return
			}
			err = errors.Wrap(err, "hex decode file content fail")
		}
		return
	}
	if os.IsNotExist(err) {
		bsFileContent, err = nil, nil
		return
	}
	err = errors.Wrap(err, "read file fail")
	return
}

func (s *fileSlaveUint8) localStorageWrite(fileDir, filePath string, bsFileContent []byte) (n int, err error) {

	// mkdir all dir
	if err = os.MkdirAll(fileDir, 0755); err != nil {
		err = errors.Wrap(err, "mkdir all fail")
		return
	}
	// open/create file
	var file *os.File
	if file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		err = errors.Wrap(err, "open file fail")
		return
	}
	defer file.Close()
	// write
	var encodeBs = make([]byte, hex.EncodedLen(len(bsFileContent)))
	hex.Encode(encodeBs, bsFileContent)
	if n, err = file.Write(encodeBs); err != nil {
		err = errors.Wrap(err, "write file content fail")
	}
	return
}
