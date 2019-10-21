package bind

import (
	"fmt"
	"github.com/therecipe/qt/core"
)

type SchResult struct {
	Title   string
	BibText string
}

type ResultListModel struct {
	core.QAbstractListModel

	_ func()               `constructor:"init"`
	_ func()               `signal:"remove,auto"`
	_ func(item SchResult) `signal:"add,auto"`

	modelData []SchResult
}

func (m *ResultListModel) init() {
	m.modelData = []SchResult{}

	m.ConnectRowCount(m.rowCount)
	m.ConnectData(m.data)
}

func (m *ResultListModel) rowCount(*core.QModelIndex) int {
	return len(m.modelData)
}

func (m *ResultListModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if role != int(core.Qt__DisplayRole) {
		return core.NewQVariant()
	}

	item := m.modelData[index.Row()]
	return core.NewQVariant1(fmt.Sprintf("%v\n\n\n%v", item.Title, item.BibText))
}

func (m *ResultListModel) remove() {
	if len(m.modelData) == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), 0, len(m.modelData)-1)
	m.modelData = m.modelData[:0]
	m.EndRemoveRows()
}

func (m *ResultListModel) add(item SchResult) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	m.modelData = append(m.modelData, item)
	m.EndInsertRows()
}

func (m *ResultListModel) GetData(index int) SchResult {
	if index > len(m.modelData) {
		return SchResult{}
	}
	return m.modelData[index]
}
func (m *ResultListModel) DataLen() int {
	return len(m.modelData)
}
