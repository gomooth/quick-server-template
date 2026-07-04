package system

import (
	"context"
	"time"

	"server-api/app/openapi/internal/helper"
	"server-api/app/openapi/internal/helper/ecode"
	"server-api/repository/platform/pcache"

	"github.com/gomooth/xerror"
)

type service struct{}

func (s *service) CurrentTime(ctx context.Context) (*currentTimeEntity, error) {
	now := time.Now()
	return &currentTimeEntity{
		Zone:      now.In(time.Local).Location().String(),
		Time:      now.Format("2006-01-02 15:04:05.999"),
		Timestamp: now.UnixMilli(),
	}, nil
}

func (s *service) Sign(ctx context.Context, h *helper.Header) (*helper.SignResult, error) {
	app, err := pcache.NewOpenAPP().FirstByAppID(ctx, h.AppID)
	if err != nil {
		return nil, xerror.WrapWithXCode(err, ecode.BusinessError)
	}

	path := helper.GetRequestPath(ctx)
	qs, body := helper.ExtractRequestParams(ctx)
	result := helper.Sign(h.AppID, app.AppSecret, h.Timestamp, path, qs, body)
	result.Input = h.Sign
	result.Success = result.Input == result.Signature
	return result, nil
}
