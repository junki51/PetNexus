import 'package:flutter/material.dart';
import '../../../shared/widgets/app_dropdown.dart';

class GenderDropdown extends StatelessWidget {
  final String? value;
  final ValueChanged<String?> onChanged;

  const GenderDropdown({
    super.key,
    required this.value,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return AppDropdown<String>(
      label: 'เพศ',
      value: value,
      items: const ['เลือกเพศ', 'ชาย', 'หญิง', 'ไม่ระบุเพศ'],
      onChanged: onChanged,
    );
  }
}
