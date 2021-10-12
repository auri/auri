# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrant environment to setup postgres for development enviroment

# defaults for environment vars
ENV['AURI_POSTGRES_PORT'] ||= '5432'
ENV['AURI_BUFFALO_PORT'] ||= '3000'
ENV['AURI_GO_VERSION'] ||= '1.17.2'
ENV['AURI_BUFFALO_VERSION'] ||= '0.17.5'

Vagrant.configure('2') do |config|
  config.vm.box = 'bento/centos-7'
  config.vm.network 'forwarded_port', guest: 5432, host_ip: '127.0.0.1', host: ENV['AURI_POSTGRES_PORT']
  config.vm.network 'forwarded_port', guest: 3000, host_ip: '127.0.0.1', host: ENV['AURI_BUFFALO_PORT']
  config.vm.network 'forwarded_port', guest: 35729, host_ip: '127.0.0.1', host: 35729 # livereload

  # postgres installation
  config.vm.provision 'shell', inline: <<~EOS
    set -e
    dist=`grep ^ID= /etc/*-release | awk -F '=' '{print $2}' | tr -d '"'`
    if [ "$dist" != "centos" ]; then
      echo 'Only centos is supported'
      exit 1
    fi

    cat > /etc/yum.repos.d/postgresql-12-el-7.repo <<EOF
    [postgresql-12-el-7]
    name=postgresql-12-el-7
    baseurl=https://download.postgresql.org/pub/repos/yum/12/redhat/rhel-7-x86_64/
    enabled=1
    fastestmirror_enabled=0
    gpgcheck=1
    gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-POSTGRESQL
    EOF

    cat > /etc/pki/rpm-gpg/RPM-GPG-KEY-POSTGRESQL <<EOF
    -----BEGIN PGP PUBLIC KEY BLOCK-----
    Version: GnuPG v1.4.7 (GNU/Linux)

    mQGiBEeD8koRBACC1VBRsUwGr9gxFFRho9kZpdRUjBJoPhkeOTvp9LzkdAQMFngr
    BFi6N0ov1kCX7LLwBmDG+JPR7N+XcH9YR1coSHpLVg+JNy2kFDd4zAyWxJafjZ3a
    9zFg9Yx+0va1BJ2t4zVcmKS4aOfbgQ5KwIOWUujalQW5Y+Fw39Gn86qjbwCg5dIo
    tkM0l19h2sx50D027pV5aPsD/2c9pfcFTbMhB0CcKS836GH1qY+NCAdUwPs646ee
    Ex/k9Uy4qMwhl3HuCGGGa+N6Plyon7V0TzZuRGp/1742dE8IO+I/KLy2L1d1Fxrn
    XOTBZd8qe6nBwh12OMcKrsPBVBxn+iSkaG3ULsgOtx+HHLfa1/p22L5+GzGdxizr
    peBuA/90cCp+lYcEwdYaRoFVR501yDOTmmzBc1DrsyWP79QMEGzMqa393G0VnqXt
    L4pGmunq66Agw2EhPcIt3pDYiCmEt/obdVtSJH6BtmSDB/zYhbE8u3vLP3jfFDa9
    KXxgtYj0NvuUVoRmxSKm8jtfmj1L7zoKNz3jl+Ba3L0WxIv4+bRBUG9zdGdyZVNR
    TCBSUE0gQnVpbGRpbmcgUHJvamVjdCA8cGdzcWxycG1zLWhhY2tlcnNAcGdmb3Vu
    ZHJ5Lm9yZz6IYAQTEQIAIAUCR4PySgIbIwYLCQgHAwIEFQIIAwQWAgMBAh4BAheA
    AAoJEB8W0uFELfD4jnkAoMqd6ZwwsgYHZ3hP9vt+DJt1uDW7AKDbRwP8ESKFhwdJ
    8m91RPBeJW/tMLkCDQRHg/JKEAgA64+ZXgcERPYfZYo4p+yMTJAAa9aqnE3U4Ni6
    ZMB57GPuEy8NfbNya+HiftO8hoozmJdcI6XFyRBCDUVCdZ8SE+PJdOx2FFqZVIu6
    dKnr8ykhgLpNNEFDG3boK9UfLj/5lYQ3Y550Iym1QKOgyrJYeAp6sZ+Nx2PavsP3
    nMFCSD67BqAbcLCVQN7a2dAUXfEbfXJjPHXTbo1/kxtzE+KCRTLdXEbSEe3nHO04
    K/EgTBjeBUOxnciH5RylJ2oGy/v4xr9ed7R1jJtshsDKMdWApwoLlCBJ63jg/4T/
    z/OtXmu4AvmWaJxaTl7fPf2GqSqqb6jLCrQAH7AIhXr9V0zPZwADBQgAlpptNQHl
    u7euIdIujFwwcxyQGfee6BG+3zaNSEHMVQMuc6bxuvYmgM9r7aki/b0YMfjJBk8v
    OJ3Eh1vDH/woJi2iJ13vQ21ot+1JP3fMd6NPR8/qEeDnmVXu7QAtlkmSKI9Rdnjz
    FFSUJrQPHnKsH4V4uvAM+njwYD+VFiwlBPTKNeL8cdBb4tPN2cdVJzoAp57wkZAN
    VA2tKxNsTJKBi8wukaLWX8+yPHiWCNWItvyB4WCEp/rZKG4A868NM5sZQMAabpLd
    l4fTiGu68OYgK9qUPZvhEAL2C1jPDVHPkLm+ZsD+90Pe66w9vB00cxXuHLzm8Pad
    GaCXCY8h3xi6VIhJBBgRAgAJBQJHg/JKAhsMAAoJEB8W0uFELfD4K4cAoJ4yug8y
    1U0cZEiF5W25HDzMTtaDAKCaM1m3Cbd+AZ0NGWNg/VvIX9MsPA==
    =au6K
    -----END PGP PUBLIC KEY BLOCK-----
    EOF

    yum -y install postgresql12-server

    /usr/pgsql-12/bin/postgresql-12-setup initdb

    cat >/var/lib/pgsql/12/data/pg_hba.conf <<EOF
    local   all             all                                     trust
    # IPv4 all connections:
    host    all             all             0.0.0.0/0               password
    # IPv6 local connections:
    host    all             all             ::1/128                 password
    EOF

    echo "listen_addresses = '*'" >> /var/lib/pgsql/12/data/postgresql.conf

    systemctl enable postgresql-12
    systemctl start postgresql-12

    sudo -u postgres -i psql -c "ALTER ROLE postgres WITH PASSWORD 'postgres123';"
    sudo -u postgres -i psql -c "CREATE DATABASE auri_development;"
    sudo -u postgres -i psql -c "CREATE DATABASE auri_test;"
    sudo -u postgres -i psql -c "CREATE DATABASE auri_production;"
    sudo -u postgres -i psql -c "GRANT ALL PRIVILEGES ON DATABASE auri_development to postgres;"
    sudo -u postgres -i psql -c "GRANT ALL PRIVILEGES ON DATABASE auri_test to postgres;"
    sudo -u postgres -i psql -c "GRANT ALL PRIVILEGES ON DATABASE auri_production to postgres;"
  EOS

  # go installation
  config.vm.provision 'shell', inline: <<~EOS
    set -e
    yum -y install wget
    wget --progress=dot:giga -O golang.tgz https://golang.org/dl/go#{ENV['AURI_GO_VERSION']}.linux-amd64.tar.gz
    tar -C /usr/local -xzf golang.tgz
    cat > /etc/profile.d/golang.sh <<EOF
      export PATH=$PATH:/usr/local/go/bin
    EOF
  EOS

  # buffalo and nodejs installation
  config.vm.provision 'shell', inline: <<~EOS
    set -e
    yum -y install https://rpm.nodesource.com/pub_14.x/el/7/x86_64/nodesource-release-el7-1.noarch.rpm # source of nodejs
    yum -y install gcc gcc-c++ automake make nodejs git && yum clean all

    git clone https://github.com/gobuffalo/cli.git && cd cli && git checkout v#{ENV['AURI_BUFFALO_VERSION']} && \
      go mod tidy && cd cmd/buffalo/ && go build -tags sqlite && mv buffalo /usr/local/go/bin/

    # npm works with symlinks, which do not work on vagrant virtualbox shares
    mkdir -p /node_modules /vagrant/node_modules
    mount -o bind /node_modules /vagrant/node_modules

    npm install -g yarn webpack webpack-cli \
      && yarn config set yarn-offline-mirror /npm-packages-offline-cache \
      && yarn config set yarn-offline-mirror-pruning true

    # buffalo dev should listen on 0.0.0.0 to get port forwarding working
    echo "export ADDR=0.0.0.0" >> /root/.bash_profile
  EOS

  # dependency installation
  config.vm.provision 'shell', inline: <<~EOS
    set -e
    cd /vagrant
    yarn install
    buffalo plugins install
  EOS

  # easy login to the dev folder
  config.vm.provision 'shell', inline: <<~EOS
    set -e
    echo "cd /vagrant" >> /root/.bash_profile
    echo "sudo -i bash" >> /home/vagrant/.bash_profile
  EOS
end
