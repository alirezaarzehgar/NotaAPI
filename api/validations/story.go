package validations

import (
	"net/url"
	"time"

	"github.com/Asrez/NotaAPI/models"
)

func GetWrongStoryField(story models.Story) string {
	if story.Type != models.STORY_TYPE_NORMAL && story.Type != models.STORY_TYPE_EXPLORE {
		return "type"
	}

	if story.Type == models.STORY_TYPE_NORMAL {
		if time.Until(story.To) < 0 || story.To.Sub(story.From) <= 0 {
			return "to"
		}
		if _, err := url.ParseRequestURI(story.AttachedWebpage); err != nil {
			return "attached_webpage"
		}
	}

	assets := map[string]string{
		"final_image":            story.FinalImageUrl,
		"background_url":         story.BackgroundUrl,
		"main_background_url":    story.MainBackgroundUrl,
		"cropped_background_url": story.CroppedBackgroundUrl,
		"attached_file_url":      story.AttachedFileUrl,
		"logo_url":               story.LogoUrl,
	}
	for field, asset := range assets {
		if !IsValidAsset(asset) {
			return field
		}
	}

	return ""
}
