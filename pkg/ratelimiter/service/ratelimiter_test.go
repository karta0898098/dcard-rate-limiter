package service

import (
	"context"
	"testing"

	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/domain"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/repository"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/repository/mocks"

	"github.com/karta0898098/kara/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimiterServiceRequireResource(t *testing.T) {
	var (
		mockMaxCount       int64 = 10
		mockTimeUintSecond int64 = 10
		mockRemoteAddr           = "127.0.0.1"
		mockURL                  = "/proctored"
	)
	type fields struct {
		maxCount       int64
		timeUintSecond int64
		repo           func() repository.RateLimiterRepository
	}
	type args struct {
		ctx  context.Context
		addr string
		url  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		claims *domain.Claims
		err    *errors.Exception
	}{
		{
			name: "AddrIsNil",
			fields: fields{
				maxCount:       mockMaxCount,
				timeUintSecond: mockTimeUintSecond,
				repo: func() repository.RateLimiterRepository {
					repo := new(mocks.RateLimiterRepository)
					return repo
				},
			},
			args: args{
				ctx:  context.Background(),
				addr: "",
				url:  "",
			},
			claims: nil,
			err:    errors.ErrInternal,
		},
		{
			name: "GetRequestCountOccurError",
			fields: fields{
				maxCount:       mockMaxCount,
				timeUintSecond: mockTimeUintSecond,
				repo: func() repository.RateLimiterRepository {
					repo := new(mocks.RateLimiterRepository)
					repo.
						On("GetRequestCount", mock.Anything, mock.Anything).
						Return(int64(0), errors.ErrInternal)
					return repo
				},
			},
			args: args{
				ctx:  context.Background(),
				addr: mockRemoteAddr,
				url:  mockURL,
			},
			claims: nil,
			err:    errors.ErrInternal,
		},
		{
			name: "TooManyRequest",
			fields: fields{
				maxCount:       mockMaxCount,
				timeUintSecond: mockTimeUintSecond,
				repo: func() repository.RateLimiterRepository {
					repo := new(mocks.RateLimiterRepository)
					repo.
						On("GetRequestCount", mock.Anything, mock.Anything).
						Return(int64(10), nil)
					return repo
				},
			},
			args: args{
				ctx:  context.Background(),
				addr: mockRemoteAddr,
				url:  mockURL,
			},
			claims: nil,
			err:    errors.ErrTooManyRequests,
		},
		{
			name: "AddRequestCountOccurError",
			fields: fields{
				maxCount:       mockMaxCount,
				timeUintSecond: mockTimeUintSecond,
				repo: func() repository.RateLimiterRepository {
					repo := new(mocks.RateLimiterRepository)
					repo.
						On("GetRequestCount", mock.Anything, mock.Anything).
						Return(int64(0), nil)
					repo.
						On("AddRequestCount", mock.Anything, mock.Anything, mock.Anything).
						Return(errors.ErrInternal)
					return repo
				},
			},
			args: args{
				ctx:  context.Background(),
				addr: mockRemoteAddr,
				url:  mockURL,
			},
			claims: nil,
			err:    errors.ErrInternal,
		},
		{
			name: "Success",
			fields: fields{
				maxCount:       mockMaxCount,
				timeUintSecond: mockTimeUintSecond,
				repo: func() repository.RateLimiterRepository {
					repo := new(mocks.RateLimiterRepository)
					repo.
						On("GetRequestCount", mock.Anything, mock.Anything).
						Return(int64(9), nil)
					repo.
						On("AddRequestCount", mock.Anything, mock.Anything, mock.Anything).
						Return(nil)
					return repo
				},
			},
			args: args{
				ctx:  context.Background(),
				addr: mockRemoteAddr,
				url:  mockURL,
			},
			claims: &domain.Claims{
				URL:          mockURL,
				RemoteAddr:   mockRemoteAddr,
				RequestCount: 10,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &rateLimiterService{
				maxCount:       tt.fields.maxCount,
				timeUintSecond: tt.fields.timeUintSecond,
				repo:           tt.fields.repo(),
			}
			claims, err := srv.RequireResource(tt.args.ctx, tt.args.addr, tt.args.url)
			if err != nil {
				assert.True(t, tt.err.Is(err), "error type not equal")
				return
			}
			assert.Equal(t, tt.claims, claims)
		})
	}
}
