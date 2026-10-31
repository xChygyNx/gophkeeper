package function

import "github.com/xChygyNx/gophkeeper/internal/client/storage/labels"

type DataTables struct {
	DataTblLoginPassword [][]string
	DataTblText          [][]string
	DataTblCard          [][]string
	DataTblBinary        [][]string
}

func InitDataTables() *DataTables {
	return &DataTables{
		DataTblLoginPassword: [][]string{{labels.NameItem, labels.DescriptionItem, labels.LoginItem, labels.PasswordItem,
			labels.CreatedAtItem, labels.UpdatedAtItem}},
		DataTblText: [][]string{{labels.NameItem, labels.DescriptionItem, labels.DataItem, labels.CreatedAtItem, labels.UpdatedAtItem}},
		DataTblCard: [][]string{{labels.NameItem, labels.DescriptionItem, labels.PaymentSystemItem, labels.NumberItem, labels.HolderItem, labels.CVCItem,
			labels.EndDateItem, labels.CreatedAtItem, labels.UpdatedAtItem}},
		DataTblBinary: [][]string{{labels.NameItem, labels.CreatedAtItem}},
	}
}
