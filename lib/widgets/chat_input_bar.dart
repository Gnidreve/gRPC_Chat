import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

/// Eingabezeile am unteren Bildschirmrand.
/// Rein visuelles Template: kein State, kein Senden-Callback.
/// TODO: TextEditingController + onSend-Callback ergänzen, sobald
/// die eigentliche Chat-Logik angebunden wird.
class ChatInputBar extends StatelessWidget {
  const ChatInputBar({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.only(
        left: 12,
        right: 12,
        top: 10,
        bottom: 10 + MediaQuery.of(context).padding.bottom,
      ),
      decoration: const BoxDecoration(
        color: AppColors.bgApp,
        border: Border(top: BorderSide(color: AppColors.border)),
      ),
      child: Row(
        children: [
          Expanded(
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              decoration: BoxDecoration(
                color: AppColors.bgSurface,
                borderRadius: BorderRadius.circular(100),
                border: Border.all(color: AppColors.border),
              ),
              child: const TextField(
                decoration: InputDecoration(
                  hintText: 'Nachricht',
                  hintStyle: TextStyle(
                    color: AppColors.textSecondary,
                    fontSize: 14.5,
                  ),
                  border: InputBorder.none,
                  isDense: true,
                  contentPadding: EdgeInsets.symmetric(vertical: 10),
                ),
                style: TextStyle(
                  color: AppColors.textPrimary,
                  fontSize: 14.5,
                ),
              ),
            ),
          ),
          const SizedBox(width: 8),
          _SendButton(onPressed: () {
            // TODO: Sende-Logik anbinden.
          }),
        ],
      ),
    );
  }
}

class _SendButton extends StatelessWidget {
  final VoidCallback onPressed;

  const _SendButton({required this.onPressed});

  @override
  Widget build(BuildContext context) {
    return Material(
      color: AppColors.bubbleOwn,
      shape: const CircleBorder(),
      child: InkWell(
        onTap: onPressed,
        customBorder: const CircleBorder(),
        child: const SizedBox(
          width: 34,
          height: 34,
          child: Icon(
            Icons.arrow_upward_rounded,
            size: 17,
            color: AppColors.bubbleOwnText,
          ),
        ),
      ),
    );
  }
}
