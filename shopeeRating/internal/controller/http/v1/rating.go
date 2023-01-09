package v1

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/TranQuocToan1996/shopeerating/internal/usecase"
	"github.com/TranQuocToan1996/shopeerating/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ratingRoutes struct {
	u usecase.Rating
	l logger.Interface
}

func newRatingRoutes(handler *gin.RouterGroup, u usecase.Rating, l logger.Interface) {
	r := &ratingRoutes{u, l}

	h := handler.Group("/rating")
	{
		h.POST("/csv", r.getOne)
	}
}

func (r *ratingRoutes) getOne(c *gin.Context) {
	req := &struct {
		RawURL string `json:"url"`
	}{}
	err := c.ShouldBindJSON(req)
	if len(req.RawURL) == 0 || err != nil {
		errorResponse(c, http.StatusInternalServerError, "URL empty")

		return
	}

	ratingObj, err := r.u.GetRatings(c.Request.Context(), req.RawURL)
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		errorResponse(c, http.StatusInternalServerError, "some problems")

		return
	}

	records := [][]string{}
	for _, data := range ratingObj.Data.Ratings {
		records = append(records, []string{
			data.AuthorUsername,
			fmt.Sprint(data.RatingStar),
			data.Comment,
		})
	}

	name := time.Now().Unix()

	f, err := os.Create(fmt.Sprintf("%v.csv", name))
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		errorResponse(c, http.StatusInternalServerError, "some problems")
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			r.l.Error(err, "http - v1 - history")
			errorResponse(c, http.StatusInternalServerError, "some problems")
		}
	}
	c.FileAttachment(fmt.Sprintf("./%v.csv", name), fmt.Sprintf("%v.csv", name))
	c.Writer.Header().Set("attachment", "filename=tasks.csv")
}
