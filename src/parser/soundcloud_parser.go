package parser

import (
	"encoding/json"
	"fmt"
	models "github.com/AlienStream/Shared-Go/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const SOUNDCLOUD_CLIENT_ID = "ff43d208510d35ce49ed972b01f116ab"

type SoundcloudUser struct {
	Permalink   string `json:"permalink_url"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Thumbnail   string `json:"avatar_url"`
}

type SoundcloudTrack struct {
	Title        string         `json:"title"`
	Permalink    string         `json:"permalink_url"`
	Thumbnail    string         `json:"artwork_url"`
	Submitted_by string         `json:"username"`
	Created_at   string         `json:"created_at"`
	Comments     int            `json:"comment_count"`
	Favorites    int            `json:"favoritings_count"`
	User         SoundcloudUser `json:"user"`
}

func getSoundcloudChannelData(source_data DataObject) DataObject {
	fmt.Printf("Updating %s \n", source_data.Source.Title)
	info := getRawSoundcloudChannelMeta(source_data.Source.Url)
	source_data.Source.Title = info.Username
	source_data.Source.Description = info.Description
	source_data.Source.Thumbnail = info.Thumbnail

	// TODO: Multithread this into a queueable worker that respects the reddit limits
	raw_posts := getRawSoundcloudTracks(source_data.Source.Url)

	for _, raw_post := range raw_posts {
		post := models.Post{
			Id:                 0,
			Source_id:          source_data.Source.Id,
			Title:              raw_post.Title,
			Number_of_comments: raw_post.Comments,
			Permalink:          raw_post.Permalink,
			Thumbnail:          raw_post.Thumbnail,
			Embed_url:          raw_post.Permalink,
			Likes:              raw_post.Favorites,
			Dislikes:           0,
			Submitter:          raw_post.User.Username,
		}
		post.Posted_at, _ = time.Parse("2006/01/02 15:04:05 -0700", raw_post.Created_at)

		source_data.Posts = append(source_data.Posts, post)
	}

	return source_data
}

func getRawSoundcloudChannelMeta(source_url string) SoundcloudUser {
	client := &http.Client{}
	request_url := strings.Replace(source_url, "://soundcloud.com/", "://api.soundcloud.com/users/", -1)
	request_url += "?client_id=" + SOUNDCLOUD_CLIENT_ID
	request, _ := http.NewRequest("GET", request_url, nil)
	request.Header.Set("User-Agent", "AlienStream Master Server v. 2.0")

	response, request_err := client.Do(request)
	if request_err != nil {
		panic(request_err)
	}

	defer response.Body.Close()
	data := SoundcloudUser{}
	temp, _ := ioutil.ReadAll(response.Body)

	parse_err := json.Unmarshal(temp, &data)

	if parse_err != nil {
		panic("Requester failed to fetch from source")
	}

	return data
}

func getRawSoundcloudTracks(source_url string) []SoundcloudTrack {
	client := &http.Client{}
	request_url := strings.Replace(source_url, "://soundcloud.com/", "://api.soundcloud.com/users/", -1)
	request_url += "/tracks?client_id=" + SOUNDCLOUD_CLIENT_ID
	request, _ := http.NewRequest("GET", request_url, nil)
	request.Header.Set("User-Agent", "AlienStream Master Server v. 2.0")

	response, request_err := client.Do(request)
	if request_err != nil {
		panic(request_err)
	}

	defer response.Body.Close()
	data := []SoundcloudTrack{}
	temp, _ := ioutil.ReadAll(response.Body)

	parse_err := json.Unmarshal(temp, &data)

	if parse_err != nil {
		panic("Requester failed to fetch from source")
	}

	return data
}

func getSoundcloudPlaylistData(source_data DataObject) DataObject {
	fmt.Printf("Updating %s \n", source_data.Source.Title)
	return DataObject{}
}
