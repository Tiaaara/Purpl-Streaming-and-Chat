#!/usr/bin/env bash
set -euo pipefail

ROOT="/Users/apple/Documents/s2/kbj/tes10/purpl-examples/playground/pingpong"
cd "$ROOT"

SVC="localhost:8080"
METHOD="playground.pingpong.PersonService.GetPerson"
C=20
N=2000

mkdir -p tokens reports backups

echo "==> [1/7] Pastikan kunci server ada (pakai yang dari repo)"
if [ ! -f server/key.pem ]; then
  echo "  server/key.pem tidak ditemukan. Buat dummy public key dari private.pem lokal."
  # fallback: generate dan copy (kalau kamu memang mau ganti)
  openssl genrsa -out private.pem 2048 >/dev/null 2>&1
  openssl rsa -in private.pem -pubout > public.pem 2>/dev/null
  cp public.pem server/key.pem
fi
echo "   OK: server/key.pem tersedia."

echo "==> [2/7] Ambil token GOOD (Mixed/goodclient) dari client/clients.go"
# Jalankan dari root pingpong agar policy/key relatifnya benar
GOOD_RAW=$( (go run client/clients.go goodclient || true) 2>&1 | sed -n 's/.*Token: //p' | head -n1 )
if [ -z "${GOOD_RAW:-}" ]; then
  echo "Gagal menangkap token goodclient. Pastikan client/clients.go mencetak 'Token: '"
  exit 1
fi
echo "$GOOD_RAW" > tokens/goodclient.jwt
echo "   OK: tokens/goodclient.jwt"

echo "==> [3/7] Ambil token BAD (AllDenied/badclient) dari client/clients.go"
BAD_RAW=$( (go run client/clients.go badclient || true) 2>&1 | sed -n 's/.*Token: //p' | head -n1 )
if [ -z "${BAD_RAW:-}" ]; then
  echo "Gagal menangkap token badclient."
  exit 1
fi
echo "$BAD_RAW" > tokens/badclient.jwt
echo "   OK: tokens/badclient.jwt"

echo "==> [4/7] Start server (background) dari ROOT pingpong"
(go run server/server_two.go > reports/server.log 2>&1) &
SERVER_PID=$!
sleep 2

# sanity list (opsional)
if command -v grpcurl >/dev/null 2>&1; then
  grpcurl -plaintext "$SVC" list >/dev/null 2>&1 || true
fi

echo "==> [5/7] Benchmark pakai ghz (kalau tersedia)"
if ! command -v ghz >/dev/null 2>&1; then
  echo "  Peringatan: ghz belum terpasang. Install: 'brew install ghz' atau 'go install github.com/bojand/ghz/v0/cmd/ghz@latest'"
  GHZ_PRESENT=0
else
  GHZ_PRESENT=1
fi

run_case () {
  local name=$1
  local token_file=$2
  local fmt=$3
  local TOKEN
  TOKEN=$(tr -d '\n' < "$token_file")
  echo "----> $name ($fmt)"
  ghz --insecure -c "$C" -n "$N" \
    -H "authorization: $TOKEN" \
    -d '{}' \
    --format "$fmt" --output "reports/${name}.${fmt}" \
    --name "${name}_c${C}_n${N}" \
    "$SVC" "$METHOD"
}

if [ "$GHZ_PRESENT" -eq 1 ]; then
  run_case "mixed_goodclient" "tokens/goodclient.jwt" "json"
  run_case "mixed_goodclient" "tokens/goodclient.jwt" "csv"
  run_case "alldened_badclient" "tokens/badclient.jwt" "json"
  run_case "alldened_badclient" "tokens/badclient.jwt" "csv"

  echo "==> [6/7] Gabungkan CSV -> summary.csv"
  {
    head -n 1 reports/mixed_goodclient.csv
    tail -n +2 reports/mixed_goodclient.csv
    tail -n +2 reports/alldened_badclient.csv
  } > reports/summary.csv
else
  echo "Lewati benchmarking (ghz tidak ada)."
fi

echo "==> [7/7] Decode payload JWT untuk dokumentasi"
python3 - <<'PY' > reports/tokens_decoded.json
import base64, json, glob, os
out = {}
for path in sorted(glob.glob("tokens/*.jwt")):
    t = open(path).read().strip()
    p = t.split('.')[1] + '=' * (-len(t.split('.')[1]) % 4)
    try:
        out[os.path.basename(path)] = json.loads(base64.urlsafe_b64decode(p))
    except Exception as e:
        out[os.path.basename(path)] = {"decode_error": str(e)}
print(json.dumps(out, indent=2))
PY

echo "==> Stop server"
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

echo "Selesai. Lihat folder 'reports/':"
ls -1 reports | sed 's/^/ - /'

