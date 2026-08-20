import 'dart:convert';
import 'dart:io';

import 'package:http/http.dart' as http;
import 'package:open_filex/open_filex.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:path_provider/path_provider.dart';

/// A newer build found on GitHub than the one currently installed.
class AvailableUpdate {
  final int buildNumber;
  final String downloadUrl;

  const AvailableUpdate({required this.buildNumber, required this.downloadUrl});
}

/// Checks GitHub releases for a newer build than the one currently
/// installed, and can download + hand it to the system installer.
///
/// Release tags are "android-`<github-run-number>`" (see
/// .github/workflows/build-app.yml), and the app's own build number is set
/// to that same run number at build time — so comparing them is just an
/// integer comparison, no semver needed.
class UpdateChecker {
  static const _latestReleaseUrl =
      'https://api.github.com/repos/Gnidreve/gRPC_Chat/releases/latest';
  static const _apkAssetName = 'app-release.apk';

  /// Returns null if no update is available or the check itself failed
  /// (e.g. no network) — a failed check must never block the app.
  Future<AvailableUpdate?> check() async {
    try {
      final response = await http
          .get(Uri.parse(_latestReleaseUrl))
          .timeout(const Duration(seconds: 10));
      if (response.statusCode != 200) return null;

      final json = jsonDecode(response.body) as Map<String, dynamic>;
      final tagName = json['tag_name'] as String?;
      final remoteBuild = _parseBuildNumber(tagName);
      if (remoteBuild == null) return null;

      final assets = json['assets'] as List<dynamic>? ?? [];
      final apkAsset = assets.cast<Map<String, dynamic>>().firstWhere(
            (a) => a['name'] == _apkAssetName,
            orElse: () => const {},
          );
      final downloadUrl = apkAsset['browser_download_url'] as String?;
      if (downloadUrl == null) return null;

      final packageInfo = await PackageInfo.fromPlatform();
      final localBuild = int.tryParse(packageInfo.buildNumber) ?? 0;

      if (remoteBuild <= localBuild) return null;
      return AvailableUpdate(buildNumber: remoteBuild, downloadUrl: downloadUrl);
    } catch (_) {
      return null;
    }
  }

  static int? _parseBuildNumber(String? tagName) {
    if (tagName == null) return null;
    final match = RegExp(r'(\d+)$').firstMatch(tagName);
    if (match == null) return null;
    return int.tryParse(match.group(1)!);
  }

  /// Downloads the update APK, reporting progress in [0, 1], then hands it
  /// to the system package installer.
  Future<void> downloadAndInstall(
    AvailableUpdate update, {
    required void Function(double progress) onProgress,
  }) async {
    final client = http.Client();
    try {
      final request = http.Request('GET', Uri.parse(update.downloadUrl));
      final response = await client.send(request);
      if (response.statusCode != 200) {
        throw HttpException(
          'update download failed with status ${response.statusCode}',
          uri: Uri.parse(update.downloadUrl),
        );
      }
      final total = response.contentLength ?? 0;

      final dir = await getTemporaryDirectory();
      final file = File('${dir.path}/proximity-chat-update.apk');
      final sink = file.openWrite();

      var received = 0;
      await for (final chunk in response.stream) {
        sink.add(chunk);
        received += chunk.length;
        if (total > 0) onProgress(received / total);
      }
      await sink.close();

      await OpenFilex.open(
        file.path,
        type: 'application/vnd.android.package-archive',
      );
    } finally {
      client.close();
    }
  }
}
