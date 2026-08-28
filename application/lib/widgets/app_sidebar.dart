import 'package:flutter/material.dart';
import '../theme/app_theme.dart';

/// Rein kosmetische Sidebar — kein Inhalt, keine Funktion. Öffnen/Schließen
/// (per Hamburger-Icon oder Edge-Swipe-Geste), die Slide-Animation und das
/// Overlay hinter der Sidebar kommen von Flutters eingebautem
/// Scaffold/Drawer-Mechanismus (Material Best Practice) — hier wird nur das
/// Aussehen angepasst (abgerundete obere Ecken).
class AppSidebar extends StatelessWidget {
  const AppSidebar({super.key});

  @override
  Widget build(BuildContext context) {
    return Drawer(
      backgroundColor: AppColors.bgApp,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.only(
          topLeft: Radius.circular(20),
          topRight: Radius.circular(20),
        ),
      ),
      child: const SizedBox.expand(),
    );
  }
}
