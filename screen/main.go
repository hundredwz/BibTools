package screen

import (
	"github.com/hundredwz/BibTools/cmn"
	"github.com/hundredwz/BibTools/lib"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
)

type MainWin struct {
	app       *widgets.QApplication
	window    *widgets.QMainWindow
	clipBoard *gui.QClipboard
	mainmenu  *widgets.QMenuBar
	engineMap map[string]lib.Lib
	proxyUrl  string
}

func NewMainWin() *MainWin {
	mw := &MainWin{
		app:      widgets.NewQApplication(len(os.Args), os.Args),
		window:   widgets.NewQMainWindow(nil, 0),
		proxyUrl: "",
	}
	mw.MakeMenu()
	mw.engineMap = map[string]lib.Lib{
		cmn.IEEE:   lib.NewIEEE(mw.proxyUrl),
		cmn.Google: lib.NewGoogle(mw.proxyUrl),
	}
	return mw
}

func (mw *MainWin) Show() {
	mw.clipBoard = mw.app.Clipboard()
	mw.window.SetMinimumSize2(250, 200)
	mw.window.SetWindowTitle("Bib Tools")

	schLayout := NewSearchLayout(mw.clipBoard, mw.engineMap)
	schWidget := schLayout.MakeWindow()

	expLayout := NewExportLayout(mw.clipBoard, mw.engineMap)
	expWidget := expLayout.MakeWindow()

	tab := widgets.NewQTabWidget(nil)
	tab.AddTab(schWidget, "Search")
	tab.AddTab(expWidget, "Export")
	mw.window.SetCentralWidget(tab)

	mw.window.Show()

	mw.app.Exec()

}
