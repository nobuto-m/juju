// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/agent/provisioner (interfaces: Machine,BridgePolicy,Unit,Application,Charm)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	set "github.com/juju/collections/set"
	provisioner "github.com/juju/juju/apiserver/facades/agent/provisioner"
	constraints "github.com/juju/juju/core/constraints"
	instance "github.com/juju/juju/core/instance"
	network "github.com/juju/juju/network"
	containerizer "github.com/juju/juju/network/containerizer"
	state "github.com/juju/juju/state"
	charm_v6 "gopkg.in/juju/charm.v6"
	names_v3 "gopkg.in/juju/names.v3"
	reflect "reflect"
)

// MockMachine is a mock of Machine interface
type MockMachine struct {
	ctrl     *gomock.Controller
	recorder *MockMachineMockRecorder
}

// MockMachineMockRecorder is the mock recorder for MockMachine
type MockMachineMockRecorder struct {
	mock *MockMachine
}

// NewMockMachine creates a new mock instance
func NewMockMachine(ctrl *gomock.Controller) *MockMachine {
	mock := &MockMachine{ctrl: ctrl}
	mock.recorder = &MockMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMachine) EXPECT() *MockMachineMockRecorder {
	return m.recorder
}

// AllAddresses mocks base method
func (m *MockMachine) AllAddresses() ([]containerizer.Address, error) {
	ret := m.ctrl.Call(m, "AllAddresses")
	ret0, _ := ret[0].([]containerizer.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllAddresses indicates an expected call of AllAddresses
func (mr *MockMachineMockRecorder) AllAddresses() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllAddresses", reflect.TypeOf((*MockMachine)(nil).AllAddresses))
}

// AllLinkLayerDevices mocks base method
func (m *MockMachine) AllLinkLayerDevices() ([]containerizer.LinkLayerDevice, error) {
	ret := m.ctrl.Call(m, "AllLinkLayerDevices")
	ret0, _ := ret[0].([]containerizer.LinkLayerDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllLinkLayerDevices indicates an expected call of AllLinkLayerDevices
func (mr *MockMachineMockRecorder) AllLinkLayerDevices() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllLinkLayerDevices", reflect.TypeOf((*MockMachine)(nil).AllLinkLayerDevices))
}

// AllSpaces mocks base method
func (m *MockMachine) AllSpaces() (set.Strings, error) {
	ret := m.ctrl.Call(m, "AllSpaces")
	ret0, _ := ret[0].(set.Strings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSpaces indicates an expected call of AllSpaces
func (mr *MockMachineMockRecorder) AllSpaces() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSpaces", reflect.TypeOf((*MockMachine)(nil).AllSpaces))
}

// Constraints mocks base method
func (m *MockMachine) Constraints() (constraints.Value, error) {
	ret := m.ctrl.Call(m, "Constraints")
	ret0, _ := ret[0].(constraints.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Constraints indicates an expected call of Constraints
func (mr *MockMachineMockRecorder) Constraints() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Constraints", reflect.TypeOf((*MockMachine)(nil).Constraints))
}

// ContainerType mocks base method
func (m *MockMachine) ContainerType() instance.ContainerType {
	ret := m.ctrl.Call(m, "ContainerType")
	ret0, _ := ret[0].(instance.ContainerType)
	return ret0
}

// ContainerType indicates an expected call of ContainerType
func (mr *MockMachineMockRecorder) ContainerType() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerType", reflect.TypeOf((*MockMachine)(nil).ContainerType))
}

// Id mocks base method
func (m *MockMachine) Id() string {
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

// Id indicates an expected call of Id
func (mr *MockMachineMockRecorder) Id() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockMachine)(nil).Id))
}

// InstanceId mocks base method
func (m *MockMachine) InstanceId() (instance.Id, error) {
	ret := m.ctrl.Call(m, "InstanceId")
	ret0, _ := ret[0].(instance.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceId indicates an expected call of InstanceId
func (mr *MockMachineMockRecorder) InstanceId() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceId", reflect.TypeOf((*MockMachine)(nil).InstanceId))
}

// IsManual mocks base method
func (m *MockMachine) IsManual() (bool, error) {
	ret := m.ctrl.Call(m, "IsManual")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsManual indicates an expected call of IsManual
func (mr *MockMachineMockRecorder) IsManual() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsManual", reflect.TypeOf((*MockMachine)(nil).IsManual))
}

// MachineTag mocks base method
func (m *MockMachine) MachineTag() names_v3.MachineTag {
	ret := m.ctrl.Call(m, "MachineTag")
	ret0, _ := ret[0].(names_v3.MachineTag)
	return ret0
}

// MachineTag indicates an expected call of MachineTag
func (mr *MockMachineMockRecorder) MachineTag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MachineTag", reflect.TypeOf((*MockMachine)(nil).MachineTag))
}

