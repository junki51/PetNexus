import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../pet/models/pet_model.dart';

class StatusBarWidget extends StatelessWidget {
  final PetModel? pet;

  const StatusBarWidget({super.key, this.pet});

  @override
  Widget build(BuildContext context) {
    // Mock status values (0.0 – 1.0)
    const statuses = [
      _StatusItem(icon: Icons.restaurant_rounded, label: 'อาหาร', value: 0.7, color: Color(0xFFFF8A65)),
      _StatusItem(icon: Icons.favorite_rounded, label: 'สุขภาพ', value: 0.9, color: Color(0xFFEF5350)),
      _StatusItem(icon: Icons.favorite_border_rounded, label: 'ความรัก', value: 0.5, color: Color(0xFFEC407A)),
      _StatusItem(icon: Icons.water_drop_rounded, label: 'อาบน้ำ', value: 0.3, color: Color(0xFF42A5F5)),
    ];

    return Padding(
      padding: EdgeInsets.symmetric(horizontal: context.nw(24)),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: statuses
            .map((s) => _StatusButton(item: s, enabled: pet != null))
            .toList(),
      ),
    );
  }
}

class _StatusItem {
  final IconData icon;
  final String label;
  final double value;
  final Color color;

  const _StatusItem({
    required this.icon,
    required this.label,
    required this.value,
    required this.color,
  });
}

class _StatusButton extends StatelessWidget {
  final _StatusItem item;
  final bool enabled;

  const _StatusButton({required this.item, required this.enabled});

  @override
  Widget build(BuildContext context) {
    final size = context.nw(56);
    return GestureDetector(
      onTap: enabled
          ? () => ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text('${item.label}: Coming soon! 🐾'),
                  duration: const Duration(seconds: 1),
                ),
              )
          : null,
      child: Column(
        children: [
          SizedBox(
            width: size,
            height: size,
            child: Stack(
              alignment: Alignment.center,
              children: [
                // Background circle fill
                SizedBox.expand(
                  child: CircularProgressIndicator(
                    value: enabled ? item.value : 0,
                    strokeWidth: 4,
                    backgroundColor: AppColors.border,
                    valueColor: AlwaysStoppedAnimation<Color>(item.color),
                  ),
                ),
                // Inner circle
                Container(
                  width: size * 0.72,
                  height: size * 0.72,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: enabled
                        ? item.color.withValues(alpha: 0.12)
                        : AppColors.border.withValues(alpha: 0.4),
                  ),
                  child: Icon(
                    item.icon,
                    color: enabled ? item.color : AppColors.textSecondary,
                    size: context.icon(22),
                  ),
                ),
              ],
            ),
          ),
          SizedBox(height: context.nh(4)),
          Text(
            item.label,
            style: AppTextStyles.caption(context).copyWith(
              fontSize: context.nf(11),
              color: AppColors.textSecondary,
              fontWeight: FontWeight.normal,
            ),
          ),
        ],
      ),
    );
  }
}
