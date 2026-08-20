import 'package:flutter/material.dart';
import '../models/chat_message.dart';
import '../theme/app_theme.dart';
import '../widgets/chat_input_bar.dart';
import '../widgets/message_bubble.dart';
import '../widgets/online_indicator.dart';

/// Einziger Screen der App: Top Bar, Online-Indikator, Nachrichtenliste,
/// Input-Bar. Nutzt statische Demo-Daten aus chat_message.dart.
class ChatScreen extends StatelessWidget {
  const ChatScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgApp,
      body: SafeArea(
        child: Column(
          children: [
            _TopBar(),
            const OnlineIndicator(count: 14),
            Expanded(
              child: ListView.builder(
                padding: const EdgeInsets.fromLTRB(16, 18, 16, 10),
                itemCount: demoMessages.length,
                itemBuilder: (context, index) {
                  return MessageBubble(message: demoMessages[index]);
                },
              ),
            ),
            const ChatInputBar(),
          ],
        ),
      ),
    );
  }
}

class _TopBar extends StatelessWidget {
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
