import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import '../../api_config.dart';

enum RegisterState { initial, loading, success, error }

class RegisterController extends ChangeNotifier {
  RegisterState _state = RegisterState.initial;
  RegisterState get state => _state;

  // แยก State การมองเห็นรหัสผ่านออกจากกัน
  bool _isPasswordVisible = false;
  bool get isPasswordVisible => _isPasswordVisible;

  bool _isConfirmPasswordVisible = false;
  bool get isConfirmPasswordVisible => _isConfirmPasswordVisible;

  bool _isAcceptedTerms = false;
  bool get isAcceptedTerms => _isAcceptedTerms;

  void togglePasswordVisibility() {
    _isPasswordVisible = !_isPasswordVisible;
    notifyListeners();
  }

  void toggleConfirmPasswordVisibility() {
    _isConfirmPasswordVisible = !_isConfirmPasswordVisible;
    notifyListeners();
  }

  void toggleAcceptedTerms(bool? value) {
    _isAcceptedTerms = value ?? false;
    notifyListeners();
  }

  Future<void> register(String email, String password, String confirmPassword) async {
    if (!_isAcceptedTerms) return; // หรือแสดง Error

    _state = RegisterState.loading;
    notifyListeners();

    try {
      final response = await http.post(
        Uri.parse(ApiConfig.register), // เรียกใช้ผ่าน ApiConfig.register
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'email': email,
          'password': password,
          'password_confirmation': confirmPassword,
        }),
      );
      await Future.delayed(const Duration(seconds: 2));
      _state = RegisterState.success;
    } catch (e) {
      _state = RegisterState.error;
    }
    notifyListeners();
  }
}