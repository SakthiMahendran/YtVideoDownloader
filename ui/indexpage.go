package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/SakthiMahendran/YtDownloader/yt"
)

type IndexPage struct {
	size      fyne.Size
	name      string
	newWindow bool
	window    fyne.Window
	mainFrame *MainFrame
	urlEntry  *widget.Entry
	findBtn   *widget.Button
}

func (ip *IndexPage) pageSize() fyne.Size {
	return ip.size
}

func (ip *IndexPage) pageName() string {
	return ip.name
}

func (ip *IndexPage) isNewWindow() bool {
	return ip.newWindow
}

func (ip *IndexPage) setMainFrame(mf *MainFrame) {
	ip.mainFrame = mf
}

func (ip *IndexPage) setWindow(w fyne.Window) {
	ip.window = w
}

func (ip *IndexPage) renderUi() fyne.CanvasObject {
	ip.size = fyne.NewSize(300, 300)
	ip.name = "Main Window"
	ip.newWindow = false

	ip.urlEntry = ip.renderUrlEntry()
	ip.findBtn = ip.renderFindButton()

	return container.NewVBox(
		ip.urlEntry,
		container.NewCenter(ip.findBtn),
	)
}

func (ip *IndexPage) renderUrlEntry() *widget.Entry {
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter the video URL")

	urlEntry.OnChanged = func(s string) {
		go func(s string) {
			s = strings.TrimSpace(s)

			if s != "" && ip.findBtn.Disabled() {
				ip.findBtn.Enable()
			} else if s == "" && !ip.findBtn.Disabled() {
				ip.findBtn.Disable()
			}
		}(s)
	}

	return urlEntry
}

func (ip *IndexPage) renderFindButton() *widget.Button {
	findBtn := widget.NewButton("Find", func() {
		go func() {
			video, err := yt.GetVideo(ip.urlEntry.Text)

			if err != nil {
				ip.urlEntry.SetText("")
				return
			}

			dp := &DownloadPage{}
			dp.video = video

			ip.mainFrame.RenderUiPage(dp)
		}()
	})

	findBtn.SetIcon(theme.SearchIcon())
	findBtn.Disable()

	return findBtn
}
