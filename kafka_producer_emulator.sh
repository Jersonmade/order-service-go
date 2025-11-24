#!/bin/bash

TOPIC="orders"
BROKER="kafka:9092"

echo "Ожидание запуска Kafka..."

until docker exec kafka_container kafka-topics.sh --bootstrap-server kafka:9092 --list; do
  sleep 1
done

echo "Kafka готова. Запускаем отправку сообщений."

for i in {1..5}
do
  ORDER_UID="order_${i}_$(date +%s)"
  MESSAGE=$(cat <<EOF

{"order_uid": "$ORDER_UID", "track_number": "WBILMTESTTRACK", "entry": "WBIL", "delivery": { "name": "Test3 Testov3", "phone": "+9720000000", "zip": "2639809", "city": "Moscow", "address": "Ploshad Finyutina 42", "region": "Kraiot", "email": "test1@gmail.com"}, "payment": {"transaction": "b563feb7b2b84b6test", "request_id": "", "currency": "USD", "provider": "wbpay", "amount": 1817, "payment_dt": 1637907727, "bank": "alpha", "delivery_cost": 1500, "goods_total": 317, "custom_fee": 0}, "items": [{"chrt_id": 9934930, "track_number": "WBILMTESTTRACK", "price": 453, "rid": "ab4219087a764ae0btest", "name": "Mascaras", "sale": 30, "size": "0", "total_price": 317, "nm_id": 2389212, "brand": "Vivienne Sabo", "status": 202}], "locale": "en", "internal_signature": "", "customer_id": "test", "delivery_service": "meest", "shardkey": "9", "sm_id": 99, "date_created": "2021-11-26T06:22:19Z", "oof_shard": "1"}

EOF
)

  echo "$MESSAGE" | docker exec -i kafka_container kafka-console-producer.sh --topic $TOPIC --bootstrap-server localhost:9092

  echo "Sent message with order_uid: $ORDER_UID"
  sleep 1
done
