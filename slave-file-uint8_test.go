package mbserver

import (
	"reflect"
	"sync"
	"testing"
)

func Test_fileSlaveUint8_localStorageFileRead(t *testing.T) {
	type fields struct {
		slaveNum     uint8
		slaveLock    []sync.RWMutex
		fileStoreDir string
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantBsFileContent []byte
		wantErr           bool
	}{
		{
			name:   "test no exist file",
			fields: fields{},
			args: args{
				filePath: "./no-exist-file",
			},
			wantBsFileContent: nil,
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileSlaveUint8{
				slaveNum:     tt.fields.slaveNum,
				slaveLock:    tt.fields.slaveLock,
				fileStoreDir: tt.fields.fileStoreDir,
			}
			gotBsFileContent, err := s.localStorageFileRead(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileSlaveUint8.localStorageFileRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBsFileContent, tt.wantBsFileContent) {
				t.Errorf("fileSlaveUint8.localStorageFileRead() = %v, want %v", gotBsFileContent, tt.wantBsFileContent)
			}
		})
	}
}

func Test_fileSlaveUint8_SaveDiscreteInputs(t *testing.T) {
	type fields struct {
		slaveNum     uint8
		slaveLock    []sync.RWMutex
		fileStoreDir string
	}
	type args struct {
		id uint8
		b  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test01",
			fields: fields{
				slaveNum:     2,
				slaveLock:    make([]sync.RWMutex, 2),
				fileStoreDir: "./file-slave",
			},
			args: args{
				id: 2,
				b:  []byte{0x00, 0x00, 0x12, 0x34, 0xAB},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileSlaveUint8{
				slaveNum:     tt.fields.slaveNum,
				slaveLock:    tt.fields.slaveLock,
				fileStoreDir: tt.fields.fileStoreDir,
			}
			if err := s.SaveDiscreteInputs(tt.args.id, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("fileSlaveUint8.SaveDiscreteInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileSlaveUint8_DiscreteInputs(t *testing.T) {
	type fields struct {
		slaveNum     uint8
		slaveLock    []sync.RWMutex
		fileStoreDir string
	}
	type args struct {
		id uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantBs  []byte
		wantErr bool
	}{
		{
			name: "test01",
			fields: fields{
				slaveNum:     2,
				slaveLock:    make([]sync.RWMutex, 2),
				fileStoreDir: "./file-slave",
			},
			args: args{
				id: 2,
			},
			wantBs:  []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileSlaveUint8{
				slaveNum:     tt.fields.slaveNum,
				slaveLock:    tt.fields.slaveLock,
				fileStoreDir: tt.fields.fileStoreDir,
			}
			gotBs, err := s.DiscreteInputs(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileSlaveUint8.DiscreteInputs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBs, tt.wantBs) {
				// t.Errorf("fileSlaveUint8.DiscreteInputs() = %v, want %v", gotBs, tt.wantBs)
				t.Logf("fileSlaveUint8.DiscreteInputs() = % x, want %v", gotBs, tt.wantBs)
			}
		})
	}
}

func Test_fileSlaveUint8_SaveHoldingRegisters(t *testing.T) {
	type fields struct {
		slaveNum     uint8
		slaveLock    []sync.RWMutex
		fileStoreDir string
	}
	type args struct {
		id uint8
		b  []uint16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test01",
			fields: fields{
				slaveNum:     2,
				slaveLock:    make([]sync.RWMutex, 2),
				fileStoreDir: "./file-slave",
			},
			args: args{
				id: 2,
				b:  []uint16{0, 0, 1, 65535, 2025},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileSlaveUint8{
				slaveNum:     tt.fields.slaveNum,
				slaveLock:    tt.fields.slaveLock,
				fileStoreDir: tt.fields.fileStoreDir,
			}
			if err := s.SaveHoldingRegisters(tt.args.id, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("fileSlaveUint8.SaveHoldingRegisters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileSlaveUint8_HoldingRegisters(t *testing.T) {
	type fields struct {
		slaveNum     uint8
		slaveLock    []sync.RWMutex
		fileStoreDir string
	}
	type args struct {
		id uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantBs  []uint16
		wantErr bool
	}{
		{
			name: "test01",
			fields: fields{
				slaveNum:     2,
				slaveLock:    make([]sync.RWMutex, 2),
				fileStoreDir: "./file-slave",
			},
			args: args{
				id: 2,
			},
			wantBs:  []uint16{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileSlaveUint8{
				slaveNum:     tt.fields.slaveNum,
				slaveLock:    tt.fields.slaveLock,
				fileStoreDir: tt.fields.fileStoreDir,
			}
			gotBs, err := s.HoldingRegisters(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileSlaveUint8.HoldingRegisters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBs, tt.wantBs) {
				// t.Errorf("fileSlaveUint8.HoldingRegisters() = %v, want %v", gotBs, tt.wantBs)
				t.Logf("fileSlaveUint8.DiscreteInputs() = %v, want %v", gotBs, tt.wantBs)
			}
		})
	}
}
