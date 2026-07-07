class BreedModel {
  final String id;
  final String species;
  final String name;
  final String? nameTh;

  const BreedModel({
    required this.id,
    required this.species,
    required this.name,
    this.nameTh,
  });

  factory BreedModel.fromJson(Map<String, dynamic> json) {
    return BreedModel(
      id: json['id'] as String? ?? '',
      species: json['species'] as String? ?? 'dog',
      name: json['name'] as String? ?? '',
      nameTh: json['name_th'] as String?,
    );
  }

  String get displayName => nameTh ?? name;
}
