package manga

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Item struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	AltTitle    string `json:"alt_title"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	Year        string `json:"year"`
}

type kitsuResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			CanonicalTitle string `json:"canonicalTitle"`
			Description    string `json:"description"`
			Synopsis       string `json:"synopsis"`
			StartDate      string `json:"startDate"`
			Titles         struct {
				EnJP string `json:"en_jp"`
				JaJP string `json:"ja_jp"`
				EnUS string `json:"en_us"`
			} `json:"titles"`
			PosterImage struct {
				Medium   string `json:"medium"`
				Large    string `json:"large"`
				Original string `json:"original"`
			} `json:"posterImage"`
		} `json:"attributes"`
	} `json:"data"`
}

var httpClient = &http.Client{
	Timeout: 7 * time.Second,
}

// Search mencari manga berdasarkan nama melalui Kitsu API publik
func Search(query string) ([]Item, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []Item{}, nil
	}

	apiURL := fmt.Sprintf("https://kitsu.io/api/edge/manga?filter[text]=%s&page[limit]=6", url.QueryEscape(query))
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("User-Agent", "ComicReaderTV/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi API manga: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API manga merespon status: %s", resp.Status)
	}

	var parsed kitsuResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("gagal parse data manga: %w", err)
	}

	results := make([]Item, 0, len(parsed.Data))
	for _, d := range parsed.Data {
		title := d.Attributes.CanonicalTitle
		if title == "" {
			title = d.Attributes.Titles.EnJP
		}
		if title == "" {
			title = d.Attributes.Titles.EnUS
		}

		alt := d.Attributes.Titles.JaJP
		if alt == "" && d.Attributes.Titles.EnUS != title {
			alt = d.Attributes.Titles.EnUS
		}

		desc := d.Attributes.Synopsis
		if desc == "" {
			desc = d.Attributes.Description
		}
		desc = strings.TrimSpace(desc)
		if len(desc) > 350 {
			desc = desc[:347] + "..."
		}

		cover := d.Attributes.PosterImage.Large
		if cover == "" {
			cover = d.Attributes.PosterImage.Medium
		}
		if cover == "" {
			cover = d.Attributes.PosterImage.Original
		}

		year := ""
		if len(d.Attributes.StartDate) >= 4 {
			year = d.Attributes.StartDate[:4]
		}

		results = append(results, Item{
			ID:          d.ID,
			Title:       title,
			AltTitle:    alt,
			Description: desc,
			CoverURL:    cover,
			Year:        year,
		})
	}

	return results, nil
}
