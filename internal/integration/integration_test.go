package integration

import (
	orm "github.com/jrmarcco/easy-orm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"time"
)

type Suite struct {
	suite.Suite

	driver  string
	dsn     string
	db      *orm.DB
	dialect orm.Dialect
}

func (s *Suite) SetupSuite() {
	t := s.T()

	db, err := orm.Open(s.driver, s.dsn, orm.DBWithDialect(s.dialect))
	require.NoError(t, err)

	err = db.Wait(10 * time.Second)
	require.NoError(t, err)

	s.db = db
}
