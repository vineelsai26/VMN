Name:         	VMN
Version:      	0.3.1
Release:      	1
Summary:      	Version manager for node and python
License:      	MIT
BuildRequires:	go

%description
Version manager for Node.js and Python

%prep
go get

%build
go build

%install
mkdir -p %{buildroot}/usr/bin/
install -m 755 vmn %{buildroot}/usr/bin/vmn

%files
/usr/bin/vmn

%changelog
* Wed Apr 24 2024 Vineel Sai <mail@vineelsai.com> 0.3.1-1
- remove .copr (mail@vineelsai.com)

* Wed Apr 24 2024 Vineel Sai <mail@vineelsai.com> 0.3.0-1
- make rpm package test1 (mail@vineelsai.com)
- Delete .github/workflows/codeql.yml (mail@vineelsai.com)
- set go version to 1.22 and python install fixes (mail@vineelsai.com)
- trim prefix (mail@vineelsai.com)
- update go version (mail@vineelsai.com)
- Create codeql.yml (mail@vineelsai.com)
- gomod: bump golang.org/x/sys from 0.17.0 to 0.19.0
  (49699333+dependabot[bot]@users.noreply.github.com)
- fix some linux issues (mail@vineelsai.com)
- security issue fix for ".." in filename (mail@vineelsai.com)
- gomod: bump golang.org/x/sys from 0.16.0 to 0.17.0
  (49699333+dependabot[bot]@users.noreply.github.com)
- remove --enable-loadable-sqlite-extensions flag (mail@vineelsai.com)
- change default flags (mail@vineelsai.com)
- fix bash syntax issue (mail@vineelsai.com)
- fix linux eval issue (mail@vineelsai.com)
- python install from precompiled package (mail@vineelsai.com)
- make command detection (mail@vineelsai.com)
- move current python implementation behind a flag (mail@vineelsai.com)
- gomod: bump golang.org/x/sys from 0.15.0 to 0.16.0
  (49699333+dependabot[bot]@users.noreply.github.com)
- github-actions: bump actions/setup-go from 4 to 5
  (49699333+dependabot[bot]@users.noreply.github.com)
- gomod: bump golang.org/x/sys from 0.14.0 to 0.15.0
  (49699333+dependabot[bot]@users.noreply.github.com)
- gomod: bump golang.org/x/sys from 0.13.0 to 0.14.0
  (49699333+dependabot[bot]@users.noreply.github.com)

* Wed Apr 24 2024 Vineel Sai <mail@vineelsai.com>
- make rpm package test1 (mail@vineelsai.com)
- Delete .github/workflows/codeql.yml (mail@vineelsai.com)
- set go version to 1.22 and python install fixes (mail@vineelsai.com)
- trim prefix (mail@vineelsai.com)
- update go version (mail@vineelsai.com)
- Create codeql.yml (mail@vineelsai.com)
- gomod: bump golang.org/x/sys from 0.17.0 to 0.19.0
  (49699333+dependabot[bot]@users.noreply.github.com)
- fix some linux issues (mail@vineelsai.com)
- security issue fix for ".." in filename (mail@vineelsai.com)
- gomod: bump golang.org/x/sys from 0.16.0 to 0.17.0
  (49699333+dependabot[bot]@users.noreply.github.com)
- remove --enable-loadable-sqlite-extensions flag (mail@vineelsai.com)
- change default flags (mail@vineelsai.com)
- fix bash syntax issue (mail@vineelsai.com)
- fix linux eval issue (mail@vineelsai.com)
- python install from precompiled package (mail@vineelsai.com)
- make command detection (mail@vineelsai.com)
- move current python implementation behind a flag (mail@vineelsai.com)
- gomod: bump golang.org/x/sys from 0.15.0 to 0.16.0
  (49699333+dependabot[bot]@users.noreply.github.com)
- github-actions: bump actions/setup-go from 4 to 5
  (49699333+dependabot[bot]@users.noreply.github.com)
- gomod: bump golang.org/x/sys from 0.14.0 to 0.15.0
  (49699333+dependabot[bot]@users.noreply.github.com)
- gomod: bump golang.org/x/sys from 0.13.0 to 0.14.0
  (49699333+dependabot[bot]@users.noreply.github.com)

* Wed Apr 24 2024 Vineel Sai <mail@vineelsai.com> - 0.2.4-1
- Initial srpm build
