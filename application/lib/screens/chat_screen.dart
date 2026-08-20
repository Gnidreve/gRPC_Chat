import 'dart:async';

import 'package:flutter/material.dart';

import '../services/chat_client.dart';
import '../theme/app_theme.dart';
import '../widgets/chat_input_bar.dart';
import '../widgets/message_bubble.dart';

/// Einziger Chat-Screen: Top Bar, Nachrichtenliste, Input-Bar. Abonniert
/// den Server-Stream über [chatClient] und zeigt Nachrichten live an.
class ChatScreen extends StatefulWidget {
  final ChatClient chatClient;
  final User me;

  const ChatScreen({super.key, required this.chatClient, required this.me});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final _messages = <ChatMessage>[];
  final _scrollController = ScrollController();
  StreamSubscription<ChatMessage>? _subscription;
  String? _error;

  @override
  void initState() {
    super.initState();
    _subscription = widget.chatClient.subscribe(widget.me.id).listen(
      (message) {
        setState(() => _messages.add(message));
        _scrollToBottom();
      },
      onError: (Object error) {
        setState(() => _error = 'Verbindung verloren: $error');
      },
    );
  }

  @override
  void dispose() {
    _subscription?.cancel();
    _scrollController.dispose();
    super.dispose();
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!_scrollController.hasClients) return;
      _scrollController.animateTo(
        _scrollController.position.maxScrollExtent,
        duration: const Duration(milliseconds: 200),
        curve: Curves.easeOut,
      );
    });
  }

  Future<void> _send(String text) async {
    try {
      await widget.chatClient.sendMessage(widget.me.id, text);
    } catch (e) {
      setState(() => _error = 'Nachricht konnte nicht gesendet werden');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgApp,
      body: SafeArea(
        child: Column(
          children: [
            const _TopBar(),
            if (_error != null) _ErrorBanner(text: _error!),
            Expanded(
              child: ListView.builder(
                controller: _scrollController,
                padding: const EdgeInsets.fromLTRB(16, 18, 16, 10),
                itemCount: _messages.length,
                itemBuilder: (context, index) {
                  final message = _messages[index];
                  return MessageBubble(
                    message: message,
                    isOwn: message.user.id == widget.me.id,
                  );
                },
              ),
            ),
            ChatInputBar(onSend: _send),
          ],
        ),
      ),
    );
  }
}

class _TopBar extends StatelessWidget {
  const _TopBar();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.fromLTRB(20, 2, 20, 14),
      decoration: const BoxDecoration(
        border: Border(bottom: BorderSide(color: AppColors.border)),
      ),
      child: const Text(
        'Chat',
        textAlign: TextAlign.center,
        style: TextStyle(
          fontSize: 17,
          fontWeight: FontWeight.w600,
          letterSpacing: -0.2,
          color: AppColors.textPrimary,
        ),
      ),
    );
  }
}

class _ErrorBanner extends StatelessWidget {
  final String text;

  const _ErrorBanner({required this.text});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: AppColors.bgSurface,
      child: Text(
        text,
        style: const TextStyle(fontSize: 11.5, color: AppColors.textSecondary),
      ),
    );
  }
}
