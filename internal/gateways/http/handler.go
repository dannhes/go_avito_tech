package http

import (
	"github.com/labstack/echo/v4"
	"go_avito_tech/api/gen"
	"go_avito_tech/internal/domain"
)

type Handler struct {
	useCases UseCases
}

func NewHandler(useCases UseCases) *Handler {
	return &Handler{useCases: useCases}
}

func (h *Handler) PostTeamAdd(ctx echo.Context) error {
	var body gen.PostTeamAddJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	for _, member := range body.Members {
		user := domain.User{
			ID:       member.UserId,
			Username: member.Username,
			TeamName: body.TeamName,
			IsActive: member.IsActive,
		}
		if err := h.useCases.Users.Save(ctx.Request().Context(), user); err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
	}
	if err := h.useCases.Teams.Save(ctx.Request().Context(), body.TeamName); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(201, body)
}

func (h *Handler) GetTeamGet(ctx echo.Context, params gen.GetTeamGetParams) error {
	team, err := h.useCases.Teams.FindByName(ctx.Request().Context(), params.TeamName)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, team)
}

func (h *Handler) PostUsersSetIsActive(ctx echo.Context) error {
	var body gen.PostUsersSetIsActiveJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := h.useCases.Users.SetActive(ctx.Request().Context(), body.UserId, body.IsActive); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, body)
}

func (h *Handler) GetUsersGetReview(ctx echo.Context, params gen.GetUsersGetReviewParams) error {
	prs, err := h.useCases.PullRs.FindByReviewer(ctx.Request().Context(), params.UserId)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, prs)
}

func (h *Handler) PostPullRequestCreate(ctx echo.Context) error {
	var body gen.PostPullRequestCreateJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	pr := domain.PullRequest{
		ID:       body.PullRequestId,
		Name:     body.PullRequestName,
		AuthorID: body.AuthorId,
		Status:   domain.StatusOpen,
	}
	if err := h.useCases.PullRs.Save(ctx.Request().Context(), pr); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(201, pr)
}

func (h *Handler) PostPullRequestMerge(ctx echo.Context) error {
	var body gen.PostPullRequestMergeJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	pr, err := h.useCases.PullRs.Merge(ctx.Request().Context(), body.PullRequestId)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return ctx.JSON(200, pr)
}

func (h *Handler) PostPullRequestReassign(ctx echo.Context) error {
	var body gen.PostPullRequestReassignJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	revs, err := h.useCases.PullRs.ReassignReviewer(ctx.Request().Context(), body.PullRequestId, body.OldUserId)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, revs)
}

func (h *Handler) GetStats(ctx echo.Context) error {
	var body gen.PostPullRequestCreateJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	stats, err := h.useCases.Stats.GetStats(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, stats)
}
