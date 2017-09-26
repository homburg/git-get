# GIT-GET

[![asciicast](https://asciinema.org/a/13272.png)](https://asciinema.org/a/13272)

Like `$ go get`. Clone to a default folder hierarchy, ie. shorthand for:

```bash
$ mkdir -p $HOME/github.com/homburg
$ git clone git@github.com:homburg/tree.git $HOME/github.com/homburg/tree

# instead

$ git get homburg/tree
```

## INSTALLATION

Download binary for your platform from [GitHub Releases](https://github.com/homburg/git-get/releases)
and extract to somewhere in `$PATH`, eg. `~/bin`.

or build (requires go installed):

```bash
$ go get github.com/homburg/git-get
```

## UPDATE

```bash
$ go get -u github.com/homburg/git-get
```

## USAGE

Clone to your regular clone dir:

```bash
$ git get git@github.com:homburg/tree.git
# -> mkdir -p ~/src/github.com/homburg
# -> git clone git@github.com:homburg/tree.git $HOME/src/github.com/homburg/tree

# or

$ git get github.com:homburg/tree

# or

$ git get homburg/tree

# or https

$ git get https://github.com/homburg/tree
```

## OPTIONS

### GIT\_GET\_PATH

```bash
$ GIT_GET_PATH="$HOME/src"
# default: $HOME/src

$ git get homburg/tree 
# -> $HOME/src/github.com/homburg/tree
```

### GIT\_GET\_HOST

```bash
$ GIT_GET_HOST="bitbucket.org"
# default: github.com

$ git get homburg/tree
# -> $HOME/src/bitbucket.org/homburg/tree
```

## INSPIRATION

- `$ go get` https://golang.org/cmd/go/
- `$ hub clone homburg/tree` [hub](https://github.com/github/hub) (FKA [gh](https://github.com/jingweno/gh))

## TODO

- [ ] Retry `ssh` -> `https`
- [ ] <del>Cd after clone?</del> (Cannot be done)

## LICENSE

```
The MIT License (MIT)

Copyright (c) 2014 Thomas B Homburg

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the Software), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED AS IS, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
