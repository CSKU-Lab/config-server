package requests

type GetPagination struct {
	Page      int
	PageSize  int
	SortOrder string
	Search    string
}
