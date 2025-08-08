package gui_elements

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/xChygyNx/gophkeeper/internal/client/cmp/function"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
)

type Tabs struct {
	TblLoginPassword *widget.Table
	TblText          *widget.Table
	TblCard          *widget.Table
	TblBinary        *widget.Table
	TabLoginPassword *container.TabItem
	TabText          *container.TabItem
	TabCard          *container.TabItem
	TabBinary        *container.TabItem
}

func GetTabs() *Tabs {
	return &Tabs{}
}

func (t *Tabs) InitLoginPasswordTable(dataTables *DataTables) {
	t.TblLoginPassword = widget.NewTable(
		func() (int, int) {
			return len(dataTables.DataTblLoginPassword), len(dataTables.DataTblLoginPassword[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(labels.TblLabel)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(dataTables.DataTblLoginPassword[i.Row][i.Col])
		})
	function.SetDefaultColumnsWidthLoginPassword(t.TblLoginPassword)
}

func (t *Tabs) InitTextTable(dataTables *DataTables) {
	t.TblText = widget.NewTable(
		func() (int, int) {
			return len(dataTables.DataTblText), len(dataTables.DataTblText[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(labels.TblLabel)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(dataTables.DataTblText[i.Row][i.Col])
		})
	function.SetDefaultColumnsWidthText(t.TblText)
}

func (t *Tabs) InitCardTable(dataTables *DataTables) {
	t.TblCard = widget.NewTable(
		func() (int, int) {
			return len(dataTables.DataTblCard), len(dataTables.DataTblCard[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(labels.TblLabel)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(dataTables.DataTblCard[i.Row][i.Col])
		})
	function.SetDefaultColumnsWidthCard(t.TblCard)
}

func (t *Tabs) InitBinaryTable(dataTables *DataTables) {
	t.TblBinary = widget.NewTable(
		func() (int, int) {
			return len(dataTables.DataTblBinary), len(dataTables.DataTblBinary[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(labels.TblLabel)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(dataTables.DataTblBinary[i.Row][i.Col])
		})
	function.SetDefaultColumnsWidthBinary(t.TblBinary)
}
