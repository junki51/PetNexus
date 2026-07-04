import 'package:flutter/material.dart';

import '../../core/constants/app_text_styles.dart';
import '../../layout/responsive_layout.dart';

class AppSectionTitle extends StatelessWidget {
  final String title;

  final String? subtitle;

  const AppSectionTitle({
    super.key,
    required this.title,
    this.subtitle,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(
        vertical: context.nh(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: AppTextStyles.heading(context),
          ),

          if (subtitle != null) ...[
            SizedBox(
              height: context.nh(4),
            ),
            Text(
              subtitle!,
              style: AppTextStyles.caption(context),
            ),
          ]
        ],
      ),
    );
  }
}