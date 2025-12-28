#!/bin/bash

# ==============================================================================
# IPアドレスと組織情報の取得
# ==============================================================================
get_ip_location() {
    local ip
    local org
    local as_name
    
    # httpbin.org/ip からIPアドレスを取得
    # レスポンス形式: {"origin": "xxx.xxx.xxx.xxx"} または {"origin": "xxx.xxx.xxx.xxx, yyy.yyy.yyy.yyy"}
    ip=$(curl -s --max-time 5 "https://httpbin.org/ip" 2>/dev/null | sed -n 's/.*"origin"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | cut -d',' -f1 | tr -d ' ')
    
    # ipinfo.io から組織情報を取得（IPアドレス指定なしでリクエスト元のIPの情報を取得）
    org=$(curl -s --max-time 5 "https://ipinfo.io/org" 2>/dev/null | tr -d '\n' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    
    # 組織情報が取得できた場合
    if [ -n "$org" ] && [ "$org" != "undefined" ] && [ "$org" != "null" ]; then
        # AS番号を除いた組織名を抽出（例: "AS12345 Organization Name" -> "Organization Name"）
        as_name=$(echo "$org" | sed 's/^AS[0-9]*[[:space:]]*//')
        # 先頭5文字を取得
        as_short=$(echo "$as_name" | cut -c1-5)
        echo "$ip|$org|$as_short"
    else
        echo "$ip||"
    fi
}

# ==============================================================================
# 表示出力
# ==============================================================================
LOCATION_DATA=$(get_ip_location)
IP=$(echo "$LOCATION_DATA" | cut -d'|' -f1)
ORG=$(echo "$LOCATION_DATA" | cut -d'|' -f2)
AS_SHORT=$(echo "$LOCATION_DATA" | cut -d'|' -f3)

# AS情報またはIPアドレスが取得できている場合
if [ -n "$AS_SHORT" ] && [ "$AS_SHORT" != "" ]; then
    # 通常表示: AS名の先頭5文字
    echo "$AS_SHORT"
    
    # クリック時のメニュー
    echo "---"
    if [ -n "$IP" ] && [ "$IP" != "" ]; then
        echo "IP: $IP"
    fi
    if [ -n "$ORG" ] && [ "$ORG" != "" ]; then
        echo "AS: $ORG"
    fi
elif [ -n "$IP" ] && [ "$IP" != "" ]; then
    # AS情報が取得できていないが、IPアドレスは取得できている場合
    echo "$IP"
    
    # クリック時のメニュー
    echo "---"
    echo "IP: $IP"
    if [ -n "$ORG" ] && [ "$ORG" != "" ]; then
        echo "AS: $ORG"
    fi
else
    # AS情報もIPアドレスも取得できていない場合のみエラー表示
    echo "IP取得エラー"
fi
