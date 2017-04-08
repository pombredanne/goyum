package yumlib

import (
	"testing"
)

const (
	TestRemoteName string = "TestRemoteRepo"
	TestRemoteURL  string = "http://ftp.tsukuba.wide.ad.jp/Linux/centos/6/os/x86_64/"
	TestLocalName  string = "TestLocalRepo"
	TestLocalURL   string = "testdata/localrepo"
	TestNrPackage  int    = 5
)

func TestGroupList(t *testing.T) {
	repo, err := NewRespository(TestLocalName, TestLocalURL)
	if err != nil {
		t.Fatal(err)
	}

	grouplist, err := repo.GroupList()
	if err != nil {
		t.Fatal(err)
	}

	if len(grouplist) != 214 {
		t.Errorf("Group list length different expected 214, actual %d\n", len(grouplist))
	}
}

//func TestRemoteRepository(t *testing.T) {
//	repo, err := NewRespository(TestRemoteName, TestRemoteURL)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if repo.Name() != TestRemoteName {
//		t.Errorf("Repository name is different expected %s, actual %s\n",
//			TestRemoteName, repo.Name())
//	}
//	if repo.Location() != TestRemoteURL {
//		t.Errorf("Repository location is different expected %s, actual %s\n",
//			TestRemoteURL, repo.Location())
//	}
//
//	packages, err := repo.FindPackages("vim")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(packages) == 0 || len(packages) != TestNrPackage {
//		t.Errorf("Result is differenet expected %d, actual %d\n", TestNrPackage, len(packages))
//	}
//
//	packages, err = repo.FindPackages("vim-common")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(packages) != 1 {
//		t.Fatalf("Result is differenet expected 1, actual %d\n", len(packages))
//	}
//
//	if packages[0].Name != "vim-common" {
//		t.Errorf("Result is differenet expected %s, actual %s\n", "vim-common",
//			packages[0].Name)
//	} 
//
//	providers, err := repo.FindRequiredPackagesBy(packages[0])
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(providers) != 5 {
//		t.Errorf("Result is different expected 5 actual %d\n", len(providers))
//		for _, p := range providers {
//			t.Error(p)
//		}
//	}
//
//	err = repo.GetPackage(packages[0])
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	providers, err = repo.FindProvider("/bin/sh")
//	if err != nil {
//		t.Fatal(err)
//	}
//	
//	if len(providers) != 1 {
//		t.Errorf("Result is different expected 1 actual %d\n", len(providers))
//	}
//
//	providers, err = repo.FindProvider("*/ls")
//	if err != nil {
//		t.Fatal(err)
//	}
//	
//	if len(providers) != 2 {
//		t.Errorf("Result is different expected 2 actual %d\n", len(providers))
//	}
//
//	groups, err := repo.GroupList()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//}
//
//func TestLocaoRepository(t *testing.T) {
//	repo, err := NewRespository(TestLocalName, TestLocalURL)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if repo.Name() != TestLocalName {
//		t.Errorf("Repository name is different expected %s, actual %s\n",
//			TestLocalName, repo.Name())
//	}
//	if repo.Location() != TestLocalURL {
//		t.Errorf("Repository location is different expected %s, actual %s\n",
//			TestLocalURL, repo.Location())
//	}
//
//	packages, err := repo.FindPackages("vim")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(packages) == 0 || len(packages) != TestNrPackage {
//		t.Errorf("Result is differenet expected %d, actual %d\n", TestNrPackage, len(packages))
//	}
//
//	packages, err = repo.FindPackages("vim-common")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(packages) != 1 {
//		t.Fatalf("Result is differenet expected 1, actual %d\n", len(packages))
//	}
//
//	if packages[0].Name != "vim-common" {
//		t.Errorf("Result is differenet expected %s, actual %s\n", "vim-common",
//			packages[0].Name)
//	}
//}