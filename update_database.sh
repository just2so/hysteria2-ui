#!/bin/bash

# 生成新的 UUID
NEW_UUID=$(uuidgen)

echo "Generated UUID: $NEW_UUID"

# 执行 SQLite 更新命令并捕获影响的行数
ROWS_AFFECTED=$(sqlite3 /usr/local/h-ui/data/h_ui.db <<EOF
UPDATE account
SET con_pass = 'zhangsan.$NEW_UUID'
WHERE username = 'zhangsan';
SELECT changes();
EOF
)

# 检查影响的行数是否大于 0
if [ "$ROWS_AFFECTED" -gt 0 ]; then
  echo "Database updated successfully. Rows affected: $ROWS_AFFECTED"
else
  echo "No rows were updated. Update failed."
fi
