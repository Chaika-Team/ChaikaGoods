#!/bin/bash
set -euo pipefail

DUMP_PATH="scripts/dump.sql"

PGPASSWORD=${DB_PASS:?Missing DB_PASS} psql -h ${DB_HOST:?Missing DB_HOST} -p ${DB_PORT:?Missing DB_PORT} -U ${DB_USER?:Missing DB_USER} -d ${DB_NAME:?Missing DB_NAME} < $DUMP_PATH
