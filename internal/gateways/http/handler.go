package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go_avito_tech/api/gen"
	"go_avito_tech/internal/domain"
	"net/http"
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.useCases.Teams.Save(ctx.Request().Context(), body.TeamName); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, member := range body.Members {
		user := domain.User{
			ID:       member.UserId,
			Username: member.Username,
			TeamName: body.TeamName,
			IsActive: member.IsActive,
		}
		if err := h.useCases.Users.Save(ctx.Request().Context(), user); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusCreated, body)
}

func (h *Handler) GetTeamGet(ctx echo.Context, params gen.GetTeamGetParams) error {
	team, err := h.useCases.Teams.FindByName(ctx.Request().Context(), params.TeamName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "team not found")
	}
	return ctx.JSON(http.StatusOK, team)
}

func (h *Handler) PostUsersSetIsActive(ctx echo.Context) error {
	var body gen.PostUsersSetIsActiveJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := h.useCases.Users.FindByID(ctx.Request().Context(), body.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	if err := h.useCases.Users.SetActive(ctx.Request().Context(), user.ID, body.IsActive); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	user.IsActive = body.IsActive
	return ctx.JSON(http.StatusOK, user)
}

func (h *Handler) GetUsersGetReview(ctx echo.Context, params gen.GetUsersGetReviewParams) error {
	_, err := h.useCases.Users.FindByID(ctx.Request().Context(), params.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	prs, err := h.useCases.PullRs.FindByReviewer(ctx.Request().Context(), params.UserId)
	if err != nil {
		fmt.Println("TYT")
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, prs)
}

func (h *Handler) PostPullRequestCreate(ctx echo.Context) error {
	var body gen.PostPullRequestCreateJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	pr := &domain.PullRequest{
		ID:                body.PullRequestId,
		Name:              body.PullRequestName,
		AuthorID:          body.AuthorId,
		Status:            domain.StatusOpen,
		NeedMoreReviewers: true,
	}
	if err := h.useCases.PullRs.Save(ctx.Request().Context(), *pr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if _, err := h.useCases.PullRs.AssignReviewers(ctx.Request().Context(), pr.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	apiPr := gen.PullRequest{
		PullRequestId:   pr.ID,
		PullRequestName: pr.Name,
		Status:          gen.PullRequestStatusOPEN,
		CreatedAt:       pr.CreatedAt,
	}
	return ctx.JSON(http.StatusCreated, apiPr)
}

func (h *Handler) PostPullRequestMerge(ctx echo.Context) error {
	var body gen.PostPullRequestMergeJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	pr, err := h.useCases.PullRs.Merge(ctx.Request().Context(), body.PullRequestId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	apiPr := gen.PullRequest{
		PullRequestId:   pr.ID,
		PullRequestName: pr.Name,
		Status:          gen.PullRequestStatusMERGED,
		CreatedAt:       pr.CreatedAt,
		MergedAt:        pr.MergedAt,
	}
	return ctx.JSON(http.StatusOK, apiPr)
}

func (h *Handler) PostPullRequestReassign(ctx echo.Context) error {
	var body gen.PostPullRequestReassignJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err := h.useCases.Users.FindByID(ctx.Request().Context(), body.OldUserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "old reviewer not found")
	}
	revs, err := h.useCases.PullRs.ReassignReviewer(ctx.Request().Context(), body.PullRequestId, body.OldUserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, revs)
}

func (h *Handler) GetStats(ctx echo.Context) error {
	stats, err := h.useCases.Stats.GetStats(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, stats)
}
