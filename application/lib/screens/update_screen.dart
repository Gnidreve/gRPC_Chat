import 'package:flutter/material.dart';

import '../services/update_checker.dart';
import '../theme/app_theme.dart';

/// Mandatory, non-dismissible screen shown while a newer build is being
/// downloaded and handed off to the system installer. There is no skip —
/// the back button is intercepted and does nothing.
class UpdateScreen extends StatefulWidget {
  final AvailableUpdate update;

  const UpdateScreen({super.key, required this.update});

  @override
  State<UpdateScreen> createState() => _UpdateScreenState();
}

class _UpdateScreenState extends State<UpdateScreen> {
  final _checker = UpdateChecker();
  double _progress = 0;
  bool _readyToInstall = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _run();
  }

  Future<void> _run() async {
    setState(() {
      _error = null;
      _readyToInstall = false;
      _progress = 0;
    });
    try {
      await _checker.downloadAndInstall(
        widget.update,
        onProgress: (p) {
          if (mounted) setState(() => _progress = p);
        },
      );
      if (mounted) setState(() => _readyToInstall = true);
    } catch (e) {
      if (mounted) setState(() => _error = '$e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return PopScope(
      canPop: false,
      child: Scaffold(
        backgroundColor: AppColors.bgApp,
        body: SafeArea(
          child: Center(
            child: Padding(
              padding: const EdgeInsets.all(32),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Text(
                    'Update wird installiert',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w700,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    _error != null
                        ? 'Fehlgeschlagen: $_error'
                        : _readyToInstall
                            ? 'Bitte im Installations-Dialog bestätigen, um fortzufahren.'
                            : 'Ein Pflicht-Update wird heruntergeladen …',
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      fontSize: 13,
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 28),
                  if (_error == null && !_readyToInstall)
                    ClipRRect(
                      borderRadius: BorderRadius.circular(100),
                      child: LinearProgressIndicator(
                        value: _progress > 0 ? _progress : null,
                        minHeight: 6,
                        backgroundColor: AppColors.bgSurface,
                        color: AppColors.accentGreen,
                      ),
                    ),
                  if (_readyToInstall || _error != null) ...[
                    const SizedBox(height: 8),
                    FilledButton(
                      onPressed: _run,
                      style: FilledButton.styleFrom(
                        backgroundColor: AppColors.bubbleOwn,
                      ),
                      child: Text(
                        _readyToInstall
                            ? 'Installationsdialog erneut öffnen'
                            : 'Erneut versuchen',
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
