import 'package:flutter/material.dart';
import '../models/pet_model.dart';
import '../models/breed_model.dart';
import '../repositories/pet_repository.dart';

enum PetState {
  initial,
  loading,
  success,
  error,
}

class PetController extends ChangeNotifier {
  final PetRepository _repository = PetRepository();

  PetState _state = PetState.initial;
  PetState get state => _state;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  PetModel? _registeredPet;
  PetModel? get registeredPet => _registeredPet;

  // Selected pet type ('dog' or 'cat')
  String _selectedSpecies = 'dog';
  String get selectedSpecies => _selectedSpecies;

  // Dynamic breed catalog fetched from backend
  List<BreedModel> _breeds = [];
  List<BreedModel> get breeds => _breeds;

  BreedModel? _selectedBreed;
  BreedModel? get selectedBreed => _selectedBreed;

  // My pets list fetched from backend
  List<PetModel> _myPets = [];
  List<PetModel> get myPets => _myPets;

  PetModel? _selectedPet;
  PetModel? get selectedPet => _selectedPet;

  // Form temporary fields
  String _name = '';
  String? _gender;
  String? _dateOfBirth;
  int? _age;
  double? _weightKg;
  String? _microchipId;
  String? _avatarUrl;

  // Getters for form fields
  String get name => _name;
  String? get gender => _gender;
  String? get dateOfBirth => _dateOfBirth;
  int? get age => _age;
  double? get weightKg => _weightKg;
  String? get microchipId => _microchipId;
  String? get avatarUrl => _avatarUrl;

  void setSpecies(String species) {
    if (_selectedSpecies != species) {
      _selectedSpecies = species;
      _selectedBreed = null; // Clear selected breed when species changes
      _breeds = [];
      notifyListeners();
    }
  }

  Future<void> fetchBreeds(String species) async {
    _breeds = await _repository.getBreeds(species);
    notifyListeners();
  }

  Future<void> fetchMyPets() async {
    _myPets = await _repository.listMyPets();
    notifyListeners();
  }

  Future<void> fetchPetDetail(String id) async {
    try {
      _selectedPet = await _repository.getMyPet(id);
      notifyListeners();
    } catch (_) {}
  }

  void setSelectedBreed(BreedModel? breed) {
    _selectedBreed = breed;
    notifyListeners();
  }

  void setTemporaryInfo({
    required String name,
    String? gender,
    String? dateOfBirth,
    int? age,
    BreedModel? breed,
    double? weightKg,
    String? microchipId,
  }) {
    _name = name;
    _gender = gender;
    _dateOfBirth = dateOfBirth;
    _age = age;
    _selectedBreed = breed;
    _weightKg = weightKg;
    _microchipId = microchipId;
    notifyListeners();
  }

  void setAvatarUrl(String? url) {
    _avatarUrl = url;
    notifyListeners();
  }

  Future<bool> createPetProfile() async {
    _setLoading();
    try {
      _registeredPet = await _repository.createPet(
        species: _selectedSpecies,
        name: _name,
        breedId: _selectedBreed?.id,
        gender: _gender,
        dateOfBirth: _dateOfBirth,
        weightKg: _weightKg,
        microchipId: _microchipId,
        avatarUrl: _avatarUrl,
      );

      _state = PetState.success;
      notifyListeners();
      return true;
    } catch (e) {
      _setError(e.toString().replaceAll("Exception: ", ""));
      return false;
    }
  }

  void reset() {
    _state = PetState.initial;
    _errorMessage = null;
    _registeredPet = null;
    _selectedSpecies = 'dog';
    _name = '';
    _gender = null;
    _dateOfBirth = null;
    _age = null;
    _selectedBreed = null;
    _breeds = [];
    _weightKg = null;
    _microchipId = null;
    _avatarUrl = null;
    notifyListeners();
  }

  void _setLoading() {
    _errorMessage = null;
    _state = PetState.loading;
    notifyListeners();
  }

  void _setError(String message) {
    _errorMessage = message;
    _state = PetState.error;
    notifyListeners();
  }
}
