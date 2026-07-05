import 'package:flutter/material.dart';
import '../models/owner_profile_model.dart';
import '../repositories/owner_profile_repository.dart';

enum OwnerProfileState {
  initial,
  loading,
  success,
  error,
}

class OwnerProfileController extends ChangeNotifier {
  final OwnerProfileRepository _repository = OwnerProfileRepository();

  OwnerProfileState _state = OwnerProfileState.initial;
  OwnerProfileState get state => _state;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  OwnerProfileModel? _profile;
  OwnerProfileModel? get profile => _profile;

  //======================
  // Create Profile
  //======================
  Future<bool> createProfile({
    required String firstName,
    required String lastName,
    required String phoneNumber,
    required int age,
    String? gender,
    String? avatarUrl,
    String? address,
    String? province,
  }) async {
    _setLoading();

    try {
      // Calculate date of birth based on age: e.g. if age is 25, DOB year = current_year - 25.
      // Day and month are kept as today's date.
      final now = DateTime.now();
      final dob = DateTime(now.year - age, now.month, now.day);
      final dobString = "${dob.year}-${dob.month.toString().padLeft(2, '0')}-${dob.day.toString().padLeft(2, '0')}";

      _profile = await _repository.createProfile(
        firstName: firstName,
        lastName: lastName,
        phoneNumber: phoneNumber,
        gender: gender,
        dateOfBirth: dobString,
        avatarUrl: avatarUrl,
        addressLine1: address,
        province: province,
      );

      _state = OwnerProfileState.success;
      notifyListeners();
      return true;
    } catch (e) {
      _setError(e.toString().replaceAll("Exception: ", ""));
      return false;
    }
  }

  //======================
  // Fetch Profile
  //======================
  Future<void> fetchProfile() async {
    _setLoading();
    try {
      _profile = await _repository.getProfile();
      _state = OwnerProfileState.success;
      notifyListeners();
    } catch (e) {
      // If error is 404/not found, set to initial state so user can create profile
      _profile = null;
      if (e.toString().contains("not found") || e.toString().contains("404")) {
        _state = OwnerProfileState.initial;
      } else {
        _setError(e.toString().replaceAll("Exception: ", ""));
      }
      notifyListeners();
    }
  }

  void clearProfile() {
    _profile = null;
    _state = OwnerProfileState.initial;
    _errorMessage = null;
    notifyListeners();
  }

  //======================
  // Helpers
  //======================
  void _setLoading() {
    _errorMessage = null;
    _state = OwnerProfileState.loading;
    notifyListeners();
  }

  void _setError(String message) {
    _errorMessage = message;
    _state = OwnerProfileState.error;
    notifyListeners();
  }

  void resetState() {
    _errorMessage = null;
    _state = OwnerProfileState.initial;
    notifyListeners();
  }
}
