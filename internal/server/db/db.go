package db

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/Snakder/Mon_go/internal/utils"
)

func New() *DB {
	db := new(DB)
	db.Metrics = utils.NewMetricsStorage()
	db.mut = new(sync.Mutex)
	return db
}

type DB struct {
	mut     *sync.Mutex
	Metrics utils.MetricsStorage
}

func (db *DB) Set(t, name, val string) error {
	var m *utils.Metrics
	switch t {
	case "counter":
		d, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		m = utils.NewMetrics(name, t)
		db.mut.Lock()
		if db.Metrics[name] != nil && db.Metrics[name].Delta != nil {
			*m.Delta = d + *db.Metrics[name].Delta
			db.Metrics[name].Value = nil
		} else {
			*m.Delta = d
		}
		db.mut.Unlock()
	case "gauge":
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		m = utils.NewMetrics(name, t)
		*m.Value = v
	default:
		return errors.New("invalid type")
	}
	db.mut.Lock()
	db.Metrics[name] = m
	db.mut.Unlock()
	return nil

}

func (db *DB) Get(t, name string) (utils.SysGather, error) {
	db.mut.Lock()
	defer db.mut.Unlock()
	if m, ok := db.Metrics[name]; ok {
		switch strings.ToLower(t) {
		case "gauge", "counter":
			if m.MType == strings.ToLower(t) {
				return m, nil
			}
		default:
			return nil, errors.New("invalid type")
		}
	}
	return nil, errors.New("unknown metric")

}

func (db *DB) GetAll() map[string]utils.SysGather {
	fullMap := make(map[string]utils.SysGather, len(db.Metrics))
	for name, m := range db.Metrics {
		fullMap[name] = m
	}
	return fullMap
}
