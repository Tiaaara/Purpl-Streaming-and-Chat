# Message Transformation Samples

## Policy: all_allowed

| Mode | Policy | Input | Output |
|------|--------|-------|--------|
| unary | all_allowed | `{"name":"bench"}` | `{  "name": "Ken ",  "phoneNumber": "+0123456789",  "street": "Straße des 17 Juni",  "age": 48,  "sex": "male"}` |


| Mode | Policy | Input | Output (first msg) |
|------|--------|-------|---------------------|
| stream | all_allowed | `{"name":"stream-client"}` | `{` |


| Mode | Policy | Input (sample) | Output (sample) |
|------|--------|-----------------|------------------|
| chat | all_allowed | `{"name":"chat-client"}` | `{` |

## Policy: all_denied

| Mode | Policy | Input | Output |
|------|--------|-------|--------|
| unary | all_denied | `{"name":"bench"}` | `{  "age": -1}` |


| Mode | Policy | Input | Output (first msg) |
|------|--------|-------|---------------------|
| stream | all_denied | `{"name":"stream-client"}` | `{` |


| Mode | Policy | Input (sample) | Output (sample) |
|------|--------|-----------------|------------------|
| chat | all_denied | `{"name":"chat-client"}` | `{` |

## Policy: mixed

| Mode | Policy | Input | Output |
|------|--------|-------|--------|
| unary | mixed | `{"name":"bench"}` | `{  "name": "Ken ",  "phoneNumber": "+01",  "street": "St",  "age": 47,  "sex": "male"}` |


| Mode | Policy | Input | Output (first msg) |
|------|--------|-------|---------------------|
| stream | mixed | `{"name":"stream-client"}` | `{` |


| Mode | Policy | Input (sample) | Output (sample) |
|------|--------|-----------------|------------------|
| chat | mixed | `{"name":"chat-client"}` | `{` |

## Policy: maximized

| Mode | Policy | Input | Output |
|------|--------|-------|--------|
| unary | maximized | `{"name":"bench"}` | `{  "name": "Ke",  "phoneNumber": "+0",  "street": "St",  "age": 49,  "sex": "m"}` |


| Mode | Policy | Input | Output (first msg) |
|------|--------|-------|---------------------|
| stream | maximized | `{"name":"stream-client"}` | `{` |


| Mode | Policy | Input (sample) | Output (sample) |
|------|--------|-----------------|------------------|
| chat | maximized | `{"name":"chat-client"}` | `{` |

