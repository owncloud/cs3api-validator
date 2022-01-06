package main

import (
	"context"

	linkv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	"github.com/cucumber/godog"
)

// deletePublicShares empties the trash for all users after running the scenario
func (f *FeatureContext) deletePublicShares(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {

	for u, _ := range f.Users {
		reqctx, err := f.getAuthContext(u)
		if err != nil {
			continue
		}

		publicSharesRes, err := f.Client.ListPublicShares(
			reqctx,
			&linkv1beta1.ListPublicSharesRequest{},
		)

		if err != nil {
			continue
		}

		for _, share := range publicSharesRes.Share {

			_, _ = f.Client.RemovePublicShare(
				reqctx,
				&linkv1beta1.RemovePublicShareRequest{
					Ref: &linkv1beta1.PublicShareReference{
						Spec: &linkv1beta1.PublicShareReference_Id{
							Id: share.Id,
						},
					},
				},
			)
		}

	}
	return ctx, nil
}
