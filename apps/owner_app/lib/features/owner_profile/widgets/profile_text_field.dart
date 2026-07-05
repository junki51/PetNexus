import 'package:flutter/material.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';

class ProfileTextField extends StatelessWidget {
  final TextEditingController controller;
  final String label;
  final String hintText;
  final IconData? prefixIcon;
  final Widget? prefix;
  final TextInputType keyboardType;
  final String? errorText;
  final bool isRequired;
  final int maxLines;

  const ProfileTextField({
    super.key,
    required this.controller,
    required this.label,
    required this.hintText,
    this.prefixIcon,
    this.prefix,
    this.keyboardType = TextInputType.text,
    this.errorText,
    this.isRequired = false,
    this.maxLines = 1,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // RichText for label with optional red asterisk
        RichText(
          text: TextSpan(
            text: label,
            style: AppTextStyles.body(context).copyWith(
              fontWeight: FontWeight.w600,
            ),
            children: [
              if (isRequired)
                TextSpan(
                  text: ' *',
                  style: TextStyle(
                    color: AppColors.error,
                    fontWeight: FontWeight.bold,
                  ),
                ),
            ],
          ),
        ),

        SizedBox(
          height: context.nh(8),
        ),

        TextField(
          controller: controller,
          keyboardType: keyboardType,
          maxLines: maxLines,
          style: AppTextStyles.body(context),
          decoration: InputDecoration(
            hintText: hintText,
            hintStyle: AppTextStyles.hint(context),
            prefixIcon: prefixIcon != null
                ? Icon(
                    prefixIcon,
                    size: context.icon(22),
                  )
                : null,
            prefix: prefix, // Custom prefix widget support (e.g. 🇹🇭 +66)
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
              borderSide: const BorderSide(
                color: AppColors.border,
              ),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(
                context.radius(18),
              ),
              borderSide: const BorderSide(
                color: AppColors.border,
              ),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(
                context.radius(18),
              ),
              borderSide: const BorderSide(
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
