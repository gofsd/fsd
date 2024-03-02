import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'allIn.dart';

class Translation {
  Map<String, Map<String, String>> langMap;
  Map<String, Map<String, String>> newLangMap;
  Map<String, Map<String, String>> toTranslate =
      Map<String, Map<String, String>>();
  String from = 'en';
  Map<String, String> apiEx = {"iw": "he"};
  Map<String, bool> skip = {
    "el": true,
    "he": true,
    "fi": true,
    "hi": true,
    "th": true
  };

  void setMap() {
    langMap = map;
    newLangMap = langMap;
  }

  Future<File> writeNewMap() async {
    final file = await _localFile;

    // Write the file
    return file.writeAsString("Map<String, Map<String, String>> map = " +
        jsonEncode(newLangMap) +
        ";");
  }

  Future<String> get _localPath async {
    final directory = "./";

    return directory;
  }

  Future<File> get _localFile async {
    final path = await _localPath;
    return File('$path/newAllIn.dart');
  }

  void prepareNewMap() async {
    Map<String, String> rootLang;
    for (var langMap in newLangMap.entries) {
      if (langMap.key == from) {
        rootLang = langMap.value;
      }
    }

    for (var langMap in newLangMap.entries) {
      var value = langMap.value;
      for (var langItem in rootLang.entries) {
        if (value[langItem.key] == null) {
          value[langItem.key] = "";
        }
      }
    }

    for (var langMap in newLangMap.entries) {
      if (skip[langMap.key] != null && skip[langMap.key]) {
        continue;
      }
      var value = langMap.value;
      for (var langItem in rootLang.entries) {
        if (value[langItem.key] == "" &&
            langItem.value != "" &&
            langMap.key != from) {
          if (toTranslate[langMap.key] == null) {
            toTranslate[langMap.key] = Map();
          }
          toTranslate[langMap.key][langItem.key] = langItem.value;
          //value[langItem.key] = await getTranslate(langItem.value);
        }
      }
    }

    for (var itemsByLang in toTranslate.entries) {
      String lang = itemsByLang.key;
      if (apiEx[itemsByLang.key] != null) {
        lang = apiEx[itemsByLang.key];
      }
      Map<String, String> translChunk = Map();
      for (var itToTranslate in itemsByLang.value.entries) {
        if (translChunk.length < 9) {
          translChunk[itToTranslate.key] = itToTranslate.value;
        } else {
          await getTranslate(lang, translChunk, itemsByLang.key);
          translChunk = Map();
        }
      }

      if (translChunk.length > 0) {
        await getTranslate(lang, translChunk, itemsByLang.key);
        translChunk = Map();
      }
    }
    // getTranslate("hello");
  }

  Future<String> getTranslate(
      String lang, Map<String, String> toTranslate, String lk) async {
    var client = HttpClient();

    HttpClientRequest request = await client.postUrl(Uri.parse(
        'https://microsoft-translator-text.p.rapidapi.com/translate?to=$lang&api-version=3.0&profanityAction=NoAction&textType=plain&from=$from'));
    request.headers
        .add(HttpHeaders.contentTypeHeader, "application/json; charset=utf-8");
    request.headers
        .add("x-rapidapi-host", "microsoft-translator-text.p.rapidapi.com");
    request.headers.add(
        "x-rapidapi-key", "81b08e4222msh4a977efabc8cf53p105e20jsn720a3c104861");
    List<Map<String, String>> tr = List();
    for (var t in toTranslate.entries) {
      Map<String, String> m = Map();
      m["Text"] = t.value;
      tr.add(m);
    }
    request.write(jsonEncode(tr));
    // [{"detectedLanguage":{"language":"en","score":1.0},"translations":[{"text":"Мне бы очень хотелось несколько раз проехать на вашей машине по кварталу.","to":"ru"}]},{"detectedLanguage":{"language":"en","score":1.0},"translations":[{"text":"Всем привет.","to":"ru"}]}]
    var response = await request.close();
    final stringData = await response.transform(utf8.decoder).join();
    print(stringData);

    List<dynamic> respDecoded = jsonDecode(stringData);
    print(stringData);
    List<MapEntry<String, String>> l = toTranslate.entries.toList();
    for (var i = 0; i < l.length; i++) {
      List<dynamic> entry = respDecoded[i]['translations'] as List<dynamic>;

      newLangMap[lk][l[i].key] = entry[0]["text"];
    }
    return "";
  }

  void run() async {
    setMap();
    try {
      await prepareNewMap();
    } catch (e) {
      print(e.toString());
    }

    writeNewMap();
  }
}

void main() {
  Translation().run();
}
