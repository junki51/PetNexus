import 'package:flutter/material.dart';

extension ResponsiveContext on BuildContext {
  // ดึงความกว้างและความสูงจริงของหน้าจอเครื่องผู้ใช้
  double get screenWidth => MediaQuery.sizeOf(this).width;
  double get screenHeight => MediaQuery.sizeOf(this).height;

  // สมมติว่าขนาดหน้าจอใน Figma ดีไซน์เป็น Base อยู่ที่ 375 x 812 (มาตรฐานโมบายทั่วไป)
  // ฟังก์ชันนี้จะแปลงค่าความกว้างให้ยืดหยุ่นตามสัดส่วนจอจริง
  double nw(double width) => (screenWidth / 375) * width;

  // ฟังก์ชันแปลงค่าความสูงให้ยืดหยุ่นตามสัดส่วนจอจริง
  double nh(double height) => (screenHeight / 812) * height;

  // สำหรับฟอนต์ เพื่อให้สเกลตามขนาดจอได้เหมาะสม
  double nf(double fontSize) => (screenWidth / 375) * fontSize;
}