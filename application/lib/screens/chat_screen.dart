import 'dart:async';

import 'package:flutter/material.dart';

import '../services/chat_client.dart';
import '../services/location_service.dart';
import '../theme/app_theme.dart';
import '../widgets/chat_input_bar.dart';
import '../widgets/message_bubble.dart';
import '../widgets/online_indicator.dart';

/// Einziger Chat-Screen: Top Bar, Nachrichtenliste, Input-Bar. Abonniert
/// den Server-Stream über [chatClient] und zeigt Nachrichten live an.
class ChatScreen extends StatefulWidget {
  final ChatClient chatClient;
  final User me;

  const ChatScreen({super.key, required this.chatClient, required this.me});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

// Reconnect delays, in order — a mobile connection commonly drops after the
// app sits backgrounded for a while, so a killed Subscribe stream needs to
// come back on its own instead of leaving the chat stuck on an error.
const _reconnectDelays = [
  Duration(milliseconds: 300),
  Duration(seconds: 1),
  Duration(seconds: 2),
  Duration(seconds: 4),
  Duration(seconds: 8),
];

class _ChatScreenState extends State<ChatScreen> with WidgetsBindingObserver {
  final _messages = <ChatMessage>[];
  final _scrollController = ScrollController();
  final _locationService = LocationService();
  StreamSubscription<ChatEvent>? _subscription;
  Timer? _reconnectTimer;
  int _reconnectAttempt = 0;
  int _onlineCount = 0;
  String? _error;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _subscribeToEvents();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    // Resuming from background re-subscribes with a fresh location, even
    // if the underlying stream survived backgrounding — the distance shown
    // for new messages should reflect where the user is now, not where
    // they were when the app was last opened.
    if (state == AppLifecycleState.resumed) {
      _subscribeToEvents();
    }
  }

  Future<void> _subscribeToEvents() async {
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
    _subscription?.cancel();

    final LatLng location;
    try {
      location = await _locationService.current();
    } catch (_) {
      if (!mounted) return;
      setState(() => _error = 'Standort nicht verfügbar, versuche erneut …');
      _scheduleReconnect();
      return;
    }
    if (!mounted) return;

    _subscription =
        widget.chatClient.subscribe(widget.me.id, location).listen(
      (event) {
        if (!mounted) return;
        setState(() {
          _reconnectAttempt = 0;
          _error = null;
          if (event.hasMessage()) {
            _messages.add(event.message);
          } else if (event.hasPresence()) {
            _onlineCount = event.presence.onlineCount;
          }
        });
        if (event.hasMessage()) _scrollToBottom();
      },
      onError: (Object error) {
        if (!mounted) return;
        setState(() => _error = 'Verbindung verloren, verbinde neu …');
        _scheduleReconnect();
      },
      onDone: () {
        if (!mounted) return;
        _scheduleReconnect();
      },
    );
  }

  void _scheduleReconnect() {
    if (_reconnectTimer != null) return; // already scheduled
    final delay = _reconnectDelays[
        _reconnectAttempt.clamp(0, _reconnectDelays.length - 1)];
    _reconnectAttempt++;
    _reconnectTimer = Timer(delay, () {
      _reconnectTimer = null;
      if (!mounted) return;
      // The next Subscribe call replays the full history again, so avoid
      // duplicating what's already shown.
      setState(() => _messages.clear());
      _subscribeToEvents();
    });
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    _reconnectTimer?.cancel();
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
      final location = await _locationService.current();
      await widget.chatClient.sendMessage(widget.me.id, text, location);
    } catch (e) {
      if (!mounted) return;
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
            OnlineIndicator(count: _onlineCount),
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
