package rest

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunnyyssh/designing-software-cw2/file-storage/internal/errs"
	"github.com/sunnyyssh/designing-software-cw2/file-storage/internal/model"
)

type FileService interface {
	Upload(context.Context, string) (*model.FileMeta, error)
	Get(context.Context, int64) (*model.FileMeta, error)
	ListByHash(context.Context, model.MD5) ([]model.FileMeta, error)
}

type FileHandler struct {
	svc FileService
}

func NewFileHandler(svc FileService) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	if ty := c.ContentType(); !strings.HasPrefix(ty, "text/") {
		c.String(400, "invalid content type. We are waiting text/* from you. You sent %s to us. What for?", ty)
		return
	}

	defer c.Request.Body.Close()
	text, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(400, "we cannot read what you've sent: %s", err)
		return
	}

	meta, err := h.svc.Upload(c.Request.Context(), string(text))
	if err != nil {
		sendErr(c, err)
		return
	}

	c.JSON(200, meta)
}

func (h *FileHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(400, "you should pass int64 id")
		return
	}

	meta, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		sendErr(c, err)
		return
	}

	c.JSON(200, meta)
}

func (h *FileHandler) ListByHash(c *gin.Context) {
	hashRaw := c.Query("hash")
	if hashRaw == "" {
		c.String(400, "you should pass base64-encoded md5 hash")
		return
	}
	hash, err := url.PathUnescape(hashRaw)
	if err != nil {
		c.String(400, "path escaping troubles: %s", err)
		return
	}

	metas, err := h.svc.ListByHash(c.Request.Context(), model.MD5(hash))
	if err != nil {
		sendErr(c, err)
		return
	}

	c.JSON(200, metas)
}

func sendErr(c *gin.Context, err error) {
	var httpErr errs.HTTPError
	if errors.As(err, &httpErr) {
		c.String(httpErr.Code, httpErr.Message)
		return
	}

	c.String(500, "internal server error: %s", err)
}
