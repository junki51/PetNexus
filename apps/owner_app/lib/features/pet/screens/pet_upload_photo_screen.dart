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
import '../controllers/pet_controller.dart';
import '../widgets/step_tracker.dart';

class PetUploadPhotoScreen extends StatefulWidget {
  const PetUploadPhotoScreen({super.key});

  @override
  State<PetUploadPhotoScreen> createState() => _PetUploadPhotoScreenState();
}

class _PetUploadPhotoScreenState extends State<PetUploadPhotoScreen> {
  String? _selectedMockUrl;

  // Mock list of premium pet avatars for the picker dialog
  final List<String> _dogAvatars = [
    'https://images.unsplash.com/photo-1543466835-00a7907e9de1?w=200',
    'https://images.unsplash.com/photo-1583511655857-d19b40a7a54e?w=200',
    'https://images.unsplash.com/photo-1537151608828-ea2b117b6b86?w=200',
  ];

  final List<String> _catAvatars = [
    'https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=200',
    'https://images.unsplash.com/photo-1495360010541-f48722b34f7d?w=200',
    'https://images.unsplash.com/photo-1533738363-b7f9aef128ce?w=200',
  ];

  void _showMockPhotoPickerDialog(bool isDog) {
    final list = isDog ? _dogAvatars : _catAvatars;

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(
            'เลือกรูปภาพสัตว์เลี้ยง (Mock)',
            style: AppTextStyles.title(context),
          ),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: list.map((url) => _buildAvatarOption(url)).toList(),
              ),
              if (_selectedMockUrl != null) ...[
                const SizedBox(height: 16),
                TextButton(
                  onPressed: () {
                    setState(() => _selectedMockUrl = null);
                    Navigator.pop(context);
                  },
                  child: const Text('ล้างรูปภาพ', style: TextStyle(color: Colors.red)),
                )
              ]
            ],
          ),
        );
      },
    );
  }

  Widget _buildAvatarOption(String url) {
    return GestureDetector(
      onTap: () {
        setState(() => _selectedMockUrl = url);
        Navigator.pop(context);
      },
      child: CircleAvatar(
        radius: 35,
        backgroundImage: NetworkImage(url),
      ),
    );
  }

  Future<void> _createDigitalId(PetController controller) async {
    controller.setAvatarUrl(_selectedMockUrl);

    final success = await controller.createPetProfile();

    if (!mounted) return;

    if (success) {
      Navigator.pushReplacementNamed(context, AppRoutes.petSuccess);
    } else {
      AppDialog.showMessage(
        context: context,
        title: 'เกิดข้อผิดพลาด',
        message: controller.errorMessage ?? 'ไม่สามารถสร้าง ID Card ได้',
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final petController = context.watch<PetController>();
    final isDog = petController.selectedSpecies == 'dog';

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
                  const StepTracker(currentStep: 3),
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
              child: Padding(
                padding: EdgeInsets.symmetric(
                  horizontal: context.nw(32),
                  vertical: context.nh(40),
                ),
                child: Column(
                  children: [
                    // Stacks of cards with placeholder graphic
                    Expanded(
                      child: Center(
                        child: GestureDetector(
                          onTap: () => _showMockPhotoPickerDialog(isDog),
                          child: Stack(
                            alignment: Alignment.center,
                            children: [
                              // Stack Card 1 (Back rotated card)
                              Transform.rotate(
                                angle: 0.12,
                                child: Container(
                                  width: context.nw(160),
                                  height: context.nw(220),
                                  decoration: BoxDecoration(
                                    color: Colors.grey.shade200,
                                    borderRadius: BorderRadius.circular(16),
                                    border: Border.all(color: Colors.grey.shade300, width: 2),
                                  ),
                                ),
                              ),
                              // Stack Card 2 (Front main card)
                              Container(
                                width: context.nw(160),
                                height: context.nw(220),
                                decoration: BoxDecoration(
                                  color: Colors.white,
                                  borderRadius: BorderRadius.circular(16),
                                  border: Border.all(color: Colors.grey.shade300, width: 2),
                                  boxShadow: [
                                    BoxShadow(
                                      color: Colors.black.withValues(alpha: 0.08),
                                      blurRadius: 12,
                                      offset: const Offset(0, 4),
                                    )
                                  ],
                                  image: _selectedMockUrl != null
                                      ? DecorationImage(
                                          image: NetworkImage(_selectedMockUrl!),
                                          fit: BoxFit.cover,
                                        )
                                      : null,
                                ),
                                child: _selectedMockUrl == null
                                    ? Center(
                                        child: Icon(
                                          Icons.image_outlined,
                                          size: context.icon(44),
                                          color: Colors.grey.shade400,
                                        ),
                                      )
                                    : null,
                              ),
                              // Upload round button overlapping at bottom right of the cards
                              Positioned(
                                bottom: context.nw(8),
                                right: context.nw(8),
                                child: CircleAvatar(
                                  radius: context.nw(24),
                                  backgroundColor: AppColors.primary,
                                  child: Icon(
                                    Icons.file_upload_outlined,
                                    color: Colors.white,
                                    size: context.icon(26),
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                    SizedBox(height: context.nh(24)),
                    Text(
                      'เพิ่มภาพสัตว์เลี้ยง',
                      style: AppTextStyles.title(context).copyWith(
                        fontSize: context.nf(22),
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    SizedBox(height: context.nh(24)),
                    // Security Banner
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
                    // Create Digital ID Card Button
                    AppButton.primary(
                      text: 'สร้าง Digital ID Card',
                      icon: Icons.pets,
                      loading: petController.state == PetState.loading,
                      onPressed: () => _createDigitalId(petController),
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
}
