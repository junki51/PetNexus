import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';

class SettingsTile extends StatelessWidget {
  final IconData icon;
  final Color iconColor;
  final String label;
  final Widget? trailing;
  final VoidCallback onTap;

  const SettingsTile({
    super.key,
    required this.icon,
    required this.iconColor,
    required this.label,
    this.trailing,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(context.radius(16)),
      child: Padding(
        padding: EdgeInsets.symmetric(
            horizontal: context.nw(16), vertical: context.nh(14)),
        child: Row(
          children: [
            Container(
              width: context.nw(32),
              height: context.nw(32),
              decoration: BoxDecoration(
                color: iconColor.withValues(alpha: 0.12),
                borderRadius: BorderRadius.circular(context.radius(8)),
              ),
              child: Icon(icon, color: iconColor, size: context.icon(18)),
            ),
            SizedBox(width: context.nw(12)),
            Expanded(
              child: Text(
                label,
                style: AppTextStyles.body(context)
                    .copyWith(fontSize: context.nf(14)),
              ),
            ),
            if (trailing != null) ...[
              trailing!,
              SizedBox(width: context.nw(8)),
            ],
            Icon(Icons.chevron_right_rounded,
                color: AppColors.textSecondary, size: context.icon(20)),
          ],
        ),
      ),
    );
  }
}
