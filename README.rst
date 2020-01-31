=======
Patcher
=======

.. image:: https://img.shields.io/github/tag/klmitch/patcher.svg
    :target: https://github.com/klmitch/patcher/tags
.. image:: https://img.shields.io/hexpm/l/plug.svg
    :target: https://github.com/klmitch/patcher/blob/master/LICENSE
.. image:: https://travis-ci.org/klmitch/patcher.svg?branch=master
    :target: https://travis-ci.org/klmitch/patcher
.. image:: https://coveralls.io/repos/github/klmitch/patcher/badge.svg?branch=master
    :target: https://coveralls.io/github/klmitch/patcher?branch=master
.. image:: https://godoc.org/github.com/klmitch/patcher?status.svg
    :target: http://godoc.org/github.com/klmitch/patcher
.. image:: https://img.shields.io/github/issues/klmitch/patcher.svg
    :target: https://github.com/klmitch/patcher/issues
.. image:: https://img.shields.io/github/issues-pr/klmitch/patcher.svg
    :target: https://github.com/klmitch/patcher/pulls
.. image:: https://goreportcard.com/badge/github.com/klmitch/patcher
    :target: https://goreportcard.com/report/github.com/klmitch/patcher

This repository contains Patcher.  Patcher is a testing tool intended
to aid in the construction of tests which "patch" the source in some
fashion.  This is not monkey patching, like is performed with dynamic
languages like Python; Patcher requires that the source code be
written with the patching in mind.  For instance, a particular
function called by the source may be assigned to a variable, then that
variable used in the call.

Using
=====

The primary API in Patcher is the ``Patcher`` interface.  Something
that implements the ``Patcher`` interface has two idempotent methods:
``Install()`` and ``Restore()``; for convenience, both return the same
``Patcher``.  This means that a patch may be placed with code such
as::

    defer MakeThePatcher().Install().Restore()

The ``Patcher`` will be created and told to install itself, and the
``Restore()`` will be evaluated at the end of the test function.

Patcher provides three implementations of ``Patcher``.  The first is
``MockPatcher``, which is provided for testing code that manipulates
``Patcher``; most users of Patcher will not find this type useful.
The more useful ``Patcher`` implementations are created with
``SetVar()`` and ``NewPatchMaster()``.

``SetVar()``
------------

The ``SetVar()`` function creates an instance of a ``VariableSetter``
struct, which implements ``Patcher``.  The ``SetVar()`` function is
called with an address of a variable and a desired value to assign to
that variable while the patch is active; when the ``Patcher`` is
installed, the current value of that variable is saved for later
restoration and the desired value is assigned.  Note that ``SetVar()``
performs various sanity checks on its arguments, such as verifying
that the variable and the value have compatible types, and will
``panic()`` if those checks fail.

The ``SetVar()`` patcher is the primary use for Patcher in test
suites; it allows the temporary substitution of another function to be
used in code, allowing test cases that would ordinarily be difficult
to test reliably.  The below is an example of how ``SetVar()`` may be
used::

    var readFile func(string) ([]byte, error) = ioutil.ReadFile

    func DoSomething(filename string) error {
    	data, err := readFile(filename)
    	if err != nil {
    		return err
    	}

    	// Do something...

    	return nil
    }

    func TestDoSomething(t *testing.T) {
    	defer SetVar(&readFile, func(filename string) ([]byte, error) {
    		return []byte("hello"), nil
    	}).Install().Restore()

    	err := DoSomething("some-filename")

    	if err != nil {
    		t.Fail("non-nil error!")
    	}
    }

``Log()``
---------

The ``Log()`` function creates an instance of a ``LogPatcher`` struct,
which implements ``Patcher``.  The ``Log()`` function is called with
an ``io.Writer`` to which the output of the default logger from the
``log`` package should be redirected while the patch is active; when
the ``Patcher`` is installed, all output to the default logger will be
redirected to that ``io.Writer`` and the original ``io.Writer`` will
be saved for later restoration.  For instance::

    func DoSomething(filename string) error {
    	data, err := io.ReadFile(filename)
    	if err != nil {
    		log.Printf("Error reading file")
    		return err
    	}

    	// Do something...

    	return nil
    }

    func TestDoSomething(t *testing.T) {
    	logStream := &bytes.Buffer{}
    	defer Log(logStream).Install().Restore()

    	err := DoSomething("some-filename")

    	if logStream.String() != "Error reading file" {
    		t.Fail("failed to log!")
    	}
    }

``NewPatchMaster()``
--------------------

