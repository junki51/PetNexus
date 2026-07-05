import 'package:flutter/material.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';

class SectionTitle extends StatelessWidget {
  final String title;
  final VoidCallback? onBack;

  const SectionTitle({
    super.key,
    required this.title,
    this.onBack,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        if (onBack != null) ...[
          GestureDetector(
            onTap: onBack,
            child: CircleAvatar(
              radius: context.nw(20),
              backgroundColor: AppColors.primaryLight,
              child: Icon(
                Icons.chevron_left,
                color: AppColors.primary,
                size: context.icon(28),
              ),
            ),
          ),
          SizedBox(width: context.nw(12)),
        ],
        Text(
          title,
          style: AppTextStyles.title(context).copyWith(
            fontSize: context.nf(24),
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
      ],
    );
  }
}
