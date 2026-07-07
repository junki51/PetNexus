class PetModel {
  final String id;
  final String ownerProfileId;
  final String? breedId;
  final String species;
  final String name;
  final String? gender;
  final String? dateOfBirth;
  final double? weightKg;
  final String? microchipId;
  final String? avatarUrl;
  final String createdAt;
  final String updatedAt;

  const PetModel({
    required this.id,
    required this.ownerProfileId,
    this.breedId,
    required this.species,
    required this.name,
    this.gender,
    this.dateOfBirth,
    this.weightKg,
    this.microchipId,
    this.avatarUrl,
    required this.createdAt,
    required this.updatedAt,
  });

  factory PetModel.fromJson(Map<String, dynamic> json) {
    return PetModel(
      id: json['id'] as String? ?? '',
      ownerProfileId: json['owner_profile_id'] as String? ?? '',
      breedId: json['breed_id'] as String?,
      species: json['species'] as String? ?? 'dog',
      name: json['name'] as String? ?? '',
      gender: json['gender'] as String?,
      dateOfBirth: json['date_of_birth'] as String?,
      weightKg: (json['weight_kg'] as num?)?.toDouble(),
      microchipId: json['microchip_id'] as String?,
      avatarUrl: json['avatar_url'] as String?,
      createdAt: json['created_at'] as String? ?? '',
      updatedAt: json['updated_at'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      if (breedId != null) 'breed_id': breedId,
      'species': species,
      'name': name,
      if (gender != null) 'gender': gender,
      if (dateOfBirth != null) 'date_of_birth': dateOfBirth,
      if (weightKg != null) 'weight_kg': weightKg,
      if (microchipId != null) 'microchip_id': microchipId,
      if (avatarUrl != null) 'avatar_url': avatarUrl,
    };
  }
}
