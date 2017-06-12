// Code generated by mockery v1.0.0
package architecture

import mock "github.com/stretchr/testify/mock"

// PriorityQueue is an autogenerated mock type for the PriorityQueue type
type MockPriorityQueue struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *MockPriorityQueue) Delete(id string) PriorityQueueItem {
	ret := _m.Called(id)

	var r0 PriorityQueueItem
	if rf, ok := ret.Get(0).(func(string) PriorityQueueItem); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PriorityQueueItem)
		}
	}

	return r0
}

// Dequeue provides a mock function with given fields:
func (_m *MockPriorityQueue) Dequeue() PriorityQueueItem {
	ret := _m.Called()

	var r0 PriorityQueueItem
	if rf, ok := ret.Get(0).(func() PriorityQueueItem); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PriorityQueueItem)
		}
	}

	return r0
}

// Enqueue provides a mock function with given fields: item
func (_m *MockPriorityQueue) Enqueue(item PriorityQueueItem) {
	_m.Called(item)
}

// Find provides a mock function with given fields: id
func (_m *MockPriorityQueue) Find(id string) PriorityQueueItem {
	ret := _m.Called(id)

	var r0 PriorityQueueItem
	if rf, ok := ret.Get(0).(func(string) PriorityQueueItem); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PriorityQueueItem)
		}
	}

	return r0
}

// Init provides a mock function with given fields:
func (_m *MockPriorityQueue) Init(tubeName string) {
	_m.Called()
}

// Peek provides a mock function with given fields:
func (_m *MockPriorityQueue) Peek() PriorityQueueItem {
	ret := _m.Called()

	var r0 PriorityQueueItem
	if rf, ok := ret.Get(0).(func() PriorityQueueItem); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PriorityQueueItem)
		}
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *MockPriorityQueue) Size() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
