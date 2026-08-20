import 'package:flutter_test/flutter_test.dart';

import 'package:grpc_chat/main.dart';

void main() {
  testWidgets('ChatApp shows the chat screen with demo messages',
      (WidgetTester tester) async {
    await tester.pumpWidget(const ChatApp());

    expect(find.text('Chat'), findsOneWidget);
    expect(find.text('14 online'), findsOneWidget);
    expect(find.text('Hey, seid ihr schon am Set?'), findsOneWidget);
  });
}
