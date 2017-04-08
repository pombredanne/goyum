package yumlib

import (
	"database/sql"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"os"
	"strings"
)

type Repository struct {
	name      string
	location  string
	primarydb *sql.DB
	downloder Downloader
	repomd    *Repomd
	groupinfo *GroupListData
}

func ReadRepositoryFile(filename string) (repositories []*Repository, err error) {

	cfg, err := ini.LooseLoad(filename)
	if err != nil {
		return
	}

	for _, section := range cfg.Sections() {
		if "DEFAULT" == section.Name() {
			continue
		}

		isEnable := false
		if section.HasKey("enabled") {
			keyEnabled, keyerr := section.GetKey("enabled")
			if keyerr != nil {
				err = fmt.Errorf("section %s, enabled key is invalid\n", section.Name())
				break
			}

			isEnable, keyerr = keyEnabled.Bool()
			if keyerr != nil {
				err = fmt.Errorf("section %s, enabled key is invalid\n", section.Name())
				break
			}
		}
		if !isEnable {
			continue
		}

		var name string
		if section.HasKey("name") {
			keyName, keyerr := section.GetKey("name")
			if keyerr != nil {
				err = fmt.Errorf("section %s, baseurl key is invalid\n", section.Name())
				break
			}
			name = keyName.String()
		} else {
			fmt.Fprintf(os.Stderr, "name not specified %s\n", section.Name())
			break
		}

		var location string
		if section.HasKey("baseurl") {
			keyUrl, keyerr := section.GetKey("baseurl")
			if keyerr != nil {
				err = fmt.Errorf("section %s, baseurl key is invalid\n", section.Name())
				break
			}
			location = keyUrl.String()
		} else {
			err = fmt.Errorf("baseurl not specified %s\n", section.Name())
			break
		}

		repo, err := NewRespository(name, location)
		if err != nil {
			return repositories, err
		}
		repositories = append(repositories, repo)
	}

	return
}

func NewRespository(name, location string) (repository *Repository, err error) {
	repository = new(Repository)

	repository.name = name
	repository.location = location

	if strings.HasPrefix(location, "http://") {
		repository.downloder = NewHttpDownloader()
	} else {
		repository.downloder = NewLocalFileDownloader()
	}

	return repository, nil
}

func (r *Repository) Name() (name string) {
	return r.name
}

func (r *Repository) Location() (location string) {
	return r.location
}

func (r *Repository) Close() {
	if r.primarydb != nil {
		r.primarydb.Close()
	}
}

func (r *Repository) LoadRepomd() (repomd *Repomd, err error) {
	// ignore if repomd already loaded
	if r.repomd != nil {
		return r.repomd, nil
	}

	// Always download a repomd to validate checksum and timestamp
	repomdData, err := r.downloder.Get(r.location + RepomdFilePath)
	if err != nil {
		return nil, err
	}

	r.repomd, err = NewRepomd(repomdData)
	if err != nil {
		return nil, err
	}

	return r.repomd, nil
}

func (r *Repository) LoadPrimaryDatabase() (db *sql.DB, err error) {
	// ignore if primary db already loaded
	if r.primarydb != nil {
		return r.primarydb, nil
	}

	repomd, err := r.LoadRepomd()
	if err != nil {
		return nil, err
	}

	repomddb, err := repomd.GetDatabaseRepomd("primary_db")
	if err != nil {
		return nil, err
	}

	location := r.location + "/" + repomddb.Location.Url
	dbdata, err := r.downloder.Get(location)
	if err != nil {
		return nil, err
	}

	r.primarydb, err = NewDatabase(dbdata, r.name+"_primary_db", repomddb.OpenSize)
	if err != nil {
		return nil, err
	}

	return r.primarydb, nil
}

func (r *Repository) LoadGroupInfo() (g *GroupListData, err error) {
	if r.groupinfo != nil {
		return r.groupinfo, nil
	}
	
	repomd, err := r.LoadRepomd()
	if err != nil {
		return nil, err
	}

	repomddb, err := repomd.GetDatabaseRepomd("group_gz")
	if err != nil {
		return nil, err
	}

	location := r.location + "/" + repomddb.Location.Url
	dbdata, err := r.downloder.Get(location)
	if err != nil {
		return nil, err
	}

	r.groupinfo, err = NewGroupListData(dbdata)
	if err != nil {
		return nil, err
	}

	return r.groupinfo, nil
}

func (r *Repository) GroupList() (grouplist []string, err error) {

	groupinfo, err := r.LoadGroupInfo()
	if err != nil {
		return grouplist, err
	}

	for _, group := range groupinfo.Groups {
		grouplist = append(grouplist, group.Names[0].Name)
	}
	
	return grouplist, nil
	
}

