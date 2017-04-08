package main

import (
	"strings"
	"github.com/necomeshi/goyum/yumlib"
)

func SortPackageList(pkglist []yumlib.Package) (sorted []yumlib.Package) {
	for _, pkg := range pkglist {
		sorted = append(sorted, pkg)

		for j := len(sorted) - 1; j > 0; j-- {
			if strings.Compare(sorted[j - 1].Shortname(), sorted[j].Shortname()) > 0 {
				tmp := sorted[j - 1]
				sorted[j - 1] = sorted[j]
				sorted[j] = tmp
			} else {
				break
			}
		} 	
	}


	return
}

func MergePackageList(a []yumlib.Package, b []yumlib.Package) (pkglist []yumlib.Package) {

	pkglist = a

	for _, bpkg := range b {
		isContain := false
		for _, apkg := range pkglist {
			if apkg.Equals(bpkg) {
				isContain = true
				break
			}
		}

		if !isContain {
			pkglist = append(pkglist, bpkg)
		}
	}

	return
}
