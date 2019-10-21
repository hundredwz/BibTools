package screen

import "github.com/therecipe/qt/widgets"

func (mw *MainWin) MakeMenu() {
	mw.mainmenu = widgets.NewQMenuBar(nil)

	filemenu := widgets.NewQMenu2("&File", nil)

	// Options action.
	proxy := filemenu.AddAction("&Proxy")
	proxy.SetMenuRole(widgets.QAction__PreferencesRole)
	proxy.ConnectTriggered(func(checked bool) {
		win := widgets.NewQMainWindow(mw.window, 0)
		widget := widgets.NewQWidget(win, 0)
		widget.SetLayout(widgets.NewQVBoxLayout())
		proxyInput := widgets.NewQTextEdit(nil)
		if mw.proxyUrl == "" {
			proxyInput.SetPlaceholderText("set proxy for the tool to get access to the libraries.\ne.g. socks5://127.0.0.1:1080")
		} else {
			proxyInput.SetText(mw.proxyUrl)
		}
		widget.Layout().AddWidget(proxyInput)
		setBtn := widgets.NewQPushButton2("save", nil)
		setBtn.ConnectClicked(func(checked bool) {
			mw.proxyUrl=proxyInput.ToPlainText()
			for _, v := range mw.engineMap {
				v.SetProxy(proxyInput.ToPlainText())
			}
			win.Hide()
		})
		widget.Layout().AddWidget(setBtn)
		win.SetCentralWidget(widget)
		win.Show()
	})

	// Separator :)
	filemenu.AddSeparator()

	// Exit URTrator.
	exit := filemenu.AddAction("&Exit")
	exit.SetMenuRole(widgets.QAction__QuitRole)
	exit.ConnectTriggered(func(checked bool) {
		mw.app.Quit()
	})

	mw.mainmenu.AddMenu(filemenu)


	aboutmenu := widgets.NewQMenu2("&Help", nil)

	aboutTool := aboutmenu.AddAction("&About This Tool")
	aboutTool.SetMenuRole(widgets.QAction__AboutRole)
	aboutTool.ConnectTriggered(func(checked bool) {
		win := widgets.NewQMainWindow(mw.window, 0)
		widget := widgets.NewQWidget(win, 0)
		widget.SetLayout(widgets.NewQVBoxLayout())
		text := widgets.NewQLabel2("This tool is a small tool to help researchers easily \n search papers and download citations", win, 0)
		widget.Layout().AddWidget(text)
		closeBtn := widgets.NewQPushButton2("close", win)
		closeBtn.ConnectClicked(func(checked bool) {
			win.Hide()
		})
		widget.Layout().AddWidget(closeBtn)
		win.SetCentralWidget(widget)
		win.Show()
	})

	mw.mainmenu.AddMenu(aboutmenu)

	mw.window.SetMenuBar(mw.mainmenu)
}
