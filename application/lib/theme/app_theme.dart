import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// Zentrale Farb- und Formtokens, 1:1 aus dem Web-Mockup übernommen.
class AppColors {
  AppColors._();

  static const bgApp = Color(0xFFFFFFFF);
  static const bgSurface = Color(0xFFF4F4F5);
  static const border = Color(0xFFE5E5E5);

  static const textPrimary = Color(0xFF111113);
  static const textSecondary = Color(0xFF8A8A8F);

  static const accentGreen = Color(0xFF12B76A);

  static const bubbleOwn = Color(0xFF111113);
  static const bubbleOwnText = Color(0xFFFFFFFF);
  static const bubbleOther = Color(0xFFF4F4F5);
  static const bubbleOtherText = Color(0xFF111113);

  /// Fixed palette users pick their own color from on first launch.
  static const userColors = [
    Color(0xFF12B76A),
    Color(0xFFF04438),
    Color(0xFFF79009),
    Color(0xFF7A5AF8),
    Color(0xFF0BA5EC),
    Color(0xFFEE46BC),
    Color(0xFF84CC16),
    Color(0xFFF97066),
  ];
}

extension AppColorHex on Color {
  /// "#RRGGBB", the wire format sent to and received from the server.
  String toHex() =>
      '#${(toARGB32() & 0xFFFFFF).toRadixString(16).padLeft(6, '0').toUpperCase()}';
}

/// Parses a "#RRGGBB" string as sent by the server. Falls back to
/// [AppColors.textSecondary] for anything malformed.
Color colorFromHex(String hex) {
  final cleaned = hex.replaceFirst('#', '');
  final value = int.tryParse('ff$cleaned', radix: 16);
  return value == null ? AppColors.textSecondary : Color(value);
}

class AppRadius {
  AppRadius._();

  static const lg = 22.0;
  static const md = 14.0;
}

/// Fertiges ThemeData, das an MaterialApp übergeben wird.
class AppTheme {
  AppTheme._();

  static ThemeData get light {
    final base = ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,
      scaffoldBackgroundColor: AppColors.bgApp,
      colorScheme: ColorScheme.fromSeed(
        seedColor: AppColors.textPrimary,
        brightness: Brightness.light,
        surface: AppColors.bgApp,
      ),
    );

    return base.copyWith(
      textTheme: GoogleFonts.interTextTheme(base.textTheme).apply(
        bodyColor: AppColors.textPrimary,
        displayColor: AppColors.textPrimary,
      ),
      dividerColor: AppColors.border,
    );
  }
}
