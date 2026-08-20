import 'package:flutter/material.dart';
import '../models/chat_message.dart';
import '../theme/app_theme.dart';

/// Stellt eine einzelne Chat-Nachricht dar (eigen oder fremd).
/// Reines Präsentations-Widget – bekommt alles über [message] rein.
class MessageBubble extends StatelessWidget {
  final ChatMessage message;

  const MessageBubble({super.key, required this.message});

  @override
  Widget build(BuildContext context) {
    final isOwn = message.isOwn;

    return Align(
      alignment: isOwn ? Alignment.centerRight : Alignment.centerLeft,
      child: ConstrainedBox(
        constraints: BoxConstraints(
          maxWidth: MediaQuery.of(context).size.width * 0.78,
        ),
        child: Padding(
          padding: const EdgeInsets.only(bottom: 10),
          child: Column(
            crossAxisAlignment:
                isOwn ? CrossAxisAlignment.end : CrossAxisAlignment.start,
            children: [
              if (!isOwn && message.senderName != null)
                Padding(
                  padding: const EdgeInsets.only(left: 2, bottom: 4),
                  child: Text(
                    message.senderName!,
                    style: const TextStyle(
                      fontSize: 11.5,
                      fontWeight: FontWeight.w600,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ),
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 13, vertical: 9),
                decoration: BoxDecoration(
                  color: isOwn ? AppColors.bubbleOwn : AppColors.bubbleOther,
                  borderRadius: BorderRadius.only(
                    topLeft: const Radius.circular(AppRadius.lg),
                    topRight: const Radius.circular(AppRadius.lg),
                    bottomLeft: Radius.circular(isOwn ? AppRadius.lg : 5),
                    bottomRight: Radius.circular(isOwn ? 5 : AppRadius.lg),
                  ),
                ),
                child: Text(
                  message.text,
                  style: TextStyle(
                    fontSize: 14.5,
                    height: 1.4,
                    color: isOwn
                        ? AppColors.bubbleOwnText
                        : AppColors.bubbleOtherText,
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.only(top: 4, left: 4, right: 4),
                child: Text(
                  message.time,
                  style: const TextStyle(
                    fontSize: 10.5,
                    color: AppColors.textSecondary,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
