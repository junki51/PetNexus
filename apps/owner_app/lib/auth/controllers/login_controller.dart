import 'package:flutter/material.dart';

// กำหนดสถานะของหน้าจอ Login
enum AuthState { initial, loading, success, error }

class LoginController extends ChangeNotifier {
  AuthState _state = AuthState.initial;
  AuthState get state => _state;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  // ฟังก์ชันจำลองการเข้าสู่ระบบ (เชื่อม API จริงตรงนี้)
  Future<void> loginWithEmail() async {
    _setState(AuthState.loading);
    
    try {
      // TODO: ใส่ Logic เชื่อม Backend หรือ Firebase Auth
      await Future.delayed(const Duration(seconds: 2)); // จำลองดีเลย์ API
      
      _setState(AuthState.success);
      // นำทางไปหน้า Home ต่อไป
    } catch (e) {
      _errorMessage = e.toString();
      _setState(AuthState.error);
    }
  }

  Future<void> loginWithSocial(String provider) async {
    _setState(AuthState.loading);
    try {
      // TODO: ใส่ Logic ของ Google/Apple/Facebook Sign-in
      await Future.delayed(const Duration(seconds: 2));
      _setState(AuthState.success);
    } catch (e) {
      _errorMessage = "Social Login Error: $e";
      _setState(AuthState.error);
    }
  }

  void _setState(AuthState newState) {
    _state = newState;
    notifyListeners(); // แจ้งเตือน UI ให้ Re-build
  }
}