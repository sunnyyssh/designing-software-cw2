package rest

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/errs"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/model"
)

type AnalysisService interface {
	Analyze(ctx context.Context, id int64) (*model.AnalysisResult, error)
}

type AnalysisHandler struct {
	svc AnalysisService
}

func NewAnalysisService(svc AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{
		svc: svc,
	}
}

func (h *AnalysisHandler) Analyze(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(400, "you should pass int64 id")
		return
	}

	res, err := h.svc.Analyze(c.Request.Context(), id)
	if err != nil {
		sendErr(c, err)
		return
	}

	c.JSON(200, res)
}

func sendErr(c *gin.Context, err error) {
	var httpErr errs.HTTPError
	if errors.As(err, &httpErr) {
		c.String(httpErr.Code, httpErr.Message)
		return
	}

	c.String(500, "internal server error: %s", err)
}
