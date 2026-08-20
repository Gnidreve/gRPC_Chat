import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:grpc_chat/main.dart';

void main() {
  testWidgets('ChatApp asks for a nickname and color on first launch',
      (WidgetTester tester) async {
    SharedPreferences.setMockInitialValues({});

    await tester.pumpWidget(const ChatApp());
    await tester.pumpAndSettle();

    expect(find.text('Willkommen'), findsOneWidget);
    expect(find.text('Nickname'), findsOneWidget);
  });
}
