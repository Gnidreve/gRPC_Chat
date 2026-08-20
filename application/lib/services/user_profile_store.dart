import 'dart:math';

import 'package:shared_preferences/shared_preferences.dart';

/// The id, nickname and color a user has for themselves, persisted on the
/// device. [id] is generated once and kept stable across app restarts, so
/// the server recognizes returning users. There is currently no UI to
/// change nickname/color later — reinstalling the app is the only way to
/// pick different ones (which also resets id).
class UserProfile {
  final String id;
  final String nickname;
  final String color;

  const UserProfile({
    required this.id,
    required this.nickname,
    required this.color,
  });
}

class UserProfileStore {
  static const _idKey = 'user_profile_id';
  static const _nicknameKey = 'user_profile_nickname';
  static const _colorKey = 'user_profile_color';

  Future<UserProfile?> load() async {
    final prefs = await SharedPreferences.getInstance();
    final id = prefs.getString(_idKey);
    final nickname = prefs.getString(_nicknameKey);
    final color = prefs.getString(_colorKey);
    if (id == null || nickname == null || color == null) return null;
    return UserProfile(id: id, nickname: nickname, color: color);
  }

  Future<void> save(UserProfile profile) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_idKey, profile.id);
    await prefs.setString(_nicknameKey, profile.nickname);
    await prefs.setString(_colorKey, profile.color);
  }
}

/// A random 32-character hex id, generated once on first launch and then
/// persisted via [UserProfileStore] — this app's identity for the user,
/// independent of the server's own (ephemeral, in-memory) state.
String generateUserId() {
  final random = Random.secure();
  final bytes = List<int>.generate(16, (_) => random.nextInt(256));
  return bytes.map((b) => b.toRadixString(16).padLeft(2, '0')).join();
}
