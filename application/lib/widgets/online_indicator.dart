import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

/// Zeile mittig unter der Top Bar: grüner Punkt + "X online". Der Wert wird
/// als Parameter übergeben, keine eigene Logik. Ein Tap auf die ganze Zeile
/// (unsichtbarer Hit-Bereich, keine visuelle Änderung) ruft [onTap] auf.
class OnlineIndicator extends StatelessWidget {
  final int count;
  final VoidCallback? onTap;

  const OnlineIndicator({super.key, required this.count, this.onTap});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      behavior: HitTestBehavior.opaque,
      onTap: onTap,
      child: Container(
        width: double.infinity,
        padding: const EdgeInsets.symmetric(vertical: 9),
        decoration: const BoxDecoration(
          color: AppColors.bgApp,
          border: Border(bottom: BorderSide(color: AppColors.border)),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const _PulsingDot(),
            const SizedBox(width: 7),
            Text.rich(
              TextSpan(
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: AppColors.textSecondary,
                ),
                children: [
                  TextSpan(
                    text: '$count',
                    style: const TextStyle(color: AppColors.textPrimary),
                  ),
                  const TextSpan(text: ' online'),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// Rein visueller Pulse-Effekt, keine fachliche Logik.
class _PulsingDot extends StatefulWidget {
  const _PulsingDot();

  @override
  State<_PulsingDot> createState() => _PulsingDotState();
}

class _PulsingDotState extends State<_PulsingDot>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller = AnimationController(
    vsync: this,
    duration: const Duration(milliseconds: 2200),
  )..repeat();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 24,
      height: 24,
      child: AnimatedBuilder(
        animation: _controller,
        builder: (context, _) {
          final t = _controller.value;
          return Stack(
            alignment: Alignment.center,
            children: [
              Opacity(
                opacity: (1 - t) * 0.35,
                child: Transform.scale(
                  scale: 0.6 + t * 1.3,
                  child: Container(
                    width: 10,
                    height: 10,
                    decoration: const BoxDecoration(
                      color: AppColors.accentGreen,
                      shape: BoxShape.circle,
                    ),
                  ),
                ),
              ),
              Container(
                width: 10,
                height: 10,
                decoration: const BoxDecoration(
                  color: AppColors.accentGreen,
                  shape: BoxShape.circle,
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}
