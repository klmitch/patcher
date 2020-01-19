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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatchMasterImplementsPatcher(t *testing.T) {
	assert.Implements(t, (*Patcher)(nil), &PatchMaster{})
}

func TestNewPatchMaster(t *testing.T) {
	p1 := &MockPatcher{}
	p2 := &MockPatcher{}

	pm := NewPatchMaster(p1, p2)

	assert.Len(t, pm.patches, 2)
	assert.Same(t, p1, pm.patches[0])
	assert.Same(t, p2, pm.patches[1])
}

type OrderPatcher struct {
	ordering *[]string
	name     string
}

func (op OrderPatcher) Install() Patcher {
	*op.ordering = append(*op.ordering, fmt.Sprintf("Install %s", op.name))

	return op
}

func (op OrderPatcher) Restore() Patcher {
	*op.ordering = append(*op.ordering, fmt.Sprintf("Restore %s", op.name))

	return op
}

func TestPatchMasterInstall(t *testing.T) {
	ordering := []string{}
	pm := NewPatchMaster(
		OrderPatcher{
			ordering: &ordering,
			name:     "patch1",
		},
		OrderPatcher{
			ordering: &ordering,
			name:     "patch2",
		},
	)

	result := pm.Install()

	assert.Equal(t, []string{"Install patch1", "Install patch2"}, ordering)
	assert.Same(t, pm, result)
}

func TestPatchMasterRestore(t *testing.T) {
	ordering := []string{}
	pm := NewPatchMaster(
		OrderPatcher{
			ordering: &ordering,
			name:     "patch1",
		},
		OrderPatcher{
			ordering: &ordering,
			name:     "patch2",
		},
	)

	result := pm.Restore()

	assert.Equal(t, []string{"Restore patch2", "Restore patch1"}, ordering)
	assert.Same(t, pm, result)
}

func TestPatchMasterAdd(t *testing.T) {
	p1 := &MockPatcher{}
	p2 := &MockPatcher{}
	pm := NewPatchMaster(p1)

	result := pm.Add(p2)

	assert.Same(t, p2, result)
	assert.Len(t, pm.patches, 2)
	assert.Same(t, p1, pm.patches[0])
	assert.Same(t, p2, pm.patches[1])
}
