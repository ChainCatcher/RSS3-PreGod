package handler

import (
	"errors"
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GetBackLinkListRequest struct {
	Type          string   `form:"type"`
	Limit         int      `form:"limit"`
	LastInstance  string   `form:"last_instance"`
	Instance      string   `form:"instance"`
	LinkSources   []string `form:"link_sources"`
	ProfileSource string   `form:"profile_source"`
}

func GetBackLinkListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	request := GetBackLinkListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(errors.New("invalid params"))

		return
	}

	linkModels, err := database.QueryLinksByTo(
		database.DB, constants.LinkTypeFollowing.Int(), instance.Identity, constants.ProfileSourceIDCrossbell.Int(),
	)
	if err != nil {
		_ = c.Error(errors.New("invalid params"))

		return
	}

	var links []protocol.Link
	for _, linkModel := range linkModels {
		links = append(links, protocol.Link{
			DateCreated: linkModel.CreatedAt,
			From:        linkModel.From,
			To:          linkModel.To,
			Source:      constants.ProfileSourceID(linkModel.Source).Name().String(),
			Metadata:    protocol.LinkMetadata{
				// TODO
			},
		})
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier:  fmt.Sprintf("%s/backlinks", rss3uri.New(instance).String()),
		DateUpdated: time.Now(),
		Total:       len(links),
		List:        links,
	})
}
