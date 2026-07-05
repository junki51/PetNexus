import 'package:flutter/cupertino.dart';
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
    final size = context.nw(60).clamp(55.0, 63.0).toDouble();
    final isGoogle = icon == CupertinoIcons.globe && color == AppColors.google;

    return Container(
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: AppColors.shadow.withValues(alpha: 0.12),
            blurRadius: context.nw(12),
            spreadRadius: context.nw(1),
            offset: Offset(0, context.nh(4)),
          ),
        ],
      ),
      child: Material(
        color: AppColors.surface,
        shape: const CircleBorder(),
        elevation: 0,
        shadowColor: AppColors.shadow,
        child: InkWell(
          onTap: onTap,
          customBorder: const CircleBorder(),
          child: SizedBox(
            width: size,
            height: size,
            child: Center(
              child: isGoogle
                  ? CustomPaint(
                      size: Size(context.icon(50), context.icon(50)),
                      painter: GoogleLogoPainter(),
                    )
                  : Icon(
                      icon,
                      color: color,
                      size: context.icon(55,), // 40 is more proportional for standard icons
                    ),
            ),
          ),
        ),
      ),
    );
  }
}

class GoogleLogoPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final double w = size.width;
    final double h = size.height;
    final double radius = w / 2;

    final paint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = w * 0.22
      ..strokeCap = StrokeCap.square;

    final center = Offset(w / 2, h / 2);
    final rect = Rect.fromCircle(
      center: center,
      radius: radius - paint.strokeWidth / 2,
    );

    // Red arc: Top segment
    paint.color = const Color(0xFFEA4335);
    canvas.drawArc(rect, -2.356, 2.094, false, paint);

    // Yellow arc: Left segment
    paint.color = const Color(0xFFFBBC05);
    canvas.drawArc(rect, 2.443, 1.396, false, paint);

    // Green arc: Bottom segment
    paint.color = const Color(0xFF34A853);
    canvas.drawArc(rect, 0.698, 1.745, false, paint);

    // Blue arc: Right segment
    paint.color = const Color(0xFF4285F4);
    canvas.drawArc(rect, -0.349, 1.047, false, paint);

    // Horizontal bar: Blue
    final barPaint = Paint()
      ..color = const Color(0xFF4285F4)
      ..style = PaintingStyle.fill;
    final barRect = Rect.fromLTRB(
      center.dx,
      center.dy - paint.strokeWidth / 2,
      w,
      center.dy + paint.strokeWidth / 2,
    );
    canvas.drawRect(barRect, barPaint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
