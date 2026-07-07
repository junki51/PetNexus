import 'package:dio/dio.dart';
import '../../../../core/network/api_client.dart';
import '../models/pet_model.dart';
import '../models/breed_model.dart';

class PetRepository {
  Future<List<BreedModel>> getBreeds(String species) async {
    try {
      final response = await ApiClient.instance.dio.get(
        "/breeds",
        queryParameters: {"species": species},
      );
      final list = response.data['data'] as List<dynamic>;
      return list.map((item) => BreedModel.fromJson(item as Map<String, dynamic>)).toList();
    } on DioException catch (_) {
      // Fallback mocks
      if (species == 'dog') {
        return const [
          BreedModel(id: 'mock-dog-1', species: 'dog', name: 'Golden Retriever'),
          BreedModel(id: 'mock-dog-2', species: 'dog', name: 'Labrador Retriever'),
          BreedModel(id: 'mock-dog-3', species: 'dog', name: 'Poodle'),
          BreedModel(id: 'mock-dog-4', species: 'dog', name: 'Shiba Inu'),
          BreedModel(id: 'mock-dog-5', species: 'dog', name: 'Siberian Husky'),
          BreedModel(id: 'mock-dog-6', species: 'dog', name: 'Chihuahua'),
          BreedModel(id: 'mock-dog-7', species: 'dog', name: 'Pomeranian'),
          BreedModel(id: 'mock-dog-8', species: 'dog', name: 'Thai Bangkaew'),
        ];
      } else {
        return const [
          BreedModel(id: 'mock-cat-1', species: 'cat', name: 'Persian'),
          BreedModel(id: 'mock-cat-2', species: 'cat', name: 'Scottish Fold'),
          BreedModel(id: 'mock-cat-3', species: 'cat', name: 'British Shorthair'),
          BreedModel(id: 'mock-cat-4', species: 'cat', name: 'Siamese'),
          BreedModel(id: 'mock-cat-5', species: 'cat', name: 'Maine Coon'),
          BreedModel(id: 'mock-cat-6', species: 'cat', name: 'Ragdoll'),
          BreedModel(id: 'mock-cat-7', species: 'cat', name: 'Sphynx'),
          BreedModel(id: 'mock-cat-8', species: 'cat', name: 'Domestic Shorthair'),
        ];
      }
    }
  }

  Future<PetModel> createPet({
    required String species,
    required String name,
    String? breedId,
    String? gender,
    String? dateOfBirth,
    double? weightKg,
    String? microchipId,
    String? avatarUrl,
  }) async {
    try {
      final Map<String, dynamic> postData = {
        "species": species,
        "name": name,
      };
      if (breedId != null) postData["breed_id"] = breedId;
      if (gender != null) postData["gender"] = gender;
      if (dateOfBirth != null) postData["date_of_birth"] = dateOfBirth;
      if (weightKg != null) postData["weight_kg"] = weightKg;
      if (microchipId != null) postData["microchip_id"] = microchipId;
      if (avatarUrl != null) postData["avatar_url"] = avatarUrl;

      final response = await ApiClient.instance.dio.post(
        "/pets",
        data: postData,
      );
      final data = response.data['data'] as Map<String, dynamic>;
      return PetModel.fromJson(data);
    } on DioException catch (e) {
      // Graceful fallback to Mock response if the backend route doesn't exist yet (404)
      if (e.response?.statusCode == 404) {
        final mockId = "mock-pet-${DateTime.now().millisecondsSinceEpoch}";
        return PetModel(
          id: mockId,
          ownerProfileId: "mock-owner-profile-id",
          species: species,
          name: name,
          breedId: breedId,
          gender: gender,
          dateOfBirth: dateOfBirth,
          weightKg: weightKg,
          microchipId: microchipId,
          avatarUrl: avatarUrl,
          createdAt: DateTime.now().toIso8601String(),
          updatedAt: DateTime.now().toIso8601String(),
        );
      }
      throw Exception(e.response?.data["message"] ?? "Failed to register pet");
    }
  }

  Future<List<PetModel>> listMyPets() async {
    try {
      final response = await ApiClient.instance.dio.get("/pets");
      final list = response.data['data'] as List<dynamic>;
      return list.map((item) => PetModel.fromJson(item as Map<String, dynamic>)).toList();
    } on DioException catch (_) {
      return []; // Empty list on error — no crash
    }
  }

  Future<PetModel> getMyPet(String id) async {
    try {
      final response = await ApiClient.instance.dio.get("/pets/$id");
      final data = response.data['data'] as Map<String, dynamic>;
      return PetModel.fromJson(data);
    } on DioException catch (e) {
      throw Exception(e.response?.data["message"] ?? "Failed to load pet");
    }
  }

  Future<PetModel> updateMyPet(String id, Map<String, dynamic> updates) async {
    try {
      final response =
          await ApiClient.instance.dio.patch("/pets/$id", data: updates);
      final data = response.data['data'] as Map<String, dynamic>;
      return PetModel.fromJson(data);
    } on DioException catch (e) {
      throw Exception(e.response?.data["message"] ?? "Failed to update pet");
    }
  }
}
