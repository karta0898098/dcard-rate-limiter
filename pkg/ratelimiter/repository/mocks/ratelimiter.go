// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// RateLimiterRepository is an autogenerated mock type for the RateLimiterRepository type
type RateLimiterRepository struct {
	mock.Mock
}

// AddRequestCount provides a mock function with given fields: ctx, key, timeUintSecond
func (_m *RateLimiterRepository) AddRequestCount(ctx context.Context, key string, timeUintSecond int64) error {
	ret := _m.Called(ctx, key, timeUintSecond)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) error); ok {
		r0 = rf(ctx, key, timeUintSecond)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRequestCount provides a mock function with given fields: ctx, key
func (_m *RateLimiterRepository) GetRequestCount(ctx context.Context, key string) (int64, error) {
	ret := _m.Called(ctx, key)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}