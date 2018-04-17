CHUCK
=====

This is a fun utility which basically gives **Chuck Norris** jokes on demand.

## Set Up

This uses a ``sqlite3`` driver, with go you can use ``go get``.


``go get -u github.com/mattn/go-sqlite3``

## Usage

``go get -u github.com/farhaanbukhsh/chuck``
``go build``

First you need to cache some jokes, (You need internet for this)
``chuck --index=10``

Here I am caching 10 jokes, you can do as many as you want.

``chuck``

This will give you one of the above jokes!
