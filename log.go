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
	"io"
	"log"
)

// LogPatcher is a patcher that, given a io.Writer, will update the
// output of the default logger in the log package.
type LogPatcher struct {
	value    io.Writer
	original io.Writer
	applied  bool
}

// Log constructs a LogPatcher, storing the desired io.Writer to use
// with the default logger.  It could be used in a test function like
// so:
//
//	func TestDoSomething(t *testing.T) {
//		logStream := &bytes.Buffer{}
//		defer Log(logStream).Install().Restore()
//
//		err := DoSomething("some-filename")
//
//		if logStream.String() != "Error reading file" {
//			t.Fail("failed to log!")
//		}
//	}
func Log(value io.Writer) *LogPatcher {
	return &LogPatcher{
		value: value,
	}
}

// Install installs the patch.  It should store metadata sufficient to
// allow Restore to restore the original data.  This method must be
// idempotent.
func (lp *LogPatcher) Install() Patcher {
	// Be idempotent
	if lp.applied {
		return lp
	}

	// Save the current value of the log output
	lp.original = log.Writer()

	// Set the patch value
	log.SetOutput(lp.value)
	lp.applied = true

	return lp
}

// Restore uses the metadata stored by Install to restore the patch to
// its original value.  This method must be idempotent.
func (lp *LogPatcher) Restore() Patcher {
	// Be idempotent
	if !lp.applied {
		return lp
	}

	// Restore the original log output
	log.SetOutput(lp.original)
	lp.applied = false

	return lp
}
