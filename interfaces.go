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

// Patcher is a testing tool intended to aid in the construction of
// tests which "patch" the source in some fashion.  This is not monkey
// patching, like is performed with dynamic languages like Python;
// Patcher requires that the source code be written with the patching
// in mind.  For instance, a particular function called by the source
// may be assigned to a variable, then that variable used in the call.
//
// The primary API in this package is the Patcher interface.
// Something that implements the Patcher interface has two idempotent
// methods: Install and Restore; for convenience, both return the same
// Patcher.
//
// This package provides three implementations of Patcher.  The first
// is MockPatcher, which is provided for testing code that manipulates
// a Patcher; most users of this package will not find this type
// useful.  The more useful Patcher implementations are created with
// SetVar and NewPatchMaster.

package patcher

// Patcher is an interface for patchers.  Patchers have Install and
// Restore methods.
type Patcher interface {
	// Install installs the patch.  It should store metadata
	// sufficient to allow Restore to restore the original data.
	// This method must be idempotent.
	Install() Patcher

	// Restore uses the metadata stored by Install to restore the
	// patch to its original value.  This method must be
	// idempotent.
	Restore() Patcher
}
