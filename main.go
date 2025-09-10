package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"syscall/js"

	id3v2 "github.com/bogem/id3v2"
	"github.com/dhowden/tag"
)

type MP3Metadata struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Year        string `json:"year"`
	Genre       string `json:"genre"`
	Format      string `json:"format"`
	Track       int    `json:"track"`
	TotalTracks int    `json:"total"`
	Cover       []byte `json:"cover"`
}

func updateMetadata(this js.Value, p []js.Value) any {
	data := make([]byte, p[0].Length())
	js.CopyBytesToGo(data, p[0])

	var mp3metadata MP3Metadata
	err := json.Unmarshal([]byte(p[1].String()), &mp3metadata)
	if err != nil {
		panic(err)
	}

	newTag := id3v2.NewEmptyTag()

	newTag.SetTitle(mp3metadata.Title)
	newTag.SetAlbum(mp3metadata.Album)
	newTag.SetYear(mp3metadata.Year)
	newTag.SetGenre(mp3metadata.Genre)
	newTag.SetArtist(mp3metadata.Artist)
	newTag.AddTextFrame("TRCK", id3v2.EncodingUTF8, fmt.Sprintf("%d/%d", mp3metadata.Track, mp3metadata.TotalTracks))

	if len(mp3metadata.Cover) > 0 {
		pic := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/jpeg",
			Picture:     mp3metadata.Cover,
			Description: "Portada",
			PictureType: id3v2.PTFrontCover,
		}
		newTag.AddAttachedPicture(pic)
	}
	// Guardar nueva cabecera + audio original
	var buf bytes.Buffer
	if _, err := newTag.WriteTo(&buf); err != nil {
		panic(err)
	}

	audioStart := newTag.Size()
	if audioStart > len(data) {
		audioStart = len(data)
	}
	buf.Write(data[audioStart:])

	out := js.Global().Get("Uint8Array").New(buf.Len())
	js.CopyBytesToJS(out, buf.Bytes())
	return out
}

func parseMP3(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf("error: no file provided")
	}
	uint8Array := args[0]
	length := uint8Array.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, uint8Array)
	reader := bytes.NewReader(data)
	meta, err := tag.ReadFrom(reader)
	if err != nil {
		return js.ValueOf("error: " + err.Error())
	}
	track, totalTracks := meta.Track()
	result := map[string]interface{}{
		"title":  meta.Title(),
		"artist": meta.Artist(),
		"album":  meta.Album(),
		"year":   meta.Year(),
		"genre":  meta.Genre(),
		"format": meta.Format(),
		"track":  track,
		"total":  totalTracks,
	}
	jsonBytes, _ := json.Marshal(result)
	return js.ValueOf(string(jsonBytes))
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("playSong", js.FuncOf(playSong))
	js.Global().Set("parseMP3", js.FuncOf(parseMP3))
	js.Global().Set("updateMetadata", js.FuncOf(updateMetadata))
	<-c
}

func playSong(this js.Value, args []js.Value) interface{} {
	audio := js.Global().Get("Audio").New("music/song.mp3")
	audio.Call("play")
	return nil
}
