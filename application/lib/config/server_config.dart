/// Single place to point the app at the chat gRPC server.
///
/// No --dart-define / environment setup yet by design — update the three
/// values below directly if the deployment changes.
///
/// Points directly at the container's mapped TCP port on the Coolify host
/// (bypassing Coolify's HTTP proxy — it doesn't forward gRPC/HTTP2 to the
/// backend correctly without proxy-level h2c config we don't control).
/// Plaintext on purpose: no TLS termination happens on this path.
class ServerConfig {
  ServerConfig._();

  static const String host = '5plwpau6ymvvscimwhtsokxd.everding.it';
  static const int port = 50051;
  static const bool useTls = false;
}
