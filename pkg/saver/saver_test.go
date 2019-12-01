package saver

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/matryer/is"
)

var tweetId int64 = 1234567890

func init() {
	os.Remove("./output/1234567890.json")
}

func TestTweetSaverSaveJson(t *testing.T) {
	is := is.New(t)

	tweet := twitter.Tweet{
		ID:        tweetId,
		CreatedAt: "2019-01-01 01:01:01",
		Text:      "This is a test tweet. @helloworld http://www.example.com",
	}

	ts := TweetSaver{
		SaveDir:  "./output",
		SaveJson: true,
	}

	err := ts.Save(tweet)
	is.NoErr(err)

	fbytes, _ := ioutil.ReadFile("./fixtures/1234567890.json")
	obytes, _ := ioutil.ReadFile("./output/1234567890.json")
	is.True(bytes.Equal(fbytes, obytes))
}

func TestTweetSaverSaveMedia(t *testing.T) {
	is := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./fixtures/omg.gif")
	}))
	defer ts.Close()

	url := ts.URL + "/omg.gif" // Add filename to URL.
	tweet := twitter.Tweet{
		ID:        tweetId,
		CreatedAt: "2019-01-01 01:01:01",
		Text:      "This is a test tweet. @helloworld http://www.example.com",
		ExtendedEntities: &twitter.ExtendedEntity{
			Media: []twitter.MediaEntity{
				twitter.MediaEntity{
					MediaURLHttps: url,
					Type:          "animated_gif",
					VideoInfo: twitter.VideoInfo{
						Variants: []twitter.VideoVariant{
							twitter.VideoVariant{
								URL: url,
							},
						},
					},
				},
			},
		},
	}

	saver := TweetSaver{
		SaveDir:   "./output",
		SaveMedia: true,
	}

	err := saver.Save(tweet)
	is.NoErr(err)
}