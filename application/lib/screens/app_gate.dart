import 'package:flutter/material.dart';

import '../services/chat_client.dart';
import '../services/location_service.dart';
import '../services/update_checker.dart';
import '../services/user_profile_store.dart';
import '../theme/app_theme.dart';
import 'chat_screen.dart';
import 'location_permission_screen.dart';
import 'profile_setup_screen.dart';
import 'update_screen.dart';

/// App entry point: checks for a mandatory update first (on every launch
/// and every resume from background), then loads the locally stored
/// profile (if any) and joins the chat server with it, showing
/// [ChatScreen]. If no profile exists yet, [ProfileSetupScreen] collects
/// one first.
class AppGate extends StatefulWidget {
  const AppGate({super.key});

  @override
  State<AppGate> createState() => _AppGateState();
}

enum _Status {
  loading,
  updateRequired,
  locationRequired,
  needsProfile,
  ready,
  error,
}

class _AppGateState extends State<AppGate> with WidgetsBindingObserver {
  final _profileStore = UserProfileStore();
  final _chatClient = ChatClient();
  final _updateChecker = UpdateChecker();
  final _locationService = LocationService();

  _Status _status = _Status.loading;
  AvailableUpdate? _update;
  User? _me;
  String? _error;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _checkForUpdateThenBootstrap();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state != AppLifecycleState.resumed) return;
    if (_status == _Status.locationRequired) {
      // Most likely returning from the system settings screen after
      // granting the permission there — retry immediately.
      _ensureLocationThenBootstrap();
    } else {
      _checkForUpdate();
    }
  }

  /// Runs on every launch: only proceeds to the location check (then
  /// profile/join) once no update is pending.
  Future<void> _checkForUpdateThenBootstrap() async {
    final update = await _updateChecker.check();
    if (update != null) {
      setState(() {
        _update = update;
        _status = _Status.updateRequired;
      });
      return;
    }
    await _ensureLocationThenBootstrap();
  }

  /// Runs on resume from background: only interrupts the current screen if
  /// an update is actually found, otherwise leaves whatever's showing
  /// alone.
  Future<void> _checkForUpdate() async {
    if (_status == _Status.updateRequired) return;
    final update = await _updateChecker.check();
    if (update != null && mounted) {
      setState(() {
        _update = update;
        _status = _Status.updateRequired;
      });
    }
  }

  /// Location is mandatory (Proximity Chat computes distance-to-message
  /// server-side from it) — only proceeds to profile/join once granted.
  Future<void> _ensureLocationThenBootstrap() async {
    if (!await _locationService.requestPermission()) {
      if (mounted) setState(() => _status = _Status.locationRequired);
      return;
    }
    await _bootstrap();
  }

  Future<void> _bootstrap() async {
    final profile = await _profileStore.load();
    if (profile == null) {
      setState(() => _status = _Status.needsProfile);
      return;
    }
    await _join(profile.id, profile.nickname, colorFromHex(profile.color));
  }

  /// Used by [ProfileSetupScreen] for a brand-new profile: no persisted id
  /// yet, so generate one now.
  Future<void> _joinNew(String nickname, Color color) {
    return _join(generateUserId(), nickname, color);
  }

  Future<void> _join(String id, String nickname, Color color) async {
    setState(() {
      _status = _Status.loading;
      _error = null;
    });
    try {
      final user = await _chatClient.join(id, nickname, color.toHex());
      await _profileStore.save(
        UserProfile(id: id, nickname: nickname, color: color.toHex()),
      );
      setState(() {
        _me = user;
        _status = _Status.ready;
      });
    } catch (e) {
      setState(() {
        _status = _Status.error;
        _error = '$e';
      });
    }
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    _chatClient.shutdown();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    switch (_status) {
      case _Status.loading:
        return const Scaffold(
          backgroundColor: AppColors.bgApp,
          body: Center(child: CircularProgressIndicator()),
        );
      case _Status.updateRequired:
        return UpdateScreen(update: _update!);
      case _Status.locationRequired:
        return LocationPermissionScreen(onGranted: _bootstrap);
      case _Status.needsProfile:
        return ProfileSetupScreen(onSubmit: _joinNew);
      case _Status.ready:
        return ChatScreen(chatClient: _chatClient, me: _me!);
      case _Status.error:
        return Scaffold(
          backgroundColor: AppColors.bgApp,
          body: Center(
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Text(
                    'Verbindung zum Server fehlgeschlagen',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    _error ?? '',
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 20),
                  FilledButton(
                    onPressed: _bootstrap,
                    style: FilledButton.styleFrom(
                      backgroundColor: AppColors.bubbleOwn,
                    ),
                    child: const Text('Erneut versuchen'),
                  ),
                ],
              ),
            ),
          ),
        );
    }
  }
}
