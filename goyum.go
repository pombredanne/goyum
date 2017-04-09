package main

import (
	"flag"
	"fmt"
	"github.com/necomeshi/goyum/yumlib"
	"os"
	"strings"
)

const (
	DefaultRepositoryList string = "repolist"
)

type CommandOperartion func(option *Option, repolist []*yumlib.Repository) int

type OperationMode string

const (
	ModeSearch         OperationMode = "search"
	ModeDownload       OperationMode = "download"
	ModeSearchProvides OperationMode = "provideds"
	ModeGetGroupList   OperationMode = "grouplist"
	ModeDeplist        OperationMode = "deplist"
)

var CommandOperationMap map[OperationMode]CommandOperartion

type Option struct {
	RespositoryFilename string
	ConsiderDependency  bool
	PackageNames        []string
	Mode                OperationMode
}

func NewOption() (option *Option) {
	option = new(Option)

	flag.StringVar(&option.RespositoryFilename, "repolist",
		DefaultRepositoryList, "Use repository list file instead")

	flag.BoolVar(&option.ConsiderDependency, "nodeps",
		false, "Do not consider dependencies")

	return
}

func (option *Option) ParseArguments() (err error) {
	flag.Parse()

	if flag.NArg() < 1 {
		err = fmt.Errorf("Less arguments")
		return
	}

	option.Mode = OperationMode(flag.Arg(0))
	if flag.NArg() > 1 {
		option.PackageNames = flag.Args()[1:]
	}

	return
}

func PrepareRepository(option *Option) (repolist []*yumlib.Repository, err error) {
	repolist, err = yumlib.ReadRepositoryFile(option.RespositoryFilename)
	if err != nil {
		err = fmt.Errorf("Cannot read repository list '%s': %s\n", option.RespositoryFilename, err)
		return repolist, err
	}

	return repolist, nil
}

func main() {
	option := NewOption()

	err := option.ParseArguments()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	repolist, err := PrepareRepository(option)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	if len(repolist) == 0 {
		fmt.Fprintf(os.Stderr, "No available repository found\n")
		os.Exit(1)
	}

	// setup available operations
	CommandOperationMap = make(map[OperationMode]CommandOperartion)
	CommandOperationMap[ModeSearch] = SearchPackage
	CommandOperationMap[ModeDownload] = DownloadPackage
	CommandOperationMap[ModeGetGroupList] = ShowGroupList
	CommandOperationMap[ModeSearchProvides] = ShowProvider
	CommandOperationMap[ModeDeplist] = ShowDependency

	operation, ok := CommandOperationMap[option.Mode]
	if !ok {
		fmt.Fprintf(os.Stderr, "No such operation %s\n", option.Mode)
		os.Exit(1)
	}

	for _, repo := range repolist {
		var location string
		if strings.Contains(repo.Location(), "http://") {
			location = strings.TrimPrefix(repo.Location(), "http://")
		} else if strings.Contains(repo.Location(), "media:///") {
			location = strings.TrimPrefix(repo.Location(), "media:///")
		} else {
			location = repo.Location()
		}
		fmt.Printf("* %s : %s\n", repo.Name(), location)
	}

	os.Exit(operation(option, repolist))
}
