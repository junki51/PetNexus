import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../core/constants/app_text_styles.dart';
import '../../layout/responsive_layout.dart';

class AppTextField extends StatelessWidget {
  final TextEditingController controller;

  final String label;

  final String hintText;

  final IconData? prefixIcon;

  final Widget? suffixIcon;

  final bool obscureText;

  final bool readOnly;

  final VoidCallback? onTap;

  final TextInputType keyboardType;

  final String? errorText;

  final int maxLines;

  const AppTextField({
    super.key,
    required this.controller,
    required this.label,
    required this.hintText,
    this.prefixIcon,
    this.suffixIcon,
    this.obscureText = false,
    this.readOnly = false,
    this.onTap,
    this.keyboardType = TextInputType.text,
    this.errorText,
    this.maxLines = 1,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: AppTextStyles.body(context),
        ),

        SizedBox(
          height: context.nh(8),
        ),

        TextField(
          controller: controller,
          obscureText: obscureText,
          readOnly: readOnly,
          onTap: onTap,
          keyboardType: keyboardType,
          maxLines: maxLines,
          decoration: InputDecoration(
            hintText: hintText,

            hintStyle: AppTextStyles.hint(context),

            prefixIcon: prefixIcon != null
                ? Icon(
                    prefixIcon,
                    size: context.icon(22),
                  )
                : null,

            suffixIcon: suffixIcon,

            filled: true,
            fillColor: Colors.white,

            errorText: errorText,

            contentPadding: EdgeInsets.symmetric(
              horizontal: context.nw(18),
              vertical: context.nh(16),
            ),

            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(
                context.radius(18),
              ),
              borderSide: BorderSide(
                color: AppColors.border,
              ),
            ),

            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(
                context.radius(18),
              ),
              borderSide: BorderSide(
                color: AppColors.border,
              ),
            ),

            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(
                context.radius(18),
              ),
              borderSide: BorderSide(
                color: AppColors.primary,
                width: 2,
              ),
            ),
          ),
        ),
      ],
    );
  }
}