import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../app/app_routes.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_card.dart';
import '../../../shared/widgets/app_dialog.dart';
import '../../../shared/widgets/app_scaffold.dart';
import '../../auth/controllers/auth_controller.dart';
import '../controllers/owner_profile_controller.dart';
import '../widgets/section_title.dart';
import '../widgets/profile_avatar.dart';
import '../widgets/profile_text_field.dart';
import '../widgets/gender_dropdown.dart';
import '../widgets/province_dropdown.dart';
import '../widgets/address_field.dart';

class OwnerProfileScreen extends StatefulWidget {
  const OwnerProfileScreen({super.key});

  @override
  State<OwnerProfileScreen> createState() => _OwnerProfileScreenState();
}

class _OwnerProfileScreenState extends State<OwnerProfileScreen> {
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();
  final _ageController = TextEditingController();
  final _phoneController = TextEditingController();
  final _addressController = TextEditingController();

  String? _selectedGender;
  String? _selectedProvince;
  String? _mockAvatarUrl;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<OwnerProfileController>().fetchProfile();
    });
  }

  @override
  void dispose() {
    _firstNameController.dispose();
    _lastNameController.dispose();
    _ageController.dispose();
    _phoneController.dispose();
    _addressController.dispose();
    super.dispose();
  }

  // Calculate age from date_of_birth (YYYY-MM-DD)
  int _calculateAge(String? dobString) {
    if (dobString == null) return 0;
    try {
      final parts = dobString.split('-');
      if (parts.length != 3) return 0;
      final year = int.parse(parts[0]);
      return DateTime.now().year - year;
    } catch (_) {
      return 0;
    }
  }

  // Opens a beautiful dialog to pick a mock profile image
  void _pickMockAvatar() {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(
            'เลือกรูปโปรไฟล์ (Mock)',
            style: AppTextStyles.title(context),
          ),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  _mockAvatarOption(
                      'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150'),
                  _mockAvatarOption(
                      'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150'),
                  _mockAvatarOption(
                      'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150'),
                ],
              ),
              if (_mockAvatarUrl != null) ...[
                const SizedBox(height: 16),
                TextButton(
                  onPressed: () {
                    setState(() => _mockAvatarUrl = null);
                    Navigator.pop(context);
                  },
                  child: const Text('ล้างรูปโปรไฟล์', style: TextStyle(color: Colors.red)),
                )
              ]
            ],
          ),
        );
      },
    );
  }

  Widget _mockAvatarOption(String url) {
    return GestureDetector(
      onTap: () {
        setState(() => _mockAvatarUrl = url);
        Navigator.pop(context);
      },
      child: CircleAvatar(
        radius: 35,
        backgroundImage: NetworkImage(url),
      ),
    );
  }

  Future<void> _submitProfile() async {
    final firstName = _firstNameController.text.trim();
    final lastName = _lastNameController.text.trim();
    final ageText = _ageController.text.trim();
    final phone = _phoneController.text.trim();
    final address = _addressController.text.trim();

    if (firstName.isEmpty || lastName.isEmpty || ageText.isEmpty || phone.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('กรุณากรอกข้อมูลที่จำเป็น (*) ให้ครบถ้วน')),
      );
      return;
    }

    final age = int.tryParse(ageText);
    if (age == null || age <= 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('กรุณากรอกอายุให้ถูกต้อง')),
      );
      return;
    }

    String? mappedGender;
    if (_selectedGender == 'ชาย') {
      mappedGender = 'male';
    } else if (_selectedGender == 'หญิง') {
      mappedGender = 'female';
    } else if (_selectedGender == 'ไม่ระบุเพศ') {
      mappedGender = 'prefer_not_to_say';
    }

    final success = await context.read<OwnerProfileController>().createProfile(
          firstName: firstName,
          lastName: lastName,
          phoneNumber: phone,
          age: age,
          gender: mappedGender,
          avatarUrl: _mockAvatarUrl,
          address: address.isNotEmpty ? address : null,
          province: _selectedProvince != 'เลือกจังหวัด' ? _selectedProvince : null,
        );

    if (!mounted) return;

    if (success) {
      AppDialog.showMessage(
        context: context,
        title: 'สำเร็จ',
        message: 'ตั้งค่าโปรไฟล์เรียบร้อยแล้ว!',
      );
    } else {
      final error = context.read<OwnerProfileController>().errorMessage;
      AppDialog.showMessage(
        context: context,
        title: 'เกิดข้อผิดพลาด',
        message: error ?? 'ไม่สามารถบันทึกโปรไฟล์ได้',
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final profileController = context.watch<OwnerProfileController>();
    final authController = context.watch<AuthController>();

    if (profileController.state == OwnerProfileState.loading) {
      return const Scaffold(
        body: Center(
          child: CircularProgressIndicator(color: AppColors.primary),
        ),
      );
    }

    // Mode 1: Home Dashboard Mode (If profile already exists)
    if (profileController.profile != null) {
      final p = profileController.profile!;
      final calculatedAge = _calculateAge(p.dateOfBirth);
      String displayGender = 'ไม่ระบุ';
      if (p.gender == 'male') displayGender = 'ชาย';
      if (p.gender == 'female') displayGender = 'หญิง';
      if (p.gender == 'prefer_not_to_say') displayGender = 'ไม่ระบุเพศ';

      return Scaffold(
        appBar: AppBar(
          title: RichText(
            text: const TextSpan(
              style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
              children: [
                TextSpan(text: 'Pet', style: TextStyle(color: AppColors.textPrimary)),
                TextSpan(text: 'Nexus', style: TextStyle(color: AppColors.primary)),
              ],
            ),
          ),
          centerTitle: true,
          backgroundColor: Colors.transparent,
          elevation: 0,
          actions: [
            IconButton(
              icon: const Icon(Icons.logout, color: AppColors.error),
              onPressed: () async {
                final navigator = Navigator.of(context);
                final confirm = await AppDialog.showConfirm(
                  context: context,
                  title: 'ออกจากระบบ',
                  message: 'คุณต้องการออกจากระบบใช่หรือไม่?',
                );
                if (confirm) {
                  profileController.clearProfile();
                  await authController.logout();
                  navigator.pushReplacementNamed(AppRoutes.first);
                }
              },
            )
          ],
        ),
        body: SafeArea(
          child: SingleChildScrollView(
            padding: EdgeInsets.all(context.nw(24)),
            child: Column(
              children: [
                Center(
                  child: CircleAvatar(
                    radius: context.nw(55),
                    backgroundColor: AppColors.primaryLight,
                    backgroundImage: p.avatarUrl != null ? NetworkImage(p.avatarUrl!) : null,
                    child: p.avatarUrl == null
                        ? Icon(Icons.person, size: context.icon(55), color: AppColors.primary)
                        : null,
                  ),
                ),
                SizedBox(height: context.nh(16)),
                Text(
                  p.displayName,
                  style: AppTextStyles.title(context).copyWith(fontSize: context.nf(24)),
                ),
                SizedBox(height: context.nh(24)),
                AppCard(
                  child: Column(
                    children: [
                      _profileInfoRow(context, 'เบอร์โทรศัพท์', p.phoneNumber),
                      _profileDivider(context),
                      _profileInfoRow(context, 'เพศ', displayGender),
                      _profileDivider(context),
                      _profileInfoRow(context, 'อายุ', calculatedAge > 0 ? '$calculatedAge ปี' : 'ไม่ระบุ'),
                      _profileDivider(context),
                      _profileInfoRow(context, 'ที่อยู่', p.addressLine1 ?? 'ไม่ระบุ'),
                      _profileDivider(context),
                      _profileInfoRow(context, 'จังหวัด', p.province ?? 'ไม่ระบุ'),
                    ],
                  ),
                ),
                SizedBox(height: context.nh(32)),
                Text(
                  'ยินดีต้อนรับสู่ PetNexus!\nโปรไฟล์ของคุณตั้งค่าเสร็จสมบูรณ์แล้ว',
                  textAlign: TextAlign.center,
                  style: AppTextStyles.body(context).copyWith(color: AppColors.textSecondary),
                )
              ],
            ),
          ),
        ),
      );
    }

    // Mode 2: Complete Profile Setup Mode (Form Mode)
    final formLayoutPadding = EdgeInsets.fromLTRB(
      context.nw(32),
      context.nh(40),
      context.nw(32),
      context.nh(40),
    );

    return AppScaffold(
      scrollable: false,
      backgroundColor: AppColors.primaryLight,
      child: Column(
        children: [
          SizedBox(
            height: context.nh(90),
            child: Padding(
              padding: EdgeInsets.symmetric(horizontal: context.nw(16)),
              child: Align(
                alignment: Alignment.centerLeft,
                child: SectionTitle(
                  title: 'ตั้งค่าโปรไฟล์',
                  onBack: () => Navigator.maybePop(context),
                ),
              ),
            ),
          ),
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
                padding: formLayoutPadding,
                child: Column(
                  children: [
                    Center(
                      child: ProfileAvatar(
                        imageFile: null, // Image selection is mocked via URLs
                        onTap: _pickMockAvatar,
                      ),
                    ),
                    if (_mockAvatarUrl != null) ...[
                      SizedBox(height: context.nh(8)),
                      CircleAvatar(
                        radius: 24,
                        backgroundImage: NetworkImage(_mockAvatarUrl!),
                      ),
                    ],
                    SizedBox(height: context.nh(24)),
                    ProfileTextField(
                      controller: _firstNameController,
                      label: 'ชื่อจริง',
                      hintText: 'กรอกข้อมูล',
                      isRequired: true,
                    ),
                    SizedBox(height: context.nh(14)),
                    ProfileTextField(
                      controller: _lastNameController,
                      label: 'นามสกุล',
                      hintText: 'กรอกข้อมูล',
                      isRequired: true,
                    ),
                    SizedBox(height: context.nh(14)),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(
                          flex: 5,
                          child: GenderDropdown(
                            value: _selectedGender,
                            onChanged: (val) {
                              setState(() => _selectedGender = val);
                            },
                          ),
                        ),
                        SizedBox(width: context.nw(16)),
                        Expanded(
                          flex: 5,
                          child: ProfileTextField(
                            controller: _ageController,
                            label: 'อายุ',
                            hintText: 'กรอกข้อมูล',
                            keyboardType: TextInputType.number,
                            isRequired: true,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: context.nh(14)),
                    ProfileTextField(
                      controller: _phoneController,
                      label: 'เบอร์โทรศัพท์',
                      hintText: 'กรอกข้อมูล',
                      keyboardType: TextInputType.phone,
                      isRequired: true,
                      prefix: Padding(
                        padding: EdgeInsets.only(right: context.nw(8)),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text('🇹🇭', style: TextStyle(fontSize: context.nf(18))),
                            SizedBox(width: context.nw(4)),
                            Text('+66', style: AppTextStyles.body(context)),
                            SizedBox(width: context.nw(8)),
                            Container(
                              width: 1,
                              height: context.nh(16),
                              color: AppColors.border,
                            )
                          ],
                        ),
                      ),
                    ),
                    SizedBox(height: context.nh(14)),
                    AddressField(controller: _addressController),
                    SizedBox(height: context.nh(14)),
                    ProvinceDropdown(
                      value: _selectedProvince,
                      onChanged: (val) {
                        setState(() => _selectedProvince = val);
                      },
                    ),
                    SizedBox(height: context.nh(32)),
                    AppButton.primary(
                      text: 'ถัดไป',
                      icon: Icons.pets,
                      loading: profileController.state == OwnerProfileState.loading,
                      onPressed: _submitProfile,
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

  Widget _profileInfoRow(BuildContext context, String label, String value) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: context.nh(8)),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: AppTextStyles.body(context).copyWith(color: AppColors.textSecondary),
          ),
          Text(
            value,
            style: AppTextStyles.body(context).copyWith(fontWeight: FontWeight.bold),
          ),
        ],
      ),
    );
  }

  Widget _profileDivider(BuildContext context) {
    return Divider(height: context.nh(16), color: AppColors.border);
  }
}