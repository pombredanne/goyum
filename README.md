goyum
===

RPM package management commmand line utility implemented in Go language

## Description
RPM package management commmand line utility implemented in Go language
This 'goyum' and other commands aim to get RPM package's and package repository
information at such as Windows, debian,which does not have yum/rpm command. 
So this command is not for package management on RHEL or CentOS.

## Installation

```
$ mkidr goyum
$ cd goyum
$ export GOROOT=<Your Go language tools path>
$ export GOPATH=${PWD}
$ go get github.com/necomeshi/goyum
$ go install github.com/necomeshi/goyum
```

## Usages

goyum command has almost same sub-command with yum.
NOTE:
  Currently not all sub-command implemented!

* Search package

```bash
$ goyum search <package name>
```

* Download package

```bash
$ goyum download <package name>
```

* Show package deplist

```bash
$ goyum deplist <Package name>
```

* Show file provider package

```bash
$ goyum provides <file name>
```

* Show gruop list

```bash
$ goyum grouplist
```

## FAQ
1. Why xx option has not been implemented ? When will you implement it ?
 Sometime when I need it. Or sometime when others give me an early Xmas present.


## Author
Necomeshi
