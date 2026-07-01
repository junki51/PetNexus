import 'package:flutter/material.dart';

extension ResponsiveContext on BuildContext {
  // ดึงความกว้างและความสูงจริงของหน้าจอเครื่องผู้ใช้
  double get screenWidth => MediaQuery.sizeOf(this).width;
  double get screenHeight => MediaQuery.sizeOf(this).height;

  /// [Senior Fix] จำกัดความกว้างสูงสุดที่จะนำไปคำนวณสัดส่วน 
  /// เพื่อไม่ให้ UI ขยายใหญ่เกินไปเวลาเปิดบน Tablet หรือ Desktop
  double get _responsiveWidth => screenWidth > 600 ? 450 : screenWidth;
  double get _responsiveHeight => screenHeight > 900 ? 812 : screenHeight;

  // ฟังก์ชันแปลงค่าความกว้าง โดยใช้ตัวแปรที่จำกัดเพดานแล้ว
  double nw(double width) => (_responsiveWidth / 375) * width;

  // ฟังก์ชันแปลงค่าความสูง โดยใช้ตัวแปรที่จำกัดเพดานแล้ว
  double nh(double height) => (_responsiveHeight / 812) * height;

  // สำหรับฟอนต์ สเกลตามความกว้างที่จำกัดเพดาน ป้องกันฟอนต์ตัวเท่าบ้านบน Desktop
  double nf(double fontSize) => (_responsiveWidth / 375) * fontSize;
}