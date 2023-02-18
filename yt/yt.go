package yt

import (
	"io"

	"fyne.io/fyne/v2"
	"github.com/kkdai/youtube/v2"
)

var ytClient = youtube.Client{}

func GetVideo(url string) (*youtube.Video, error) {
	return ytClient.GetVideo(url)
}

func Download(video *youtube.Video, format *youtube.Format, uc fyne.URIWriteCloser) <-chan float64 {
	dloadStatus := make(chan float64)

	go func(chan float64) {
		vStream, vSize, _ := ytClient.GetStream(video, format)
		vReader := NewVideoReader(vStream)

		go io.Copy(uc, &vReader)

		for {
			status := vReader.total / float64(vSize) * 100
			dloadStatus <- status
			if status == 100 {
				break
			}
		}

		close(dloadStatus)
		vStream.Close()
		uc.Close()
	}(dloadStatus)

	return dloadStatus
}
