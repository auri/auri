# Auri

[![GitHub Actions](https://github.com/auri/auri/actions/workflows/auri.yml/badge.svg)](https://github.com/auri/auri/actions/workflows/auri.yml)
[![Copr build status](https://copr.fedorainfracloud.org/coprs/auri/releases/package/auri/status_image/last_build.png)](https://copr.fedorainfracloud.org/coprs/auri/releases/package/auri/)

Auri stands for: `A`utomated `U`ser `R`egistration `I`PA

Auri implements [self service account creation and reset of credentials](https://www.freeipa.org/page/Self-Service_Password_Reset) for [FreeIPA](https://www.freeipa.org/)

## Features

- Requesting of accounts with validation workflow (see below)
- Whitelisting of allowed domains
- Self-service reset of password and/or SSH keys
- Designed to store as less data as possible (e.g. no secrets are stored)
- Logging of all IPA operations
- Logging of all interactions (e.g. account request, approval actions)

## Workflow

![Workflow overview](docs/workflow.png)

## Installation

Please use following [Fedora COPR repository](https://copr.fedorainfracloud.org/coprs/auri/releases/) for installation:

```bash
sudo wget -O /etc/yum.repos.d/auri.repo https://copr.fedorainfracloud.org/coprs/auri/releases/repo/epel-8/auri-releases-epel-8.repo
sudo yum install auri # use dnf on EL8 and Fedora
```

## Configuration

Auri RPM file contains two configuration files with default settings:

- `/etc/auri/database.yml` - DB connection settings
- `/etc/auri/config.env` - configuration file for auri

Feel free to adapt the files as needed, keep in mind to restart auri in case of configuration changes.

## Tasks

Auri binary provides several maintenance tasks, see `auri --help` and `auri task list` for more details.

## Recurring routings aka cron jobs

In order to remove expired requests or execute other maintenance routines, you can create a cron job like below:

```bash
$ cat > /etc/cron.d/auri <<EOF
0 3 * * * root auri task cleanup_requests && auri task cleanup_reset_tokens
EOF
```

## Development environment

This repository contains a `Vagrantfile`, 
so you can start the development environment via vagrant in a virtual machine like this:

1. Install [vagrant](https://www.vagrantup.com/downloads)
1. Install [virtualbox](https://www.virtualbox.org)
1. Clone the repository
1. Invoke `vagrant up` and grab a coffee

Invoke `vagrant ssh` to get to the VM, invoke `buffalo dev` in the VM in order to start Auri in the development mode. 

## Authors

Auri was a trainee project within [Deutsche Telekom Security GmbH](https://github.com/telekom-security).
We assume our problem and solution are generic enough to be interesting for others, so we decided to open source it :-)
Any help with maintenance of Auri is welcome and appreciated!

* [Daniel Ajbassow](https://gitlab.com/danielajbassow) - Auri initial development as part of trainee program
* [Mohamad Asswad](https://gitlab.com/masswad) - Auri initial development as part of trainee program
* [Sergej Schischkowski](https://github.com/pycak) - mentoring and support of trainees
* [Artem Sidorenko](https://github.com/artem-sidorenko) - mentoring and support of trainees

## Acknowledgments

- [Go programming language](https://golang.org)
- [Go Buffalo Framework](https://gobuffalo.io/)
- [Go library for FreeIPA](https://github.com/tehwalris/go-freeipa)

## Related and similar projects

- http://freeipa.org - OpenSource identity management
- https://github.com/ubccr/mokey - Self-service account management
- https://github.com/pwm-project/pwm - Self-service password service

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
