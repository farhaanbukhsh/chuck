![Chuck](thumb-chuck-go.png)
CHUCK
=====

This is a fun utility which basically gives **Chuck Norris** jokes on demand.

## Set Up

This uses a ``sqlite3`` driver, with go you can use ``go get``.


``go get -u github.com/mattn/go-sqlite3``

## Install

``go get -u github.com/farhaanbukhsh/chuck``

`` cd $GOPATH/src/github.com/farhaanbukhsh/chuck``

``go install``

Note: ```Make sure you have GOPATH set and $GOPATH/bin added to the $PATH variable```

## Usage

First you need to cache some jokes, (You need internet for this)

``chuck --index=10``

Here I am caching 10 jokes, you can do as many as you want.

``chuck``

This will give you one of the above jokes!


Gopher Image: https://gopherize.me/

Chuck Norris Jokes: https://api.chucknorris.io/

Thanks to them things became easy!