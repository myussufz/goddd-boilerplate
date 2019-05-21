package kit

// Paginate :
type Paginate struct {
	Cursor        string
	FilterColumns []string
	Filters       map[string]map[string]interface{}
	Limit         int64
}
