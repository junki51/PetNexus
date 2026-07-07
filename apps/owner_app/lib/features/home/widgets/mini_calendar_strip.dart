import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';

class MiniCalendarStrip extends StatelessWidget {
  final DateTime selectedDate;
  final ValueChanged<DateTime> onDateSelected;

  const MiniCalendarStrip({
    super.key,
    required this.selectedDate,
    required this.onDateSelected,
  });

  @override
  Widget build(BuildContext context) {
    final today = DateTime.now();
    // Show 7 days centered on today
    final days = List.generate(
        7, (i) => today.subtract(Duration(days: 3 - i)));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: EdgeInsets.symmetric(horizontal: context.nw(20)),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'กิจกรรมวันนี้',
                style: AppTextStyles.body(context).copyWith(
                  fontWeight: FontWeight.bold,
                  fontSize: context.nf(15),
                ),
              ),
              Text(
                _monthYearThai(today),
                style: AppTextStyles.caption(context).copyWith(
                  color: AppColors.textSecondary,
                  fontSize: context.nf(12),
                  fontWeight: FontWeight.normal,
                ),
              ),
            ],
          ),
        ),
        SizedBox(height: context.nh(8)),
        SizedBox(
          height: context.nh(64),
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            padding: EdgeInsets.symmetric(horizontal: context.nw(16)),
            itemCount: days.length,
            itemBuilder: (ctx, i) {
              final day = days[i];
              final isSelected = _isSameDay(day, selectedDate);
              final isToday = _isSameDay(day, today);
              return GestureDetector(
                onTap: () => onDateSelected(day),
                child: AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  margin: EdgeInsets.symmetric(horizontal: context.nw(4)),
                  width: context.nw(40),
                  decoration: BoxDecoration(
                    color: isSelected
                        ? AppColors.primary
                        : isToday
                            ? AppColors.primaryLight
                            : Colors.transparent,
                    borderRadius:
                        BorderRadius.circular(context.radius(12)),
                    border: isToday && !isSelected
                        ? Border.all(color: AppColors.primary, width: 1.5)
                        : null,
                  ),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        _dayThai(day),
                        style: TextStyle(
                          fontSize: context.nf(10),
                          color: isSelected
                              ? Colors.white
                              : AppColors.textSecondary,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      SizedBox(height: context.nh(2)),
                      Text(
                        '${day.day}',
                        style: TextStyle(
                          fontSize: context.nf(16),
                          fontWeight: FontWeight.bold,
                          color: isSelected
                              ? Colors.white
                              : AppColors.textPrimary,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  bool _isSameDay(DateTime a, DateTime b) =>
      a.year == b.year && a.month == b.month && a.day == b.day;

  String _dayThai(DateTime d) {
    const days = ['อา', 'จ', 'อ', 'พ', 'พฤ', 'ศ', 'ส'];
    return days[d.weekday % 7];
  }

  String _monthYearThai(DateTime d) {
    const months = [
      'ม.ค.', 'ก.พ.', 'มี.ค.', 'เม.ย.', 'พ.ค.', 'มิ.ย.',
      'ก.ค.', 'ส.ค.', 'ก.ย.', 'ต.ค.', 'พ.ย.', 'ธ.ค.'
    ];
    return '${months[d.month - 1]} ${d.year + 543}';
  }
}
