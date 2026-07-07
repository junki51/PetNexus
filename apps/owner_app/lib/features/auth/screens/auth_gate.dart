import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../app/app_routes.dart';
import '../controllers/auth_controller.dart';
import '../../owner_profile/controllers/owner_profile_controller.dart';

class AuthGate extends StatefulWidget {
  const AuthGate({super.key});

  @override
  State<AuthGate> createState() => _AuthGateState();
}

class _AuthGateState extends State<AuthGate> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) async {
      final authController = context.read<AuthController>();
      final ownerController = context.read<OwnerProfileController>();

      final isLoggedIn = await authController.checkAuthentication();

      if (!mounted) return;

      if (!isLoggedIn) {
        Navigator.pushReplacementNamed(context, AppRoutes.first);
        return;
      }

      // Check if owner profile exists
      await ownerController.fetchProfile();

      if (!mounted) return;

      if (ownerController.profile == null) {
        // Logged in but no profile → complete profile first
        Navigator.pushReplacementNamed(context, AppRoutes.completeProfile);
      } else {
        // Fully onboarded → go to main shell
        Navigator.pushReplacementNamed(context, AppRoutes.main);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(child: CircularProgressIndicator()),
    );
  }
}