// Raw mocks base method
func (m *MockMachine) Raw() *state.Machine {
	ret := m.ctrl.Call(m, "Raw")
	ret0, _ := ret[0].(*state.Machine)
	return ret0
}

// Raw indicates an expected call of Raw
func (mr *MockMachineMockRecorder) Raw() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Raw", reflect.TypeOf((*MockMachine)(nil).Raw))
}

// RemoveAllAddresses mocks base method
func (m *MockMachine) RemoveAllAddresses() error {
	ret := m.ctrl.Call(m, "RemoveAllAddresses")
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAllAddresses indicates an expected call of RemoveAllAddresses
func (mr *MockMachineMockRecorder) RemoveAllAddresses() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAllAddresses", reflect.TypeOf((*MockMachine)(nil).RemoveAllAddresses))
}

// SetConstraints mocks base method
func (m *MockMachine) SetConstraints(arg0 constraints.Value) error {
	ret := m.ctrl.Call(m, "SetConstraints", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConstraints indicates an expected call of SetConstraints
func (mr *MockMachineMockRecorder) SetConstraints(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConstraints", reflect.TypeOf((*MockMachine)(nil).SetConstraints), arg0)
}

// SetDevicesAddresses mocks base method
func (m *MockMachine) SetDevicesAddresses(arg0 ...state.LinkLayerDeviceAddress) error {
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetDevicesAddresses", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDevicesAddresses indicates an expected call of SetDevicesAddresses
func (mr *MockMachineMockRecorder) SetDevicesAddresses(arg0 ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDevicesAddresses", reflect.TypeOf((*MockMachine)(nil).SetDevicesAddresses), arg0...)
}

// SetLinkLayerDevices mocks base method
func (m *MockMachine) SetLinkLayerDevices(arg0 ...state.LinkLayerDeviceArgs) error {
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetLinkLayerDevices", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLinkLayerDevices indicates an expected call of SetLinkLayerDevices
func (mr *MockMachineMockRecorder) SetLinkLayerDevices(arg0 ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLinkLayerDevices", reflect.TypeOf((*MockMachine)(nil).SetLinkLayerDevices), arg0...)
}

// SetParentLinkLayerDevicesBeforeTheirChildren mocks base method
func (m *MockMachine) SetParentLinkLayerDevicesBeforeTheirChildren(arg0 []state.LinkLayerDeviceArgs) error {
	ret := m.ctrl.Call(m, "SetParentLinkLayerDevicesBeforeTheirChildren", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetParentLinkLayerDevicesBeforeTheirChildren indicates an expected call of SetParentLinkLayerDevicesBeforeTheirChildren
func (mr *MockMachineMockRecorder) SetParentLinkLayerDevicesBeforeTheirChildren(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetParentLinkLayerDevicesBeforeTheirChildren", reflect.TypeOf((*MockMachine)(nil).SetParentLinkLayerDevicesBeforeTheirChildren), arg0)
}

// Units mocks base method
func (m *MockMachine) Units() ([]provisioner.Unit, error) {
	ret := m.ctrl.Call(m, "Units")
	ret0, _ := ret[0].([]provisioner.Unit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Units indicates an expected call of Units
func (mr *MockMachineMockRecorder) Units() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Units", reflect.TypeOf((*MockMachine)(nil).Units))
}

// MockBridgePolicy is a mock of BridgePolicy interface
type MockBridgePolicy struct {
	ctrl     *gomock.Controller
	recorder *MockBridgePolicyMockRecorder
}

// MockBridgePolicyMockRecorder is the mock recorder for MockBridgePolicy
type MockBridgePolicyMockRecorder struct {
	mock *MockBridgePolicy
}

// NewMockBridgePolicy creates a new mock instance
func NewMockBridgePolicy(ctrl *gomock.Controller) *MockBridgePolicy {
	mock := &MockBridgePolicy{ctrl: ctrl}
	mock.recorder = &MockBridgePolicyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBridgePolicy) EXPECT() *MockBridgePolicyMockRecorder {
	return m.recorder
}

// FindMissingBridgesForContainer mocks base method
func (m *MockBridgePolicy) FindMissingBridgesForContainer(arg0 containerizer.Machine, arg1 containerizer.Container) ([]network.DeviceToBridge, int, error) {
	ret := m.ctrl.Call(m, "FindMissingBridgesForContainer", arg0, arg1)
	ret0, _ := ret[0].([]network.DeviceToBridge)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindMissingBridgesForContainer indicates an expected call of FindMissingBridgesForContainer
func (mr *MockBridgePolicyMockRecorder) FindMissingBridgesForContainer(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMissingBridgesForContainer", reflect.TypeOf((*MockBridgePolicy)(nil).FindMissingBridgesForContainer), arg0, arg1)
}

// PopulateContainerLinkLayerDevices mocks base method
func (m *MockBridgePolicy) PopulateContainerLinkLayerDevices(arg0 containerizer.Machine, arg1 containerizer.Container) error {
	ret := m.ctrl.Call(m, "PopulateContainerLinkLayerDevices", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PopulateContainerLinkLayerDevices indicates an expected call of PopulateContainerLinkLayerDevices
func (mr *MockBridgePolicyMockRecorder) PopulateContainerLinkLayerDevices(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopulateContainerLinkLayerDevices", reflect.TypeOf((*MockBridgePolicy)(nil).PopulateContainerLinkLayerDevices), arg0, arg1)
}

// MockUnit is a mock of Unit interface
type MockUnit struct {
	ctrl     *gomock.Controller
	recorder *MockUnitMockRecorder
}

// MockUnitMockRecorder is the mock recorder for MockUnit
type MockUnitMockRecorder struct {
	mock *MockUnit
}

// NewMockUnit creates a new mock instance
func NewMockUnit(ctrl *gomock.Controller) *MockUnit {
	mock := &MockUnit{ctrl: ctrl}
	mock.recorder = &MockUnitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUnit) EXPECT() *MockUnitMockRecorder {
	return m.recorder
}

// Application mocks base method
func (m *MockUnit) Application() (provisioner.Application, error) {
	ret := m.ctrl.Call(m, "Application")
	ret0, _ := ret[0].(provisioner.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Application indicates an expected call of Application
func (mr *MockUnitMockRecorder) Application() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockUnit)(nil).Application))
}

// Name mocks base method
func (m *MockUnit) Name() string {
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockUnitMockRecorder) Name() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockUnit)(nil).Name))
}

// MockApplication is a mock of Application interface
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// Charm mocks base method
func (m *MockApplication) Charm() (provisioner.Charm, bool, error) {
	ret := m.ctrl.Call(m, "Charm")
	ret0, _ := ret[0].(provisioner.Charm)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Charm indicates an expected call of Charm
func (mr *MockApplicationMockRecorder) Charm() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Charm", reflect.TypeOf((*MockApplication)(nil).Charm))
}

// Name mocks base method
func (m *MockApplication) Name() string {
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockApplicationMockRecorder) Name() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockApplication)(nil).Name))
}

