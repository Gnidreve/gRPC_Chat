import 'package:flutter/material.dart';

import '../services/location_service.dart';
import '../theme/app_theme.dart';

/// Mandatory, non-dismissible screen shown when location permission isn't
/// granted — Proximity Chat needs it to compute distance-to-message for
/// every recipient. There is no skip.
class LocationPermissionScreen extends StatefulWidget {
  final VoidCallback onGranted;

  const LocationPermissionScreen({super.key, required this.onGranted});

  @override
  State<LocationPermissionScreen> createState() =>
      _LocationPermissionScreenState();
}

class _LocationPermissionScreenState extends State<LocationPermissionScreen> {
  final _locationService = LocationService();
  bool _permanentlyDenied = false;
  bool _requesting = false;

  Future<void> _request() async {
    setState(() => _requesting = true);
    final granted = await _locationService.requestPermission();
    if (granted) {
      widget.onGranted();
      return;
    }
    final permanentlyDenied = await _locationService.isPermanentlyDenied();
    if (mounted) {
      setState(() {
        _requesting = false;
        _permanentlyDenied = permanentlyDenied;
      });
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
                    'Standort erforderlich',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w700,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    _permanentlyDenied
                        ? 'Die Berechtigung wurde dauerhaft abgelehnt. Bitte in den Einstellungen aktivieren.'
                        : 'Proximity Chat zeigt dir, wie weit Nachrichten entfernt wurden. '
                            'Dafür braucht die App Zugriff auf deinen ungefähren Standort — '
                            'deine genauen Koordinaten verlassen dabei nie den Server.',
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      fontSize: 13,
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 28),
                  FilledButton(
                    onPressed: _requesting
                        ? null
                        : _permanentlyDenied
                            ? _locationService.openSettings
                            : _request,
                    style: FilledButton.styleFrom(
                      backgroundColor: AppColors.bubbleOwn,
                    ),
                    child: Text(
                      _permanentlyDenied
                          ? 'Einstellungen öffnen'
                          : 'Standortzugriff erlauben',
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
