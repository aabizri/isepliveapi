package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Decode attachments and place them in the temporary directory
func (pubreq *postPublicationsRequest) downloadURIAttachments() (pictures []string, videos []string, audio []string, documents []string, err error) {
	if len(pubreq.PicturesURI) != 0 {
		pictures, err = downloadURIAttachments("PIC", pubreq.PicturesURI)
		if err != nil {
			return
		}
	}

	if len(pubreq.VideosURI) != 0 {
		videos, err = downloadURIAttachments("VID", pubreq.VideosURI)
		if err != nil {
			return
		}
	}

	if len(pubreq.AudioURI) != 0 {
		audio, err = downloadURIAttachments("AUD", pubreq.AudioURI)
		if err != nil {
			return
		}
	}

	if len(pubreq.DocumentsURI) != 0 {
		documents, err = downloadURIAttachments("DOC", pubreq.DocumentsURI)
	}
	return
}

// Download a particular attachment type
func downloadURIAttachments(prefix string, attachmentsURI []string) ([]string, error) {
	var paths []string = make([]string, len(attachmentsURI))
	for index, uriStr := range attachmentsURI {
		// Download file
		resp, err := http.Get(uriStr)
		if err != nil {
			return paths, err
		}
		defer resp.Body.Close()

		// Transform the uri into a net/url form
		url, err := url.Parse(uriStr)
		if err != nil {
			return paths, err
		}

		// Build filename
		filename := filepath.Base(url.Path)
		path := filepath.Join(tempDirPath, filename)

		// Create file
		file, err := os.Create(path)
		if err != nil {
			return paths, err
		}
		defer file.Close()

		// Make buffer, decode & write to file
		io.Copy(file, resp.Body)

		// Detect file type

		paths[index] = path
	}
	return paths, nil
}
