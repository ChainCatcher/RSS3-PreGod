package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
	"github.com/gin-gonic/gin"
)

type GetBackLinkListRequest struct {
	Type           string `form:"type"`
	Limit          int    `form:"limit"`
	LastInstance   string `form:"last_instance"`
	Instance       string `form:"instance"`
	LinkSources    []int  `form:"link_sources"`
	ProfileSources []int  `form:"profile_sources"`
}

func GetBackLinkListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	request := GetBackLinkListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(errors.New("invalid params"))

		return
	}

	linkSources := make([]int, 0)
	for _, linkSource := range request.LinkSources {
		linkSources = append(linkSources, constants.LinkSourceName(linkSource).ID().Int())
	}

	var linkType *int

	linkModels, err := database.QueryLinksByTo(
		database.DB,
		linkType,
		instance.Identity,
		request.LinkSources,
		request.Limit,
	)
	if err != nil {
		_ = c.Error(errors.New("invalid params"))

		return
	}

	links := make([]protocol.Link, 0)

	for _, linkModel := range linkModels {
		fromInstance, err := rss3uri.NewInstance(
			constants.InstanceTypeID(linkModel.FromInstanceType).String(),
			linkModel.From,
			constants.PlatformID(linkModel.FromPlatformID).Symbol().String(),
		)
		if err != nil {
			_ = c.Error(err)

			return
		}

		toInstance, err := rss3uri.NewInstance(
			constants.InstanceTypeID(linkModel.ToInstanceType).String(),
			linkModel.From,
			constants.PlatformID(linkModel.ToPlatformID).Symbol().String(),
		)
		if err != nil {
			_ = c.Error(err)

			return
		}

		links = append(links, protocol.Link{
			DateCreated: timex.Time(linkModel.CreatedAt),
			From:        rss3uri.New(fromInstance).String(),
			To:          rss3uri.New(toInstance).String(),
			Type:        constants.LinkTypeID(linkModel.Type).String(),
			Source:      constants.ProfileSourceID(linkModel.Source).Name().String(),
			Metadata: protocol.LinkMetadata{
				Network: constants.NetworkSymbolCrossbell.String(),
				Proof:   "TODO",
			},
		})
	}

	var dateUpdated *timex.Time

	for _, link := range links {
		internalTime := link.DateCreated
		if dateUpdated == nil {
			dateUpdated = &internalTime
		} else if dateUpdated.Time().Before(link.DateCreated.Time()) {
			dateUpdated = &internalTime
		}
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier:  fmt.Sprintf("%s/backlinks", rss3uri.New(instance)),
		DateUpdated: dateUpdated,
		Total:       len(links),
		List:        links,
	})
}
