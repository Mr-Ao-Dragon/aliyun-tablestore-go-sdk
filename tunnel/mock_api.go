package tunnel

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MocktunnelDataApi is a mock of tunnelDataApi interface
type MocktunnelDataApi struct {
	ctrl     *gomock.Controller
	recorder *MocktunnelDataApiMockRecorder
}

// MocktunnelDataApiMockRecorder is the mock recorder for MocktunnelDataApi
type MocktunnelDataApiMockRecorder struct {
	mock *MocktunnelDataApi
}

// NewMocktunnelDataApi creates a new mock instance
func NewMocktunnelDataApi(ctrl *gomock.Controller) *MocktunnelDataApi {
	mock := &MocktunnelDataApi{ctrl: ctrl}
	mock.recorder = &MocktunnelDataApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MocktunnelDataApi) EXPECT() *MocktunnelDataApiMockRecorder {
	return m.recorder
}

// readRecords mocks base method
func (m *MocktunnelDataApi) readRecords(tunnelId, clientId, channelId, token string) ([]*Record, string, string, int, error) {
	ret := m.ctrl.Call(m, "readRecords", tunnelId, clientId, channelId, token)
	ret0, _ := ret[0].([]*Record)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(int)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// readRecords indicates an expected call of readRecords
func (mr *MocktunnelDataApiMockRecorder) readRecords(tunnelId, clientId, channelId, token interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "readRecords", reflect.TypeOf((*MocktunnelDataApi)(nil).readRecords), tunnelId, clientId, channelId, token)
}

// MockTunnelMetaApi is a mock of TunnelMetaApi interface
type MockTunnelMetaApi struct {
	ctrl     *gomock.Controller
	recorder *MockTunnelMetaApiMockRecorder
}

// MockTunnelMetaApiMockRecorder is the mock recorder for MockTunnelMetaApi
type MockTunnelMetaApiMockRecorder struct {
	mock *MockTunnelMetaApi
}

// NewMockTunnelMetaApi creates a new mock instance
func NewMockTunnelMetaApi(ctrl *gomock.Controller) *MockTunnelMetaApi {
	mock := &MockTunnelMetaApi{ctrl: ctrl}
	mock.recorder = &MockTunnelMetaApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTunnelMetaApi) EXPECT() *MockTunnelMetaApiMockRecorder {
	return m.recorder
}

// CreateTunnel mocks base method
func (m *MockTunnelMetaApi) CreateTunnel(req *CreateTunnelRequest) (*CreateTunnelResponse, error) {
	ret := m.ctrl.Call(m, "CreateTunnel", req)
	ret0, _ := ret[0].(*CreateTunnelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTunnel indicates an expected call of CreateTunnel
func (mr *MockTunnelMetaApiMockRecorder) CreateTunnel(req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTunnel", reflect.TypeOf((*MockTunnelMetaApi)(nil).CreateTunnel), req)
}

// ListTunnel mocks base method
func (m *MockTunnelMetaApi) ListTunnel(req *ListTunnelRequest) (*ListTunnelResponse, error) {
	ret := m.ctrl.Call(m, "ListTunnel", req)
	ret0, _ := ret[0].(*ListTunnelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTunnel indicates an expected call of ListTunnel
func (mr *MockTunnelMetaApiMockRecorder) ListTunnel(req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTunnel", reflect.TypeOf((*MockTunnelMetaApi)(nil).ListTunnel), req)
}

// DescribeTunnel mocks base method
func (m *MockTunnelMetaApi) DescribeTunnel(req *DescribeTunnelRequest) (*DescribeTunnelResponse, error) {
	ret := m.ctrl.Call(m, "DescribeTunnel", req)
	ret0, _ := ret[0].(*DescribeTunnelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeTunnel indicates an expected call of DescribeTunnel
func (mr *MockTunnelMetaApiMockRecorder) DescribeTunnel(req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeTunnel", reflect.TypeOf((*MockTunnelMetaApi)(nil).DescribeTunnel), req)
}

// DeleteTunnel mocks base method
func (m *MockTunnelMetaApi) DeleteTunnel(req *DeleteTunnelRequest) (*DeleteTunnelResponse, error) {
	ret := m.ctrl.Call(m, "DeleteTunnel", req)
	ret0, _ := ret[0].(*DeleteTunnelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTunnel indicates an expected call of DeleteTunnel
func (mr *MockTunnelMetaApiMockRecorder) DeleteTunnel(req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTunnel", reflect.TypeOf((*MockTunnelMetaApi)(nil).DeleteTunnel), req)
}

// NewTunnelWorker mocks base method
func (m *MockTunnelMetaApi) NewTunnelWorker(tunnelId string, workerConfig *TunnelWorkerConfig) (TunnelWorker, error) {
	ret := m.ctrl.Call(m, "NewTunnelWorker", tunnelId, workerConfig)
	ret0, _ := ret[0].(TunnelWorker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewTunnelWorker indicates an expected call of NewTunnelWorker
func (mr *MockTunnelMetaApiMockRecorder) NewTunnelWorker(tunnelId, workerConfig interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTunnelWorker", reflect.TypeOf((*MockTunnelMetaApi)(nil).NewTunnelWorker), tunnelId, workerConfig)
}

// MockTunnelWorker is a mock of TunnelWorker interface
type MockTunnelWorker struct {
	ctrl     *gomock.Controller
	recorder *MockTunnelWorkerMockRecorder
}

// MockTunnelWorkerMockRecorder is the mock recorder for MockTunnelWorker
type MockTunnelWorkerMockRecorder struct {
	mock *MockTunnelWorker
}

// NewMockTunnelWorker creates a new mock instance
func NewMockTunnelWorker(ctrl *gomock.Controller) *MockTunnelWorker {
	mock := &MockTunnelWorker{ctrl: ctrl}
	mock.recorder = &MockTunnelWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTunnelWorker) EXPECT() *MockTunnelWorkerMockRecorder {
	return m.recorder
}

// ConnectAndWorking mocks base method
func (m *MockTunnelWorker) ConnectAndWorking() error {
	ret := m.ctrl.Call(m, "ConnectAndWorking")
	ret0, _ := ret[0].(error)
	return ret0
}

// ConnectAndWorking indicates an expected call of ConnectAndWorking
func (mr *MockTunnelWorkerMockRecorder) ConnectAndWorking() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectAndWorking", reflect.TypeOf((*MockTunnelWorker)(nil).ConnectAndWorking))
}

// Shutdown mocks base method
func (m *MockTunnelWorker) Shutdown() {
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown
func (mr *MockTunnelWorkerMockRecorder) Shutdown() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockTunnelWorker)(nil).Shutdown))
}
