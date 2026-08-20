import 'package:shared_preferences/shared_preferences.dart';

/// The nickname and color a user picked for themselves, persisted on the
/// device. There is currently no UI to change these later — reinstalling
/// the app is the only way to pick a different nickname/color.
class UserProfile {
  final String nickname;
  final String color;

  const UserProfile({required this.nickname, required this.color});
}

class UserProfileStore {
  static const _nicknameKey = 'user_profile_nickname';
  static const _colorKey = 'user_profile_color';

  Future<UserProfile?> load() async {
    final prefs = await SharedPreferences.getInstance();
    final nickname = prefs.getString(_nicknameKey);
    final color = prefs.getString(_colorKey);
    if (nickname == null || color == null) return null;
    return UserProfile(nickname: nickname, color: color);
  }

  Future<void> save(UserProfile profile) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_nicknameKey, profile.nickname);
    await prefs.setString(_colorKey, profile.color);
  }
}
