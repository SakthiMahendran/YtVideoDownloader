package yt

import (
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
		var vBuff [1024]byte
		var totReaded float64

		vStream, vSize, _ := ytClient.GetStream(video, format)

		for {
			readed, _ := vStream.Read(vBuff[:])

			if readed == 0 {
				break
			}

			uc.Write(vBuff[:])

			totReaded += float64(readed)
			percent := (totReaded / float64(vSize))

			dloadStatus <- percent
		}

		close(dloadStatus)
		vStream.Close()
		uc.Close()
	}(dloadStatus)

	return dloadStatus
}
