// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/interface/user_repo.go

// Package mockusecase is a generated GoMock package.
package mockusecase

import (
	reflect "reflect"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	entities "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	request "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/requestModels"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	gomock "github.com/golang/mock/gomock"
)

// MockIUserRepo is a mock of IUserRepo interface.
type MockIUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIUserRepoMockRecorder
}

// MockIUserRepoMockRecorder is the mock recorder for MockIUserRepo.
type MockIUserRepoMockRecorder struct {
	mock *MockIUserRepo
}

// NewMockIUserRepo creates a new mock instance.
func NewMockIUserRepo(ctrl *gomock.Controller) *MockIUserRepo {
	mock := &MockIUserRepo{ctrl: ctrl}
	mock.recorder = &MockIUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserRepo) EXPECT() *MockIUserRepoMockRecorder {
	return m.recorder
}

// AddUserAddress mocks base method.
func (m *MockIUserRepo) AddUserAddress(newAddress *entities.UserAddress) *e.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserAddress", newAddress)
	ret0, _ := ret[0].(*e.Error)
	return ret0
}

// AddUserAddress indicates an expected call of AddUserAddress.
func (mr *MockIUserRepoMockRecorder) AddUserAddress(newAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserAddress", reflect.TypeOf((*MockIUserRepo)(nil).AddUserAddress), newAddress)
}

// CreateUser mocks base method.
func (m *MockIUserRepo) CreateUser(arg0 *entities.User) *e.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0)
	ret0, _ := ret[0].(*e.Error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIUserRepoMockRecorder) CreateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIUserRepo)(nil).CreateUser), arg0)
}

// DeleteUserAddress mocks base method.
func (m *MockIUserRepo) DeleteUserAddress(id uint) *e.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserAddress", id)
	ret0, _ := ret[0].(*e.Error)
	return ret0
}

// DeleteUserAddress indicates an expected call of DeleteUserAddress.
func (mr *MockIUserRepoMockRecorder) DeleteUserAddress(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserAddress", reflect.TypeOf((*MockIUserRepo)(nil).DeleteUserAddress), id)
}

// DoAddressNameExists mocks base method.
func (m *MockIUserRepo) DoAddressNameExists(name string) (bool, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoAddressNameExists", name)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// DoAddressNameExists indicates an expected call of DoAddressNameExists.
func (mr *MockIUserRepoMockRecorder) DoAddressNameExists(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoAddressNameExists", reflect.TypeOf((*MockIUserRepo)(nil).DoAddressNameExists), name)
}

// EditProfile mocks base method.
func (m *MockIUserRepo) EditProfile(userID uint, req *request.EditProfileReq) *e.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProfile", userID, req)
	ret0, _ := ret[0].(*e.Error)
	return ret0
}

// EditProfile indicates an expected call of EditProfile.
func (mr *MockIUserRepoMockRecorder) EditProfile(userID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProfile", reflect.TypeOf((*MockIUserRepo)(nil).EditProfile), userID, req)
}

// EditUserAddress mocks base method.
func (m *MockIUserRepo) EditUserAddress(newaddress *entities.UserAddress) *e.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditUserAddress", newaddress)
	ret0, _ := ret[0].(*e.Error)
	return ret0
}

// EditUserAddress indicates an expected call of EditUserAddress.
func (mr *MockIUserRepoMockRecorder) EditUserAddress(newaddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditUserAddress", reflect.TypeOf((*MockIUserRepo)(nil).EditUserAddress), newaddress)
}

// GetAddressNameByID mocks base method.
func (m *MockIUserRepo) GetAddressNameByID(id uint) (*string, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddressNameByID", id)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetAddressNameByID indicates an expected call of GetAddressNameByID.
func (mr *MockIUserRepoMockRecorder) GetAddressNameByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddressNameByID", reflect.TypeOf((*MockIUserRepo)(nil).GetAddressNameByID), id)
}

// GetEmailByUserID mocks base method.
func (m *MockIUserRepo) GetEmailByUserID(userID uint) (*string, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailByUserID", userID)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetEmailByUserID indicates an expected call of GetEmailByUserID.
func (mr *MockIUserRepoMockRecorder) GetEmailByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailByUserID", reflect.TypeOf((*MockIUserRepo)(nil).GetEmailByUserID), userID)
}

// GetPasswordAndUserDetailsByEmail mocks base method.
func (m *MockIUserRepo) GetPasswordAndUserDetailsByEmail(arg0 string) (*entities.User, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPasswordAndUserDetailsByEmail", arg0)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetPasswordAndUserDetailsByEmail indicates an expected call of GetPasswordAndUserDetailsByEmail.
func (mr *MockIUserRepoMockRecorder) GetPasswordAndUserDetailsByEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswordAndUserDetailsByEmail", reflect.TypeOf((*MockIUserRepo)(nil).GetPasswordAndUserDetailsByEmail), arg0)
}

