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
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogPatcherImplementsPatcher(t *testing.T) {
	assert.Implements(t, (*Patcher)(nil), &LogPatcher{})
}

func TestLog(t *testing.T) {
	value := &bytes.Buffer{}

	lp := Log(value)

	assert.Equal(t, &LogPatcher{
		value: value,
	}, lp)
}

func TestLogPatcherInstallBase(t *testing.T) {
	value := &bytes.Buffer{}
	original := log.Writer()
	defer func() { log.SetOutput(original) }()
	lp := &LogPatcher{
		value: value,
	}

	result := lp.Install()

	assert.Same(t, lp, result)
	assert.Same(t, original, lp.original)
	assert.Same(t, value, log.Writer())
	assert.True(t, lp.applied)
}

func TestLogPatcherInstallIdempotent(t *testing.T) {
	value := &bytes.Buffer{}
	original := log.Writer()
	defer func() { log.SetOutput(original) }()
	lp := &LogPatcher{
		value:   value,
		applied: true,
	}

	result := lp.Install()

	assert.Same(t, lp, result)
	assert.Nil(t, lp.original)
	assert.Same(t, original, log.Writer())
	assert.True(t, lp.applied)
}

func TestLogPatcherRestoreBase(t *testing.T) {
	value := &bytes.Buffer{}
	original := log.Writer()
	defer func() { log.SetOutput(original) }()
	lp := &LogPatcher{
		value:    original,
		original: value,
		applied:  true,
	}

	result := lp.Restore()

	assert.Same(t, lp, result)
	assert.Same(t, value, log.Writer())
	assert.False(t, lp.applied)
}

func TestLogPatcherRestoreIdempotent(t *testing.T) {
	value := &bytes.Buffer{}
	original := log.Writer()
	defer func() { log.SetOutput(original) }()
	lp := &LogPatcher{
		original: value,
	}

	result := lp.Restore()

	assert.Same(t, lp, result)
	assert.Same(t, original, log.Writer())
	assert.False(t, lp.applied)
}
