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

// PatchMaster is a patcher that handles multiple patchers.  Patchers
// can be passed in to the constructor (NewPatchMaster), or can be
// added using the Add method.
type PatchMaster struct {
	patches []Patcher
}

// NewPatchMaster constructs a new PatchMaster.  It could be used in a
// test function like so:
//
//	func TestSomething(t *testing.T) {
//		pm := NewPatchMaster(
//			SetVar(&var1, "value1"),
//			SetVar(&var2, "value2"),
//		)
//		defer pm.Install().Restore()
//
//		// Do some tests
//
//		// Patch an additional variable
//		pm.Add(SetVar(&var3, "value3")).Install()
//
//		// Do some more tests
//	}
func NewPatchMaster(patches ...Patcher) *PatchMaster {
	return &PatchMaster{patches: patches}
}

// Install installs the patch.  It should store metadata sufficient to
// allow Restore to restore the original data.  This method must be
// idempotent.
func (pm *PatchMaster) Install() Patcher {
	for _, patch := range pm.patches {
		patch.Install()
	}

	return pm
}

// Restore uses the metadata stored by Install to restore the patch to
// its original value.  This method must be idempotent.
func (pm *PatchMaster) Restore() Patcher {
	// Walk the patches in reverse for the restore
	for i := len(pm.patches) - 1; i >= 0; i-- {
		pm.patches[i].Restore()
	}

	return pm
}

// Add adds a new patcher to the PatchMaster.  For convenience, it
// returns the patcher it just added.
func (pm *PatchMaster) Add(patch Patcher) Patcher {
	pm.patches = append(pm.patches, patch)

	return patch
}
