import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../layout/responsive_layout.dart';

class StepTracker extends StatelessWidget {
  final int currentStep; // 1, 2, or 3

  const StepTracker({
    super.key,
    required this.currentStep,
  });

  @override
  Widget build(BuildContext context) {
    final dotSize = context.nw(12);
    final lineWidth = context.nw(50);
    final activeColor = AppColors.primary;
    final inactiveColor = Colors.white;
    const borderSide = BorderSide(color: AppColors.border, width: 2);

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        // Step 1 Dot
        _buildDot(context, currentStep >= 1, dotSize, activeColor, inactiveColor, borderSide),
        // Line 1
        _buildLine(context, currentStep >= 2, lineWidth, activeColor),
        // Step 2 Dot
        _buildDot(context, currentStep >= 2, dotSize, activeColor, inactiveColor, borderSide),
        // Line 2
        _buildLine(context, currentStep >= 3, lineWidth, activeColor),
        // Step 3 Dot
        _buildDot(context, currentStep >= 3, dotSize, activeColor, inactiveColor, borderSide),
      ],
    );
  }

  Widget _buildDot(
    BuildContext context,
    bool isActive,
    double size,
    Color activeColor,
    Color inactiveColor,
    BorderSide borderSide,
  ) {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: isActive ? activeColor : inactiveColor,
        shape: BoxShape.circle,
        border: isActive ? null : Border.fromBorderSide(borderSide),
        boxShadow: isActive
            ? null
            : [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.05),
                  blurRadius: 4,
                  offset: const Offset(0, 2),
                )
              ],
      ),
    );
  }

  Widget _buildLine(
    BuildContext context,
    bool isActive,
    double width,
    Color activeColor,
  ) {
    return Container(
      width: width,
      height: 3,
      color: isActive ? activeColor : AppColors.border,
    );
  }
}
