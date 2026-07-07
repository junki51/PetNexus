import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../auth/controllers/auth_controller.dart';
import '../../owner_profile/controllers/owner_profile_controller.dart';
import '../widgets/settings_tile.dart';

class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final authCtrl = context.watch<AuthController>();
    final ownerCtrl = context.watch<OwnerProfileController>();
    final profile = ownerCtrl.profile;

    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: ListView(
          padding: EdgeInsets.symmetric(horizontal: context.nw(20)),
          children: [
            // Header
            Padding(
              padding: EdgeInsets.symmetric(vertical: context.nh(16)),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('Settings', style: AppTextStyles.heading(context)),
                  _NotificationBell(),
                ],
              ),
            ),

            // Profile card
            Container(
              padding: EdgeInsets.all(context.nw(16)),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(context.radius(16)),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withValues(alpha: 0.06),
                    blurRadius: 10,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: Row(
                children: [
                  // Avatar
                  CircleAvatar(
                    radius: context.nw(28),
                    backgroundColor: AppColors.primaryLight,
                    backgroundImage: profile?.avatarUrl != null
                        ? NetworkImage(profile!.avatarUrl!)
                        : null,
                    child: profile?.avatarUrl == null
                        ? Icon(Icons.person_rounded,
                            color: AppColors.primary, size: context.icon(28))
                        : null,
                  ),
                  SizedBox(width: context.nw(12)),
                  // Info
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          profile != null
                              ? '${profile.firstName} ${profile.lastName}'
                              : 'นายตัวอย่าง นามสกุล',
                          style: AppTextStyles.body(context).copyWith(
                              fontWeight: FontWeight.bold,
                              fontSize: context.nf(15)),
                        ),
                        Text(
                          authCtrl.currentUser?.email ?? 'example@gmail.com',
                          style: AppTextStyles.caption(context).copyWith(
                              color: AppColors.textSecondary,
                              fontWeight: FontWeight.normal,
                              fontSize: context.nf(13)),
                        ),
                      ],
                    ),
                  ),
                  // Edit icon
                  GestureDetector(
                    onTap: () => ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                          content: Text('แก้ไขโปรไฟล์ — Coming soon!')),
                    ),
                    child: Container(
                      padding: EdgeInsets.all(context.nw(8)),
                      decoration: BoxDecoration(
                        color: AppColors.primaryLight,
                        shape: BoxShape.circle,
                      ),
                      child: Icon(Icons.edit_rounded,
                          color: AppColors.primary, size: context.icon(16)),
                    ),
                  ),
                ],
              ),
            ),

            SizedBox(height: context.nh(24)),

            // Sharing & Authorization
            _SectionLabel('Sharing & Authorization'),
            _SettingsGroup(tiles: [
              SettingsTile(
                icon: Icons.people_rounded,
                iconColor: const Color(0xFF66BB6A),
                label: 'Family Sharing',
                trailing: _Badge('3 members'),
                onTap: () => _comingSoon(context),
              ),
              SettingsTile(
                icon: Icons.local_hospital_rounded,
                iconColor: const Color(0xFF42A5F5),
                label: 'Authorized Clinics',
                trailing: _Badge('2 clinics'),
                onTap: () => _comingSoon(context),
              ),
              SettingsTile(
                icon: Icons.pending_rounded,
                iconColor: const Color(0xFFFFA726),
                label: 'Pending Requests',
                trailing: _Badge('1 request', isAlert: true),
                onTap: () => _comingSoon(context),
              ),
            ]),

            SizedBox(height: context.nh(16)),

            // Privacy & Data
            _SectionLabel('Privacy & Data'),
            _SettingsGroup(tiles: [
              SettingsTile(
                icon: Icons.lock_rounded,
                iconColor: const Color(0xFF7E57C2),
                label: 'Privacy Settings',
                onTap: () => _comingSoon(context),
              ),
              SettingsTile(
                icon: Icons.picture_as_pdf_rounded,
                iconColor: const Color(0xFFEF5350),
                label: 'Export Data (PDF)',
                onTap: () => _comingSoon(context),
              ),
              SettingsTile(
                icon: Icons.cloud_upload_rounded,
                iconColor: const Color(0xFF26C6DA),
                label: 'Data Backup',
                onTap: () => _comingSoon(context),
              ),
            ]),

            SizedBox(height: context.nh(16)),

            // More
            _SectionLabel('More'),
            _SettingsGroup(tiles: [
              SettingsTile(
                icon: Icons.language_rounded,
                iconColor: const Color(0xFF66BB6A),
                label: 'Language',
                onTap: () => _comingSoon(context),
              ),
              SettingsTile(
                icon: Icons.help_rounded,
                iconColor: const Color(0xFFFFA726),
                label: 'Help & Support',
                onTap: () => _comingSoon(context),
              ),
              SettingsTile(
                icon: Icons.info_rounded,
                iconColor: AppColors.primary,
                label: 'About PetNexus',
                onTap: () => _comingSoon(context),
              ),
            ]),

            SizedBox(height: context.nh(24)),

            // Logout
            SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: () async {
                  await context.read<AuthController>().logout();
                  if (context.mounted) {
                    Navigator.pushNamedAndRemoveUntil(
                        context, '/first', (_) => false);
                  }
                },
                icon: const Icon(Icons.logout_rounded, color: Colors.red),
                label: const Text('ออกจากระบบ',
                    style: TextStyle(color: Colors.red)),
                style: OutlinedButton.styleFrom(
                  side: const BorderSide(color: Colors.red),
                  padding: EdgeInsets.symmetric(vertical: context.nh(14)),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(context.radius(14)),
                  ),
                ),
              ),
            ),

            SizedBox(height: context.nh(32)),
          ],
        ),
      ),
    );
  }

  void _comingSoon(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Coming soon! 🐾')),
    );
  }
}

class _SectionLabel extends StatelessWidget {
  final String text;
  const _SectionLabel(this.text);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(bottom: context.nh(8)),
      child: Text(
        text,
        style: AppTextStyles.caption(context).copyWith(
          color: AppColors.textSecondary,
          fontWeight: FontWeight.w600,
          fontSize: context.nf(13),
        ),
      ),
    );
  }
}

class _SettingsGroup extends StatelessWidget {
  final List<SettingsTile> tiles;
  const _SettingsGroup({required this.tiles});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(context.radius(16)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.04),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        children: tiles.asMap().entries.map((e) {
          final isLast = e.key == tiles.length - 1;
          return Column(
            children: [
              e.value,
              if (!isLast)
                Divider(
                  height: 1,
                  color: AppColors.border,
                  indent: context.nw(52),
                ),
            ],
          );
        }).toList(),
      ),
    );
  }
}

class _Badge extends StatelessWidget {
  final String text;
  final bool isAlert;
  const _Badge(this.text, {this.isAlert = false});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.symmetric(
          horizontal: context.nw(8), vertical: context.nh(3)),
      decoration: BoxDecoration(
        color: isAlert
            ? Colors.orange.withValues(alpha: 0.12)
            : AppColors.primaryLight,
        borderRadius: BorderRadius.circular(context.radius(12)),
      ),
      child: Text(
        text,
        style: TextStyle(
          fontSize: context.nf(12),
          color: isAlert ? Colors.orange : AppColors.primary,
          fontWeight: FontWeight.w600,
        ),
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
                color: Colors.red, shape: BoxShape.circle),
          ),
        ),
      ],
    );
  }
}
