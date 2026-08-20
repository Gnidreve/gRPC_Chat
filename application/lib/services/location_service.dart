import 'package:geolocator/geolocator.dart';

/// A plain (lat, lng) pair — deliberately not geolocator's [Position] type,
/// so callers (ChatClient, screens) don't depend on that package directly.
class LatLng {
  final double lat;
  final double lng;

  const LatLng(this.lat, this.lng);
}

/// Thrown when location can't be obtained because the permission is denied
/// or the device's location service is off. Proximity Chat requires
/// location for every Subscribe/SendMessage call, so callers should surface
/// this rather than silently sending no location.
class LocationUnavailable implements Exception {
  final String message;

  const LocationUnavailable(this.message);

  @override
  String toString() => message;
}

/// Wraps `geolocator`: permission handling and fetching the current
/// position. Coarse accuracy is enough — Proximity Chat only ever exposes
/// a rounded-to-kilometer distance, never the raw coordinate.
class LocationService {
  /// True if permission is already granted (does not prompt).
  Future<bool> hasPermission() async {
    final permission = await Geolocator.checkPermission();
    return permission == LocationPermission.always ||
        permission == LocationPermission.whileInUse;
  }

  /// Requests permission if not already granted. Returns true once granted.
  /// A permanently denied permission must be fixed in system settings —
  /// callers should offer [openSettings] in that case.
  Future<bool> requestPermission() async {
    if (!await Geolocator.isLocationServiceEnabled()) return false;

    var permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
    }
    return permission == LocationPermission.always ||
        permission == LocationPermission.whileInUse;
  }

  /// True if the permission was permanently denied ("don't ask again") —
  /// only the system settings screen can fix that from here on.
  Future<bool> isPermanentlyDenied() async {
    final permission = await Geolocator.checkPermission();
    return permission == LocationPermission.deniedForever;
  }

  Future<void> openSettings() => Geolocator.openAppSettings();

  /// Returns the current position. Throws [LocationUnavailable] if
  /// permission isn't granted or the location service is off — callers must
  /// have a granted permission (via [requestPermission]) before calling.
  Future<LatLng> current() async {
    if (!await Geolocator.isLocationServiceEnabled()) {
      throw const LocationUnavailable('Standortdienst ist deaktiviert');
    }
    if (!await hasPermission()) {
      throw const LocationUnavailable('Standort-Berechtigung fehlt');
    }
    final position = await Geolocator.getCurrentPosition(
      locationSettings: const LocationSettings(accuracy: LocationAccuracy.low),
    );
    return LatLng(position.latitude, position.longitude);
  }
}
