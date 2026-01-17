package service

import (
	"context"
	"errors"
	"strings"

	"{{ .Computed.common_module_final }}/copierx"
	"{{ .Computed.common_module_final }}/jwt"
	{{- if .Computed.enable_hotspot_final }}
	"{{ .Computed.common_module_final }}/log"
	{{- end }}
	jwtV4 "github.com/golang-jwt/jwt/v4"
	{{- if or .Computed.enable_captcha_final .Computed.enable_user_lock_final .Computed.enable_hotspot_final }}
	"github.com/golang-module/carbon/v2"
	{{- end }}
	"github.com/google/wire"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/biz"
	"{{ .Computed.module_name_final }}/internal/conf"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(New{{ .Computed.service_name_capitalized }}Service)

// {{ .Computed.service_name_capitalized }}Service implements the {{ .Computed.service_name_capitalized }} gRPC/HTTP service.
type {{ .Computed.service_name_capitalized }}Service struct {
	v1.Unimplemented{{ .Computed.service_name_capitalized }}Server

	c *conf.Bootstrap

	uc         *biz.AuthUseCase
	user       *biz.UserUseCase
	role       *biz.RoleUseCase
	permission *biz.PermissionUseCase
	{{- if .Computed.enable_action_final }}
	action *biz.ActionUseCase
	{{- end }}
	{{- if .Computed.enable_user_group_final }}
	userGroup *biz.UserGroupUseCase
	{{- end }}
	{{- if .Computed.enable_whitelist_final }}
	whitelist *biz.WhitelistUseCase
	{{- end }}
	{{- if .Computed.enable_hotspot_final }}
	hotspot biz.HotspotRepo
	{{- end }}
}

// New{{ .Computed.service_name_capitalized }}Service creates a new service instance.
func New{{ .Computed.service_name_capitalized }}Service(
	c *conf.Bootstrap,
	uc *biz.AuthUseCase,
	user *biz.UserUseCase,
	role *biz.RoleUseCase,
	permission *biz.PermissionUseCase,
	{{- if .Computed.enable_action_final }}
	action *biz.ActionUseCase,
	{{- end }}
	{{- if .Computed.enable_user_group_final }}
	userGroup *biz.UserGroupUseCase,
	{{- end }}
	{{- if .Computed.enable_whitelist_final }}
	whitelist *biz.WhitelistUseCase,
	{{- end }}
	{{- if .Computed.enable_hotspot_final }}
	hotspot biz.HotspotRepo,
	{{- end }}
) *{{ .Computed.service_name_capitalized }}Service {
	return &{{ .Computed.service_name_capitalized }}Service{
		c:          c,
		uc:         uc,
		user:       user,
		role:       role,
		permission: permission,
		{{- if .Computed.enable_action_final }}
		action: action,
		{{- end }}
		{{- if .Computed.enable_user_group_final }}
		userGroup: userGroup,
		{{- end }}
		{{- if .Computed.enable_whitelist_final }}
		whitelist: whitelist,
		{{- end }}
		{{- if .Computed.enable_hotspot_final }}
		hotspot: hotspot,
		{{- end }}
	}
}

func (s *{{ .Computed.service_name_capitalized }}Service) flushCache(ctx context.Context) {
	if s.user != nil {
		s.user.FlushCache(ctx)
	}
	{{- if .Computed.enable_hotspot_final }}
	if s.hotspot != nil {
		if err := s.hotspot.Refresh(ctx); err != nil {
			log.WithContext(ctx).WithError(err).Warn("refresh hotspot failed")
		}
	}
	{{- end }}
}

func (s *{{ .Computed.service_name_capitalized }}Service) Register(ctx context.Context, req *v1.RegisterRequest) (*emptypb.Empty, error) {
	if s.user == nil {
		return nil, biz.ErrInternal(ctx, "user usecase not configured")
	}
	r := &biz.User{}
	_ = copierx.Copy(&r, req)

	{{- if .Computed.enable_captcha_final }}
	// Skip captcha verification if configured.
	if !s.c.Server.SkipCaptcha {
		if !s.user.VerifyCaptcha(ctx, req.CaptchaId, req.CaptchaAnswer) {
			return nil, biz.ErrInvalidCaptcha(ctx)
		}
	}
	{{- end }}

	if err := s.user.Create(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) Pwd(ctx context.Context, req *v1.PwdRequest) (*emptypb.Empty, error) {
	if s.user == nil {
		return nil, biz.ErrInternal(ctx, "user usecase not configured")
	}
	r := &biz.User{}
	_ = copierx.Copy(&r, req)
	if err := s.user.Pwd(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	if req == nil {
		return nil, biz.ErrIllegalParameter(ctx, "req")
	}
	if s.user == nil {
		return nil, biz.ErrInternal(ctx, "user usecase not configured")
	}

	{{- if or .Computed.enable_captcha_final .Computed.enable_user_lock_final .Computed.enable_hotspot_final }}
	r := &biz.Login{}
	_ = copierx.Copy(&r, req)

	res, err := s.user.Login(ctx, r)
	if err != nil {
		loginFailedErr := biz.ErrLoginFailed(ctx)
		loginFailed := err.Error() == loginFailedErr.Error()
		notFound := err.Error() == biz.ErrRecordNotFound(ctx).Error()
		invalidCaptcha := err.Error() == biz.ErrInvalidCaptcha(ctx).Error()
		if invalidCaptcha {
			return nil, err
		}
		if notFound {
			// Avoid username probing.
			return nil, loginFailedErr
		}
		if loginFailed {
			// Record wrong password attempts; include a timestamp to avoid races with a later successful login.
			loginTime := biz.LoginTime{
				Username: req.Username,
				LastLogin: carbon.DateTime{
					Carbon: carbon.Now(),
				},
				{{- if or .Computed.enable_captcha_final .Computed.enable_user_lock_final }}
				Wrong: res.Wrong,
				{{- end }}
			}
			_ = s.user.WrongPwd(ctx, &loginTime)
			s.flushCache(ctx)
		}
		return nil, err
	}

	// Successful login: update last_login and reset wrong/locked flags.
	_ = s.user.LastLogin(ctx, req.Username)
	s.flushCache(ctx)

	rp := &v1.LoginReply{}
	_ = copierx.Copy(&rp, res)
	return rp, nil
	{{- else }}
	// Standard preset: use AuthUseCase (no captcha/lock flows).
	if s.uc == nil {
		return nil, biz.ErrInternal(ctx, "auth usecase not configured")
	}
	loginReq := &biz.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	if req.Platform != nil {
		loginReq.Platform = *req.Platform
	}
	res, err := s.uc.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}
	_ = s.user.LastLogin(ctx, req.Username)

	rp := &v1.LoginReply{}
	_ = copierx.Copy(&rp, res)
	return rp, nil
	{{- end }}
}

func (*{{ .Computed.service_name_capitalized }}Service) Logout(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.LoginReply, error) {
	if req == nil {
		return nil, biz.ErrIllegalParameter(ctx, "token")
	}
	tokenStr := strings.TrimSpace(req.Token)
	if tokenStr == "" {
		return nil, biz.ErrJwtMissingToken(ctx)
	}
	if s.c == nil {
		return nil, biz.ErrInternal(ctx, "config not configured")
	}

	info, err := parseToken(ctx, s.c.Server.Jwt.Key, tokenStr)
	if err != nil {
		return nil, err
	}

	claims, ok := info.Claims.(jwtV4.MapClaims)
	if !ok {
		return nil, biz.ErrJwtTokenParseFail(ctx)
	}

	// Extract attrs from claims (defaults to empty strings when missing).
	code, _ := claims["code"].(string)
	platform, _ := claims["platform"].(string)
	authUser := jwt.User{
		Attrs: map[string]string{
			"code":     code,
			"platform": platform,
		},
	}
	token, expireTime := authUser.CreateToken(s.c.Server.Jwt.Key, s.c.Server.Jwt.Expires)
	return &v1.LoginReply{
		Token:   token,
		Expires: expireTime.ToDateTimeString(),
	}, nil
}

{{- if .Computed.enable_captcha_final }}
func (s *{{ .Computed.service_name_capitalized }}Service) Captcha(ctx context.Context, _ *emptypb.Empty) (*v1.CaptchaReply, error) {
	if s.user == nil {
		return nil, biz.ErrInternal(ctx, "user usecase not configured")
	}
	res := s.user.Captcha(ctx)
	rp := &v1.CaptchaReply{
		Captcha: &v1.Captcha{},
	}
	_ = copierx.Copy(&rp.Captcha, &res)
	return rp, nil
}
{{- end }}

func (s *{{ .Computed.service_name_capitalized }}Service) Status(ctx context.Context, req *v1.StatusRequest) (*v1.StatusReply, error) {
	if req == nil {
		return nil, biz.ErrIllegalParameter(ctx, "username")
	}
	if s.user == nil {
		return nil, biz.ErrInternal(ctx, "user usecase not configured")
	}
	res, err := s.user.Status(ctx, req.Username, true)
	if err != nil {
		return nil, err
	}

	rp := &v1.StatusReply{
		Locked: res.Locked,
	}
	{{- if .Computed.enable_user_lock_final }}
	rp.LockExpire = res.LockExpire
	{{- end }}
	{{- if .Computed.enable_captcha_final }}
	if res.NeedCaptcha {
		rp.Captcha = &v1.Captcha{}
		_ = copierx.Copy(&rp.Captcha, &res.Captcha)
	}
	{{- end }}
	return rp, nil
}

func parseToken(ctx context.Context, key, jwtToken string) (info *jwtV4.Token, err error) {
	info, err = jwtV4.Parse(jwtToken, func(_ *jwtV4.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		var ve *jwtV4.ValidationError
		if !errors.As(err, &ve) {
			return nil, err
		}
		switch {
		case ve.Errors&jwtV4.ValidationErrorMalformed != 0:
			return nil, biz.ErrJwtTokenInvalid(ctx)
		case ve.Errors&(jwtV4.ValidationErrorExpired|jwtV4.ValidationErrorNotValidYet) != 0:
			return nil, biz.ErrJwtTokenExpired(ctx)
		default:
			return nil, biz.ErrJwtTokenParseFail(ctx)
		}
	}
	if !info.Valid {
		return nil, biz.ErrJwtTokenParseFail(ctx)
	}
	if info.Method != jwtV4.SigningMethodHS512 {
		return nil, biz.ErrJwtUnSupportSigningMethod(ctx)
	}
	return info, nil
}
