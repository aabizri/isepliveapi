package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"gopkg.in/h2non/filetype.v1"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// Decode attachments and place them in the temporary directory
func (pubreq *postPublicationsRequest) decodeBase64Attachments() (pictures []string, videos []string, audio []string, documents []string, err error) {
	if len(pubreq.PicturesBase64) != 0 {
		pictures, err = decodeBase64Attachment("PIC", pubreq.PicturesBase64)
		if err != nil {
			return
		}
	}

	if len(pubreq.VideosBase64) != 0 {
		videos, err = decodeBase64Attachment("VID", pubreq.VideosBase64)
		if err != nil {
			return
		}
	}

	if len(pubreq.AudioBase64) != 0 {
		audio, err = decodeBase64Attachment("AUD", pubreq.AudioBase64)
		if err != nil {
			return
		}
	}

	if len(pubreq.DocumentsBase64) != 0 {
		documents, err = decodeBase64Attachment("DOC", pubreq.DocumentsBase64)
	}
	return
}

// Decode a particular attachment type
func decodeBase64Attachment(prefix string, attachments []string) ([]string, error) {
	var paths []string = make([]string, len(attachments))
	for index, encoded := range attachments {
		// Decode the string
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return paths, err
		}

		// Detect MIME type & extension
		kind, err := filetype.Match(decoded)
		if err != nil {
			return paths, err
		}

		// Make suffix (REPLACE WITH HASH ?)
		suffix := strconv.FormatInt(int64(index+len(encoded)), 10)

		// Build filename
		filename := prefix + suffix + "." + kind.Extension
		path := filepath.Join(tempDirPath, filename)

		// Create file
		file, err := os.Create(path)
		if err != nil {
			return paths, err
		}
		defer file.Close()

		// Make buffer, decode & write to file
		decodedReader := bytes.NewBuffer(decoded)
		n, err := io.Copy(file, decodedReader)
		if err == nil && n == 0 {
			err = errors.New("ERROR: 0 bytes copied")
		}
		if err != nil {
			return paths, err
		}

		// Detect file type

		paths[index] = path
	}
	return paths, nil
}
