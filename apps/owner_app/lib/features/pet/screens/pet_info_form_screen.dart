import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../app/app_routes.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_card.dart';
import '../../../shared/widgets/app_scaffold.dart';
import '../../owner_profile/widgets/profile_text_field.dart';
import '../controllers/pet_controller.dart';
import '../models/breed_model.dart';
import '../widgets/step_tracker.dart';

class PetInfoFormScreen extends StatefulWidget {
  const PetInfoFormScreen({super.key});

  @override
  State<PetInfoFormScreen> createState() => _PetInfoFormScreenState();
}

class _PetInfoFormScreenState extends State<PetInfoFormScreen> {
  final _nameController = TextEditingController();
  final _dobController = TextEditingController();
  final _ageController = TextEditingController();
  final _weightController = TextEditingController();
  final _microchipController = TextEditingController();

  String? _selectedGender; // 'male' or 'female'
  BreedModel? _selectedBreed;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final controller = context.read<PetController>();
      controller.fetchBreeds(controller.selectedSpecies);
    });
  }

  @override
  void dispose() {
    _nameController.dispose();
    _dobController.dispose();
    _ageController.dispose();
    _weightController.dispose();
    _microchipController.dispose();
    super.dispose();
  }

  void _nextStep(String species) {
    final name = _nameController.text.trim();
    final ageText = _ageController.text.trim();
    final dob = _dobController.text.trim();
    final weightText = _weightController.text.trim();
    final microchip = _microchipController.text.trim();

    if (name.isEmpty) {
      _showSnackBar('กรุณากรอกชื่อเล่นสัตว์เลี้ยง');
      return;
    }
    if (_selectedGender == null) {
      _showSnackBar('กรุณาเลือกเพศ');
      return;
    }
    if (ageText.isEmpty) {
      _showSnackBar('กรุณากรอกอายุ');
      return;
    }
    if (_selectedBreed == null) {
      _showSnackBar('กรุณาเลือกสายพันธุ์');
      return;
    }

    final age = int.tryParse(ageText);
    if (age == null || age < 0) {
      _showSnackBar('กรุณากรอกอายุให้ถูกต้อง');
      return;
    }

    final weight = double.tryParse(weightText);

    // Save temporary state to controller
    context.read<PetController>().setTemporaryInfo(
          name: name,
          gender: _selectedGender,
          dateOfBirth: dob.isNotEmpty ? dob : null,
          age: age,
          breed: _selectedBreed,
          weightKg: weight,
          microchipId: microchip.isNotEmpty ? microchip : null,
        );

    // Navigate to step 3
    Navigator.pushNamed(context, AppRoutes.petUploadPhoto);
  }

  Future<void> _selectDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: DateTime.now(),
      firstDate: DateTime(1900),
      lastDate: DateTime.now(),
      builder: (context, child) {
        return Theme(
          data: Theme.of(context).copyWith(
            colorScheme: const ColorScheme.light(
              primary: AppColors.primary,
              onPrimary: Colors.white,
              onSurface: AppColors.textPrimary,
            ),
          ),
          child: child!,
        );
      },
    );
    if (picked != null) {
      setState(() {
        _dobController.text = "${picked.year}-${picked.month.toString().padLeft(2, '0')}-${picked.day.toString().padLeft(2, '0')}";
      });
    }
  }

  void _showSnackBar(String msg) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(msg)));
  }

  @override
  Widget build(BuildContext context) {
    final petController = context.watch<PetController>();
    final species = petController.selectedSpecies;
    final isDog = species == 'dog';
    final breedsList = petController.breeds;

    return AppScaffold(
      scrollable: false,
      backgroundColor: AppColors.primaryLight,
      child: Column(
        children: [
          // Step Header
          SizedBox(
            height: context.nh(90),
            child: Padding(
              padding: EdgeInsets.symmetric(horizontal: context.nw(16)),
              child: Stack(
                alignment: Alignment.center,
                children: [
                  Positioned(
                    left: 0,
                    child: GestureDetector(
                      onTap: () => Navigator.pop(context),
                      child: CircleAvatar(
                        radius: context.nw(20),
                        backgroundColor: Colors.white,
                        child: Icon(
                          Icons.chevron_left,
                          color: AppColors.primary,
                          size: context.icon(28),
                        ),
                      ),
                    ),
                  ),
                  const StepTracker(currentStep: 2),
                ],
              ),
            ),
          ),
          // Content Card
          Expanded(
            child: AppCard(
              padding: EdgeInsets.zero,
              color: AppColors.background,
              borderRadius: BorderRadius.only(
                topLeft: Radius.elliptical(
                  context.nw(250),
                  context.nh(52),
                ),
                topRight: Radius.elliptical(
                  context.nw(250),
                  context.nh(52),
                ),
              ),
              child: SingleChildScrollView(
                padding: EdgeInsets.fromLTRB(
                  context.nw(32),
                  context.nh(40),
                  context.nw(32),
                  context.nh(40),
                ),
                child: Column(
                  children: [
                    // Pet circular avatar/icon
                    CircleAvatar(
                      radius: context.nw(44),
                      backgroundColor: Colors.grey.shade200,
                      child: Icon(
                        Icons.pets, // Using Pets icon as standard for both dog and cat
                        size: context.icon(40),
                        color: Colors.grey.shade600,
                      ),
                    ),
                    SizedBox(height: context.nh(12)),
                    Text(
                      isDog ? 'กรอกข้อมูลสุนัข' : 'กรอกข้อมูลแมว',
                      style: AppTextStyles.title(context).copyWith(
                        fontSize: context.nf(22),
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    SizedBox(height: context.nh(24)),

                    // Nickname field
                    ProfileTextField(
                      controller: _nameController,
                      label: 'ชื่อเล่น',
                      hintText: 'กรอกข้อมูล',
                      isRequired: true,
                    ),
                    SizedBox(height: context.nh(14)),

                    // Gender field (Toggle buttons ผู้ / เมีย)
                    _buildGenderToggleField(context),
                    SizedBox(height: context.nh(14)),

                    // DOB Field (DatePicker mapped to ProfileTextField style)
                    ProfileTextField(
                      controller: _dobController,
                      label: 'วันเกิด',
                      hintText: 'กรอกข้อมูล',
                      prefixIcon: Icons.calendar_month,
                      readOnly: true,
                      onTap: () => _selectDate(context),
                    ),
                    SizedBox(height: context.nh(14)),

                    // Age field
                    ProfileTextField(
                      controller: _ageController,
                      label: 'อายุ',
                      hintText: 'กรอกข้อมูล',
                      keyboardType: TextInputType.number,
                      isRequired: true,
                    ),
                    SizedBox(height: context.nh(14)),

                    // Breed dropdown field
                    _buildBreedDropdownField(context, breedsList),
                    SizedBox(height: context.nh(14)),

                    // Weight field (With kg suffix unit)
                    _buildWeightField(context),
                    SizedBox(height: context.nh(14)),

                    // Microchip field
                    ProfileTextField(
                      controller: _microchipController,
                      label: 'รหัสไมโครชิป',
                      hintText: 'กรอกข้อมูล',
                    ),
                    SizedBox(height: context.nh(24)),

                    // Banner Security
                    Container(
                      padding: EdgeInsets.symmetric(
                        horizontal: context.nw(16),
                        vertical: context.nh(12),
                      ),
                      decoration: BoxDecoration(
                        color: AppColors.primary.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(context.radius(16)),
                      ),
                      child: Row(
                        children: [
                          Icon(
                            Icons.shield_outlined,
                            color: AppColors.primary,
                            size: context.icon(28),
                          ),
                          SizedBox(width: context.nw(12)),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'ข้อมูลของน้องจะปลอดภัย',
                                  style: AppTextStyles.body(context).copyWith(
                                    fontWeight: FontWeight.bold,
                                    fontSize: context.nf(14),
                                    color: AppColors.primary,
                                  ),
                                ),
                                SizedBox(height: context.nh(2)),
                                Text(
                                  'เราปกป้องข้อมูลของสัตว์เลี้ยงตามมาตรฐานความปลอดภัยระดับสากล',
                                  style: AppTextStyles.caption(context).copyWith(
                                    fontSize: context.nf(11),
                                    color: AppColors.textSecondary,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                    SizedBox(height: context.nh(28)),

                    // Next Button
                    AppButton.primary(
                      text: 'ถัดไป',
                      icon: Icons.pets,
                      onPressed: () => _nextStep(species),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildGenderToggleField(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        RichText(
          text: TextSpan(
            text: 'เพศ',
            style: AppTextStyles.body(context).copyWith(
              fontWeight: FontWeight.w600,
            ),
            children: const [
              TextSpan(
                text: ' *',
                style: TextStyle(
                  color: AppColors.error,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ),
        SizedBox(height: context.nh(8)),
        Row(
          children: [
            // Male Button
            Expanded(
              child: _buildGenderButton(
                context,
                label: '♂ ผู้',
                isActive: _selectedGender == 'male',
                onTap: () => setState(() => _selectedGender = 'male'),
              ),
            ),
            SizedBox(width: context.nw(16)),
            // Female Button
            Expanded(
              child: _buildGenderButton(
                context,
                label: '♀ เมีย',
                isActive: _selectedGender == 'female',
                onTap: () => setState(() => _selectedGender = 'female'),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildGenderButton(
    BuildContext context, {
    required String label,
    required bool isActive,
    required VoidCallback onTap,
  }) {
    final height = context.nh(48).clamp(44.0, 52.0).toDouble();

    return GestureDetector(
      onTap: onTap,
      child: Container(
        height: height,
        decoration: BoxDecoration(
          color: isActive ? AppColors.primary.withValues(alpha: 0.12) : Colors.white,
          borderRadius: BorderRadius.circular(height / 2),
          border: Border.all(
            color: isActive ? AppColors.primary : const Color.fromARGB(255, 0, 0, 0),
            width: 1.5,
          ),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withValues(alpha: 0.05),
              blurRadius: 8,
              offset: const Offset(0, 3),
            ),
          ],
        ),
        child: Center(
          child: Text(
            label,
            style: AppTextStyles.body(context).copyWith(
              fontWeight: FontWeight.bold,
              color: isActive ? AppColors.primary : AppColors.textPrimary,
              fontSize: context.nf(16),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildBreedDropdownField(BuildContext context, List<BreedModel> breeds) {
    final height = context.nh(56).clamp(52.0, 60.0).toDouble();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        RichText(
          text: TextSpan(
            text: 'สายพันธุ์',
            style: AppTextStyles.body(context).copyWith(
              fontWeight: FontWeight.w600,
            ),
            children: const [
              TextSpan(
                text: ' *',
                style: TextStyle(
                  color: AppColors.error,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ),
        SizedBox(height: context.nh(8)),
        Container(
          height: height,
          padding: EdgeInsets.symmetric(horizontal: context.nw(18)),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(context.radius(18)),
            border: Border.all(
              color: const Color.fromARGB(255, 0, 0, 0),
              width: 1.5,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.05),
                blurRadius: 8,
                offset: const Offset(0, 3),
              ),
            ],
          ),
          child: DropdownButtonHideUnderline(
            child: DropdownButtonFormField<BreedModel>(
              initialValue: _selectedBreed,
              hint: Text(
                'กรอกข้อมูล',
                style: AppTextStyles.hint(context),
              ),
              decoration: const InputDecoration(border: InputBorder.none),
              items: breeds.map((breed) {
                return DropdownMenuItem<BreedModel>(
                  value: breed,
                  child: Text(breed.displayName, style: AppTextStyles.body(context)),
                );
              }).toList(),
              onChanged: (val) {
                setState(() => _selectedBreed = val);
              },
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildWeightField(BuildContext context) {
    final height = context.nh(56).clamp(52.0, 60.0).toDouble();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'น้ำหนัก',
          style: AppTextStyles.body(context).copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        SizedBox(height: context.nh(8)),
        Container(
          height: height,
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(context.radius(18)),
            border: Border.all(
              color: const Color.fromARGB(255, 0, 0, 0),
              width: 1.5,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.05),
                blurRadius: 8,
                offset: const Offset(0, 3),
              ),
            ],
          ),
          child: TextField(
            controller: _weightController,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            textAlignVertical: TextAlignVertical.center,
            decoration: InputDecoration(
              hintText: 'กรอกข้อมูล',
              hintStyle: AppTextStyles.hint(context),
              border: InputBorder.none,
              contentPadding: EdgeInsets.symmetric(horizontal: context.nw(18)),
              suffixIcon: Padding(
                padding: EdgeInsets.symmetric(horizontal: context.nw(18)),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    Text(
                      'kg',
                      style: AppTextStyles.body(context).copyWith(
                        color: AppColors.textSecondary,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}
