// This file was generated by counterfeiter
package fakes

import (
	"os"
	"sync"

	"github.com/IBM/ubiquity/utils"
)

type FakeExecutor struct {
	ExecuteStub        func(command string, args []string) ([]byte, error)
	executeMutex       sync.RWMutex
	executeArgsForCall []struct {
		command string
		args    []string
	}
	executeReturns struct {
		result1 []byte
		result2 error
	}
	StatStub        func(string) (os.FileInfo, error)
	statMutex       sync.RWMutex
	statArgsForCall []struct {
		arg1 string
	}
	statReturns struct {
		result1 os.FileInfo
		result2 error
	}
	MkdirStub        func(string, os.FileMode) error
	mkdirMutex       sync.RWMutex
	mkdirArgsForCall []struct {
		arg1 string
		arg2 os.FileMode
	}
	mkdirReturns struct {
		result1 error
	}
	RemoveAllStub        func(string) error
	removeAllMutex       sync.RWMutex
	removeAllArgsForCall []struct {
		arg1 string
	}
	removeAllReturns struct {
		result1 error
	}
	HostnameStub        func() (string, error)
	hostnameMutex       sync.RWMutex
	hostnameArgsForCall []struct{}
	hostnameReturns     struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeExecutor) Execute(command string, args []string) ([]byte, error) {
	var argsCopy []string
	if args != nil {
		argsCopy = make([]string, len(args))
		copy(argsCopy, args)
	}
	fake.executeMutex.Lock()
	fake.executeArgsForCall = append(fake.executeArgsForCall, struct {
		command string
		args    []string
	}{command, argsCopy})
	fake.recordInvocation("Execute", []interface{}{command, argsCopy})
	fake.executeMutex.Unlock()
	if fake.ExecuteStub != nil {
		return fake.ExecuteStub(command, args)
	}
	return fake.executeReturns.result1, fake.executeReturns.result2
}

func (fake *FakeExecutor) ExecuteCallCount() int {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return len(fake.executeArgsForCall)
}

func (fake *FakeExecutor) ExecuteArgsForCall(i int) (string, []string) {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return fake.executeArgsForCall[i].command, fake.executeArgsForCall[i].args
}

func (fake *FakeExecutor) ExecuteReturns(result1 []byte, result2 error) {
	fake.ExecuteStub = nil
	fake.executeReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeExecutor) Stat(arg1 string) (os.FileInfo, error) {
	fake.statMutex.Lock()
	fake.statArgsForCall = append(fake.statArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Stat", []interface{}{arg1})
	fake.statMutex.Unlock()
	if fake.StatStub != nil {
		return fake.StatStub(arg1)
	}
	return fake.statReturns.result1, fake.statReturns.result2
}

func (fake *FakeExecutor) StatCallCount() int {
	fake.statMutex.RLock()
	defer fake.statMutex.RUnlock()
	return len(fake.statArgsForCall)
}

func (fake *FakeExecutor) StatArgsForCall(i int) string {
	fake.statMutex.RLock()
	defer fake.statMutex.RUnlock()
	return fake.statArgsForCall[i].arg1
}

func (fake *FakeExecutor) StatReturns(result1 os.FileInfo, result2 error) {
	fake.StatStub = nil
	fake.statReturns = struct {
		result1 os.FileInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeExecutor) Mkdir(arg1 string, arg2 os.FileMode) error {
	fake.mkdirMutex.Lock()
	fake.mkdirArgsForCall = append(fake.mkdirArgsForCall, struct {
		arg1 string
		arg2 os.FileMode
	}{arg1, arg2})
	fake.recordInvocation("Mkdir", []interface{}{arg1, arg2})
	fake.mkdirMutex.Unlock()
	if fake.MkdirStub != nil {
		return fake.MkdirStub(arg1, arg2)
	}
	return fake.mkdirReturns.result1
}

func (fake *FakeExecutor) MkdirCallCount() int {
	fake.mkdirMutex.RLock()
	defer fake.mkdirMutex.RUnlock()
	return len(fake.mkdirArgsForCall)
}

func (fake *FakeExecutor) MkdirArgsForCall(i int) (string, os.FileMode) {
	fake.mkdirMutex.RLock()
	defer fake.mkdirMutex.RUnlock()
	return fake.mkdirArgsForCall[i].arg1, fake.mkdirArgsForCall[i].arg2
}

func (fake *FakeExecutor) MkdirReturns(result1 error) {
	fake.MkdirStub = nil
	fake.mkdirReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeExecutor) RemoveAll(arg1 string) error {
	fake.removeAllMutex.Lock()
	fake.removeAllArgsForCall = append(fake.removeAllArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("RemoveAll", []interface{}{arg1})
	fake.removeAllMutex.Unlock()
	if fake.RemoveAllStub != nil {
		return fake.RemoveAllStub(arg1)
	}
	return fake.removeAllReturns.result1
}

func (fake *FakeExecutor) RemoveAllCallCount() int {
	fake.removeAllMutex.RLock()
	defer fake.removeAllMutex.RUnlock()
	return len(fake.removeAllArgsForCall)
}

func (fake *FakeExecutor) RemoveAllArgsForCall(i int) string {
	fake.removeAllMutex.RLock()
	defer fake.removeAllMutex.RUnlock()
	return fake.removeAllArgsForCall[i].arg1
}

func (fake *FakeExecutor) RemoveAllReturns(result1 error) {
	fake.RemoveAllStub = nil
	fake.removeAllReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeExecutor) Hostname() (string, error) {
	fake.hostnameMutex.Lock()
	fake.hostnameArgsForCall = append(fake.hostnameArgsForCall, struct{}{})
	fake.recordInvocation("Hostname", []interface{}{})
	fake.hostnameMutex.Unlock()
	if fake.HostnameStub != nil {
		return fake.HostnameStub()
	}
	return fake.hostnameReturns.result1, fake.hostnameReturns.result2
}

func (fake *FakeExecutor) HostnameCallCount() int {
	fake.hostnameMutex.RLock()
	defer fake.hostnameMutex.RUnlock()
	return len(fake.hostnameArgsForCall)
}

func (fake *FakeExecutor) HostnameReturns(result1 string, result2 error) {
	fake.HostnameStub = nil
	fake.hostnameReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeExecutor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	fake.statMutex.RLock()
	defer fake.statMutex.RUnlock()
	fake.mkdirMutex.RLock()
	defer fake.mkdirMutex.RUnlock()
	fake.removeAllMutex.RLock()
	defer fake.removeAllMutex.RUnlock()
	fake.hostnameMutex.RLock()
	defer fake.hostnameMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeExecutor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ utils.Executor = new(FakeExecutor)