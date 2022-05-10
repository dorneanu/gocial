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
	Author          string `json:"author"`
	LifecycleState  string `json:"lifecycleState"`
	SpecificContent struct {
		ShareContent LinkedinUGCShareContent `json:"com.linkedin.ugc.ShareContent"`
	} `json:"specificContent"`
	Visibility struct {
		MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
	} `json:"visibility"`
}

// LinkedinShareRepository implements share.Repository
type LinkedinShareRepository struct {
	identity entity.IdentityProvider
	client   *http.Client
}

func NewLinkedinShareRepository(identity entity.IdentityProvider) *LinkedinShareRepository {
	return &LinkedinShareRepository{
		identity: identity,
		client:   &http.Client{},
	}
}

func (l *LinkedinShareRepository) createNewPost(article entity.ArticleShare) *LinkedinUGCSharePost {
	// Create share content information
	shareContent := LinkedinUGCShareContent{}
	shareContent.ShareCommentary.Text = article.Comment
	shareContent.ShareMediaCategory = "ARTICLE"
	shareContent.Media = []LinkedinUGCShareMedia{
		LinkedinUGCShareMedia{
			Status: "READY",
			Description: struct {
				Text string "json:\"text\""
			}{article.Title},
			OriginalURL: article.URL,
			Title: struct {
				Text string "json:\"text\""
			}{article.Title},
		},
	}

	// Create UGC share post
	sharePost := LinkedinUGCSharePost{}
	sharePost.Author = fmt.Sprintf("urn:li:person:%s", l.identity.UserID)
	sharePost.LifecycleState = "PUBLISHED"
	sharePost.SpecificContent.ShareContent = shareContent
	sharePost.Visibility = struct {
		MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
	}{"PUBLIC"}

	return &sharePost
}

func (l *LinkedinShareRepository) ShareArticle(ctx context.Context, article entity.ArticleShare) error {
	ugcPost := l.createNewPost(article)

	// Marshalize ugcPost
	jsonStr, err := json.MarshalIndent(ugcPost, "", "  ")
	if err != nil {
		return fmt.Errorf("Couldn't marshalize ugcPost: %s\n", err)
	}

	// Create new HTTP request
	req, err := http.NewRequest("POST", linkedinUGCAPI, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", l.identity.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	fmt.Printf("%v\n", req.Header)

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
