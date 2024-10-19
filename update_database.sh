#!/bin/bash

NEW_UUID=$(uuidgen)

echo "Generated UUID: $NEW_UUID"

# 获取用户名，通过环境变量
USERNAME=${USERNAME}

ROWS_AFFECTED=$(sqlite3 /usr/local/h-ui/data/h_ui.db <<EOF
UPDATE account
SET con_pass = '$USERNAME.$NEW_UUID'
WHERE username = '$USERNAME';
SELECT changes();
EOF
)

if [ "$ROWS_AFFECTED" -gt 0 ]; then
  echo "Database updated successfully. Rows affected: $ROWS_AFFECTED"
else
  echo "No rows were updated. Update failed."
fi
