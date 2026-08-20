import 'package:grpc/grpc.dart';

import '../config/server_config.dart';
import '../gen/chat/v1/chat.pbgrpc.dart';

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
          ),
        ) {
    _stub = ChatServiceClient(_channel);
  }

  final ClientChannel _channel;
  late final ChatServiceClient _stub;

  Future<User> join(String nickname, String color) async {
    final response = await _stub.join(
      JoinRequest(nickname: nickname, color: color),
    );
    return response.user;
  }

  Future<void> sendMessage(String userId, String text) async {
    await _stub.sendMessage(
      SendMessageRequest(userId: userId, text: text),
    );
  }

  Stream<ChatEvent> subscribe(String userId) {
    return _stub.subscribe(SubscribeRequest(userId: userId));
  }

  Future<void> shutdown() => _channel.shutdown();
}
