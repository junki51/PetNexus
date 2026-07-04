import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../layout/responsive_layout.dart';

class AppSocialButton extends StatelessWidget {
  final IconData icon;
  final Color color;
  final VoidCallback onTap;

  const AppSocialButton({
    super.key,
    required this.icon,
    required this.color,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final size = context.nw(55).clamp(44.0, 52.0).toDouble();

    return Container(
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            // ปรับสีเงาให้จางลง (ลองใช้สี AppColors.shadow แล้วปรับความโปร่งแสง .withValues(alpha: 0.12))
            color: AppColors.shadow.withValues(alpha: 0.12), 
            blurRadius: context.nw(12), // 💡 ความฟุ้งของเงา (ยิ่งเยอะยิ่งนุ่ม)
            spreadRadius: context.nw(1), // 💡 ระยะการกระจายของเงา
            offset: Offset(0, context.nh(4)), // 💡 ดันเงาลงด้านล่างในแนวแกน Y เพื่อให้ดูปุ่มลอยขึ้นมา
          ),
        ],
      ),
      child: Material(
        color: AppColors.surface,
        shape: const CircleBorder(),
        elevation: 0, // 💡 ตั้งค่า elevation เป็น 0 เพราะเราใช้ BoxShadow แทน
        shadowColor: AppColors.shadow,
        child: InkWell(
          onTap: onTap,
          customBorder: const CircleBorder(),
          child: SizedBox(
            width: size,
            height: size,
            child: Icon(icon, color: color, size: context.icon(50)),
          ),
        ),
      ),
    );
  }
}
