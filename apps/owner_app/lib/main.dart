import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'app/app_routes.dart';
import 'app/main_shell.dart';
import 'features/auth/controllers/auth_controller.dart';
import 'features/auth/screens/auth_gate.dart';
import 'features/auth/screens/first_screen.dart';
import 'features/auth/screens/login_screen.dart';
import 'features/auth/screens/register_screen.dart';
import 'features/calendar/controllers/calendar_controller.dart';
import 'features/home/controllers/home_controller.dart';
import 'features/owner_profile/controllers/owner_profile_controller.dart';
import 'features/owner_profile/screens/owner_profile.dart';
import 'features/pet/controllers/pet_controller.dart';
import 'features/pet/screens/select_pet_screen.dart';
import 'features/pet/screens/pet_info_form_screen.dart';
import 'features/pet/screens/pet_upload_photo_screen.dart';
import 'features/pet/screens/pet_success_screen.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthController()),
        ChangeNotifierProvider(create: (_) => OwnerProfileController()),
        ChangeNotifierProvider(create: (_) => PetController()),
        ChangeNotifierProvider(create: (_) => HomeController()),
        ChangeNotifierProvider(create: (_) => CalendarController()),
      ],
      child: const MyApp(),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: "PetNexus",
      // Start at AuthGate — checks token, then routes appropriately
      home: const AuthGate(),
      routes: {
        AppRoutes.auth: (_) => const AuthGate(),
        AppRoutes.first: (_) => const FirstScreen(),
        AppRoutes.login: (_) => const LoginScreen(),
        AppRoutes.register: (_) => const RegisterScreen(),
        AppRoutes.home: (_) => const OwnerProfileScreen(),
        AppRoutes.completeProfile: (_) => const OwnerProfileScreen(),
        AppRoutes.main: (_) => const MainShell(),
        AppRoutes.selectPet: (_) => const SelectPetScreen(),
        AppRoutes.petInfoForm: (_) => const PetInfoFormScreen(),
        AppRoutes.petUploadPhoto: (_) => const PetUploadPhotoScreen(),
        AppRoutes.petSuccess: (_) => const PetSuccessScreen(),
      },
    );
  }
}
