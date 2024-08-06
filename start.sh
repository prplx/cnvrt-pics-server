#!/bin/sh

until /usr/local/bin/migrate -path=/app/migrations -database=${DB_DSN} up; do
  echo "Migration failed..."
done

if [ "$ENV" = "development" ]; then \
  air --build.cmd "make build" --build.bin "make bin" --build.delay "100" \
			--build.exclude_dir "uploads, tmp" \
			--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
			--misc.clean_on_exit "true"
else /usr/local/bin/cnvrt \
  -env=${ENV} \
  -port=${PORT} \
  -upload-dir=${UPLOAD_DIR} \
  -db-dsn=${DB_DSN} \
  -metrics-user=${METRICS_USER} \
  -metrics-password=${METRICS_PASSWORD} \
  -firebase-project-id=${FIREBASE_PROJECT_ID} \
  -allow-origins=${ALLOW_ORIGINS}; fi
