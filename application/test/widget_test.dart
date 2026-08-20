import 'package:flutter_test/flutter_test.dart';
import 'package:geolocator_platform_interface/geolocator_platform_interface.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:grpc_chat/main.dart';

/// Stands in for the real platform plugin (no location hardware in the test
/// environment) — grants permission immediately, as a real device would
/// once the user taps through the OS prompt.
class _FakeGeolocatorPlatform extends GeolocatorPlatform {
  @override
  Future<bool> isLocationServiceEnabled() async => true;

  @override
  Future<LocationPermission> checkPermission() async =>
      LocationPermission.whileInUse;

  @override
  Future<LocationPermission> requestPermission() async =>
      LocationPermission.whileInUse;
}

void main() {
  GeolocatorPlatform.instance = _FakeGeolocatorPlatform();

  testWidgets('ChatApp asks for a nickname and color on first launch',
      (WidgetTester tester) async {
    SharedPreferences.setMockInitialValues({});

    await tester.pumpWidget(const ChatApp());
    await tester.pumpAndSettle();

    expect(find.text('Willkommen'), findsOneWidget);
    expect(find.text('Nickname'), findsOneWidget);
  });
}