// MockCharm is a mock of Charm interface
type MockCharm struct {
	ctrl     *gomock.Controller
	recorder *MockCharmMockRecorder
}

// MockCharmMockRecorder is the mock recorder for MockCharm
type MockCharmMockRecorder struct {
	mock *MockCharm
}

// NewMockCharm creates a new mock instance
func NewMockCharm(ctrl *gomock.Controller) *MockCharm {
	mock := &MockCharm{ctrl: ctrl}
	mock.recorder = &MockCharmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCharm) EXPECT() *MockCharmMockRecorder {
	return m.recorder
}

// LXDProfile mocks base method
func (m *MockCharm) LXDProfile() *charm_v6.LXDProfile {
	ret := m.ctrl.Call(m, "LXDProfile")
	ret0, _ := ret[0].(*charm_v6.LXDProfile)
	return ret0
}

// LXDProfile indicates an expected call of LXDProfile
func (mr *MockCharmMockRecorder) LXDProfile() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LXDProfile", reflect.TypeOf((*MockCharm)(nil).LXDProfile))
}

// Revision mocks base method
func (m *MockCharm) Revision() int {
	ret := m.ctrl.Call(m, "Revision")
	ret0, _ := ret[0].(int)
	return ret0
}

// Revision indicates an expected call of Revision
func (mr *MockCharmMockRecorder) Revision() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revision", reflect.TypeOf((*MockCharm)(nil).Revision))
}