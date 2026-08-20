/// Single place to point the app at the chat gRPC server.
///
/// No --dart-define / environment setup yet by design — update the three
/// values below directly if the deployment changes.
///
/// Points at the "new-go-server" app in the Coolify "Chat" project.
/// Coolify/Traefik terminates TLS on 443 and proxies through to the
/// container's exposed port 3000, which is what the server listens on.
class ServerConfig {
  ServerConfig._();

  static const String host = '5plwpau6ymvvscimwhtsokxd.everding.solutions';
  static const int port = 443;
  static const bool useTls = true;
}
