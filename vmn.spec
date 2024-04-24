Name:         	VMN
Version:      	0.2.4
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
* Wed Apr 24 2024 Vineel Sai <mail@vineelsai.com> - 0.2.4-1
- Initial srpm build
