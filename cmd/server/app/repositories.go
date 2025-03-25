package app

import (
	"context"
	"database/sql"
	"errors"
	"go-yandex-practicum-metrics/internal/repositories"
	"go-yandex-practicum-metrics/internal/types"
	"io"
)

type Repos struct {
	DBSaveRepo       SaveRepository
	DBFilterRepo     FilterRepository
	DBListRepo       ListRepository
	FileSaveRepo     SaveRepository
	FileFilterRepo   FilterRepository
	FileListRepo     ListRepository
	MemorySaveRepo   SaveRepository
	MemoryFilterRepo FilterRepository
	MemoryListRepo   ListRepository
}

type FilterRepository interface {
	Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, error)
}

type SaveRepository interface {
	Save(ctx context.Context, metric *types.Metrics) error
}

type ListRepository interface {
	List(ctx context.Context) ([]*types.Metrics, error)
}

type RepoConfig interface {
	GetDatabaseDSN() string
	GetFileStoragePath() string
}

type DBRepo interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type FileRepo interface {
	io.Reader
	io.Seeker
	Encode(v interface{}) error
	Decode(v interface{}) error
}

func NewRepositories(
	config RepoConfig,
	db DBRepo,
	file FileRepo,
	memory map[types.MetricID]*types.Metrics,
) (*Repos, error) {
	var repos Repos

	// Ensure that at least one repository is provided
	if db == nil && file == nil && memory == nil {
		return nil, errors.New("at least one repository must be provided")
	}

	// Proceed with repository initialization
	if db != nil {
		repos.DBSaveRepo = repositories.NewMetricDBSaveRepository(db)
		repos.DBFilterRepo = repositories.NewMetricDBFilterRepository(db)
		repos.DBListRepo = repositories.NewMetricDBListRepository(db)
	}

	if file != nil {
		repos.FileSaveRepo = repositories.NewMetricFileSaveRepository(file)
		repos.FileFilterRepo = repositories.NewMetricFileFilterRepository(file, file)
		repos.FileListRepo = repositories.NewMetricFileListRepository(file, file)
	}

	if db == nil && file == nil && memory != nil {
		repos.MemorySaveRepo = repositories.NewMetricMemorySaveRepository(memory)
		repos.MemoryFilterRepo = repositories.NewMetricMemoryFilterRepository(memory)
		repos.MemoryListRepo = repositories.NewMetricMemoryListRepository(memory)
	}

	return &repos, nil
}

func (r *Repos) GetSaveRepository() SaveRepository {
	if r.DBSaveRepo != nil {
		return r.DBSaveRepo
	} else if r.FileSaveRepo != nil {
		return r.FileSaveRepo
	} else {
		return r.MemorySaveRepo
	}
}

func (r *Repos) GetFilterRepository() FilterRepository {
	if r.DBFilterRepo != nil {
		return r.DBFilterRepo
	} else if r.FileFilterRepo != nil {
		return r.FileFilterRepo
	} else {
		return r.MemoryFilterRepo
	}
}

func (r *Repos) GetListRepository() ListRepository {
	if r.DBListRepo != nil {
		return r.DBListRepo
	} else if r.FileListRepo != nil {
		return r.FileListRepo
	} else {
		return r.MemoryListRepo
	}
}
