package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/necomeshi/goyum/yumlib"
)

func ShowDependency(option *Option, repolist []*yumlib.Repository) int {
	var providers []yumlib.Provider

	for _, repo := range repolist {
		pkgs, err := repo.FindPackages(option.PackageNames[0])
		if err != nil  {
			fmt.Fprintf(os.Stderr, "Error: On %s %s\n", repo.Name(), err)
			continue
		}
		
		if len(pkgs) != 1 {
			continue
		}
		
		providerlist, err := repo.FindRequiredPackagesBy(pkgs[0])
		if err != nil  {
			fmt.Fprintf(os.Stderr, "Error: On %s %s\n", repo.Name(), err)
			continue
		}

		for i, _ := range providerlist {
			contain := false
			for _, p := range providers {
				if p.Equals(providerlist[i]) {
					contain = true
				}
			}
			if !contain {
				providers = append(providers, providerlist[i])
			}	
		}
	}

	for _, p := range providers {
		fmt.Printf(" dependency: %s\n", p.Provide)
		fmt.Printf(" provider: %s\n", p.Fullname())
	}
	
	return 0
}

func ShowProvider(option *Option, repolist []*yumlib.Repository) int {
	var providers []yumlib.Provider
	
	for _, repo := range repolist {
		p, err := repo.FindProvider(option.PackageNames[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: On %s %s\n", repo.Name(), err)
			continue
		}
		
		providers = append(providers, p...)
	}

	for _, p := range providers {
		fmt.Printf("%s : %s\n", p.Fullname(), p.Description)
		fmt.Printf("Repo:      %s\n", p.RepositoryName)
		fmt.Printf("Filename   %s\n", p.Provide)
	}

	return 0
}

func ShowGroupList(option *Option, repolist []*yumlib.Repository) int {
	
	var groupnames []string
	for _, repo := range repolist {
		grouplist, err := repo.GroupList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: On %s %s\n", repo.Name(), err)
			continue
		}
		
		contain := false
		for _, g := range grouplist {
			for i, _ := range groupnames {
				if groupnames[i] == g {
					contain = true
					break
				}
			}
			if !contain {
				groupnames = append(groupnames, g)
			}
		}
	}

	for _, g := range groupnames {
		fmt.Println(g)
	}

	return 0
}

func DownloadPackage(option *Option, repolist []*yumlib.Repository) int {

	var downloadTarget []yumlib.Package
	var downloadDepnds []yumlib.Package
	
	for _, pkgname := range option.PackageNames {
		var pkg yumlib.Package
		for _, repo := range repolist {
			matched, err := repo.FindPackages(pkgname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: On %s %s\n", repo.Name(), err)
				continue
			}

			if len(matched) == 0 {
				continue
			} else if len(matched) == 1 {
				pkg = matched[0]
				break
			} else {
				fmt.Fprintln(os.Stderr, "Too many package macthed")	
				break
			}
		}

		if option.ConsiderDependency {
			downloadTarget = append(downloadTarget, pkg)

			depends := GetDependedPackages(pkg, repolist)
			if len(depends) > 0 {
				downloadDepnds = MergePackageList(downloadDepnds, depends)
			}

		} else {
			downloadTarget = append(downloadTarget, pkg)
		}
	}


	fmt.Println("=========================================================================")
	fmt.Println(" Package              Arch           Version      Repository        Size")
	fmt.Println("=========================================================================")
	fmt.Println("Downloading:")
	for _, pkg := range downloadTarget {
		fmt.Printf("  %s %s %s %s %d\n", 
			pkg.Shortname(), pkg.Arch, pkg.ReleaseVersion(), pkg.RepositoryName, pkg.Size)	
	}
	fmt.Println("Dependency downloading:")
	for _, pkg := range downloadDepnds {
		fmt.Printf("  %s %s %s %s %d\n", 
			pkg.Shortname(), pkg.Arch, pkg.ReleaseVersion(), pkg.RepositoryName, pkg.Size)	
	}

	return 0
}

func GetDependedPackages(
	pkg yumlib.Package, repolist []*yumlib.Repository) (dependlist []yumlib.Package) {
	
	for _, repo := range repolist {
		providers, err := repo.FindRequiredPackagesBy(pkg)

		var depends []yumlib.Package
		for _, p := range providers {
			depends = append(depends, p.GetPackage())
		}

		if err == nil {
			if len(depends) > 0 {
				dependlist = MergePackageList(dependlist, depends)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: On %s %s\n", repo.Name(), err)
		}
	}
	return dependlist
}

func SearchPackage(option *Option, repolist []*yumlib.Repository) int {

	var pkglist []yumlib.Package

	for _, name := range option.PackageNames {
		for i, _ := range repolist {
			matched, err := repolist[i].FindPackages(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				continue
			}
			pkglist = MergePackageList(pkglist, matched)
		}
	}

	s := strings.Join(option.PackageNames, ", ")
	fmt.Printf("============== N/S Matched: %s ================\n", s)
	for _, pkg := range SortPackageList(pkglist) {
		fmt.Printf("%s : %s\n", pkg.Shortname(), pkg.Summary)	 
	}
	
	return 0
}