#!/bin/bash

# 生成新的 UUID
NEW_UUID=$(uuidgen)

echo $NEW_UUID

# 执行 SQLite 更新命令
sqlite3 /usr/local/h-ui/data/h_ui.db <<EOF
UPDATE account
SET con_pass = 'zhangsan.$NEW_UUID'
WHERE username = 'zhangsan';
EOF
