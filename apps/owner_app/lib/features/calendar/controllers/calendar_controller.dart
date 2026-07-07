import 'package:flutter/material.dart';
import '../../home/controllers/home_controller.dart';

class CalendarController extends ChangeNotifier {
  DateTime _focusedMonth = DateTime.now();
  DateTime get focusedMonth => _focusedMonth;

  DateTime _selectedDate = DateTime.now();
  DateTime get selectedDate => _selectedDate;

  int _selectedPetIndex = 0;
  int get selectedPetIndex => _selectedPetIndex;

  // Mock activities keyed by date string "yyyy-MM-dd"
  final Map<String, List<ActivityItem>> _activityMap = {};

  List<ActivityItem> get activities {
    final key = _dateKey(_selectedDate);
    return _activityMap[key] ?? [];
  }

  void init() {
    _loadMockData();
  }

  void prevMonth() {
    _focusedMonth = DateTime(_focusedMonth.year, _focusedMonth.month - 1, 1);
    notifyListeners();
  }

  void nextMonth() {
    _focusedMonth = DateTime(_focusedMonth.year, _focusedMonth.month + 1, 1);
    notifyListeners();
  }

  void prevYear() {
    _focusedMonth = DateTime(_focusedMonth.year - 1, _focusedMonth.month, 1);
    notifyListeners();
  }

  void nextYear() {
    _focusedMonth = DateTime(_focusedMonth.year + 1, _focusedMonth.month, 1);
    notifyListeners();
  }

  void setSelectedDate(DateTime date) {
    _selectedDate = date;
    notifyListeners();
  }

  void setSelectedPetIndex(int index) {
    _selectedPetIndex = index;
    notifyListeners();
  }

  bool hasActivityOnDate(DateTime date) {
    return _activityMap.containsKey(_dateKey(date)) &&
        _activityMap[_dateKey(date)]!.isNotEmpty;
  }

  String _dateKey(DateTime d) =>
      '${d.year}-${d.month.toString().padLeft(2, '0')}-${d.day.toString().padLeft(2, '0')}';

  void _loadMockData() {
    final now = DateTime.now();
    // Today activities
    final todayKey = _dateKey(now);
    _activityMap[todayKey] = [
      ActivityItem(
        icon: Icons.medical_services_rounded,
        title: 'ไปหาหมอ',
        subtitle: 'คลินิกรักสัตว์ สุขุมวิท · Jan 11 2026 – 10:00 AM',
        time: '10:00',
        color: const Color(0xFF66BB6A),
      ),
      ActivityItem(
        icon: Icons.vaccines_rounded,
        title: 'รับยาถ่ายพยาธิ',
        subtitle: 'นัดครั้งที่ 2',
        time: '14:00',
        color: const Color(0xFFFFA726),
      ),
    ];

    // Day +3 activities
    final futureDay = now.add(const Duration(days: 3));
    final futureKey = _dateKey(futureDay);
    _activityMap[futureKey] = [
      ActivityItem(
        icon: Icons.bathtub_rounded,
        title: 'อาบน้ำ-ตัดขน',
        subtitle: 'Grooming salon',
        time: '09:00',
        color: const Color(0xFF42A5F5),
      ),
    ];
    notifyListeners();
  }
}
