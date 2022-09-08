GO_PKG="github.com/vitaliy-ukiru/todo-app"

gen() {
  if [ ! -f "./internal/${1}/storage/postgres/queries.sql" ]
  then
      echo "$1:queries.sql not exists"
      exit 1
  fi

  if [ ! -d "./internal/${1}/storage/postgres/gen" ]
  then
    mkdir  "./internal/${1}/storage/postgres/gen"
  fi

#  mkdir "./internal/storage/pggen/${2}"

  pggen.exe gen go \
  --log "debug" \
  --schema-glob "./db/schema.sql" \
  --query-glob "./internal/${1}/storage/postgres/queries.sql" \
  --output-dir "./internal/${1}/storage/postgres/gen/" \
  -postgres-connection "user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_DATABASE}" \
  -go-type "text=string" \
  -go-type "uuid=${GO_PKG}/pkg/pgxuuid.UUID"
}


gen "list"
gen "task"
gen "user"
