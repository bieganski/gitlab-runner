// Code generated by mockery v1.1.0. DO NOT EDIT.

package probe

import mock "github.com/stretchr/testify/mock"

// MockLogger is an autogenerated mock type for the Logger type
type MockLogger struct {
	mock.Mock
}

// Warningln provides a mock function with given fields: args
func (_m *MockLogger) Warningln(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}
