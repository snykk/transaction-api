// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	v1 "github.com/snykk/transaction-api/internal/business/domains/v1"
	mock "github.com/stretchr/testify/mock"
)

// TransactionRepository is an autogenerated mock type for the TransactionRepository type
type TransactionRepository struct {
	mock.Mock
}

// Deposit provides a mock function with given fields: ctx, transactionDom
func (_m *TransactionRepository) Deposit(ctx context.Context, transactionDom v1.TransactionDomain) (v1.TransactionDomain, error) {
	ret := _m.Called(ctx, transactionDom)

	if len(ret) == 0 {
		panic("no return value specified for Deposit")
	}

	var r0 v1.TransactionDomain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, v1.TransactionDomain) (v1.TransactionDomain, error)); ok {
		return rf(ctx, transactionDom)
	}
	if rf, ok := ret.Get(0).(func(context.Context, v1.TransactionDomain) v1.TransactionDomain); ok {
		r0 = rf(ctx, transactionDom)
	} else {
		r0 = ret.Get(0).(v1.TransactionDomain)
	}

	if rf, ok := ret.Get(1).(func(context.Context, v1.TransactionDomain) error); ok {
		r1 = rf(ctx, transactionDom)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *TransactionRepository) GetAll(ctx context.Context) ([]v1.TransactionDomain, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []v1.TransactionDomain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]v1.TransactionDomain, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []v1.TransactionDomain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]v1.TransactionDomain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUserId provides a mock function with given fields: ctx, userId
func (_m *TransactionRepository) GetByUserId(ctx context.Context, userId string) ([]v1.TransactionDomain, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserId")
	}

	var r0 []v1.TransactionDomain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]v1.TransactionDomain, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []v1.TransactionDomain); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]v1.TransactionDomain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Purchase provides a mock function with given fields: ctx, trasanctionDom
func (_m *TransactionRepository) Purchase(ctx context.Context, trasanctionDom v1.TransactionDomain) (v1.TransactionDomain, error) {
	ret := _m.Called(ctx, trasanctionDom)

	if len(ret) == 0 {
		panic("no return value specified for Purchase")
	}

	var r0 v1.TransactionDomain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, v1.TransactionDomain) (v1.TransactionDomain, error)); ok {
		return rf(ctx, trasanctionDom)
	}
	if rf, ok := ret.Get(0).(func(context.Context, v1.TransactionDomain) v1.TransactionDomain); ok {
		r0 = rf(ctx, trasanctionDom)
	} else {
		r0 = ret.Get(0).(v1.TransactionDomain)
	}

	if rf, ok := ret.Get(1).(func(context.Context, v1.TransactionDomain) error); ok {
		r1 = rf(ctx, trasanctionDom)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Withdraw provides a mock function with given fields: ctx, transactionDom
func (_m *TransactionRepository) Withdraw(ctx context.Context, transactionDom v1.TransactionDomain) (v1.TransactionDomain, error) {
	ret := _m.Called(ctx, transactionDom)

	if len(ret) == 0 {
		panic("no return value specified for Withdraw")
	}

	var r0 v1.TransactionDomain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, v1.TransactionDomain) (v1.TransactionDomain, error)); ok {
		return rf(ctx, transactionDom)
	}
	if rf, ok := ret.Get(0).(func(context.Context, v1.TransactionDomain) v1.TransactionDomain); ok {
		r0 = rf(ctx, transactionDom)
	} else {
		r0 = ret.Get(0).(v1.TransactionDomain)
	}

	if rf, ok := ret.Get(1).(func(context.Context, v1.TransactionDomain) error); ok {
		r1 = rf(ctx, transactionDom)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransactionRepository creates a new instance of TransactionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionRepository {
	mock := &TransactionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
