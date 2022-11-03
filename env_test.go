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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvPatcherImplementsPatcher(t *testing.T) {
	assert.Implements(t, (*Patcher)(nil), &EnvPatcher{})
}

func TestInternalSetEnvUnsetExists(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "value", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "", true
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()

	setEnv("ENV", nil)

	assert.False(t, setenvCalled)
	assert.True(t, lookupenvCalled)
	assert.True(t, unsetenvCalled)
}

func TestInternalSetEnvUnsetMissing(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "value", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "", false
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()

	setEnv("ENV", nil)

	assert.False(t, setenvCalled)
	assert.True(t, lookupenvCalled)
	assert.False(t, unsetenvCalled)
}

func TestInternalSetEnvUnsetFails(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "value", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "", true
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return assert.AnError
		}),
	).Install().Restore()

	assert.PanicsWithValue(t, assert.AnError, func() { setEnv("ENV", nil) })
	assert.False(t, setenvCalled)
	assert.True(t, lookupenvCalled)
	assert.True(t, unsetenvCalled)
}

func TestInternalSetEnvSetBase(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "value", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "", true
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()
	value := "value" //nolint:goconst

	setEnv("ENV", &value)

	assert.True(t, setenvCalled)
	assert.False(t, lookupenvCalled)
	assert.False(t, unsetenvCalled)
}

func TestInternalSetEnvSetFails(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "value", v)
			setenvCalled = true
			return assert.AnError
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "", true
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()
	value := "value"

	assert.PanicsWithValue(t, assert.AnError, func() { setEnv("ENV", &value) })
	assert.True(t, setenvCalled)
	assert.False(t, lookupenvCalled)
	assert.False(t, unsetenvCalled)
}

func TestSetEnv(t *testing.T) {
	result := SetEnv("ENV", "value")

	assert.Equal(t, "ENV", result.name)
	assert.Equal(t, "value", *result.value)
	assert.Nil(t, result.original)
	assert.False(t, result.applied)
}

func TestUnsetEnv(t *testing.T) {
	result := UnsetEnv("ENV")

	assert.Equal(t, "ENV", result.name)
	assert.Nil(t, result.value)
	assert.Nil(t, result.original)
	assert.False(t, result.applied)
}

func TestEnvPatcherInstallExists(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "value", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "original", true //nolint:goconst
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()
	value := "value"
	obj := &EnvPatcher{
		name:  "ENV",
		value: &value,
	}

	result := obj.Install()

	assert.Same(t, obj, result)
	assert.Equal(t, "original", *obj.original)
	assert.True(t, obj.applied)
	assert.True(t, setenvCalled)
	assert.True(t, lookupenvCalled)
	assert.False(t, unsetenvCalled)
}

func TestEnvPatcherInstallMissing(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "value", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "", false
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()
	value := "value"
	obj := &EnvPatcher{
		name:  "ENV",
		value: &value,
	}

	result := obj.Install()

	assert.Same(t, obj, result)
	assert.Nil(t, obj.original)
	assert.True(t, obj.applied)
	assert.True(t, setenvCalled)
	assert.True(t, lookupenvCalled)
	assert.False(t, unsetenvCalled)
}

func TestEnvPatcherInstallIdempotent(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "value", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "original", true
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()
	value := "value"
	original := "original"
	obj := &EnvPatcher{
		name:     "ENV",
		value:    &value,
		original: &original,
		applied:  true,
	}

	result := obj.Install()

	assert.Same(t, obj, result)
	assert.True(t, obj.applied)
	assert.False(t, setenvCalled)
	assert.False(t, lookupenvCalled)
	assert.False(t, unsetenvCalled)
}

func TestEnvPatcherRestoreBase(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "original", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "original", true
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()
	original := "original"
	obj := &EnvPatcher{
		name:     "ENV",
		original: &original,
		applied:  true,
	}

	result := obj.Restore()

	assert.Same(t, obj, result)
	assert.False(t, obj.applied)
	assert.True(t, setenvCalled)
	assert.False(t, lookupenvCalled)
	assert.False(t, unsetenvCalled)
}

func TestEnvPatcherRestoreIdempotent(t *testing.T) {
	setenvCalled := false
	lookupenvCalled := false
	unsetenvCalled := false
	defer NewPatchMaster(
		SetVar(&setenv, func(n, v string) error {
			assert.Equal(t, "ENV", n)
			assert.Equal(t, "original", v)
			setenvCalled = true
			return nil
		}),
		SetVar(&lookupenv, func(n string) (string, bool) {
			assert.Equal(t, "ENV", n)
			lookupenvCalled = true
			return "original", true
		}),
		SetVar(&unsetenv, func(n string) error {
			assert.Equal(t, "ENV", n)
			unsetenvCalled = true
			return nil
		}),
	).Install().Restore()
	original := "original"
	obj := &EnvPatcher{
		name:     "ENV",
		original: &original,
	}

	result := obj.Restore()

	assert.Same(t, obj, result)
	assert.False(t, obj.applied)
	assert.False(t, setenvCalled)
	assert.False(t, lookupenvCalled)
	assert.False(t, unsetenvCalled)
}
