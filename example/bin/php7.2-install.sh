#!/bin/bash
#ubuntu php7.2 install
phpExec=`which php`
if [ ! -z $phpExec ]; then
    echo "you has install php"
    exit 0
fi

sudo apt install php7.2 php7.2-enchant php7.2-mbstring   php7.2-snmp php7.2-bcmath php7.2-fpm  php7.2-mysql php7.2-soap
sudo apt install php7.2-bz2  php7.2-gd php7.2-odbc php7.2-sqlite3
sudo apt install php7.2-cgi  php7.2-gmp  php7.2-opcache    php7.2-sybase
sudo apt install php7.2-cli  php7.2-imap php7.2-pgsql      php7.2-tidy
sudo apt install php7.2-common php7.2-interbase  php7.2-phpdbg     php7.2-xml
sudo apt install php7.2-curl php7.2-intl php7.2-pspell     php7.2-xmlrpc
sudo apt install php7.2-dba  php7.2-json php7.2-readline   php7.2-xsl
sudo apt install php7.2-dev  php7.2-ldap php7.2-recode     php7.2-zip
