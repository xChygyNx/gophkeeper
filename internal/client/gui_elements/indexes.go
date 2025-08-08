package gui_elements

type Indexes struct {
	IndexTblLoginPassword       int
	SelectedRowTblLoginPassword []string
	IndexTblText                int
	SelectedRowTblText          []string
	IndexTblCard                int
	SelectedRowTblCard          []string
	IndexTblBinary              int
	SelectedRowTblBinary        []string
}

func InitIndexes() *Indexes {
	return &Indexes{}
}
