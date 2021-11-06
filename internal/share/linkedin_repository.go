package share

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dorneanu/gomation/internal/entity"
)

// POST https://api.linkedin.com/v2/ugcPosts
// Authorization: Bearer {{(verb-var token)}}
// Accept: application/json
// X-Restli-Protocol-Version: 2.0.0
// Content-Type: application/json; charset=utf-8

// {
//     "author": "urn:li:person:GJ74B4A98G",
//     "lifecycleState": "PUBLISHED",
//     "specificContent": {
//         "com.linkedin.ugc.ShareContent": {
//             "shareCommentary": {
//                 "text": "Testing"
//             },
//             "shareMediaCategory": "NONE"
//         }
//     },
//     "visibility": {
//         "com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC"
//     }
// }

const (
	// API URL for User Generated Content (UGC)
	linkedinUGCAPI = "https://api.linkedin.com/v2/ugcPosts"
)

// LinkedinUGCShareMedia describes the media to be shared
//
// Check out https://docs.microsoft.com/en-us/linkedin/marketing/integrations/community-management/shares/ugc-post-api?tabs=http#sharemedia
type LinkedinUGCShareMedia struct {
	Status      string `json:"status"`
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	OriginalURL string `json:"originalUrl"`
	Title       struct {
		Text string `json:"text"`
	} `json:"title"`
}

// LinkedinUGCShareContent defines meta data of content to be shared
//
// Also check https://docs.microsoft.com/en-us/linkedin/marketing/integrations/community-management/shares/ugc-post-api?tabs=http#sharecontent
type LinkedinUGCShareContent struct {
	ShareCommentary struct {
		Text string `json:"text"`
	} `json:"shareCommentary"`
	ShareMediaCategory string                  `json:"shareMediaCategory"`
	Media              []LinkedinUGCShareMedia `json:"media"`
}

// LinkedinUGCPost defines the schema for a UGC post
//
// Also check https://docs.microsoft.com/en-us/linkedin/marketing/integrations/community-management/shares/ugc-post-api?tabs=http#schema
type LinkedinUGCSharePost struct {
	Author          string                  `json:"author"`
	LifecycleState  string                  `json:"lifecycleState"`
	SpecificContent LinkedinUGCShareContent `json:"specificContent"`
	Visibility      struct {
		MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
	} `json:"visibility"`
}

// LinkedinShareRepository implements share.Repository
type LinkedinShareRepository struct {
	identity entity.Identity
	client   *http.Client
}

func NewLinkedinShareRepository(identity entity.Identity) *LinkedinShareRepository {
	return &LinkedinShareRepository{
		identity: identity,
		client:   &http.Client{},
	}
}

func (l *LinkedinShareRepository) createNewPost(article entity.ArticleShare) *LinkedinUGCSharePost {
	// Create share content information
	shareContent := LinkedinUGCShareContent{}
	shareContent.ShareCommentary.Text = "some text"
	shareContent.ShareMediaCategory = "ARTICLE"
	shareContent.Media = []LinkedinUGCShareMedia{
		LinkedinUGCShareMedia{
			Status: "READY",
			Description: struct {
				Text string "json:\"text\""
			}{"alles klar"},
			OriginalURL: "http://google.de",
			Title: struct {
				Text string "json:\"text\""
			}{"alles klar"},
		},
	}

	// Create UGC share post
	sharePost := LinkedinUGCSharePost{}
	sharePost.Author = l.identity.ID
	sharePost.LifecycleState = "PUBLISHED"
	sharePost.SpecificContent = shareContent
	sharePost.Visibility = struct {
		MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
	}{"PUBLIC"}

	return &sharePost
}

func (l *LinkedinShareRepository) ShareArticle(ctx context.Context, article entity.ArticleShare) error {
	ugcPost := l.createNewPost(article)

	// Marshalize ugcPost
	jsonStr, err := json.Marshal(ugcPost)
	if err != nil {
		return fmt.Errorf("Couldn't marshalize ugcPost: %s\n", err)
	}

	// Create new HTTP request
	req, err := http.NewRequest("POST", linkedinUGCAPI, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", l.identity.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	// Send request
	resp, err := l.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
