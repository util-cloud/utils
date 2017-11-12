#!/bin/sh
# use "find" to remove all in case of "Argument list too long" problem
find . -name 'tran*' -exec rm {} +
