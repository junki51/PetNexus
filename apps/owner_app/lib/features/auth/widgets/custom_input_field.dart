import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../layout/responsive_layout.dart';

class CustomInputField extends StatelessWidget {
  final TextEditingController controller;
  final String hintText;
  final IconData prefixIcon;
  final bool isPassword;
  final bool obscureText;
  final VoidCallback? onToggleVisibility;
  final TextInputType? keyboardType;
  final TextInputAction? textInputAction;

  const CustomInputField({
    super.key,
    required this.controller,
    required this.hintText,
    required this.prefixIcon,
    this.isPassword = false,
    this.obscureText = false,
    this.onToggleVisibility,
    this.keyboardType,
    this.textInputAction,
  });

  @override
  Widget build(BuildContext context) {
    final height = context.nh(56).clamp(52.0, 60.0).toDouble();

    return Container(
      height: height,
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(height / 2),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.11),
            blurRadius: context.nw(16),
            spreadRadius: context.nw(1),
            offset: Offset(0, context.nh(6)),
          ),
        ],
      ),
      child: TextField(
        controller: controller,
        obscureText: obscureText,
        keyboardType: keyboardType,
        textInputAction: textInputAction,
        textAlignVertical: TextAlignVertical.center,
        style: TextStyle(
          color: AppColors.textPrimary,
          fontSize: context.nf(17),
        ),
        decoration: InputDecoration(
          isDense: true,
          /*hintText: hintText,
          hintStyle: TextStyle(
            color: Colors.black45, // Darker and clearer hint text
            fontSize: context.nf(17),
          ),*/
          hint: Text.rich(
            TextSpan(
              // แยกข้อความคำใบ้ออก โดยตัดเครื่องหมาย * ออกไปก่อน
              text: hintText.replaceAll('*', ''), 
              style: TextStyle(
                color: Colors.black45, // สีเทาคำใบ้ตามปกติของคุณ
                fontSize: context.nf(17),
              ),
              children: [
                // 💡 2. ตรวจสอบว่าถ้าข้อความเดิมมีดอกจัน ให้เติมดอกจันสีแดงต่อท้ายแบบนี้
                if (hintText.contains('*'))
                  TextSpan(
                    text: '*',
                    style: TextStyle(
                      color: Colors.red, // 🔴 บังคับเป็นสีแดงสดเฉพาะตัวดอกจัน
                      fontSize: context.nf(17),
                      fontWeight: FontWeight.bold,
                    ),
                  ),
              ],
            ),
          ),
          prefixIcon: Padding(
            padding: EdgeInsets.symmetric(horizontal: context.nw(12)),
            child: Icon(
              prefixIcon,
              color: AppColors.textSecondary,
              size: context.icon(24), // Slightly larger icon to match field height
            ),
          ),
          prefixIconConstraints: BoxConstraints(
            minWidth: context.nw(48),
            minHeight: height,
          ),
          suffixIcon: isPassword
              ? Padding(
                  padding: EdgeInsets.only(right: context.nw(12)),
                  child: IconButton(
                    padding: EdgeInsets.zero,
                    constraints: BoxConstraints(
                      minWidth: context.nw(36),
                      minHeight: height,
                    ),
                    icon: Icon(
                      obscureText
                          ? Icons.visibility_off_outlined
                          : Icons.visibility_outlined,
                      color: Colors.black45,
                      size: context.icon(24),
                    ),
                    onPressed: onToggleVisibility,
                  ),
                )
              : null,
          suffixIconConstraints: BoxConstraints(
            minWidth: context.nw(48),
            minHeight: height,
          ),
          border: InputBorder.none,
          enabledBorder: InputBorder.none,  // 💡 บังคับปิดเส้นขอบตอนปกติ
          focusedBorder: InputBorder.none,  // 💡 บังคับปิดเส้นขอบตอนกดคลิกกรอกข้อมูล
          contentPadding: EdgeInsets.zero,
        ),
      ),
    );
  }
}
