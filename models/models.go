package models

// QueryOpts ...
type QueryOpts struct {
	Page  int
	Limit int
}

// GetOffset returns the page offset required for query
func (db *QueryOpts) GetOffset() int {
	if db.Page < 0 {
		db.Page = 0
	}

	return db.Page * db.GetLimit()
}

// GetLimit sets the default and max limit
func (db *QueryOpts) GetLimit() int {
	if db.Limit < 1 {
		db.Limit = 1
	}
	if db.Limit > 100 {
		db.Limit = 100
	}
	return db.Limit
}
