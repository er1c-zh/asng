package api

import (
	"asng/models"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type FileManager struct {
	ctx  context.Context
	cwd  string
	rwMu sync.RWMutex
}

func NewFileManager(ctx context.Context) (*FileManager, error) {
	var err error
	const dir = "ashareNg"

	fm := &FileManager{}
	fm.ctx = ctx

	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		runtime.LogFatalf(ctx, "Failed to get user home directory: %v", err)
		return nil, err
	}

	fm.cwd = filepath.Join(homeDirPath, dir)

	_, err = os.Stat(fm.cwd)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(fm.cwd, 0755)
			if err != nil {
				runtime.LogFatalf(ctx, "Failed to create directory %s: %v", fm.cwd, err)
				return nil, err
			}
		} else {
			runtime.LogFatalf(ctx, "Failed to stat directory %s: %v", fm.cwd, err)
			return nil, err
		}
	}

	return fm, nil
}

type FileMeta struct {
	Version   uint8
	Type      FileType
	UpdatedAt int64
}

type FileType uint8

func (ft FileType) String() string {
	switch ft {
	case TypeStockMeta:
		return "stock_meta"
	default:
		return "unknown"
	}
}

func (ft FileType) fileName() string {
	switch ft {
	case TypeStockMeta:
		return "stock_meta.json"
	default:
		return "unknown"
	}
}

const (
	TypeStockMeta FileType = 1
)

func SaveFile[T any](fm *FileManager, ft FileType, data T) error {
	fm.rwMu.Lock()
	defer fm.rwMu.Unlock()
	name := ft.fileName()
	f, err := os.Create(filepath.Join(fm.cwd, name))
	if err != nil {
		runtime.LogErrorf(fm.ctx, "Failed to create %s: %v", name, err)
		return err
	}
	defer f.Close()

	fileMeta := FileMeta{
		Version:   1,
		Type:      ft,
		UpdatedAt: time.Now().Unix(),
	}

	err = binary.Write(f, binary.LittleEndian, &fileMeta)
	if err != nil {
		runtime.LogErrorf(fm.ctx, "Failed to write %s: %v", name, err)
		return err
	}

	d, err := json.Marshal(data)
	if err != nil {
		runtime.LogErrorf(fm.ctx, "Failed to marshal %s: %v", name, err)
		return err
	}

	_, err = f.Write(d)
	if err != nil {
		runtime.LogErrorf(fm.ctx, "Failed to write %s: %v", name, err)
		return err
	}

	return nil
}

func ReadFile[T any](fm *FileManager, ft FileType, dt T) (*FileMeta, T, error) {
	fm.rwMu.RLock()
	defer fm.rwMu.RUnlock()
	name := ft.fileName()
	f, err := os.Open(filepath.Join(fm.cwd, name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, dt, nil
		}
		runtime.LogErrorf(fm.ctx, "Failed to open %s: %v", name, err)
		return nil, dt, err
	}
	defer f.Close()

	var fileMeta FileMeta

	err = binary.Read(f, binary.LittleEndian, &fileMeta)
	if err != nil {
		runtime.LogErrorf(fm.ctx, "Failed to read %s: %v", name, err)
		return nil, dt, err
	}

	if fileMeta.Type != ft {
		runtime.LogErrorf(fm.ctx, "Invalid file type: %d", fileMeta.Type)
		return nil, dt, err
	}

	d, err := io.ReadAll(f)
	if err != nil {
		runtime.LogErrorf(fm.ctx, "Failed to read %s: %v", name, err)
		return nil, dt, err
	}
	err = json.Unmarshal(d, &dt)
	if err != nil {
		runtime.LogErrorf(fm.ctx, "Failed to unmarshal %s: %v", name, err)
		return nil, dt, err
	}
	return &fileMeta, dt, nil
}

func (fm *FileManager) LoadStockMeta() (*FileMeta, *models.StockMetaAll, error) {
	h, d, err := ReadFile(fm, TypeStockMeta, models.StockMetaAll{})
	if err != nil {
		return nil, nil, err
	}
	if h == nil {
		return nil, nil, nil
	}
	return h, &d, nil
}

func (fm *FileManager) SaveStockMeta(d *models.StockMetaAll) error {
	return SaveFile(fm, TypeStockMeta, d)
}
