package pb

import (
	"context"
	"errors"
	"time"

	"github.com/samuelsih/revo-voting/infra"
	ag "github.com/samuelsih/revo-voting/pb/autogenerated"
)

type Status = string

const (
	serverErr Status = "error"
	expired   Status = "expired"
	notFound  Status = "not found"
	ok        Status = "ok"
)

type VotingThemeFinder interface {
	FindVotingTheme(ctx context.Context, id string) (time.Time, error)
}

type CheckerVotingService struct {
	VotingThemeFinder VotingThemeFinder

	ag.UnimplementedVoteStatusServiceServer
}

func (c *CheckerVotingService) CheckStatus(ctx context.Context, in *ag.Request) (*ag.Response, error) {
	if in.GetVoteId() == "" {
		return &ag.Response{Status: notFound}, nil
	}

	endAt, err := c.VotingThemeFinder.FindVotingTheme(ctx, in.GetVoteId())
	if err != nil {
		if !errors.Is(err, infra.ErrVotingThemeNotFound) {
			return &ag.Response{Status: serverErr}, err
		}

		return &ag.Response{Status: notFound}, nil
	}

	nowUTC := time.Now().UTC()
	endAtUTC := endAt.UTC()

	if endAtUTC.Before(nowUTC) {
		return &ag.Response{Status: expired}, nil
	}

	return &ag.Response{Status: ok}, nil
}
