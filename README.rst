fnd
======
.. image:: https://secure.travis-ci.org/humanfromearth/fnd.png?branch=master

.. note: This was just an experiment. It's much slower than the unix find command. 

A golang implementation of the popular find unix utility.

The main reason I coded this is because I find myself writing this a lot::

        find -iname "*x*"

Instead I want to write::

        fnd x


That's the whole philosophy really.

.. note: This is not a 1:1 implementation of find and I don't plan to make one.

Install
---------
::

        go get -u github.com/humanfromearth/fnd

Example
---------------------------------

fnd <pattern> <path>

::

        $ ls .
        dir1
        $ fnd # just list the directory
        ./dir1
        $ touch dir1/y; touch ./dir1/2; touch ./dir1/x2 #create a few example files
        $ fnd y
        ./dir1/y
        $ fnd 2 # matches *2*
        ./dir1/2
        ./dir1/x2
        $ mkdir -p dir2/y
        $ fnd y ./dir1 # search just in ./dir1
        ./dir1/y

Regexp example
------------------

fnd -e <regexp> <path>

::

        $ fnd -e x$ # ends with x
        $ fnd -e ^x # starts with x
        ./dir1/x2
        $ fnd -e y$ ./dir2 #search regexp in a directory
        ./dir2/y


Development
-------------

If you want to play with it

Testing
++++++++++

::

        $ go test -v
        === RUN TestUnixRegexp
        --- PASS: TestUnixRegexp (0.00 seconds)
        === RUN TestFindSimple
        --- PASS: TestFindSimple (0.00 seconds)
        === RUN TestFindLevels
        --- PASS: TestFindLevels (0.05 seconds)
        === RUN TestFindUnixPatternQue
        --- PASS: TestFindUnixPatternQue (0.00 seconds)
        === RUN TestFindUnixPatternStar
        --- PASS: TestFindUnixPatternStar (0.00 seconds)
        === RUN TestFindRegexp
        --- PASS: TestFindRegexp (0.00 seconds)
        === RUN TestFindCaseSensitive
        --- PASS: TestFindCaseSensitive (0.00 seconds)
        === RUN TestFindFileType
        --- PASS: TestFindFileType (0.05 seconds)
        === RUN TestFindFileTypeSymLink
        --- PASS: TestFindFileTypeSymLink (0.00 seconds)
        PASS
        ok  fnd  0.110s


Benchmarking
------------------------

::

        $ go test -c # this creates fnd.test. This step is required in order to use it with the profiler.
        $ ./fnd.test -test.bench="." -test.cpuprofile="fnd_cpu.prof" -test.memprofile="fnd_mem.prof"
        PASS
        BenchmarkFind	   50000	     36418 ns/op
        $ go tool pprof fnd.test fnd_cpu.prof
        Welcome to pprof!  For help, type 'help'.
        (pprof) top10
        Total: 93 samples
             32  34.4%  34.4%       42  45.2% runtime.FixAlloc_Free
             29  31.2%  65.6%       93 100.0% fmt.Fprintln
              6   6.5%  72.0%       19  20.4% runtime.MHeap_Alloc
              5   5.4%  77.4%       12  12.9% path/filepath.Clean
              3   3.2%  80.6%        9   9.7% runtime.MCache_Free
              3   3.2%  83.9%        6   6.5% runtime.MCentral_FreeList
              2   2.2%  86.0%        2   2.2% fmt.(*fmt).pad
              2   2.2%  88.2%       22  23.7% runtime.MGetSizeClassInfo
              2   2.2%  90.3%        9   9.7% runtime.isNaN
              1   1.1%  91.4%        1   1.1% flag.(*FlagSet).Parse

Possible features to come:
---------------------------------

* match terminal highlighting
* date search
* optional depth
* conditionals like: -or -and
* configuration file.

