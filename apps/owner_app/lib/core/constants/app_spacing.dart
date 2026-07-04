import 'package:flutter/material.dart';
import '../../layout/responsive_layout.dart';

class AppSpacing {
  AppSpacing._();

  static EdgeInsets page(BuildContext context) {
    return EdgeInsets.all(context.nw(24));
  }

  static EdgeInsets horizontal(BuildContext context) {
    return EdgeInsets.symmetric(
      horizontal: context.nw(20),
    );
  }

  static EdgeInsets input(BuildContext context) {
    return EdgeInsets.symmetric(
      horizontal: context.nw(18),
      vertical: context.nh(16),
    );
  }

  static SizedBox h(BuildContext context, double value) {
    return SizedBox(
      height: context.nh(value),
    );
  }

  static SizedBox w(BuildContext context, double value) {
    return SizedBox(
      width: context.nw(value),
    );
  }
}