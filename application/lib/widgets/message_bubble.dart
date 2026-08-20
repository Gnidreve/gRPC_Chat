import 'package:flutter/material.dart';
import '../gen/chat/v1/chat.pb.dart';
import '../theme/app_theme.dart';

/// Stellt eine einzelne Chat-Nachricht dar (eigen oder fremd).
/// Reines Präsentations-Widget – bekommt alles über [message] rein.
class MessageBubble extends StatelessWidget {
  final ChatMessage message;
  final bool isOwn;

  const MessageBubble({super.key, required this.message, required this.isOwn});

  @override
  Widget build(BuildContext context) {
    return Align(
      alignment: isOwn ? Alignment.centerRight : Alignment.centerLeft,
      child: ConstrainedBox(
        constraints: BoxConstraints(
          maxWidth: MediaQuery.sizeOf(context).width * 0.78,
        ),
        child: Padding(
          padding: const EdgeInsets.only(bottom: 10),
          child: Column(
            crossAxisAlignment:
                isOwn ? CrossAxisAlignment.end : CrossAxisAlignment.start,
            children: [
              if (!isOwn)
                Padding(
                  padding: const EdgeInsets.only(left: 2, bottom: 4),
                  child: Text(
                    message.user.nickname,
                    style: TextStyle(
                      fontSize: 11.5,
                      fontWeight: FontWeight.w600,
                      color: colorFromHex(message.user.color),
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
                  _formatFooter(message.sentAt.toDateTime(toLocal: true)),
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

  /// "Uhrzeit", or "Uhrzeit - Nkm entfernt" when the server computed a
  /// distance for this message (never for the recipient's own messages —
  /// the server omits distance_km for those, see chat.proto).
  String _formatFooter(DateTime time) {
    final hour = time.hour.toString().padLeft(2, '0');
    final minute = time.minute.toString().padLeft(2, '0');
    final timeText = '$hour:$minute';
    if (message.hasDistanceKm()) {
      return '$timeText - ${message.distanceKm}km entfernt';
    }
    return timeText;
  }
}
