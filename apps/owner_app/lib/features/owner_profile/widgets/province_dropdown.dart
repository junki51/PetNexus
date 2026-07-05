import 'package:flutter/material.dart';
import '../../../shared/widgets/app_dropdown.dart';

class ProvinceDropdown extends StatelessWidget {
  final String? value;
  final ValueChanged<String?> onChanged;

  const ProvinceDropdown({
    super.key,
    required this.value,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return AppDropdown<String>(
      label: 'จังหวัด',
      value: value,
      items: const [
        'เลือกจังหวัด',
        'กรุงเทพมหานคร',
        'นนทบุรี',
        'ปทุมธานี',
        'สมุทรปราการ',
        'เชียงใหม่',
        'เชียงราย',
        'ภูเก็ต',
        'ชลบุรี',
        'นครราชสีมา',
        'ขอนแก่น',
        'สงขลา',
        'พิษณุโลก',
        'สุราษฎร์ธานี',
        'ประจวบคีรีขันธ์',
      ],
      onChanged: onChanged,
    );
  }
}
