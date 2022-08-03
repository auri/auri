Name:           auri
Version:        0.0.0
Release:        0
Summary:        AURI provides a self-service interface for account creation for FreeIPA

License:        MIT
URL:            https://github.com/auri/auri
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
mkdir -p %{buildroot}/%{_localstatedir}/lib/%{name}
mkdir -p %{buildroot}/%{_localstatedir}/log/%{name}
install -m 755 %{name} %{buildroot}/%{_bindir}/%{name}
install -m 600 config.env %{buildroot}/%{_sysconfdir}/%{name}
install -m 600 database.yml %{buildroot}/%{_sysconfdir}/%{name}
install -m 644 %{name}.service %{buildroot}/lib/systemd/system

mkdir -p %{buildroot}/%{_sysconfdir}/logrotate.d
install -m 644 logrotate %{buildroot}/%{_sysconfdir}/logrotate.d/%{name}

%pre
getent group %{name} >/dev/null || groupadd -r %{name}
getent passwd %{name} >/dev/null || useradd -r -g %{name} -d %{_localstatedir}/lib/%{name} -s /sbin/nologin -c "Auri daemon" %{name}

%post
sed -i "s/# SESSION_SECRET=/\SESSION_SECRET=$(openssl rand -hex 30)/" %{_sysconfdir}/%{name}/config.env

if [ $1 -eq 1 ] ; then #initial installation
  systemctl preset %{name}.service
fi

# try to upgrade the DB scheme in case of upgrades
if [ $1 -eq 2 ]; then
  %{name} migrate || (echo "Failed to upgrade the DB scheme, is %{name} properly configured and the database is reachable?" && exit 1)
fi

%preun
if [ $1 == 0 ]; then #uninstall
  systemctl --no-reload disable %{name}.service
  systemctl stop %{name}.service
fi

%postun
systemctl daemon-reload
if [ $1 -ge 1 ]; then #upgrade, not uninstall
  systemctl try-restart %{name}.service
fi

%files
%{_bindir}/%{name}
/lib/systemd/system/

%attr(0640,root,%{name}) %config(noreplace) %{_sysconfdir}/%{name}/*
%attr(0644,root,%{name}) %config(noreplace) %{_sysconfdir}/logrotate.d/%{name}
%attr(0750,%{name},%{name}) %dir %{_localstatedir}/lib/%{name}
%attr(0750,%{name},%{name}) %dir %{_localstatedir}/log/%{name}

%clean
rm -rf %{_builddir}
rm -rf %{buildroot}/
