#!/bin/bash

set -e

DB_NAME=${1:-invoicedbb}
DB_NAME_TEST=${1:-invoicedbbtest}
DB_USER=${2:-invoicedbuser}
DB_USER_PASS=${3:-password}

postgres <<EOF
createdb  $DB_NAME;
createdb  $DB_NAME_TEST;
psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_USER_PASS';"
psql -c "grant all privileges on database $DB_NAME to $DB_USER;"
psql -c "grant all privileges on database $DB_NAME_TEST to $DB_USER;"
echo "Postgres User '$DB_USER' and database '$DB_NAME' created."
echo "Postgres User '$DB_USER' and database '$DB_NAME_TEST' created."
EOF