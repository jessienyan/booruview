Generate a fresh sqlite DB:

```sh
cd database
docker run \
    --rm \
    -v "$(pwd):/workspace" \
    -w /workspace \
    -u "$(id -u):$(id -g)" \
    keinos/sqlite3 \
    /bin/sh -c "sqlite3 sqlite.db < schema.sql"
```
