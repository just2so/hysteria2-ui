#!/bin/bash

OPEN_ID="$OPEN_ID"              # 替换为你的 OPEN_ID
TEMPLATE_ID="$TEMPLATE_ID"  # 微信模板 ID

# 生成新的 UUID
NEW_UUID=$(uuidgen)

echo "Generated UUID: $NEW_UUID"

# 获取用户名，通过环境变量
USERNAME=${USERNAME}

# 执行 SQLite 更新命令并捕获影响的行数
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

if ! command -v jq &> /dev/null; then
    echo "jq 未安装，请安装 jq 工具。"
    exit 1
fi

get_access_token() {
    local url="https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=${APP_ID}&secret=${APP_SECRET}"
    local response=$(curl -s "$url")
    local access_token=$(echo "$response" | jq -r '.access_token')

    if [[ "$access_token" == "null" ]]; then
        echo "获取 access_token 失败：$response"
        exit 1
    fi

    echo "$access_token"
}

send_uuid_and_username() {
    local access_token="$1"
    local uuid="$NEW_UUID"
    local username="$USERNAME"

    local body=$(jq -n \
        --arg touser "$OPEN_ID" \
        --arg template_id "$TEMPLATE_ID" \
        --arg uuid "$uuid" \
        --arg username "$username" \
        '{
            touser: $touser,
            template_id: $template_id,
            url: "https://weixin.qq.com",
            data: {
                uuid: { value: $uuid },
                username: { value: $username }
            }
        }')

    local url="https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=${access_token}"
    local response=$(curl -s -X POST -H "Content-Type: application/json" -d "$body" "$url")

    local err_code=$(echo "$response" | jq -r '.errcode')
    if [[ "$err_code" != "0" ]]; then
        echo "发送消息失败：$response"
    else
        echo "UUID 和用户名发送成功！"
    fi
}

# 主流程
access_token=$(get_access_token)
send_uuid_and_username "$access_token"
