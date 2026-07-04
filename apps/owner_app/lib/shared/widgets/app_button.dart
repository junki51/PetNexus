import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../core/constants/app_text_styles.dart';
import '../../layout/responsive_layout.dart';

class AppButton extends StatelessWidget {
  final String text;

  final VoidCallback? onPressed;

  final bool loading;

  final IconData? icon;

  final Color? backgroundColor;

  final Color? foregroundColor;

  final double? width;

  final double? height;

  const AppButton({
    super.key,
    required this.text,
    this.onPressed,
    this.loading = false,
    this.icon,
    this.backgroundColor,
    this.foregroundColor,
    this.width,
    this.height,
  });

  const AppButton.primary({
    super.key,
    required this.text,
    this.onPressed,
    this.loading = false,
    this.icon,
    this.width,
    this.height,
  }) : backgroundColor = AppColors.primary,
       foregroundColor = Colors.white;

  const AppButton.secondary({
    super.key,
    required this.text,
    this.onPressed,
    this.loading = false,
    this.icon,
    this.width,
    this.height,
  }) : backgroundColor = AppColors.surface,
       foregroundColor = AppColors.textPrimary;

  @override
  Widget build(BuildContext context) {
    final effectiveForegroundColor = foregroundColor ?? Colors.white;
    final isSecondary = backgroundColor == AppColors.surface;

    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(context.radius(18)),
        boxShadow: isSecondary
            ? [
                BoxShadow(
                  color: AppColors.shadow.withValues(alpha: 0.12), // เงานุ่ม ๆ จาง ๆ เข้ากับพื้นหลังขาว
                  blurRadius: context.nw(16),  // ความฟุ้งของเงา
                  spreadRadius: context.nw(1), // ระยะกระจายเงา
                  offset: Offset(0, context.nh(6)), // ดันเงาลงมาด้านล่างให้ปุ่มดูลอยขึ้น
                ),
              ]
            : null, // ปุ่ม Primary จะไม่มีเงาซ้อนตรงนี้ เพราะใช้ elevation ของปุ่มเอง
      ),
      child: SizedBox(
        width: width ?? double.infinity,
        height: height ?? context.nh(56),
        child: ElevatedButton(
          onPressed: loading ? null : onPressed,
          style: ElevatedButton.styleFrom(
            elevation: 3,
            shadowColor: AppColors.shadow,
            backgroundColor: backgroundColor ?? AppColors.primary,
            foregroundColor: effectiveForegroundColor, 
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(context.radius(18)),
            ),
          ),
          child: loading
              ? SizedBox(
                  width: context.nw(24),
                  height: context.nw(24),
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: effectiveForegroundColor,
                  ),
                )
              :  Row(
                  mainAxisAlignment: MainAxisAlignment.start,
                  children: [
                    if (icon != null) ...[
                      Icon(icon, size: context.icon(30)),
                      SizedBox(width: context.nw(8)),
                    ],
                    Expanded(
                      child: Text(
                        text,
                        style: AppTextStyles.button(
                          context,
                        ).copyWith(color: effectiveForegroundColor, fontSize: context.nf(20)
                        ),
                        textAlign: TextAlign.center,
                      ),
                    ),
                    if (icon != null) 
                    SizedBox(width: context.icon(30) + context.nw(8)), 
                  ],
                ),
        ),
      ),
    );
  }
}
