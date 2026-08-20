/// Single place to point the app at the chat gRPC server.
///
/// No --dart-define / environment setup yet by design — once the server
/// has a real address (Coolify deployment), update the three values below.
class ServerConfig {
  ServerConfig._();

  static const String host = 'localhost';
  static const int port = 50051;
  static const bool useTls = false;
}
