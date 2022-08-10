package main

import "github.com/SakthiMahendran/YtDownloader/ui"

func main() {
	mainFrame := ui.NewMainFrame("YtDownloader")

	mainFrame.RenderUiPage(&ui.IndexPage{})
	mainFrame.ShowAndRun()
}
