import 'dart:io' as Platrofm;

import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';

import 'package:get/get.dart';

import 'package:frontend/getx/controllers/settings.dart';

import '../db/cache.dart';
import 'allIn.dart';
import 'en_us.dart';
import 'pt_br.dart';
import 'ru_ru.dart';

class TranslationService extends Translations {
  static String get locale => Get.deviceLocale.toString().split("_")[0];
  static SettingsController set = Get.find<SettingsController>();

  // Get.updateLocale(Locale('ru', 'RU'));
  static Locale getLocale() {
    String lang;
    if (set.set != null && set.set.lang != null) {
      lang = set.set.lang;
    } else {
      lang = TranslationService.locale;
    }
    return Locale(lang);
  }

  static Locale setLocale(String lang) {
    set.setLang(lang);
    var locale = Locale(lang);
    Get.updateLocale(locale);
    return Locale(lang);
  }

  @override
  Map<String, Map<String, String>> get keys => map;
}
