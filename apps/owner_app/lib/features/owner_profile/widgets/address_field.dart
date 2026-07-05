import 'package:flutter/material.dart';
import 'profile_text_field.dart';

class AddressField extends StatelessWidget {
  final TextEditingController controller;

  const AddressField({
    super.key,
    required this.controller,
  });

  @override
  Widget build(BuildContext context) {
    return ProfileTextField(
      controller: controller,
      label: 'ที่อยู่',
      hintText: 'กรอกข้อมูล',
      maxLines: 3,
    );
  }
}
