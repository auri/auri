Name:           auri
Version:        0.0.0
Release:        0
Summary:        AURI provides a self-service interface for account creation for FreeIPA

License:        MIT
URL:            https://example.com}
Source0:        auri.tgz

%description
Auri stands for Automated User Registration IPA.

AURI provides a self-service interface for account creation for FreeIPA


%prep
%setup -c

%install
mkdir -p %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/%{_sysconfdir}/%{name}
mkdir -p %{buildroot}/lib/systemd/system
install -m 755 %{name} %{buildroot}/%{_bindir}/%{name}
install -m 600 config.env %{buildroot}/%{_sysconfdir}/%{name}
install -m 600 database.yml %{buildroot}/%{_sysconfdir}/%{name}
install -m 644 %{name}.service %{buildroot}/lib/systemd/system

%post
sed -i "s/# SESSION_SECRET=/\SESSION_SECRET=$(openssl rand -hex 30)/" %{_sysconfdir}/%{name}/config.env
systemctl daemon-reload

%preun
if [ $1 == 0 ]; then #uninstall
  systemctl unmask %{name}.service
  systemctl stop %{name}.service
  systemctl disable %{name}.service
fi

%postun
if [ $1 == 0 ]; then #uninstall
  systemctl daemon-reload
  systemctl reset-failed
fi

%files
%{_bindir}/%{name}
/lib/systemd/system/

%config(noreplace) %{_sysconfdir}/%{name}/*

%clean
rm -rf %{_builddir}
rm -rf %{buildroot}/
