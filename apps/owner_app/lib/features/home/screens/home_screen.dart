import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../owner_profile/controllers/owner_profile_controller.dart';
import '../../pet/controllers/pet_controller.dart';
import '../../pet/models/pet_model.dart';
import '../widgets/pet_avatar_widget.dart';
import '../widgets/status_bar_widget.dart';
import '../widgets/mini_calendar_strip.dart';
import '../widgets/activity_list_tile.dart';
import '../controllers/home_controller.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<HomeController>().init();
      context.read<PetController>().fetchMyPets();
    });
  }

  @override
  Widget build(BuildContext context) {
    final ownerCtrl = context.watch<OwnerProfileController>();
    final petCtrl = context.watch<PetController>();
    final homeCtrl = context.watch<HomeController>();
    final pets = petCtrl.myPets;
    final pet = pets.isNotEmpty ? pets[homeCtrl.selectedPetIndex] : null;

    return Scaffold(
      backgroundColor: AppColors.primaryLight,
      body: SafeArea(
        child: Column(
          children: [
            // ── Header ──
            _buildHeader(context, ownerCtrl),
            // ── Pet Avatar Section ──
            Expanded(
              flex: 5,
              child: _buildAvatarSection(context, pet, homeCtrl, pets.length),
            ),
            // ── Bottom Card ──
            Expanded(
              flex: 6,
              child: _buildBottomCard(context, homeCtrl, pet),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader(BuildContext context, OwnerProfileController ctrl) {
    final name = ctrl.profile?.firstName ?? 'เจ้าของ';
    return Padding(
      padding: EdgeInsets.symmetric(
          horizontal: context.nw(20), vertical: context.nh(12)),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            'สวัสดี, $name 👋',
            style: AppTextStyles.body(context).copyWith(
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          _NotificationBell(),
        ],
      ),
    );
  }

  Widget _buildAvatarSection(
      BuildContext context, PetModel? pet, HomeController ctrl, int petCount) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        PetAvatarWidget(pet: pet),
        if (pet != null) ...[
          SizedBox(height: context.nh(8)),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                padding: EdgeInsets.symmetric(
                    horizontal: context.nw(16), vertical: context.nh(6)),
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.8),
                  borderRadius: BorderRadius.circular(context.radius(20)),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.pets,
                        color: AppColors.primary, size: context.icon(16)),
                    SizedBox(width: context.nw(6)),
                    Text(
                      pet.name,
                      style: AppTextStyles.body(context).copyWith(
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ],
        SizedBox(height: context.nh(8)),
        // Page indicator dots
        if (petCount > 1)
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: List.generate(petCount, (i) {
              return AnimatedContainer(
                duration: const Duration(milliseconds: 200),
                margin: EdgeInsets.symmetric(horizontal: context.nw(3)),
                width: i == ctrl.selectedPetIndex
                    ? context.nw(20)
                    : context.nw(6),
                height: context.nw(6),
                decoration: BoxDecoration(
                  color: i == ctrl.selectedPetIndex
                      ? AppColors.primary
                      : AppColors.primary.withValues(alpha: 0.3),
                  borderRadius: BorderRadius.circular(context.radius(3)),
                ),
              );
            }),
          ),
        SizedBox(height: context.nh(8)),
        // Mood text
        Container(
          margin: EdgeInsets.symmetric(horizontal: context.nw(32)),
          padding: EdgeInsets.symmetric(
              horizontal: context.nw(16), vertical: context.nh(8)),
          decoration: BoxDecoration(
            color: Colors.white.withValues(alpha: 0.7),
            borderRadius: BorderRadius.circular(context.radius(12)),
          ),
          child: Text(
            pet != null
                ? '${pet.name} ดีใจที่เห็นคุณมาก! 🐾'
                : 'เพิ่มสัตว์เลี้ยงเพื่อเริ่มใช้งาน',
            textAlign: TextAlign.center,
            style: AppTextStyles.caption(context).copyWith(
              color: AppColors.textSecondary,
              fontSize: context.nf(13),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildBottomCard(
      BuildContext context, HomeController ctrl, PetModel? pet) {
    return Container(
      width: double.infinity,
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.only(
          topLeft: Radius.circular(context.radius(32)),
          topRight: Radius.circular(context.radius(32)),
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.06),
            blurRadius: 12,
            offset: const Offset(0, -3),
          ),
        ],
      ),
      child: Column(
        children: [
          SizedBox(height: context.nh(16)),
          // Status bar icons
          StatusBarWidget(pet: pet),
          Divider(
              color: AppColors.border,
              height: context.nh(24),
              indent: context.nw(20),
              endIndent: context.nw(20)),
          // Mini calendar strip
          MiniCalendarStrip(
            selectedDate: ctrl.selectedDate,
            onDateSelected: ctrl.setSelectedDate,
          ),
          SizedBox(height: context.nh(8)),
          // Activities
          Expanded(
            child: ctrl.todayActivities.isEmpty
                ? Center(
                    child: Text(
                      'ไม่มีกิจกรรมในวันนี้',
                      style: AppTextStyles.caption(context),
                    ),
                  )
                : ListView.builder(
                    padding: EdgeInsets.symmetric(horizontal: context.nw(16)),
                    itemCount: ctrl.todayActivities.length,
                    itemBuilder: (ctx, i) =>
                        ActivityListTile(activity: ctrl.todayActivities[i]),
                  ),
          ),
        ],
      ),
    );
  }
}

class _NotificationBell extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Icon(Icons.notifications_outlined,
            color: AppColors.textPrimary, size: context.icon(26)),
        Positioned(
          right: 0,
          top: 0,
          child: Container(
            width: context.nw(8),
            height: context.nw(8),
            decoration: const BoxDecoration(
              color: Colors.red,
              shape: BoxShape.circle,
            ),
          ),
        ),
      ],
    );
  }
}
