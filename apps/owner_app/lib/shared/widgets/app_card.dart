import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../layout/responsive_layout.dart';

class AppCard extends StatelessWidget {
  final Widget child;

  final EdgeInsets? padding;

  final VoidCallback? onTap;

  final Color? color;

  final BorderRadiusGeometry? borderRadius;

  const AppCard({
    super.key,
    required this.child,
    this.padding,
    this.onTap,
    this.color,
    this.borderRadius,
  });

  @override
  Widget build(BuildContext context) {
    final effectiveBorderRadius =
        borderRadius ?? BorderRadius.circular(context.radius(18));

    final card = Container(
      width: double.infinity,

      padding: padding ?? EdgeInsets.all(context.nw(18)),

      decoration: BoxDecoration(
        color: color ?? Colors.white,

        borderRadius: effectiveBorderRadius,

        boxShadow: [
          BoxShadow(
            color: AppColors.shadow,
            blurRadius: context.nw(10),
            offset: Offset(0, context.nh(4)),
          ),
        ],
      ),

      child: child,
    );

    if (onTap == null) {
      return card;
    }

    return InkWell(
      customBorder: RoundedRectangleBorder(borderRadius: effectiveBorderRadius),
      onTap: onTap,
      child: card,
    );
  }
}
