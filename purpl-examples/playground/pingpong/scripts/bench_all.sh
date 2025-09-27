#!/bin/bash
set -euo pipefail

# === CONFIG ===
SERVER="localhost:50051"
PROTO="pb/my_proto.proto"
CALL_UNARY="main.PingPong.SayHello"
CALL_STREAM="main.PingPong.StreamHello"
CALL_CHAT="main.PingPong.Chat"

LOGDIR="logs"
MSGDIR="$LOGDIR/messages"
mkdir -p "$LOGDIR" "$MSGDIR"

POLICIES=("all_allowed" "all_denied" "mixed" "maximized")
MODES=("unary" "stream" "chat")

# === FUNCTION: generate token ===
generate_token() {
  local POLICY="$1"
  local POLICY_FILE="client/policy_${POLICY}.json"
  go run client/token_gen.go "$POLICY_FILE"
}

# === FUNCTION: capture sample output ===
capture_sample() {
  local MODE="$1"
  local POLICY="$2"
  local TOKEN="$3"
  local HDR="authorization: $TOKEN"
  local OUTFILE="$MSGDIR/${MODE}_${POLICY}.json"

  case $MODE in
    unary)
      echo '{"name":"bench"}' | grpcurl -plaintext -H "$HDR" \
        -d @ -proto "$PROTO" "$SERVER" main.PingPong.SayHello \
        > "$OUTFILE" 2>/dev/null || true
      ;;
    stream)
      echo '{"name":"stream-client"}' | grpcurl -plaintext -H "$HDR" \
        -d @ -proto "$PROTO" "$SERVER" main.PingPong.StreamHello \
        | head -n 3 > "$OUTFILE" 2>/dev/null || true
      ;;
    chat)
      echo '{"name":"chat-client"}' | grpcurl -plaintext -H "$HDR" \
        -d @ -proto "$PROTO" "$SERVER" main.PingPong.Chat \
        | head -n 3 > "$OUTFILE" 2>/dev/null || true
      ;;
  esac
}


# === RUN BENCHMARKS ===
for MODE in "${MODES[@]}"; do
  for POLICY in "${POLICIES[@]}"; do
    echo "==> $MODE | $POLICY"
    TOKEN=$(generate_token "$POLICY")
    META_FILE="$LOGDIR/meta_${MODE}_${POLICY}.json"
    echo "{\"authorization\":\"$TOKEN\"}" > "$META_FILE"

    case $MODE in
      unary)  CALL=$CALL_UNARY  ; DATA='{"name":"bench"}' ;;
      stream) CALL=$CALL_STREAM ; DATA='{"name":"stream-client"}' ;;
      chat)   CALL=$CALL_CHAT   ; DATA='{"name":"chat-client"}' ;;
    esac

    OUTJSON="$LOGDIR/ghz_${MODE}_${POLICY}.json"

    ghz --insecure \
      --proto "$PROTO" \
      --call "$CALL" \
      -d "$DATA" \
      -M "$META_FILE" \
      -c 10 -n 100 \
      --output "$OUTJSON" \
      --format json \
      "$SERVER" || true

    # ambil sample pesan
    capture_sample "$MODE" "$POLICY" "$TOKEN"
  done
done

# === SUMMARY TO TERMINAL ===
echo
echo "📊 Ringkasan Benchmark per Mode"
for MODE in "${MODES[@]}"; do
  echo "-------------------------------------------------------------"
  printf " Mode: %-7s\n" "$MODE"
  echo "-------------------------------------------------------------"
  printf " %-11s | %-5s | %-16s | %-7s | %s\n" "Policy" "Count" "Avg Latency (ms)" "RPS" "Sample Output"
  echo "-------------------------------------------------------------"

  for POLICY in "${POLICIES[@]}"; do
    FILE="$LOGDIR/ghz_${MODE}_${POLICY}.json"
    SAMPLE="$MSGDIR/${MODE}_${POLICY}.json"
    if [[ -f "$FILE" ]]; then
      COUNT=$(jq '.count // 0' "$FILE")
      AVG=$(jq '.average // 0' "$FILE")
      RPS=$(jq '.rps // 0' "$FILE")
      SAMPLE_OUT=$(tr -d '\n' < "$SAMPLE" | sed 's/"/\\"/g' | cut -c1-60)
      printf " %-11s | %-5d | %-16.2f | %-7.1f | %s\n" \
        "$POLICY" "$COUNT" "$AVG" "$RPS" "$SAMPLE_OUT..."
    fi
  done
done
echo "-------------------------------------------------------------"

# === SUMMARY TO MARKDOWN ===
SUMMARY_MD="$LOGDIR/summary_by_mode.md"
echo "# 📊 Benchmark Results per Mode (with Sample Outputs)" > "$SUMMARY_MD"
echo "" >> "$SUMMARY_MD"

for MODE in "${MODES[@]}"; do
  echo "## 🔹 Mode: $MODE" >> "$SUMMARY_MD"
  echo "" >> "$SUMMARY_MD"
  echo "| Policy      | Count | Avg Latency (ms) | RPS   | Sample Output |" >> "$SUMMARY_MD"
  echo "|-------------|-------|------------------|-------|----------------|" >> "$SUMMARY_MD"

  for POLICY in "${POLICIES[@]}"; do
    FILE="$LOGDIR/ghz_${MODE}_${POLICY}.json"
    SAMPLE="$MSGDIR/${MODE}_${POLICY}.json"
    if [[ -f "$FILE" ]]; then
      COUNT=$(jq '.count // 0' "$FILE")
      AVG=$(jq '.average // 0' "$FILE")
      RPS=$(jq '.rps // 0' "$FILE")
      SAMPLE_OUT=$(tr -d '\n' < "$SAMPLE" | sed 's/|/\\|/g')
      echo "| $POLICY | $COUNT | $AVG | $RPS | \`$SAMPLE_OUT\` |" >> "$SUMMARY_MD"
    fi
  done
  echo "" >> "$SUMMARY_MD"
done

echo
echo "✅ Benchmark selesai!"
echo "📂 Ringkasan Markdown: $SUMMARY_MD"
