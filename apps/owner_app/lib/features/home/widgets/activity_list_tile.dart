import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../controllers/home_controller.dart';

class ActivityListTile extends StatelessWidget {
  final ActivityItem activity;

  const ActivityListTile({super.key, required this.activity});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: EdgeInsets.only(bottom: context.nh(8)),
      padding: EdgeInsets.symmetric(
          horizontal: context.nw(12), vertical: context.nh(10)),
      decoration: BoxDecoration(
        color: activity.color.withValues(alpha: 0.08),
        borderRadius: BorderRadius.circular(context.radius(12)),
        border: Border.all(
          color: activity.color.withValues(alpha: 0.2),
          width: 1,
        ),
      ),
      child: Row(
        children: [
          Container(
            width: context.nw(36),
            height: context.nw(36),
            decoration: BoxDecoration(
              color: activity.color.withValues(alpha: 0.15),
              shape: BoxShape.circle,
            ),
            child: Icon(activity.icon, color: activity.color, size: context.icon(18)),
          ),
          SizedBox(width: context.nw(12)),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  activity.title,
                  style: AppTextStyles.body(context).copyWith(
                    fontWeight: FontWeight.w600,
                    fontSize: context.nf(14),
                  ),
                ),
                if (activity.subtitle.isNotEmpty)
                  Text(
                    activity.subtitle,
                    style: AppTextStyles.caption(context).copyWith(
                      fontSize: context.nf(12),
                      fontWeight: FontWeight.normal,
                    ),
                  ),
              ],
            ),
          ),
          Text(
            activity.time,
            style: TextStyle(
              fontSize: context.nf(11),
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}
