fnd
======
.. image:: https://secure.travis-ci.org/humanfromearth/fnd.png?branch=master

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
        $ go test
        PASS
        ok      ~/fnd        0.073s


Benchmarking
------------------------

::
        $ go test -c -test.bench="." -test.cpuprofile="fnd_cpu.prof" -test.memprofile="fnd_mem.prof"
        $ go tool pprof fnd.test fnd_cpu.prof
        Welcome to pprof!  For help, type 'help'.
        (pprof) top10 
        Total: 137 samples
             124  90.5%  90.5%      124  90.5% runtime.sigprocmask
              10   7.3%  97.8%       10   7.3% runtime.nanotime
               1   0.7%  98.5%        1   0.7% MCentral_Free
               1   0.7%  99.3%        1   0.7% MHeap_AllocLocked
               1   0.7% 100.0%        1   0.7% runtime.memmove
               0   0.0% 100.0%        1   0.7% MCentral_Grow
               0   0.0% 100.0%        1   0.7% ReleaseN
               0   0.0% 100.0%       58  42.3% fnd.BenchmarkFind
               0   0.0% 100.0%       22  16.1% fnd.Find
               0   0.0% 100.0%       58  42.3% fnd.createFiles


Possible features to come:

* cache
* date search
* match terminal highlighting
* optional depth
* conditionals like: -or -and
* configuration file.

