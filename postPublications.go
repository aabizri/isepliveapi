package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/aabizri/goil.v0"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// A PublishRequest sent to POST /publications
type postPublicationsRequest struct {
	Message  string `json:"message"`
	Category uint8  `json:"category"`
	Group    uint8  `json:"group"`
	Official bool   `json:"official,omitempty"`
	Private  bool   `json:"private"`
	Dislike  bool   `json:"dislike"`

	// Attachments directly in the request(string instead of []byte because encoding/json is retarded, see https://stackoverflow.com/questions/31449610/illegal-base64-data-error-message)
	PicturesBase64  []string `json:"picturesBase64,omitempty"`
	VideosBase64    []string `json:"videosBase64,omitempty"`
	AudioBase64     []string `json:"audioBase64,omitempty"`
	DocumentsBase64 []string `json:"documentsBase64,omitempty"`

	// Attachments by URL
	PicturesURI  []string `json:"picturesURI"`
	VideosURI    []string `json:"videosURI"`
	AudioURI     []string `json:"audioURI"`
	DocumentsURI []string `json:"documentsURI"`

	// Event & Survey
	Event  EventReq
	Survey SurveyReq
}

type EventReq struct {
	Title  string    `json:"title"`
	Starts time.Time `json:"starts"`
	Ends   time.Time `json:"ends"`
}

func (e EventReq) Populated() bool {
	return (e.Title != "")
}

type SurveyReq struct {
	Question string    `json:"question"`
	Ends     time.Time `json:"ends"`
	Answers  []string  `json:"answers"`
	Multiple bool      `json:"multiple,omitempty"`
}

func (s SurveyReq) Populated() bool {
	return (s.Question != "")
}

// Publish
func postPublications(w http.ResponseWriter, req *http.Request) {
	var pubreq postPublicationsRequest

	// Read the body, limiting
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1024*1024*1024*3)) // 3 GiB Limit
	if err != nil {
		err := NewError(InternalErrorReadingRequestFailed, "", fmt.Sprintf("Error: %s", err)).JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	// Close the body
	if err := req.Body.Close(); err != nil {
		err := NewError(InternalError, "Closing request body failed", fmt.Sprintf("Error: %s", err.Error())).JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	// Unmarshal
	err = json.Unmarshal(body, &pubreq)
	if err != nil {
		err := NewError(InternalErrorJSONUnmarshallingFailed, "", fmt.Sprintf("Error: %s", err.Error())).JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	//fmt.Printf("%#v\n",pubreq)

	// Retrieve goil session
	raw := req.Context().Value("session")
	gs, ok := raw.(*goil.Session)
	if !ok {
		err := NewError(InternalErrorTypeAssertionFailed, "", "Type assertion failure while casting req.Context().Value(\"session\") to (*goil.Session)").JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	// Publish
	err = pubreq.publish(gs)
	if err != nil {
		err := NewError(InternalError, "While calling pubreq.publish()", fmt.Sprintf("Error: %s", err.Error())).JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}
}

// Publish a publication request
func (pubreq *postPublicationsRequest) publish(s *goil.Session) error {

	// Prepare the post
	post := goil.NewPost(pubreq.Message, goil.Category(pubreq.Category), pubreq.Private, pubreq.Dislike)
	if pubreq.Group != 0 {
		post.PostAs(goil.Group(pubreq.Group), pubreq.Official)
	}

	// Attach attachments
	err := pubreq.attach(post)
	if err != nil {
		return err
	}

	// Publish
	err = s.PublishPost(post)

	return err
}

func (pubreq *postPublicationsRequest) attach(post *goil.Post) error {
	// Prepare slices holding paths to files
	var (
		pics  []string
		vids  []string
		audio []string
		docs  []string
	)

	var err error
	// Download base64 attachments
	if len(pubreq.PicturesBase64)+len(pubreq.VideosBase64)+len(pubreq.AudioBase64)+len(pubreq.DocumentsBase64) != 0 {
		// Download attachments to the temp directory
		pics, vids, audio, docs, err = pubreq.decodeBase64Attachments()
		if err != nil {
			return err
		}
	}

	// Download URI attachments
	if len(pubreq.PicturesURI)+len(pubreq.VideosURI)+len(pubreq.AudioURI)+len(pubreq.DocumentsURI) != 0 {
		// Create variables for storage of the paths coming from URI attachments
		var (
			picsURI  []string
			vidsURI  []string
			audioURI []string
			docsURI  []string
		)

		// Download attachments to the temp directory
		picsURI, vidsURI, audioURI, docsURI, err = pubreq.downloadURIAttachments()
		if err != nil {
			return err
		}

		// Append to the main slice
		pics = append(pics, picsURI...)
		vids = append(vids, vidsURI...)
		audio = append(audio, audioURI...)
		docs = append(docs, docsURI...)
	}

	// Add them to request
	for _, path := range pics {
		post.AttachPhoto(path)
	}
	for _, path := range vids {
		post.AttachVideo(path)
	}
	for _, path := range audio {
		post.AttachAudio(path)
	}
	for _, path := range docs {
		post.AttachDocument(path)
	}

	return nil
}
