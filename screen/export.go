package screen

import (
	"bytes"
	"fmt"
	"github.com/hundredwz/BibTools/bind"
	"github.com/hundredwz/BibTools/cmn"
	"github.com/hundredwz/BibTools/lib"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"strings"
)

type ExportLayout struct {
	engine       map[string]lib.Lib
	resultChan   chan []string
	progressChan chan int
	stopChan     chan bool
	model        *bind.ResultListModel
	clipBoard    *gui.QClipboard
}

func NewExportLayout(cb *gui.QClipboard, engine map[string]lib.Lib) *ExportLayout {
	return &ExportLayout{
		engine:       engine,
		resultChan:   make(chan []string, 10),
		progressChan: make(chan int, 10),
		stopChan:     make(chan bool),
		clipBoard:    cb,
	}
}

func (s *ExportLayout) search(engine lib.Lib, titles string) {

	titleArr := strings.Split(titles, "\n")
	titleLen := len(titleArr)
	for i, v := range titleArr {
		records, err := engine.Search(v)
		data := make([]string, 0)
		if err != nil {
			data = append(data, err.Error(), err.Error())
		} else if len(records) > 0 {
			record := records[0]
			paper := fmt.Sprintf("%s", record.Title)
			data = append(data, paper)
			bibText, err := engine.BibTex(record.ID)
			if err != nil {
				data = append(data, err.Error())
			} else {
				data = append(data, bibText)
			}
		}
		s.progressChan <- 100 * (i + 1) / titleLen
		s.resultChan <- data
	}
	s.stopChan <- true

}

func (s *ExportLayout) MakeWindow() *widgets.QWidget {

	widget := widgets.NewQWidget(nil, 0)

	widget.SetLayout(widgets.NewQHBoxLayout())

	left := widgets.NewQWidget(nil, 0)
	left.SetLayout(widgets.NewQVBoxLayout())

	libBox := widgets.NewQGroupBox2("Library", nil)
	libRow := widgets.NewQHBoxLayout()
	ieeeRadio := widgets.NewQRadioButton2(cmn.IEEE, nil)
	ieeeRadio.SetAutoExclusive(true)
	ieeeRadio.SetChecked(true)
	libRow.Layout().AddWidget(ieeeRadio)
	googleRadio := widgets.NewQRadioButton2(cmn.Google, nil)
	googleRadio.SetAutoExclusive(true)
	libRow.Layout().AddWidget(googleRadio)
	libBox.SetLayout(libRow)
	left.Layout().AddWidget(libBox)

	titleInput := widgets.NewQTextEdit(nil)
	titleInput.SetPlaceholderText("input the paper title")
	left.Layout().AddWidget(titleInput)

	schBtn := widgets.NewQPushButton2("search", nil)
	schBtn.ConnectClicked(func(checked bool) {

		s.model.Remove()

		progressDialog := widgets.NewQProgressDialog(nil, core.Qt__Dialog)
		progressDialog.SetAttribute(core.Qt__WA_DeleteOnClose, true)
		progressDialog.SetCancelButtonText("Cancel")
		progressDialog.SetLabelText("Searching")
		progressDialog.SetRange(0, 100)
		progressDialog.SetAutoReset(false)

		var eng lib.Lib
		if ieeeRadio.IsChecked() {
			eng = s.engine[cmn.IEEE]
		} else {
			eng = s.engine[cmn.Google]
		}
		go s.search(eng, titleInput.ToPlainText())
		go func() {
			for {
				select {
				case <-s.stopChan:
					progressDialog.SetValue(0)
					progressDialog.Hide()
					return
				case data := <-s.resultChan:
					if len(data)!=2 {
						progressDialog.SetValue(0)
						progressDialog.Hide()
						return
					}
					s.model.Add(bind.SchResult{
						Title:   data[0],
						BibText: data[1],
					})
				case progress := <-s.progressChan:
					progressDialog.SetValue(progress)
				}
			}
		}()
	})

	expBtn := widgets.NewQPushButton2("Copy", nil)
	expBtn.ConnectClicked(func(checked bool) {
		var buf bytes.Buffer
		for i := 0; i < s.model.DataLen(); i++ {
			buf.WriteString(s.model.GetData(i).BibText + "\n")
		}
		s.clipBoard.SetText(buf.String(), gui.QClipboard__Clipboard)
	})

	btnRow := widgets.NewQWidget(nil, 0)
	btnRow.SetLayout(widgets.NewQHBoxLayout())
	btnRow.Layout().AddWidget(schBtn)
	btnRow.Layout().AddWidget(expBtn)

	left.Layout().AddWidget(btnRow)
	widget.Layout().AddWidget(left)

	listview := widgets.NewQListView(nil)
	s.model = bind.NewResultListModel(nil)
	listview.SetModel(s.model)
	listview.ConnectDoubleClicked(func(index *core.QModelIndex) {
		s.clipBoard.SetText(s.model.GetData(index.Row()).BibText, gui.QClipboard__Clipboard)
	})
	widget.Layout().AddWidget(listview)
	return widget

}
