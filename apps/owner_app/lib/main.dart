import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'features/auth/controllers/auth_controller.dart';
import 'features/auth/screens/auth_gate.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => AuthController(),
        ),
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

      home: const AuthGate(),

      routes: {
        "/auth": (_) => const AuthGate(),
        // "/login": (_) => const LoginScreen(),
        // "/register": (_) => const RegisterScreen(),
        // "/home": (_) => const HomeScreen(),
      },
    );
  }
}