// GetProfile mocks base method.
func (m *MockIUserRepo) GetProfile(userID uint) (*entities.UserDetails, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", userID)
	ret0, _ := ret[0].(*entities.UserDetails)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockIUserRepoMockRecorder) GetProfile(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockIUserRepo)(nil).GetProfile), userID)
}

// GetUserAddress mocks base method.
func (m *MockIUserRepo) GetUserAddress(addressID uint) (*entities.UserAddress, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAddress", addressID)
	ret0, _ := ret[0].(*entities.UserAddress)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetUserAddress indicates an expected call of GetUserAddress.
func (mr *MockIUserRepoMockRecorder) GetUserAddress(addressID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAddress", reflect.TypeOf((*MockIUserRepo)(nil).GetUserAddress), addressID)
}

// GetUserAddresses mocks base method.
func (m *MockIUserRepo) GetUserAddresses(userId uint) (*[]entities.UserAddress, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAddresses", userId)
	ret0, _ := ret[0].(*[]entities.UserAddress)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetUserAddresses indicates an expected call of GetUserAddresses.
func (mr *MockIUserRepoMockRecorder) GetUserAddresses(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAddresses", reflect.TypeOf((*MockIUserRepo)(nil).GetUserAddresses), userId)
}

// GetUserBasicInfoByID mocks base method.
func (m *MockIUserRepo) GetUserBasicInfoByID(id uint) (*response.UserInfoForInvoice, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBasicInfoByID", id)
	ret0, _ := ret[0].(*response.UserInfoForInvoice)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetUserBasicInfoByID indicates an expected call of GetUserBasicInfoByID.
func (mr *MockIUserRepoMockRecorder) GetUserBasicInfoByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBasicInfoByID", reflect.TypeOf((*MockIUserRepo)(nil).GetUserBasicInfoByID), id)
}

// GetUserByEmail mocks base method.
func (m *MockIUserRepo) GetUserByEmail(email string) (*entities.User, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockIUserRepoMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockIUserRepo)(nil).GetUserByEmail), email)
}

// GetUserDetailsByEmail mocks base method.
func (m *MockIUserRepo) GetUserDetailsByEmail(email string) (*entities.UserDetails, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDetailsByEmail", email)
	ret0, _ := ret[0].(*entities.UserDetails)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetUserDetailsByEmail indicates an expected call of GetUserDetailsByEmail.
func (mr *MockIUserRepoMockRecorder) GetUserDetailsByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDetailsByEmail", reflect.TypeOf((*MockIUserRepo)(nil).GetUserDetailsByEmail), email)
}

// GetUserIDFromAddressID mocks base method.
func (m *MockIUserRepo) GetUserIDFromAddressID(id uint) (uint, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDFromAddressID", id)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetUserIDFromAddressID indicates an expected call of GetUserIDFromAddressID.
func (mr *MockIUserRepoMockRecorder) GetUserIDFromAddressID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDFromAddressID", reflect.TypeOf((*MockIUserRepo)(nil).GetUserIDFromAddressID), id)
}

// GetWalletBalance mocks base method.
func (m *MockIUserRepo) GetWalletBalance(userID uint) (float32, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletBalance", userID)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// GetWalletBalance indicates an expected call of GetWalletBalance.
func (mr *MockIUserRepoMockRecorder) GetWalletBalance(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletBalance", reflect.TypeOf((*MockIUserRepo)(nil).GetWalletBalance), userID)
}

// IsEmailRegistered mocks base method.
func (m *MockIUserRepo) IsEmailRegistered(arg0 string) (bool, *e.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmailRegistered", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*e.Error)
	return ret0, ret1
}

// IsEmailRegistered indicates an expected call of IsEmailRegistered.
func (mr *MockIUserRepoMockRecorder) IsEmailRegistered(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmailRegistered", reflect.TypeOf((*MockIUserRepo)(nil).IsEmailRegistered), arg0)
}

// ResetPassword mocks base method.
func (m *MockIUserRepo) ResetPassword(id uint, newPassword *string) *e.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", id, newPassword)
	ret0, _ := ret[0].(*e.Error)
	return ret0
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockIUserRepoMockRecorder) ResetPassword(id, newPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockIUserRepo)(nil).ResetPassword), id, newPassword)
}

// UpdateUserStatus mocks base method.
func (m *MockIUserRepo) UpdateUserStatus(email, newStatus string) *e.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserStatus", email, newStatus)
	ret0, _ := ret[0].(*e.Error)
	return ret0
}

// UpdateUserStatus indicates an expected call of UpdateUserStatus.
func (mr *MockIUserRepoMockRecorder) UpdateUserStatus(email, newStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserStatus", reflect.TypeOf((*MockIUserRepo)(nil).UpdateUserStatus), email, newStatus)
}
