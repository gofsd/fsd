var unirest = require("unirest");
var _ = require('lodash');

var fs = require("fs");
translations = [];

var langs = {
              "de_DE": {
                "lang": "German",
                "region": "Germany"
              },
              "zh_CN": {
                "lang": "Chinese",
                "region": "PRC"
              },
              "cs_CZ": {
                "lang": "Czech",
                "region": "Czech Republic"
              },
              "nl_BE": {
                "lang": "Dutch",
                "region": "Belgium"
              },
              "fr_FR": {
                "lang": "French",
                "region": "France"
              },
              "it_IT": {
                "lang": "Italian",
                "region": "Italy"
              },
              "ja_JP": {
                "lang": "Japanese",
                "region": "Japanese"
              },
              "ko_KR": {
                "lang": "Korean",
                "region": "Korean"
              },
              "pl_PL": {
                "lang": "Polish",
                "region": "Polish"
              },
              "ru_RU": {
                "lang": "Russian",
                "region": "Russian"
              },
              "es_ES": {
                "lang": "Spanish",
                "region": "Spanish"
              },
              "bg_BG": {
                "lang": "Bulgarian",
                "region": "Bulgaria"
              },
              "hr_HR": {
                "lang": "Croatian",
                "region": "Croatia"
              },
              "da_DK": {
                "lang": "Danish",
                "region": "Denmark"
              },
              "fi_FI": {
                "lang": "Finnish",
                "region": "Finland"
              },
              "el_GR": {
                "lang": "Greek",
                "region": "Greece"
              },
              "he_IL": {
                "lang": "Hebrew",
                "region": "Israel"
              },
              "hi_IN": {
                "lang": "Hindi",
                "region": "India"
              },
              "hu_HU": {
                "lang": "Hungarian",
                "region": "Hungary"
              },
              "id_ID": {
                "lang": "Indonesian",
                "region": "Indonesia"
              },
              "lv_LV": {
                "lang": "Latvian",
                "region": "Latvia"
              },
              "lt_LT": {
                "lang": "Lithuanian",
                "region": "Lithuania"
              },
              "nb_NO": {
                "lang": "Norwegian-Bokmål",
                "region": "Norway"
              },
              "pt_PT": {
                "lang": "Portuguese",
                "region": "Portugal"
              },
              "ro_RO": {
                "lang": "Romanian",
                "region": "Romania"
              },
              "sr_RS": {
                "lang": "Serbian",
                "region": "Serbian"
              },
              "sk_SK": {
                "lang": "Slovak",
                "region": "Slovakia"
              },
              "sl_SI": {
                "lang": "Slovenian",
                "region": "Slovenia"
              },
              "sv_SE": {
                "lang": "Swedish",
                "region": "Sweden"
              },
              "tl_PH": {
                "lang": "Tagalog",
                "region": "Philippines"
              },
              "th_TH": {
                "lang": "Thai",
                "region": "Thailand"
              },
              "tr_TR": {
                "lang": "Turkish",
                "region": "Turkey"
              },
              "vi_VN": {
                "lang": "Vietnamese",
                "region": "Vietnam"
              },
              "uk_UA": {
                "lang": "Ukrainian",
                "region": "Ukraine"
              },
            };

fs.readFile('./input_data/gen_from_parse_ox5000_and_wordCol.json', 'utf8', async function (err,str) {
    ox5000withWordCollection = JSON.parse(str);
    objKeys = Object.keys(langs);
    for(var keyIdx = 0; keyIdx < objKeys.length; keyIdx++) {
        lang = objKeys[keyIdx].split("_")[0];
        console.log("current lang", lang)
        await (new Promise((r) => {
                    fs.readFile(`./input_data/lookup/${lang}.json`, 'utf8',async function (err,str) {
                        if(err && err.ENOENT == fs.ENOENT) {
                            fs.writeFile(`./input_data/lookup/${lang}.json`, JSON.stringify([]), "utf8", (err)=> console.log(err));
                        }
                        translations = err ? [] : JSON.parse(str);
                        while(true){
                                        translations = _.uniqBy(translations, "normalizedSource");
                                        mapToTranslate = [];
                                        for (var i = 0; i < 10; i++) {
                                            itemToTranslate = ox5000withWordCollection.find((transl, idx) => {
                                               translated = translations.find(item => {
                                                    return String(item.normalizedSource).toLowerCase() == String(transl.name).toLowerCase();
                                               });

                                               excluded = mapToTranslate.find(item => {
                                                    return String(transl.name).toLowerCase() == String(item.name).toLowerCase();
                                               });
                                               if (translated == undefined && excluded == undefined) {
                                                   return true;
                                               } else {
                                                   return false;
                                               }
                                            });
                                            itemToTranslate ? mapToTranslate.push(itemToTranslate) : null;
                                        }

                                        preparedReqData = mapToTranslate.map((item) => {
                                            return {
                                                "Text": item.name
                                            };
                                        });
                                        if(preparedReqData.length == 0){
                                            break;
                                        }
                                        var req = unirest("POST", "https://microsoft-translator-text.p.rapidapi.com/Dictionary/Lookup");

                                        req.headers({
                                            "content-type": "application/json",
                                            "x-rapidapi-key": "81b08e4222msh4a977efabc8cf53p105e20jsn720a3c104861",
                                            "x-rapidapi-host": "microsoft-translator-text.p.rapidapi.com",
                                            "useQueryString": true
                                        });

                                        req.type("json");
                                        req.query({
                                            "from": "en",
                                            "api-version": "3.0",
                                            "to": lang
                                        });
                                        console.log(preparedReqData, translations.length, ox5000withWordCollection.length);
                                        req.send(preparedReqData);
                                        await (new Promise((resolve, reject)=> {
                                                    req.end(function (res) {
                                                    	if (res.error) reject(res.error);
                                                    	console.log(res.body, 'from  last promise')
                                                    	translations = [...translations, ...res.body];
                                                    	resolve()
                                                    	console.log(res.body);
                                                    });
                                        }));
                        }
                                        fs.writeFile(`./input_data/lookup/${lang}.json`, JSON.stringify(translations), "utf8", r);


                    });
        }));
    }
});

