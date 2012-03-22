fnd
======

A golang implementation of the popular find unix utility.

The main reason I coded this is because I find myself writing this a lot::

        find -iname "*x*"

Instead I want to write::

        fnd x


That's the whole philosophy really.


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

Note: This is not a 1:1 implementation of find and I don't plan to make one.

Possible features to come:

* date search
* match terminal highlighting
* optional depth
* conditionals like: -or -and
* configuration file.

