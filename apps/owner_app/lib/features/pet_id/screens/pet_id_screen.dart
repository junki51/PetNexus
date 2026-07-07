import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:qr_flutter/qr_flutter.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../pet/controllers/pet_controller.dart';
import '../../pet/models/pet_model.dart';

class PetIdScreen extends StatefulWidget {
  const PetIdScreen({super.key});

  @override
  State<PetIdScreen> createState() => _PetIdScreenState();
}

class _PetIdScreenState extends State<PetIdScreen> {
  int _currentPage = 0;
  final _pageController = PageController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<PetController>().fetchMyPets();
    });
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final petCtrl = context.watch<PetController>();
    final pets = petCtrl.myPets;

    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: Column(
          children: [
            // Header
            Padding(
              padding: EdgeInsets.symmetric(
                  horizontal: context.nw(20), vertical: context.nh(16)),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('Pet ID Card',
                      style: AppTextStyles.heading(context)
                          .copyWith(color: AppColors.textPrimary)),
                  _NotificationBell(),
                ],
              ),
            ),

            if (pets.isEmpty)
              Expanded(child: _buildEmpty(context))
            else ...[
              // PageView of ID Cards
              Expanded(
                child: PageView.builder(
                  controller: _pageController,
                  itemCount: pets.length,
                  onPageChanged: (i) => setState(() => _currentPage = i),
                  itemBuilder: (ctx, i) => Padding(
                    padding: EdgeInsets.symmetric(horizontal: context.nw(20)),
                    child: _IdCard(pet: pets[i]),
                  ),
                ),
              ),

              // Page dots
              if (pets.length > 1) ...[
                SizedBox(height: context.nh(8)),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: List.generate(pets.length, (i) {
                    return AnimatedContainer(
                      duration: const Duration(milliseconds: 200),
                      margin: EdgeInsets.symmetric(horizontal: context.nw(3)),
                      width: i == _currentPage ? context.nw(20) : context.nw(6),
                      height: context.nw(6),
                      decoration: BoxDecoration(
                        color: i == _currentPage
                            ? AppColors.primary
                            : AppColors.primary.withValues(alpha: 0.3),
                        borderRadius: BorderRadius.circular(context.radius(3)),
                      ),
                    );
                  }),
                ),
              ],

              // Bottom actions
              SizedBox(height: context.nh(16)),
              _buildActions(context, pets.isNotEmpty ? pets[_currentPage] : null),
              SizedBox(height: context.nh(20)),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildEmpty(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.pets, size: context.icon(64), color: AppColors.border),
          SizedBox(height: context.nh(16)),
          Text('ยังไม่มีสัตว์เลี้ยง',
              style: AppTextStyles.body(context)
                  .copyWith(color: AppColors.textSecondary)),
          SizedBox(height: context.nh(8)),
          Text('เพิ่มสัตว์เลี้ยงใหม่ได้เลย!',
              style: AppTextStyles.caption(context)),
        ],
      ),
    );
  }

  Widget _buildActions(BuildContext context, PetModel? pet) {
    final actions = [
      _ActionItem(Icons.edit_rounded, 'แก้ไข', () {}),
      _ActionItem(Icons.add_circle_rounded, 'เพิ่ม',
          () => Navigator.pushNamed(context, '/select-pet')),
      _ActionItem(Icons.list_alt_rounded, 'ดูทั้งหมด', () {}),
      _ActionItem(Icons.share_rounded, 'แชร์', () {}),
    ];
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: actions.map((a) {
        return GestureDetector(
          onTap: a.onTap,
          child: Column(
            children: [
              Container(
                padding: EdgeInsets.all(context.nw(12)),
                decoration: BoxDecoration(
                  color: AppColors.primary.withValues(alpha: 0.1),
                  shape: BoxShape.circle,
                ),
                child: Icon(a.icon, color: AppColors.primary, size: context.icon(22)),
              ),
              SizedBox(height: context.nh(4)),
              Text(a.label,
                  style: TextStyle(
                      fontSize: context.nf(12),
                      color: AppColors.textSecondary)),
            ],
          ),
        );
      }).toList(),
    );
  }
}

