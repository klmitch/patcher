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
	"fmt"
	"reflect"
)

// VariableSetter is a patcher that, given a pointer to a variable and
// the desired patch value, will set that variable to that value.
type VariableSetter struct {
	variable reflect.Value
	value    reflect.Value
	original reflect.Value
	applied  bool
}

// SetVar constructs a VariableSetter, storing the variable and its
// desired new value.  It could be used in a test function like so:
//
//	func TestDoSomething(t *testing.T) {
//		defer SetVar(&readFile, func(filename string) ([]byte, error) {
//			return []byte("hello"), nil
//		}).Install().Restore()
//
//		err := DoSomething("some-filename")
//
//		if err != nil {
//			t.Fail("non-nil error!")
//		}
//	}
func SetVar(variable, value interface{}) *VariableSetter {
	// Select the variable and validate it's a settable object
	varReflect := reflect.ValueOf(variable)
	if varReflect.Type().Kind() != reflect.Ptr {
		panic("cannot set variable passed to SetVar!")
	}
	v := varReflect.Elem()

	// Convert the desired value and check that it can be assigned
	// to the variable
	val := reflect.ValueOf(value)
	if !val.Type().AssignableTo(v.Type()) {
		panic(fmt.Sprintf("cannot assign %s type to variable type %s", val.Type(), v.Type()))
	}

	return &VariableSetter{
		variable: v,
		value:    val,
	}
}

// Install installs the patch.  It should store metadata sufficient to
// allow Restore to restore the original data.  This method must be
// idempotent.
func (vs *VariableSetter) Install() Patcher {
	// Be idempotent
	if vs.applied {
		return vs
	}

	// Save the current value of the variable
	vs.original = reflect.ValueOf(vs.variable.Interface())

	// Set the new value and store that it's applied
	vs.variable.Set(vs.value)
	vs.applied = true

	return vs
}

// Restore uses the metadata stored by Install to restore the patch to
// its original value.  This method must be idempotent.
func (vs *VariableSetter) Restore() Patcher {
	// Be idempotent
	if !vs.applied {
		return vs
	}

	// Restore the variable's original value and clear the applied
	// flag
	vs.variable.Set(vs.original)
	vs.applied = false

	return vs
}
