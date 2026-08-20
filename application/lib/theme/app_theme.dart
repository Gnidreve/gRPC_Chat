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
