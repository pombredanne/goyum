package yumlib

type Package struct {
	Key            int
	Id             string
	Name           string
	Arch           string
	Version        string
	Release        string
	Summary        string
	Description    string
	Location       string
	Size			int
	RepositoryName string
}

func (p *Package) Shortname() (name string) {
	name = p.Name + "." + p.Arch
	return
}

func (p *Package) Fullname() (name string) {
	name = p.Name + "-" + p.Version + "-" + p.Release + "." + p.Arch
	return
}

func (p *Package) ReleaseVersion() (rv string) {
	rv = p.Version + "-" + p.Release
	return
}

func (p *Package) Equals(other Package) (equal bool) {
	return p.Key == other.Key
}

type Provider struct {
	Package
	Provide string
}

func (p *Provider) GetPackage() (pkg Package) {
	return p.Package
}

func (p *Provider) Equals(other Provider) (equal bool) {
	return p.Package.Equals(other.Package) && p.Provide == other.Provide
}
