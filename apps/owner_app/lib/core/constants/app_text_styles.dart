import 'package:flutter/material.dart';

import '../../layout/responsive_layout.dart';
import 'app_colors.dart';

class AppTextStyles {
  AppTextStyles._();

  static TextStyle title(BuildContext context) {
    return TextStyle(
      fontSize: context.nf(24),
      fontWeight: FontWeight.bold,
      color: AppColors.textPrimary,
    );
  }

  static TextStyle heading(BuildContext context) {
    return TextStyle(
      fontSize: context.nf(20),
      fontWeight: FontWeight.w700,
      color: AppColors.textPrimary,
    );
  }

  static TextStyle body(BuildContext context) {
    return TextStyle(
      fontSize: context.nf(16),
      color: AppColors.textPrimary,
    );
  }

  static TextStyle hint(BuildContext context) {
    return TextStyle(
      fontSize: context.nf(15),
      color: AppColors.hint,
    );
  }

  static TextStyle button(BuildContext context) {
    return TextStyle(
      fontSize: context.nf(17),
      color: Colors.white,
      fontWeight: FontWeight.bold,
    );
  }

  static TextStyle caption(BuildContext context) {
    return TextStyle(
      fontSize: context.nf(13),
      color: AppColors.textSecondary,
      fontWeight: FontWeight.bold,
    );
  }
}