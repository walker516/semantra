package ctrl

import (
	"api/internal/core"
	"api/internal/flow/command"
	"api/internal/flow/query"
	"api/pkg/errmsg"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FeatureCtrl struct {
	cmdFlow command.FeatureCommandFlow
	qryFlow query.FeatureQueryFlow
}

func NewFeatureCtrl(cmd command.FeatureCommandFlow, qry query.FeatureQueryFlow) *FeatureCtrl {
	return &FeatureCtrl{cmdFlow: cmd, qryFlow: qry}
}

func (h *FeatureCtrl) GetByID(c echo.Context) error {
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		return errmsg.Respond(c, errmsg.New("ERR_BAD_REQUEST", "Invalid id format"))
	}
	feat, err := h.qryFlow.GetByID(c.Request().Context(), id)
	if err != nil {
		return errmsg.Respond(c, err)
	}
	return c.JSON(http.StatusOK, feat)
}

// POST /features
func (h *FeatureCtrl) Create(c echo.Context) error {
	var req core.Feature
	if err := c.Bind(&req); err != nil {
		return errmsg.Respond(c, errmsg.New("ERR_BAD_REQUEST", "Invalid request body"))
	}
	ctx := c.Request().Context()
	id, err := h.cmdFlow.Create(ctx, &req)
	if err != nil {
		return errmsg.Respond(c, err)
	}
	return c.JSON(http.StatusCreated, map[string]any{"id": id})
}

// PATCH /features/:id
func (h *FeatureCtrl) Update(c echo.Context) error {
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		return errmsg.Respond(c, errmsg.New("ERR_BAD_REQUEST", "Invalid id format"))
	}
	var req core.Feature
	if err := c.Bind(&req); err != nil {
		return errmsg.Respond(c, errmsg.New("ERR_BAD_REQUEST", "Invalid request body"))
	}
	req.ID = id

	ctx := c.Request().Context()
	err = h.cmdFlow.Update(ctx, &req)
	if err != nil {
		return errmsg.Respond(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"updated": true})
}

// DELETE /features/:id
func (h *FeatureCtrl) Delete(c echo.Context) error {
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		return errmsg.Respond(c, errmsg.New("ERR_BAD_REQUEST", "Invalid id format"))
	}
	ctx := c.Request().Context()
	err = h.cmdFlow.Delete(ctx, id)
	if err != nil {
		return errmsg.Respond(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"deleted": true})
}
