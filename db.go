package orm

type Db struct {
	registry *registry
}

type DbOpt func(db *Db)

func NewDB(opts ...DbOpt) (*Db, error) {

	db := &Db{
		registry: newRegistry(),
	}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}
