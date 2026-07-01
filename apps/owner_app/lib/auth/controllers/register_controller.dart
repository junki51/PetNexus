import 'package:flutter/material.dart';
import '../../services/register_service.dart';

enum RegisterState { initial, loading, success, error }

class RegisterController extends ChangeNotifier {

  final RegisterService _registerService = RegisterService();
  RegisterState _state = RegisterState.initial;
  RegisterState get state => _state;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

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

  Future<bool> register(String email, String password, String confirmPassword) async {
    if (password != confirmPassword) {
      _errorMessage = "รหัสผ่านไม่ตรงกัน";
      _setState(RegisterState.error);
      return false;
    }

    _setState(RegisterState.loading);;

    try {
      final result = await _registerService.register(
        email: email,
        password: password,
        confirmPassword: confirmPassword,
      );
      if (result["success"]) {
        _setState(RegisterState.success);
        return true;
      }

      _errorMessage = result["message"];
      _setState(RegisterState.error);
      return false;

    } catch (e) {
      _errorMessage = "เชื่อมต่อเซิร์ฟเวอร์ไม่ได้: $e";
      _setState(RegisterState.error);
      return false;
    }
  }
  void _setState(RegisterState newState) {
    _state = newState;
    notifyListeners(); // แจ้งเตือน UI ให้ Re-build
  }
}