import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../layout/responsive_layout.dart';

class AppLogo extends StatelessWidget {
  final double _fontSize;

  const AppLogo({super.key}) : _fontSize = 40;

  const AppLogo.large({super.key}) : _fontSize = 60;

  @override
  Widget build(BuildContext context) {
    final fontSize = _fontSize;

    return RichText(
      text: TextSpan(
        style: TextStyle(
          fontSize: context.nf(fontSize),
          fontWeight: FontWeight.bold,
        ),
        children: const [
          TextSpan(
            text: 'Pet',
            style: TextStyle(color: AppColors.textPrimary),
          ),
          TextSpan(
            text: 'Nexus',
            style: TextStyle(color: AppColors.primary),
          ),
        ],
      ),
    );
  }
}
