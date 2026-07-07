import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../app/app_routes.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_card.dart';
import '../../../shared/widgets/app_scaffold.dart';
import '../controllers/pet_controller.dart';
import '../widgets/step_tracker.dart';

class SelectPetScreen extends StatefulWidget {
  const SelectPetScreen({super.key});

  @override
  State<SelectPetScreen> createState() => _SelectPetScreenState();
}

class _SelectPetScreenState extends State<SelectPetScreen> {
  final PageController _pageController = PageController();
  int _currentPage = 0; // 0 for Dog, 1 for Cat

  // Mock illustration URLs (User will replace with static local assets later)
  final String _dogImageUrl = 'https://images.unsplash.com/photo-1543466835-00a7907e9de1?w=300';
  final String _catImageUrl = 'https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=300';

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  void _onPageChanged(int index) {
    setState(() {
      _currentPage = index;
    });
    // Set species in controller ('dog' or 'cat')
    context.read<PetController>().setSpecies(index == 0 ? 'dog' : 'cat');
  }

  void _nextPage() {
    Navigator.pushNamed(context, AppRoutes.petInfoForm);
  }

  void _skip() {
    Navigator.pushNamedAndRemoveUntil(context, AppRoutes.home, (route) => false);
  }

  @override
  Widget build(BuildContext context) {
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
                  const StepTracker(currentStep: 1),
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
                    Text(
                      'เลือกสัตว์เลี้ยง',
                      style: AppTextStyles.title(context).copyWith(
                        fontSize: context.nf(24),
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    SizedBox(height: context.nh(24)),
                    // Swipe Carousel Area with indicator arrows
                    Expanded(
                      child: Stack(
                        alignment: Alignment.center,
                        children: [
                          PageView(
                            controller: _pageController,
                            onPageChanged: _onPageChanged,
                            children: [
                              _buildPetCarouselItem(_dogImageUrl),
                              _buildPetCarouselItem(_catImageUrl),
                            ],
                          ),
                          // Arrow indicators on sides
                          if (_currentPage == 0)
                            Positioned(
                              right: 0,
                              child: GestureDetector(
                                onTap: () => _pageController.nextPage(
                                  duration: const Duration(milliseconds: 300),
                                  curve: Curves.easeIn,
                                ),
                                child: Icon(
                                  Icons.play_arrow,
                                  color: AppColors.primary,
                                  size: context.icon(36),
                                ),
                              ),
                            ),
                          if (_currentPage == 1)
                            Positioned(
                              left: 0,
                              child: GestureDetector(
                                onTap: () => _pageController.previousPage(
                                  duration: const Duration(milliseconds: 300),
                                  curve: Curves.easeIn,
                                ),
                                child: Transform.rotate(
                                  angle: 3.141592653589793,
                                  child: Icon(
                                    Icons.play_arrow,
                                    color: AppColors.primary,
                                    size: context.icon(36),
                                  ),
                                ),
                              ),
                            ),
                        ],
                      ),
                    ),
                    SizedBox(height: context.nh(8)),
                    // Species Text Indicator (e.g., สุนัข | แมว)
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          'สุนัข',
                          style: AppTextStyles.body(context).copyWith(
                            fontWeight: _currentPage == 0 ? FontWeight.bold : FontWeight.normal,
                            color: _currentPage == 0 ? AppColors.textPrimary : AppColors.textSecondary.withValues(alpha: 0.5),
                            fontSize: context.nf(18),
                          ),
                        ),
                        SizedBox(width: context.nw(24)),
                        Text(
                          'แมว',
                          style: AppTextStyles.body(context).copyWith(
                            fontWeight: _currentPage == 1 ? FontWeight.bold : FontWeight.normal,
                            color: _currentPage == 1 ? AppColors.textPrimary : AppColors.textSecondary.withValues(alpha: 0.5),
                            fontSize: context.nf(18),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: context.nh(12)),
                    // Dot Page Indicator
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        _buildDotIndicator(_currentPage == 0),
                        SizedBox(width: context.nw(8)),
                        _buildDotIndicator(_currentPage == 1),
                      ],
                    ),
                    SizedBox(height: context.nh(32)),
                    // Info Security Banner
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
                      onPressed: _nextPage,
                    ),
                    SizedBox(height: context.nh(12)),
                    // Skip Button Link
                    TextButton(
                      onPressed: _skip,
                      child: Text(
                        'ข้ามไปก่อน',
                        style: AppTextStyles.body(context).copyWith(
                          color: AppColors.primary,
                          fontWeight: FontWeight.bold,
                          fontSize: context.nf(16),
                        ),
                      ),
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

  Widget _buildPetCarouselItem(String url) {
    return Center(
      child: Container(
        width: context.nw(220),
        height: context.nw(220),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(20),
          image: DecorationImage(
            image: NetworkImage(url),
            fit: BoxFit.cover,
          ),
        ),
      ),
    );
  }

  Widget _buildDotIndicator(bool isActive) {
    return Container(
      width: context.nw(8),
      height: context.nw(8),
      decoration: BoxDecoration(
        color: isActive ? AppColors.textPrimary : AppColors.border,
        shape: BoxShape.circle,
      ),
    );
  }
}
