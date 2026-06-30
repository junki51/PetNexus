import 'package:flutter/material.dart';
import "package:owner_app/auth/screens/first_screen.dart";

void main() {
  // บังคับให้ Flutter ผูก Widget Tree ให้เสร็จก่อนรันแอป (จำเป็นถ้ามี Initialized Service อื่นๆ)
  WidgetsFlutterBinding.ensureInitialized();
  
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'PetNexus',
      debugShowCheckedModeBanner: false, // ปิดป้ายกำกับ Debug มุมขวาบน
      theme: ThemeData(
        // ใช้สี Teal ของคุณเป็นสีหลักของแอป
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFF38A3A5),
        ),
        useMaterial3: true, // เปิดใช้งาน Material 3
        // แนะนำให้เพิ่มฟอนต์ภาษาไทย เช่น 'Kanit' หรือ 'Prompt' ใน pubspec.yaml แล้วมากำหนดตรงนี้
        // fontFamily: 'Kanit', 
      ),
      // กำหนดให้ LoginScreen เป็นหน้าแรกที่แอปจะแสดงผล
      home: const FirstScreen(), 
    );
  }
}
