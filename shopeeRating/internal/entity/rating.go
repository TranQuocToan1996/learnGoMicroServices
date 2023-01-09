// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type RatingResp struct {
	Error    int         `json:"error"`
	ErrorMsg interface{} `json:"error_msg"`
	Data     Data        `json:"data"`
}

type Videos struct {
	Duration     int         `json:"duration"`
	UploadTime   interface{} `json:"upload_time"`
	ID           string      `json:"id"`
	Cover        string      `json:"cover"`
	URL          string      `json:"url"`
	MmsExt       string      `json:"mms_ext"`
	CoverImageID string      `json:"cover_image_id"`
}
type ProductItems struct {
	Shopid     int      `json:"shopid"`
	IsSnapshot int      `json:"is_snapshot"`
	Itemid     int64    `json:"itemid"`
	Snapshotid int64    `json:"snapshotid"`
	Modelid    int64    `json:"modelid"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	ModelName  string   `json:"model_name"`
	Options    []string `json:"options"`
}
type ItemRatingReply struct {
	IsHidden   bool        `json:"is_hidden"`
	Ctime      int         `json:"ctime"`
	Userid     int         `json:"userid"`
	Mtime      int         `json:"mtime"`
	Orderid    int64       `json:"orderid"`
	Cmtid      int64       `json:"cmtid"`
	Itemid     interface{} `json:"itemid"`
	Rating     interface{} `json:"rating"`
	Shopid     interface{} `json:"shopid"`
	RatingStar interface{} `json:"rating_star"`
	Status     interface{} `json:"status"`
	Editable   interface{} `json:"editable"`
	Opt        interface{} `json:"opt"`
	Filter     interface{} `json:"filter"`
	Mentioned  interface{} `json:"mentioned"`
	// Comment    string      `json:"comment"`
}
type DetailedRating struct {
	ProductQuality  int         `json:"product_quality"`
	SellerService   interface{} `json:"seller_service"`
	DeliveryService interface{} `json:"delivery_service"`
}
type SipInfo struct {
	IsOversea    bool   `json:"is_oversea"`
	OriginRegion string `json:"origin_region"`
}
type ImageData struct {
	ImageID      string `json:"image_id"`
	CoverImageID string `json:"cover_image_id"`
}
type KeyMedia struct {
	KeyMediaCode int    `json:"key_media_code"`
	KeyMediaID   string `json:"key_media_id"`
}
type Ratings struct {
	Anonymous                    bool            `json:"anonymous"`
	ShowReply                    bool            `json:"show_reply"`
	Liked                        bool            `json:"liked"`
	SyncToSocial                 bool            `json:"sync_to_social"`
	ExcludeScoringDueLowLogistic bool            `json:"exclude_scoring_due_low_logistic"`
	HasTemplateTag               bool            `json:"has_template_tag"`
	IsRepeatedPurchase           bool            `json:"is_repeated_purchase"`
	IsNormalItem                 bool            `json:"is_normal_item"`
	IsHidden                     bool            `json:"is_hidden"`
	Ctime                        int             `json:"ctime"`
	Rating                       int             `json:"rating"`
	Userid                       int             `json:"userid"`
	Shopid                       int             `json:"shopid"`
	RatingStar                   int             `json:"rating_star"`
	Status                       int             `json:"status"`
	Mtime                        int             `json:"mtime"`
	Editable                     int             `json:"editable"`
	Opt                          int             `json:"opt"`
	Filter                       int             `json:"filter"`
	AuthorShopid                 int             `json:"author_shopid"`
	EditableDate                 int             `json:"editable_date"`
	LikeCount                    int             `json:"like_count"`
	OverallFit                   int             `json:"overall_fit"`
	Orderid                      int64           `json:"orderid"`
	Itemid                       int64           `json:"itemid"`
	Cmtid                        int64           `json:"cmtid"`
	Comment                      string          `json:"comment"`
	AuthorUsername               string          `json:"author_username"`
	AuthorPortrait               string          `json:"author_portrait"`
	SizeInfoAbt                  string          `json:"size_info_abt"`
	Images                       []string        `json:"images"`
	Videos                       []Videos        `json:"videos"`
	ProductItems                 []ProductItems  `json:"product_items"`
	Mentioned                    []interface{}   `json:"mentioned"`
	DeleteReason                 interface{}     `json:"delete_reason"`
	DeleteOperator               interface{}     `json:"delete_operator"`
	Tags                         interface{}     `json:"tags"`
	DetailedRating               DetailedRating  `json:"detailed_rating"`
	SipInfo                      SipInfo         `json:"sip_info"`
	LoyaltyInfo                  interface{}     `json:"loyalty_info"`
	SyncToSocialToggle           interface{}     `json:"sync_to_social_toggle"`
	DisplayVariationFilter       interface{}     `json:"display_variation_filter"`
	Viewed                       interface{}     `json:"viewed"`
	ShowView                     interface{}     `json:"show_view"`
	SyncToSocialDetail           interface{}     `json:"sync_to_social_detail"`
	Profile                      interface{}     `json:"profile"`
	SizeInfoTags                 interface{}     `json:"size_info_tags"`
	TemplateTags                 []interface{}   `json:"template_tags"`
	ImageData                    []ImageData     `json:"image_data"`
	KeyMedia                     KeyMedia        `json:"key_media"`
	ItemRatingReply              ItemRatingReply `json:"ItemRatingReply"`
}
type ItemRatingSummary struct {
	RatingTotal           int   `json:"rating_total"`
	RcountWithContext     int   `json:"rcount_with_context"`
	RcountWithImage       int   `json:"rcount_with_image"`
	RcountWithMedia       int   `json:"rcount_with_media"`
	RcountLocalReview     int   `json:"rcount_local_review"`
	RcountRepeatPurchase  int   `json:"rcount_repeat_purchase"`
	RcountOverallFitSmall int   `json:"rcount_overall_fit_small"`
	RcountOverallFitFit   int   `json:"rcount_overall_fit_fit"`
	RcountOverallFitLarge int   `json:"rcount_overall_fit_large"`
	RcountOverseaReview   int   `json:"rcount_oversea_review"`
	RcountFolded          int   `json:"rcount_folded"`
	RatingCount           []int `json:"rating_count"`
}
type Data struct {
	IsSipItem               bool              `json:"is_sip_item"`
	DowngradeSwitch         bool              `json:"downgrade_switch"`
	HasMore                 bool              `json:"has_more"`
	ShowLocalReview         bool              `json:"show_local_review"`
	EnableBuyerGalleryMedia bool              `json:"enable_buyer_gallery_media"`
	ResizeImageAbt          bool              `json:"resize_image_abt"`
	RcmdAlgo                string            `json:"rcmd_algo"`
	BrowsingUI              string            `json:"browsing_ui"`
	PurchaseBarAbt          string            `json:"purchase_bar_abt"`
	SizeInfoAbt             string            `json:"size_info_abt"`
	UserLatestRating        interface{}       `json:"user_latest_rating"`
	TopRatings              []interface{}     `json:"top_ratings"`
	Ratings                 []Ratings         `json:"ratings"`
	ItemRatingSummary       ItemRatingSummary `json:"item_rating_summary"`
}
