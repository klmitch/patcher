// Copyright (c) 2020 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package patcher

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariableSetterImplementsPatcher(t *testing.T) {
	assert.Implements(t, (*Patcher)(nil), &VariableSetter{})
}

func TestSetVarBase(t *testing.T) {
	variable := "unpatched" //nolint:goconst

	vs := SetVar(&variable, "patched")

	assert.Equal(t, reflect.ValueOf(&variable).Elem(), vs.variable)
	assert.Equal(t, "patched", vs.value.Interface())
	assert.False(t, vs.original.IsValid())
	assert.False(t, vs.applied)
}

func TestSetVarUnsettable(t *testing.T) {
	variable := "unpatched"

	assert.PanicsWithValue(t, "cannot set variable passed to SetVar!", func() {
		SetVar(variable, "patched")
	})
}

func TestSetVarUnassignable(t *testing.T) {
	variable := "unpatched"

	assert.PanicsWithValue(t, "cannot assign int type to variable type string", func() {
		SetVar(&variable, 12345)
	})
}

func TestSetVarFunc(t *testing.T) {
	variable := func() error {
		return nil
	}

	vs := SetVar(&variable, func() error {
		return assert.AnError
	})

	assert.Equal(t, reflect.ValueOf(&variable).Elem(), vs.variable)
	assert.False(t, vs.original.IsValid())
	assert.False(t, vs.applied)
}

func TestSetVarIncompatibleFunc(t *testing.T) {
	variable := func() error {
		return nil
	}

	assert.Panics(t, func() {
		SetVar(&variable, func() (int, error) {
			return 12345, assert.AnError
		})
	})
}

func TestVariableSetterInstallBase(t *testing.T) {
	variable := "unpatched"
	vs := SetVar(&variable, "patched")

	result := vs.Install()

	assert.Same(t, vs, result)
	assert.Equal(t, "patched", variable)
	assert.Equal(t, "unpatched", vs.original.Interface())
	assert.True(t, vs.applied)
}

func TestVariableSetterInstallIdempotent(t *testing.T) {
	variable := "unpatched"
	vs := SetVar(&variable, "patched")
	vs.applied = true

	result := vs.Install()

	assert.Same(t, vs, result)
	assert.Equal(t, "unpatched", variable)
	assert.True(t, vs.applied)
}

func TestVariableSetterRestoreBase(t *testing.T) {
	variable := "patched"
	vs := SetVar(&variable, "patched")
	vs.original = reflect.ValueOf("unpatched")
	vs.applied = true

	result := vs.Restore()

	assert.Same(t, vs, result)
	assert.Equal(t, "unpatched", variable)
	assert.False(t, vs.applied)
}

func TestVariableSetterRestoreIdempotent(t *testing.T) {
	variable := "patched"
	vs := SetVar(&variable, "patched")
	vs.original = reflect.ValueOf("unpatched")

	result := vs.Restore()

	assert.Same(t, vs, result)
	assert.Equal(t, "patched", variable)
	assert.False(t, vs.applied)
}

var testingVar = "unpatched"

func testVariableSetterFunctionInner(t *testing.T) {
	vs := SetVar(&testingVar, "patched")
	defer vs.Restore()

	assert.Equal(t, "unpatched", testingVar)

	vs.Install()

	assert.Equal(t, "patched", testingVar)
}

func TestVariableSetterFunction(t *testing.T) {
	assert.Equal(t, "unpatched", testingVar)

	testVariableSetterFunctionInner(t)

	assert.Equal(t, "unpatched", testingVar)
}
