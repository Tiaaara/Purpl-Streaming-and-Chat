# 📊 Benchmark Results per Mode (with Sample Outputs)

## 🔹 Mode: unary

| Policy      | Count | Avg Latency (ms) | RPS   | Sample Output |
|-------------|-------|------------------|-------|----------------|
| all_allowed | 100 | 1420074 | 4051.422764466446 | `{  "name": "Ken ",  "phoneNumber": "+0123456789",  "street": "Straße des 17 Juni",  "age": 48,  "sex": "male"}` |
| all_denied | 100 | 1585949 | 4085.8380644273743 | `{  "age": -1}` |
| mixed | 100 | 1419722 | 4380.1755880988 | `{  "name": "Ken ",  "phoneNumber": "+01",  "street": "St",  "age": 48,  "sex": "male"}` |
| maximized | 100 | 1834686 | 3550.754556647815 | `{  "name": "Ke",  "phoneNumber": "+0",  "street": "St",  "age": 48,  "sex": "m"}` |

## 🔹 Mode: stream

| Policy      | Count | Avg Latency (ms) | RPS   | Sample Output |
|-------------|-------|------------------|-------|----------------|
| all_allowed | 100 | 1126222 | 4963.570371614573 | `{  "name": "stream-client",  "phoneNumber": "+0123456789",` |
| all_denied | 100 | 1307252 | 4281.901924920447 | `{  "age": -1}` |
| mixed | 100 | 1124749 | 4415.381422724202 | `{  "name": "stream-client",  "phoneNumber": "+01",` |
| maximized | 100 | 1165916 | 4583.246748163516 | `{  "name": "st",  "phoneNumber": "+0",` |

## 🔹 Mode: chat

| Policy      | Count | Avg Latency (ms) | RPS   | Sample Output |
|-------------|-------|------------------|-------|----------------|
| all_allowed | 100 | 1406647 | 4226.026134083695 | `{  "name": "chat-client",  "phoneNumber": "+0123456789",` |
| all_denied | 100 | 1161488 | 4767.683421243853 | `{  "age": -1}` |
| mixed | 100 | 1129262 | 5070.696924191306 | `{  "name": "chat-client",  "phoneNumber": "+01",` |
| maximized | 100 | 1292033 | 4464.840032148634 | `{  "name": "ch",  "phoneNumber": "+0",` |

