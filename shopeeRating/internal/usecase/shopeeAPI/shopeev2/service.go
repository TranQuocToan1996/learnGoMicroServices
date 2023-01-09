package shopeev2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TranQuocToan1996/shopeerating/config"
	"github.com/TranQuocToan1996/shopeerating/internal/entity"
)

const (
	version      = "v2"
	ratingFormat = "%v/%v/item/get_ratings?filter=0&flag=1&itemid=%v&limit=%v&offset=%v&shopid=%v&type=0"
)

type V2 struct {
	cfg                   config.Config
	version, ratingFormat string
}

func New(config config.Config) *V2 {
	return &V2{cfg: config, version: version, ratingFormat: ratingFormat}
}

func (api *V2) GetRatings(ctx context.Context, rawURL string) (*entity.RatingResp, error) {

	itemID, shopID := api.getID(rawURL)
	if len(itemID) == 0 || len(shopID) == 0 {
		return nil, errors.New("wrong link")
	}

	total, err := api.getTotal(ctx, itemID, shopID)
	if err != nil {
		return nil, err
	}

	resp := &entity.RatingResp{}
	current := 0
	attemp := 10_000
	for current < total {
		// TODO: Check real running, if can't accept, using fanin fanout
		respElem, err := api.getRating(ctx, itemID, shopID, current)
		if respElem.Error != 0 || respElem.ErrorMsg != nil {
			resp.Error = respElem.Error
			resp.ErrorMsg = respElem.ErrorMsg
		}
		if err != nil {
			resp.Error = -1
			resp.ErrorMsg = err.Error()
		}

		resp.Data.Ratings = append(resp.Data.Ratings, respElem.Data.Ratings...)
		if current == 0 {
			resp.Data.ItemRatingSummary = respElem.Data.ItemRatingSummary
			resp.Data.IsSipItem = respElem.Data.IsSipItem
			resp.Data.RcmdAlgo = respElem.Data.RcmdAlgo
			resp.Data.DowngradeSwitch = respElem.Data.DowngradeSwitch
			resp.Data.HasMore = respElem.Data.HasMore
			resp.Data.ShowLocalReview = respElem.Data.ShowLocalReview
			resp.Data.BrowsingUI = respElem.Data.BrowsingUI
			resp.Data.EnableBuyerGalleryMedia = respElem.Data.EnableBuyerGalleryMedia
			resp.Data.UserLatestRating = respElem.Data.UserLatestRating
			resp.Data.SizeInfoAbt = respElem.Data.SizeInfoAbt
			resp.Data.TopRatings = respElem.Data.TopRatings
			resp.Data.ResizeImageAbt = respElem.Data.ResizeImageAbt
			resp.Data.PurchaseBarAbt = respElem.Data.PurchaseBarAbt
		}
		attemp--
		if attemp < 0 {
			return nil, errors.New("too large query")
		}
		current = len(resp.Data.Ratings)

	}

	return resp, nil
}

// baseURLAndEndpoint.itemID.shopID?querystring
func (api *V2) getID(rawURL string) (itemID, shopID string) {
	if len(rawURL) == 0 ||
		!strings.Contains(rawURL, "?") || !strings.Contains(rawURL, ".") {
		return "", ""
	}
	saparateByQuestion := strings.Split(rawURL, "?")
	saparateByDot := strings.Split(saparateByQuestion[0], ".")
	return saparateByDot[len(saparateByDot)-1], saparateByDot[len(saparateByDot)-2]
}

func (api *V2) getRating(ctx context.Context, itemID, shopID string, offset int) (*entity.RatingResp, error) {
	url := fmt.Sprintf(ratingFormat,
		api.cfg.BaseURL,
		api.version,
		itemID,
		api.cfg.Limit,
		offset,
		shopID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// The shopee api also limit the data resp on each call, so it is safe to readall
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rating := &entity.RatingResp{}
	err = json.Unmarshal(body, rating)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

func (api *V2) getTotal(ctx context.Context, itemID, shopID string) (int, error) {
	const (
		limit  = 1
		offset = 0
	)
	url := fmt.Sprintf(ratingFormat,
		api.cfg.BaseURL, api.version, itemID, limit, offset, shopID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0, err
	}

	// The shopee api also limit the data resp on each call, so it is safe to readall
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	rating := &entity.RatingResp{}
	err = json.Unmarshal(body, rating)
	if err != nil {
		return 0, err
	}

	return rating.Data.ItemRatingSummary.RatingTotal, nil
}