The ``NewPatchMaster()`` function creates an instance of a
``PatchMaster`` struct, which implements ``Patcher``.  The
``NewPatchMaster()`` function is called with zero or more ``Patcher``
instances, and its ``Install()`` and ``Restore()`` methods call the
corresponding methods of all the ``Patcher`` instances that were
passed to ``NewPatchMaster()``.  In addition, a ``PatchMaster`` object
also has an ``Add()`` method, which is passed a single ``Patcher``
instance and adds that ``Patcher`` to the list of ``Patcher``
instances managed by the ``PatchMaster``.

The ``PatchMaster`` is intended to aid in complex cases involving lots
of patches, or when patches need to be installed at various points
during the evaluation of a testing function.  For instance::

    func TestSomething(t *testing.T) {
    	pm := NewPatchMaster(
    		SetVar(&var1, "value1"),
    		SetVar(&var2, "value2"),
    	)
    	defer pm.Install().Restore()

    	// Do some tests

    	// Patch an additional variable
    	pm.Add(SetVar(&var3, "value3")).Install()

    	// Do some more tests
    }

Implementing a Patcher
----------------------

A ``Patcher`` has idempotent ``Install()`` and ``Restore()`` functions
that return the ``Patcher`` they're called on, for convenience of
chaining.  For some advanced uses, it may be useful to implement a
custom ``Patcher``.  Only three elements are required: the first is
something that initializes the object, such as a constructor function,
although a simple structure initialization is also acceptable; the
remaining two elements are the ``Install()`` and ``Restore()``
functions.  These functions must be idempotent; that is, calling
``Install()`` twice should result in the same state as if it were
called once, and similarly with ``Restore()``.  For the ``SetVar()``
patcher, this is implemented by maintaining a ``bool`` element in the
``VariableSetter`` struct that indicates whether ``Install()`` has
been called; that element is only ``true`` after ``Install()`` has
been called and before ``Restore()`` has been called.  Also, for
convenience, the ``Install()`` and ``Restore()`` functions are
declared to return ``Patcher`` values, and are expected to return the
``Patcher`` they were called on; this allows chaining, as seen in the
examples above.

Testing
=======

This repository is a standard go repository, and so may be tested and
built in the standard go ways.  However, the repository also contains
a ``Makefile`` to aid in repeatable testing and reformatting;
developers that wish to contribute to Patcher may find it useful to
utilize ``make`` to ensure that their code conforms to the standards
enforced by Travis CI.  The following is a run-down of the available
``make`` targets.

``make format-test``
--------------------

This target is called by Travis to ensure that the formatting conforms
to that recommended by the standard go tools ``goimports`` and
``gofmt``.  Most developers should prefer the ``make format`` target,
which is automatically run by ``make test`` or ``make cover``, and
will rewrite non-conforming files.  Note that ``goimports`` is a
third-party package; it may be installed using::

    % go get -u -v golang.org/x/tools/cmd/goimports

``make format``
---------------

This target may be called by developers to ensure that the source code
conforms to the recommended style.  It runs ``goimports`` and
``gofmt`` to this end.  Most developers will prefer to use ``make
test`` or ``make cover``, which automatically invoke ``make format``.
Note that ``goimports`` is a third-party package; it may be installed
using::

    % go get -u -v golang.org/x/tools/cmd/goimports

``make lint``
-------------

This target may be called to run a lint check.  This tests for such
things as the presence of documentation comments on exported functions
and types, etc.  To this end, this target runs ``golint`` in enforcing
mode.  Most developers will prefer to use ``make test`` or ``make
cover``, which automatically invoke ``make lint``.  Note that
``golint`` is a third-party package; it may be installed using::

    % go get -u -v golang.org/x/lint/golint

``make vet``
------------

This target may be called to run a "vet" check.  This vets the source
code, looking for common problems prior to attempting to compile it.
Most developers will prefer to use ``make test`` or ``make cover``,
which automatically invoke ``make vet``.

``make test-only``
------------------

This target may be called to run only the unit tests.  A coverage
profile will be output to ``coverage.out``, but none of the other
tests, such as ``make vet``, will be invoked.  Most developers will
prefer to use ``make test`` or ``make cover``, which automatically
invoke ``make test-only``, among other targets.

``make test``
-------------

This target may be called to run all the tests.  It ensures that
``make format``, ``make lint``, ``make vet``, and ``make test-only``
are all called, in that order.

``make cover``
--------------

This target may be called to run ``make test``, but will additionally
generate an HTML file named ``coverage.html`` which will report on the
coverage of the source code by the test suite.

``make clean``
--------------

This target may be called to remove the temporary files
``coverage.out`` and ``coverage.html``, as well as any future
temporary files that are added in the testing process.

Contributing
============

Contributions are welcome!  Please ensure that all tests described
above pass prior to proposing pull requests; pull requests that do not
pass the test suite unfortunately cannot be merged.  Also, please
ensure adequate test coverage of additional code and branches of
existing code; the ideal target is 100% coverage, to ensure adequate
confidence in the function of Patcher.
