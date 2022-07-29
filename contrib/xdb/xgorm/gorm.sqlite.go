//go:build grom.sqlite
// +build grom.sqlite

package xgorm

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zhiyunliu/glue/config"
	contribxdb "github.com/zhiyunliu/glue/contrib/xdb"
	"github.com/zhiyunliu/glue/xdb"
	"gorm.io/driver/sqlite"
)

const Proto = "grom.sqlite"

func init() {
	xdb.Register(&mssqlResolver{})
	callbackCache[Proto] = sqlite.Open
}

type mssqlResolver struct {
}

func (s *mssqlResolver) Name() string {
	return Proto
}

func (s *mssqlResolver) Resolve(setting config.Config) (interface{}, error) {
	cfg := &contribxdb.Config{}
	err := setting.Scan(cfg)
	if err != nil {
		return nil, fmt.Errorf("读取DB配置:%w", err)
	}
	gromDB, err := buildGormDB(Proto, cfg)
	if err != nil {
		return nil, err
	}
	return &dbWrap{
		gromDB: gromDB,
	}, nil
}