class _ActionItem {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  const _ActionItem(this.icon, this.label, this.onTap);
}

// ───────────────────────── ID Card ─────────────────────────

class _IdCard extends StatelessWidget {
  final PetModel pet;
  const _IdCard({required this.pet});

  String get _petIdNumber {
    final shortId = pet.id.replaceAll('-', '').substring(0, 8).toUpperCase();
    final year = DateTime.now().year;
    return 'PNX-$year-$shortId';
  }

  int get _ageMonths {
    if (pet.dateOfBirth == null) return 0;
    try {
      final dob = DateTime.parse(pet.dateOfBirth!);
      final now = DateTime.now();
      return (now.year - dob.year) * 12 + (now.month - dob.month);
    } catch (_) {
      return 0;
    }
  }

  String get _ageDisplay {
    final months = _ageMonths;
    final years = months ~/ 12;
    final rem = months % 12;
    if (years > 0 && rem > 0) return '$years ปี $rem เดือน';
    if (years > 0) return '$years ปี';
    return '$months เดือน';
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: EdgeInsets.only(bottom: context.nh(8)),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(context.radius(20)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.08),
            blurRadius: 16,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Teal header
          Container(
            padding: EdgeInsets.symmetric(
                horizontal: context.nw(16), vertical: context.nh(12)),
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: [AppColors.primary, const Color(0xFF26C6DA)],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.only(
                topLeft: Radius.circular(context.radius(20)),
                topRight: Radius.circular(context.radius(20)),
              ),
            ),
            child: Row(
              children: [
                Icon(Icons.pets, color: Colors.white, size: context.icon(16)),
                SizedBox(width: context.nw(8)),
                Text('Pet ID Card',
                    style: TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: context.nf(14))),
              ],
            ),
          ),

          // Body
          Padding(
            padding: EdgeInsets.all(context.nw(16)),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Avatar
                ClipRRect(
                  borderRadius: BorderRadius.circular(context.radius(12)),
                  child: Container(
                    width: context.nw(80),
                    height: context.nw(80),
                    color: AppColors.primaryLight,
                    child: pet.avatarUrl != null
                        ? Image.network(pet.avatarUrl!, fit: BoxFit.cover,
                            errorBuilder: (ctx, err, stack) =>
                                Center(child: Text(pet.species == 'dog' ? '🐶' : '🐱',
                                    style: TextStyle(fontSize: context.nf(36)))))
                        : Center(child: Text(pet.species == 'dog' ? '🐶' : '🐱',
                            style: TextStyle(fontSize: context.nf(36)))),
                  ),
                ),
                SizedBox(width: context.nw(12)),
                // Info
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Flexible(
                            child: Text(pet.name,
                                style: AppTextStyles.heading(context).copyWith(
                                    fontSize: context.nf(18))),
                          ),
                          SizedBox(width: context.nw(4)),
                          Icon(
                            pet.gender == 'female'
                                ? Icons.female
                                : Icons.male,
                            color: pet.gender == 'female'
                                ? Colors.pink
                                : Colors.blue,
                            size: context.icon(16),
                          ),
                        ],
                      ),
                      Text(
                        pet.species == 'dog' ? 'สุนัข' : 'แมว',
                        style: AppTextStyles.caption(context).copyWith(
                            color: AppColors.textSecondary,
                            fontWeight: FontWeight.normal,
                            fontSize: context.nf(13)),
                      ),
                      Text(
                        _ageDisplay,
                        style: AppTextStyles.caption(context).copyWith(
                            color: AppColors.textSecondary,
                            fontWeight: FontWeight.normal,
                            fontSize: context.nf(13)),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),

          // Status badges
          Padding(
            padding: EdgeInsets.symmetric(horizontal: context.nw(16)),
            child: Column(
              children: [
                _StatusRow(
                  icon: Icons.favorite_rounded,
                  iconColor: Colors.red,
                  label: 'Health Status',
                  badge: 'Good',
                  badgeColor: const Color(0xFF66BB6A),
                ),
                SizedBox(height: context.nh(8)),
                _StatusRow(
                  icon: Icons.vaccines_rounded,
                  iconColor: const Color(0xFF42A5F5),
                  label: 'Vaccination',
                  badge: 'Up-to-date',
                  badgeColor: const Color(0xFF42A5F5),
                  badgeIcon: Icons.check_rounded,
                ),
              ],
            ),
          ),

          Divider(color: AppColors.border, height: context.nh(24),
              indent: context.nw(16), endIndent: context.nw(16)),

          // ID + QR
          Padding(
            padding: EdgeInsets.fromLTRB(
                context.nw(16), 0, context.nw(16), context.nh(12)),
            child: Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('ID Number',
                          style: TextStyle(
                              fontSize: context.nf(11),
                              color: AppColors.textSecondary)),
                      SizedBox(height: context.nh(4)),
                      Text(_petIdNumber,
                          style: TextStyle(
                              fontSize: context.nf(13),
                              fontWeight: FontWeight.bold,
                              color: AppColors.textPrimary,
                              letterSpacing: 0.5)),
                      SizedBox(height: context.nh(8)),
                      Container(
                        padding: EdgeInsets.symmetric(
                            horizontal: context.nw(10), vertical: context.nh(4)),
                        decoration: BoxDecoration(
                          color: const Color(0xFF66BB6A).withValues(alpha: 0.12),
                          borderRadius: BorderRadius.circular(context.radius(20)),
                          border: Border.all(
                              color: const Color(0xFF66BB6A).withValues(alpha: 0.3)),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.verified_rounded,
                                color: const Color(0xFF66BB6A), size: context.icon(14)),
                            SizedBox(width: context.nw(4)),
                            Text('Verified',
                                style: TextStyle(
                                    fontSize: context.nf(12),
                                    color: const Color(0xFF66BB6A),
                                    fontWeight: FontWeight.w600)),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                // QR Code
                QrImageView(
                  data: pet.id,
                  version: QrVersions.auto,
                  size: context.nw(90),
                  eyeStyle: QrEyeStyle(
                    eyeShape: QrEyeShape.square,
                    color: AppColors.textPrimary,
                  ),
                  dataModuleStyle: QrDataModuleStyle(
                    dataModuleShape: QrDataModuleShape.square,
                    color: AppColors.textPrimary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _StatusRow extends StatelessWidget {
  final IconData icon;
  final Color iconColor;
  final String label;
  final String badge;
  final Color badgeColor;
  final IconData? badgeIcon;

  const _StatusRow({
    required this.icon,
    required this.iconColor,
    required this.label,
    required this.badge,
    required this.badgeColor,
    this.badgeIcon,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, color: iconColor, size: context.icon(16)),
        SizedBox(width: context.nw(8)),
        Expanded(
          child: Text(label,
              style: AppTextStyles.body(context).copyWith(
                  fontSize: context.nf(13), color: AppColors.textPrimary)),
        ),
        Container(
          padding: EdgeInsets.symmetric(
              horizontal: context.nw(10), vertical: context.nh(3)),
          decoration: BoxDecoration(
            color: badgeColor.withValues(alpha: 0.12),
            borderRadius: BorderRadius.circular(context.radius(12)),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (badgeIcon != null) ...[
                Icon(badgeIcon!, color: badgeColor, size: context.icon(12)),
                SizedBox(width: context.nw(4)),
              ],
              Text(badge,
                  style: TextStyle(
                      fontSize: context.nf(12),
                      color: badgeColor,
                      fontWeight: FontWeight.w600)),
            ],
          ),
        ),
      ],
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
                color: Colors.red, shape: BoxShape.circle),
          ),
        ),
      ],
    );
  }
}
