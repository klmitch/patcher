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

import "os"

// EnvPatcher is a patcher that, given an environment variable name,
// will set or unset that environment variable.
type EnvPatcher struct {
	name     string
	value    *string
	original *string
	applied  bool
}

// Patch points for testing the routines in this file.
var (
	setenv    = os.Setenv
	lookupenv = os.LookupEnv
	unsetenv  = os.Unsetenv
)

// setEnv is a helper for the EnvPatcher that sets or unsets an
// environment variable depending on whether the value pointer is nil
// or a string.  It will panic if there is an error.
func setEnv(name string, value *string) {
	var err error
	if value == nil {
		// Only unset if it's set
		if _, ok := lookupenv(name); ok {
			err = unsetenv(name)
		}
	} else {
		err = setenv(name, *value)
	}

	// If there was an error setting or unsetting the environment
	// variable, panic
	if err != nil {
		panic(err)
	}
}

// SetEnv constructs an EnvPatcher, storing the desired value of the
// specified environment variable.  It could be used in a test
// function like so:
//
//	func TestDoSomething(t *testing.T) {
//		defer SetEnv("VARNAME", "value").Install().Restore()
//
//		err := DoSomething("some-filename")
//
//		if err != nil {
//			t.Fail("non-nil error!")
//		}
//	}
func SetEnv(name, value string) *EnvPatcher {
	return &EnvPatcher{
		name:  name,
		value: &value,
	}
}

// UnsetEnv constructs an EnvPatcher.  The specified environment
// variable will be unset when the patch is installed.  It could be
// used in a test function like so:
//
//	func TestDoSomething(t *testing.T) {
//		defer UnsetEnv("FILENAME").Install().Restore()
//
//		err := DoSomething()
//
//		if err != nil {
//			t.Fail("non-nil error!")
//		}
//	}
func UnsetEnv(name string) *EnvPatcher {
	return &EnvPatcher{
		name: name,
	}
}

// Install installs the patch.  It should store metadata sufficient to
// allow Restore to restore the original data.  This method must be
// idempotent.
func (ep *EnvPatcher) Install() Patcher {
	// Be idempotent
	if ep.applied {
		return ep
	}

	// Save the current value of the environment variable
	if value, ok := lookupenv(ep.name); ok {
		ep.original = &value
	} else {
		ep.original = nil
	}

	// Set the environment variable to the desired value
	setEnv(ep.name, ep.value)
	ep.applied = true

	return ep
}

// Restore uses the metadata stored by Install to restore the patch to
// its original value.  This method must be idempotent.
func (ep *EnvPatcher) Restore() Patcher {
	// Be idempotent
	if !ep.applied {
		return ep
	}

	// Restore the environment variable to the original value
	setEnv(ep.name, ep.original)
	ep.applied = false

	return ep
}
