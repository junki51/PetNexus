import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../layout/responsive_layout.dart';

class AppLoading extends StatelessWidget {
  const AppLoading({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: SizedBox(
        width: context.nw(36),
        height: context.nw(36),
        child: const CircularProgressIndicator(
          color: AppColors.primary,
        ),
      ),
    );
  }
}