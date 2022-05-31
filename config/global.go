package config

type PaginationQuery struct {
	Current  int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
}
type Pagination struct {
	Current  int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
	Total    int64 `json:"total"`
}

type FuncListResult struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

type SystemConfigKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RequestQuery = map[string]interface{}

type SqlPaginationBySqlData = []interface{}
