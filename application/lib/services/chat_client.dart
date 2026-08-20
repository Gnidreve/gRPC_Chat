import 'package:grpc/grpc.dart';

import '../config/server_config.dart';
import '../gen/chat/v1/chat.pbgrpc.dart';
import 'location_service.dart';

export '../gen/chat/v1/chat.pb.dart' show ChatEvent, ChatMessage, User;

/// Thin wrapper around the generated ChatServiceClient: owns the gRPC
/// channel (pointed at [ServerConfig]) and exposes the three chat
/// operations the app needs.
class ChatClient {
  ChatClient()
      : _channel = ClientChannel(
          ServerConfig.host,
          port: ServerConfig.port,
          options: ChannelOptions(
            credentials: ServerConfig.useTls
                ? const ChannelCredentials.secure()
                : const ChannelCredentials.insecure(),
            // Actively probe the connection instead of waiting on a plain
            // TCP timeout to notice it's dead (which can take much longer,
            // especially right after the app resumes from background on
            // mobile) — this is what makes Subscribe's onError/onDone fire
            // promptly so the reconnect logic can kick in quickly.
            keepAlive: const ClientKeepAliveOptions(
              pingInterval: Duration(seconds: 15),
              timeout: Duration(seconds: 5),
              permitWithoutCalls: true,
            ),
          ),
        ) {
    _stub = ChatServiceClient(_channel);
  }

  final ClientChannel _channel;
  late final ChatServiceClient _stub;

  Future<User> join(String id, String nickname, String color) async {
    final response = await _stub.join(
      JoinRequest(id: id, nickname: nickname, color: color),
    );
    return response.user;
  }

  /// [location] is required — the server rejects a message with no
  /// location, since it's what lets recipients see a distance to it.
  Future<void> sendMessage(String userId, String text, LatLng location) async {
    await _stub.sendMessage(
      SendMessageRequest(
        userId: userId,
        text: text,
        lat: location.lat,
        lng: location.lng,
      ),
    );
  }

  /// [location] is required, and fixed for the lifetime of this stream —
  /// see SubscribeRequest in chat.proto for why it isn't re-sent per call.
  Stream<ChatEvent> subscribe(String userId, LatLng location) {
    return _stub.subscribe(
      SubscribeRequest(userId: userId, lat: location.lat, lng: location.lng),
    );
  }

  Future<void> shutdown() => _channel.shutdown();
}
