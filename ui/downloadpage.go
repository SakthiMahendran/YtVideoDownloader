package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/SakthiMahendran/YtDownloader/yt"
	"github.com/kkdai/youtube/v2"
)

type DownloadPage struct {
	size        fyne.Size
	name        string
	newWindow   bool
	downloadBtn *widget.Button
	video       *youtube.Video
	selector    *widget.Select
	window      fyne.Window
	mainFrame   *MainFrame
}

func (dp *DownloadPage) pageSize() fyne.Size {
	return dp.size
}

func (dp *DownloadPage) pageName() string {
	return dp.name
}

func (dp *DownloadPage) isNewWindow() bool {
	return dp.newWindow
}

func (dp *DownloadPage) setMainFrame(mf *MainFrame) {
	dp.mainFrame = mf
}

func (dp *DownloadPage) setWindow(w fyne.Window) {
	dp.window = w
}

func (dp *DownloadPage) renderUi() fyne.CanvasObject {
	dp.size = fyne.NewSize(400, 500)
	dp.name = "Download"
	dp.newWindow = true

	return container.NewVSplit(
		dp.renderImage(),
		container.NewVBox(
			dp.renderForm(),
			dp.renderControlButtons(),
		),
	)
}

func (dp *DownloadPage) renderImage() *canvas.Image {
	res, _ := fyne.LoadResourceFromURLString(dp.video.Thumbnails[len(dp.video.Thumbnails)-1].URL)
	return canvas.NewImageFromResource(res)
}

func (dp *DownloadPage) renderForm() *widget.Form {
	title := widget.NewFormItem("Title:", widget.NewLabel(dp.video.Title))
	author := widget.NewFormItem("Author:", widget.NewLabel(dp.video.Author))
	publishDate := widget.NewFormItem("PublishDate:", widget.NewLabel(dp.video.PublishDate.String()))
	duration := widget.NewFormItem("Duration:", widget.NewLabel(dp.video.Duration.String()))
	selector := widget.NewFormItem("Select:", dp.renderSelector())

	return widget.NewForm(title, author, publishDate, duration, selector)
}

func (dp *DownloadPage) renderControlButtons() *fyne.Container {
	dp.downloadBtn = widget.NewButton("Download", func() {
		go dp.download()
	})

	dp.downloadBtn.Disable()

	cancelBtn := widget.NewButton("Cancel", func() {
		dp.window.Close()
	})

	dp.downloadBtn.SetIcon(theme.DownloadIcon())
	cancelBtn.SetIcon(theme.CancelIcon())

	return container.NewCenter(
		container.NewHBox(cancelBtn, dp.downloadBtn),
	)
}

func (dp *DownloadPage) renderSelector() *widget.Select {
	var videoFromates []string

	for _, formate := range dp.video.Formats.WithAudioChannels().Type("video/mp4") {
		videoFromates = append(videoFromates, "Quality: "+formate.QualityLabel+", FPS: "+strconv.Itoa(formate.FPS))
	}

	dp.selector = widget.NewSelect(videoFromates, func(s string) {})
	dp.selector.OnChanged = func(s string) {
		dp.downloadBtn.Enable()
	}

	return dp.selector
}

func (dp *DownloadPage) download() {
	save := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil || uc == nil {
			return
		}

		fileIndex := dp.selector.SelectedIndex()
		format := dp.video.Formats.WithAudioChannels().Type("video/mp4")[fileIndex]

		progressPercent := yt.Download(dp.video, &format, uc)

		progress := dialog.NewProgress("", "Downloading", dp.window)
		progress.Show()

		for percent := range progressPercent {
			progress.SetValue(percent / 100)
		}

		progress.Hide()

		dialog.ShowConfirm("", "Video Downloaded", func(b bool) {}, dp.window)

	}, dp.window)

	save.Show()
}
