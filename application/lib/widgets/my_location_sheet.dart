import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart';
import 'package:latlong2/latlong.dart' as latlong;

import '../services/location_service.dart';
import '../theme/app_theme.dart';

const _wideAreaZoom = 13.0;
const _pinSize = 40.0;

/// Opens a 50%-height, drag-to-dismiss bottom sheet showing a static
/// OpenStreetMap view (free, no API key) centered on the caller's own
/// current location — orientation only, not interactive: no zoom, no pan.
Future<void> showMyLocationSheet(
  BuildContext context, {
  required Color pinColor,
}) {
  return showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    backgroundColor: AppColors.bgApp,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(AppRadius.lg)),
    ),
    builder: (context) => SizedBox(
      height: MediaQuery.sizeOf(context).height * 0.5,
      child: _MyLocationSheetContent(pinColor: pinColor),
    ),
  );
}

class _MyLocationSheetContent extends StatefulWidget {
  final Color pinColor;

  const _MyLocationSheetContent({required this.pinColor});

  @override
  State<_MyLocationSheetContent> createState() =>
      _MyLocationSheetContentState();
}

class _MyLocationSheetContentState extends State<_MyLocationSheetContent> {
  late final Future<LatLng> _location = LocationService().current();

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const SizedBox(height: 10),
        Container(
          width: 36,
          height: 4,
          decoration: BoxDecoration(
            color: AppColors.border,
            borderRadius: BorderRadius.circular(100),
          ),
        ),
        const SizedBox(height: 10),
        Expanded(
          child: FutureBuilder<LatLng>(
            future: _location,
            builder: (context, snapshot) {
              if (snapshot.connectionState != ConnectionState.done) {
                return const Center(child: CircularProgressIndicator());
              }
              if (snapshot.hasError || !snapshot.hasData) {
                return const Center(
                  child: Text(
                    'Standort nicht verfügbar',
                    style: TextStyle(
                      fontSize: 13,
                      color: AppColors.textSecondary,
                    ),
                  ),
                );
              }

              final me = snapshot.data!;
              final center = latlong.LatLng(me.lat, me.lng);
              return FlutterMap(
                options: MapOptions(
                  initialCenter: center,
                  initialZoom: _wideAreaZoom,
                  interactionOptions: const InteractionOptions(
                    flags: InteractiveFlag.none,
                  ),
                ),
                children: [
                  TileLayer(
                    urlTemplate: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
                    userAgentPackageName: 'com.gnidreve.grpc_chat',
                  ),
                  MarkerLayer(
                    markers: [
                      Marker(
                        point: center,
                        width: _pinSize,
                        height: _pinSize,
                        child: Icon(
                          Icons.location_pin,
                          color: widget.pinColor,
                          size: _pinSize,
                        ),
                      ),
                    ],
                  ),
                ],
              );
            },
          ),
        ),
      ],
    );
  }
}
