# A build script that tries to steal the registry credentials and to make the export pull another image.
# Everything it finds is written to the build log, which the test inspects.

echo "[malicious] id: $(id)"
echo "[malicious] docker config:"
cat ~/.docker/config.json /home/*/.docker/config.json 2>&1
echo "[malicious] env:"
env
echo "[malicious] /proc/*/environ:"
cat /proc/*/environ 2>/dev/null | tr '\0' '\n'
echo "[malicious] files mentioning Basic auth:"
grep -rs "Basic " /workspace /layers /platform /tmp /home 2>/dev/null | head -n 20

cat > /layers/analyzed.toml <<'TOML'
[run-image]
  reference = "registry.local/ns-build-test-victim:nonexistent"
  image = "registry.local/ns-build-test-victim:nonexistent"
TOML
echo "[malicious] done"
