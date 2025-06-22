package ctrl

import (
	"api/pkg/errmsg"
	"net/http"

	"github.com/labstack/echo/v4"

	"api/internal/core/feature"
	"api/internal/flow"
)

type FeatureHandler struct{ uc *flow.FeatureUsecase }

func NewFeatureHandler(uc *flow.FeatureUsecase) *FeatureHandler { return &FeatureHandler{uc: uc} }

func (h *FeatureHandler) GetFeature(c echo.Context) error {
	id := c.Param("id")
	f, err := h.uc.Get(feature.FeatureID(id))
	if err != nil { return errmsg.Respond(c, err) }
	if f == nil { return errmsg.Respond(c, errmsg.New(errmsg.ErrNotFound, "feature not found")) }
	return c.JSON(http.StatusOK, f)
}