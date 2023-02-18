package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	videoSize   binding.String
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
	dp.videoSize = binding.NewString()

	title := widget.NewFormItem("Title:", widget.NewLabel(dp.video.Title))
	author := widget.NewFormItem("Author:", widget.NewLabel(dp.video.Author))
	publishDate := widget.NewFormItem("PublishDate:", widget.NewLabel(dp.video.PublishDate.String()))
	duration := widget.NewFormItem("Duration:", widget.NewLabel((dp.video.Duration - time.Second).String()))
	size := widget.NewFormItem("Size:", widget.NewLabelWithData(dp.videoSize))
	selector := widget.NewFormItem("Select:", dp.renderSelector())

	return widget.NewForm(title, author, publishDate, duration, size, selector)
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
		videoFromates = append(videoFromates, "Quality: "+formate.QualityLabel+", FPS: "+fmt.Sprint(formate.FPS))
	}

	dp.selector = widget.NewSelect(videoFromates, func(s string) {})
	dp.selector.OnChanged = func(s string) {
		if dp.downloadBtn.Disabled() {
			dp.downloadBtn.Enable()
		}

		format := dp.getSelectedVideoFormat()
		vidoSize := yt.GetVideoSize(dp.video, &format)

		dp.videoSize.Set(dp.formatVideoSize(vidoSize))
	}

	return dp.selector
}

func (dp *DownloadPage) download() {
	save := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil || uc == nil {
			return
		}

		format := dp.getSelectedVideoFormat()

		progressPercent := yt.Download(dp.video, &format, uc)

		progress := dialog.NewProgress("", "Downloading", dp.window)
		progress.Show()

		for percent := range progressPercent {
			progress.SetValue(percent / 100)
		}

		progress.Hide()

		dialog.ShowInformation("", "Video Downloaded", dp.window)

	}, dp.window)

	save.Show()
}

func (dp *DownloadPage) getSelectedVideoFormat() youtube.Format {
	selectedIndex := dp.selector.SelectedIndex()
	format := dp.video.Formats.WithAudioChannels().Type("video/mp4")[selectedIndex]

	return format
}

func (dp *DownloadPage) formatVideoSize(size int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d size", size)
	}
}