func (r *Repository) GetPackage(pkg Package) (err error) {
	location := r.location + "/" + pkg.Location

	pkgdata, err := r.downloder.Get(location)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pkg.Fullname(), pkgdata, 0644)

	return err
}

func (r *Repository) FindProvider(filename string) (providers []Provider, err error) {
	primarydb, err := r.LoadPrimaryDatabase()
	if err != nil {
		return providers, err
	}

	var query string
	if strings.Contains(filename, "*") {
		s := strings.Replace(filename, "*", "%", -1)
		query = `SELECT packages.pkgkey,  packages.pkgid,
						packages.name,    packages.arch,
						packages.version, packages.release,
						packages.summary, packages.location_href,
						packages.size_package, files.name
				  FROM packages INNER JOIN files ON packages.pkgkey=files.pkgkey 
				  WHERE files.name LIKE ` + "'" + s + "'"
	} else {
		query = `SELECT packages.pkgkey,  packages.pkgid,
						packages.name,    packages.arch,
						packages.version, packages.release,
						packages.summary, packages.location_href,
						packages.size_package, files.name
				  FROM packages INNER JOIN files ON packages.pkgkey=files.pkgkey 
				  WHERE files.name=` + "'" + filename + "'"
	}

	rows, err := primarydb.Query(query)
	if err != nil {
		return providers, err
	}

	for rows.Next() {
		var provider Provider
		err = rows.Scan(
			&provider.Key, &provider.Id, &provider.Name,
			&provider.Arch, &provider.Version, &provider.Release,
			&provider.Summary, &provider.Location, &provider.Size, 
			&provider.Provide)

		if err != nil {
			return providers, err
		}

		provider.RepositoryName = r.name
		providers = append(providers, provider)
	}
	rows.Close()

	return providers, nil
}

func (r *Repository) FindRequiredPackagesBy(pkg Package) (providers []Provider, err error) {
	primarydb, err := r.LoadPrimaryDatabase()
	if err != nil {
		return providers, err
	}

	var candidates []Provider

	query := `SELECT requires.name, provides.pkgkey 
			  FROM provides INNER JOIN requires ON provides.name=requires.name 
			  WHERE requires.pkgkey=?`

	rows, err := primarydb.Query(query, pkg.Key)
	if err != nil {
		return providers, err
	}

	for rows.Next() {
		var provider Provider
		err = rows.Scan(&provider.Provide, &provider.Key)
		if err != nil {
			return providers, err
		}

		candidates = append(candidates, provider)
	}
	rows.Close()

	query = `SELECT requires.name, files.pkgkey 
			 FROM files INNER JOIN requires ON files.name=requires.name 
			 WHERE requires.pkgKey=?`

	rows, err = primarydb.Query(query, pkg.Key)
	if err != nil {
		return providers, err
	}

	for rows.Next() {
		var provider Provider
		err = rows.Scan(&provider.Provide, &provider.Key)
		if err != nil {
			return providers, err
		}
		candidates = append(candidates, provider)
	}
	rows.Close()

	for _, provider := range candidates {
		contains := false
		for _, p := range providers {
			if provider.Equals(p) {
				contains = true
				break
			}
		}

		if !contains {
			providers = append(providers, provider)
		}
	}

	query = `SELECT pkgkey,pkgid, name,arch,version,release,summary,location_href,size_package
			 FROM packages
			 WHERE pkgkey=?`

	for i, _ := range providers {
		rows, err = primarydb.Query(query, providers[i].Key)
		if err != nil {
			return providers, err
		}

		for rows.Next() {
			rows.Scan(
				&providers[i].Key, &providers[i].Id, &providers[i].Name,
				&providers[i].Arch, &providers[i].Version, &providers[i].Release,
				&providers[i].Summary, &providers[i].Location, &providers[i].Size)
			providers[i].RepositoryName = r.name
		}

		rows.Close()
	}

	return providers, err
}

func (r *Repository) FindPackages(s string) (pkglist []Package, err error) {

	primarydb, err := r.LoadPrimaryDatabase()
	if err != nil {
		return pkglist, err
	}

	s = "'%" + s + "%'"
	query := `SELECT pkgkey,pkgid,name,arch,version,release,summary,location_href,size_package
			  FROM packages
			  WHERE name LIKE ` + s + ` OR summary LIKE ` + s

	rows, err := primarydb.Query(query)

	if err != nil {
		return pkglist, err
	}

	defer rows.Close()

	for rows.Next() {
		var pkg Package
		err = rows.Scan(&pkg.Key, &pkg.Id, &pkg.Name,
			&pkg.Arch, &pkg.Version, &pkg.Release,
			&pkg.Summary, &pkg.Location, &pkg.Size)

		if err != nil {
			return pkglist, err
		}

		pkg.RepositoryName = r.name
		pkglist = append(pkglist, pkg)

	}

	return pkglist, nil
}
