// Code generated by mockery v1.0.0. DO NOT EDIT.

// This comment works around https://github.com/vektra/mockery/issues/155

package s3

import mock "github.com/stretchr/testify/mock"
import time "time"
import url "net/url"

// mockMinioClient is an autogenerated mock type for the minioClient type
type mockMinioClient struct {
	mock.Mock
}

// PresignedGetObject provides a mock function with given fields: bucketName, objectName, expires, reqParams
func (_m *mockMinioClient) PresignedGetObject(bucketName string, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	ret := _m.Called(bucketName, objectName, expires, reqParams)

	var r0 *url.URL
	if rf, ok := ret.Get(0).(func(string, string, time.Duration, url.Values) *url.URL); ok {
		r0 = rf(bucketName, objectName, expires, reqParams)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url.URL)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, time.Duration, url.Values) error); ok {
		r1 = rf(bucketName, objectName, expires, reqParams)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PresignedPutObject provides a mock function with given fields: bucketName, objectName, expires
func (_m *mockMinioClient) PresignedPutObject(bucketName string, objectName string, expires time.Duration) (*url.URL, error) {
	ret := _m.Called(bucketName, objectName, expires)

	var r0 *url.URL
	if rf, ok := ret.Get(0).(func(string, string, time.Duration) *url.URL); ok {
		r0 = rf(bucketName, objectName, expires)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url.URL)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, time.Duration) error); ok {
		r1 = rf(bucketName, objectName, expires)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}