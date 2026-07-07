import 'package:flutter/material.dart';
import '../../../app/app_routes.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_card.dart';
import '../../../shared/widgets/app_logo.dart';
import '../../../shared/widgets/app_scaffold.dart';

class PetSuccessScreen extends StatelessWidget {
  const PetSuccessScreen({super.key});

  // Mock peeking dog illustration (User can replace with local asset later)
  final String _peekingDogUrl = 'https://images.unsplash.com/photo-1543466835-00a7907e9de1?w=200';

  void _onContinue(BuildContext context) {
    Navigator.pushNamedAndRemoveUntil(context, AppRoutes.home, (route) => false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      scrollable: false,
      backgroundColor: AppColors.primaryLight,
      child: Column(
        children: [
          // Logo & Badge Area
          Expanded(
            flex: 4,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const AppLogo(),
                SizedBox(height: context.nh(24)),
                // Verified scalloped badge icon
                Icon(
                  Icons.verified,
                  color: AppColors.primary,
                  size: context.nw(120),
                ),
              ],
            ),
          ),
          // Content Card
          Expanded(
            flex: 5,
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
                padding: EdgeInsets.fromLTRB(
                  context.nw(32),
                  context.nh(40),
                  context.nw(32),
                  context.nh(32),
                ),
                child: Column(
                  children: [
                    Text(
                      'กรอกข้อมูลเสร็จสิ้น!',
                      style: AppTextStyles.title(context).copyWith(
                        fontSize: context.nf(26),
                        fontWeight: FontWeight.bold,
                        color: AppColors.primary,
                      ),
                    ),
                    SizedBox(height: context.nh(16)),
                    Text(
                      'ขอบคุณที่ไว้วางใจให้ PetNexus\nดูแลเพื่อนรักตัวน้อยของคุณ',
                      textAlign: TextAlign.center,
                      style: AppTextStyles.body(context).copyWith(
                        color: AppColors.textPrimary,
                        fontSize: context.nf(16),
                        height: 1.4,
                      ),
                    ),
                    const Spacer(),
                    // Peeking puppy illustration
                    Container(
                      width: context.nw(120),
                      height: context.nw(80),
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.only(
                          topLeft: Radius.circular(60),
                          topRight: Radius.circular(60),
                        ),
                        image: DecorationImage(
                          image: NetworkImage(_peekingDogUrl),
                          fit: BoxFit.cover,
                        ),
                        border: const Border(
                          bottom: BorderSide(color: Colors.black, width: 2),
                        ),
                      ),
                    ),
                    SizedBox(height: context.nh(24)),
                    // Continue Button
                    AppButton.primary(
                      text: 'กดเพื่อไปต่อ',
                      icon: Icons.pets,
                      onPressed: () => _onContinue(context),
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
